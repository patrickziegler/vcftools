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
	in := "BEGIN:VCARD\r\nFN:Ahmed\r\n  Kashef\r\nTEL:1\r\nEND:VCARD\r\n"
	if full, _, _ := parse(t, in)[0].Names(); full != "Ahmed Kashef" {
		t.Errorf("name=%q", full)
	}
}

func TestNamesRaw(t *testing.T) {
	c := parse(t, "BEGIN:VCARD\nFN:Ahmed Kashef\nEND:VCARD\n")[0]
	if full, fam, giv := c.Names(); full != "Ahmed Kashef" || fam != "Kashef" || giv != "Ahmed" {
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
