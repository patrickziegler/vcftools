package vcf

import (
	"bytes"
	"strings"
	"testing"
)

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

func TestNormalizePhone(t *testing.T) {
	cases := map[string]string{
		"+49 176-8112 3123": "+4917681123123",
		"0176 811 23123":    "+4917681123123",
		"(0176) 811/23123":  "+4917681123123",
		"0049 176 81123123": "+4917681123123",
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
	in := "BEGIN:VCARD\nVERSION:3.0\nFN:Thea Kummeth\nN:Kummeth;Thea;;;\n" +
		"item1.TEL:+49 171 5400000\nitem1.X-ABLabel:Мобільний\n" +
		"TEL;TYPE=X-Мобільний;TYPE=PREF:+491715400000\nCATEGORIES:myContacts\nEND:VCARD\n"
	want := "BEGIN:VCARD\nVERSION:3.0\nFN:Thea Kummeth\nN:Kummeth;Thea;;;\nTEL:+491715400000\nEND:VCARD\n"
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
