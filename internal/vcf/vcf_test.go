package vcf

import (
	"bytes"
	"strings"
	"testing"
)

const qpCard = "BEGIN:VCARD\nVERSION:2.1\n" +
	"N;CHARSET=UTF-8;ENCODING=QUOTED-PRINTABLE:=4D=C3=BC=6C=6C=65=72;=53=74=65=6C=6C=61;;;\n" +
	"FN;CHARSET=UTF-8;ENCODING=QUOTED-PRINTABLE:=53=74=65=6C=6C=61=20=4D=C3=BC=\n=6C=6C=65=72\n" +
	"TEL;CELL:+49172342342\nORG:\nEND:VCARD\n"

func parse(t *testing.T, s string) []*Card {
	t.Helper()
	cards, err := Parse(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}
	return cards
}

func render(cards []*Card) string {
	var b bytes.Buffer
	for _, c := range cards {
		c.Write(&b)
	}
	return b.String()
}

func cleanOne(t *testing.T, in string) string {
	t.Helper()
	cards, err := Parse(strings.NewReader(in))
	if err != nil || len(cards) != 1 {
		t.Fatalf("parse: %v (%d cards)", err, len(cards))
	}
	c, ok := Clean(cards[0], "+49")
	if !ok {
		return ""
	}
	var b bytes.Buffer
	c.Write(&b)
	return b.String()
}

func TestParseRoundTripPreservesRaw(t *testing.T) {
	in := qpCard + "BEGIN:VCARD\nVERSION:3.0\nFN:A B\nitem1.TEL:1\nitem1.X-ABLabel:x\nEND:VCARD\n"
	cards := parse(t, in)
	if got := render(cards); got != in {
		t.Errorf("round trip differs:\n%s", got)
	}
	if p, ok := cards[1].Get("TEL"); !ok || p.Group != "item1" {
		t.Errorf("group not parsed: %+v", p)
	}
}

func TestParseCRLFAndFolding(t *testing.T) {
	in := "BEGIN:VCARD\r\nFN:John\r\n  Doe\r\nTEL:1\r\nEND:VCARD\r\n"
	if full, _, _ := parse(t, in)[0].Names(); full != "John Doe" {
		t.Errorf("name=%q", full)
	}
}

func TestNamesRaw(t *testing.T) {
	c := parse(t, "BEGIN:VCARD\nFN:John Doe\nEND:VCARD\n")[0]
	if full, fam, giv := c.Names(); full != "John Doe" || fam != "Doe" || giv != "John" {
		t.Errorf("names: %q %q %q", full, fam, giv)
	}
}

func TestMerge(t *testing.T) {
	a := qpCard + "BEGIN:VCARD\nVERSION:3.0\nFN:Bob\nTEL:+49151111\nEND:VCARD\n"
	b := "BEGIN:VCARD\nVERSION:3.0\nFN:Other Name\nTEL:+49172342342\nEND:VCARD\n" +
		"BEGIN:VCARD\nVERSION:3.0\nFN:bob\nTEL:+49151222\nEND:VCARD\n" +
		"BEGIN:VCARD\nVERSION:3.0\nFN:Aaron Zed\nN:Zed;Aaron;;;\nEND:VCARD\n"
	m := Merge(append(parse(t, a), parse(t, b)...))
	if len(m) != 3 {
		t.Fatalf("got %d cards", len(m))
	}
	SortByLastName(m)
	got := render(m)
	if !strings.Contains(got, "FN:Bob\nTEL:+49151111\nTEL:+49151222\nEND") {
		t.Errorf("phones not combined:\n%s", got)
	}
	if !strings.Contains(got, qpCard[len("BEGIN:VCARD\n"):len(qpCard)-len("END:VCARD\n")]) {
		t.Errorf("raw QP card modified:\n%s", got)
	}
	if !strings.HasSuffix(got, "FN:Aaron Zed\nN:Zed;Aaron;;;\nEND:VCARD\n") {
		t.Errorf("sort by last name failed:\n%s", got)
	}
}

func TestNormalizePhone(t *testing.T) {
	cases := map[string]string{
		"+49 176-1234 5678": "+4917612345678",
		"0176 1234 5678":    "+4917612345678",
		"(0176) 1234/5678":  "+4917612345678",
		"0049 176 12345678": "+4917612345678",
		"+1 (555) 123-4567": "+15551234567",
		"":                  "",
	}
	for in, want := range cases {
		if got := NormalizePhone(in, "+49"); got != want {
			t.Errorf("NormalizePhone(%q)=%q want %q", in, got, want)
		}
	}
	if got := NormalizePhone("0123", "+43"); got != "+43123" {
		t.Errorf("custom prefix: got %q", got)
	}
}

func TestDecodeLatin1(t *testing.T) {
	if got := DecodeValue("", "M\xfcller"); got != "Müller" {
		t.Errorf("got %q", got)
	}
}

func TestCleanQuotedPrintable(t *testing.T) {
	want := "BEGIN:VCARD\nVERSION:3.0\nFN:Stella Müller\nN:Müller;Stella;;;\nTEL;TYPE=CELL:+49172342342\nEND:VCARD\n"
	if got := cleanOne(t, qpCard); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestCleanAppleGroupedTel(t *testing.T) {
	in := "BEGIN:VCARD\nVERSION:3.0\nFN:Jane Smith\nN:Smith;Jane;;;\n" +
		"item1.TEL:+49 171 1234567\nitem1.X-ABLabel:Мобільний\n" +
		"TEL;TYPE=X-Мобільний;TYPE=PREF:+491711234567\nCATEGORIES:myContacts\nEND:VCARD\n"
	want := "BEGIN:VCARD\nVERSION:3.0\nFN:Jane Smith\nN:Smith;Jane;;;\nTEL:+491711234567\nEND:VCARD\n"
	if got := cleanOne(t, in); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestCleanAppleLabelTypeAndPhoneInLabel(t *testing.T) {
	in := "BEGIN:VCARD\nVERSION:3.0\nFN:A B\nN:B;A;;;\n" +
		"item1.TEL:0171 111\nitem1.X-ABLabel:_$!<Mobile>!$_\n" +
		"item2.X-ABLabel:+49 172 222\nitem3.X-ABLabel:Mother\nEND:VCARD\n"
	want := "BEGIN:VCARD\nVERSION:3.0\nFN:A B\nN:B;A;;;\nTEL;TYPE=CELL:+49171111\nTEL:+49172222\nEND:VCARD\n"
	if got := cleanOne(t, in); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestCleanEmptyCard(t *testing.T) {
	if got := cleanOne(t, "BEGIN:VCARD\nVERSION:3.0\nORG:x\nEND:VCARD\n"); got != "" {
		t.Errorf("expected rejection, got %q", got)
	}
}
