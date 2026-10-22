package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	in := "BEGIN:VCARD\nVERSION:3.0\nFN:Patrick Ziegler\nN:Ziegler;Patrick;;;\nTEL;TYPE=CELL:0176 99999\nEND:VCARD\n" +
		"BEGIN:VCARD\nVERSION:3.0\nFN:Ahmed Kashef\nN:Kashef;Ahmed;;;\nTEL;TYPE=CELL:0176 811 23123\nTEL;TYPE=WORK:+49 89 12345\nEND:VCARD\n" +
		"BEGIN:VCARD\nVERSION:3.0\nFN:Thea Kummeth\nN:Kummeth;Thea;;;\nEND:VCARD\n"

	var out, errb bytes.Buffer
	if code := run(nil, strings.NewReader(in), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}

	want := "Kashef;Ahmed;0176 811 23123;+49 89 12345\n" +
		"Kummeth;Thea\n" +
		"Ziegler;Patrick;0176 99999\n"

	if out.String() != want {
		t.Errorf("stdout:\n%s\nwant:\n%s", out.String(), want)
	}

	if got := errb.String(); got != "Total processed contacts: 3\n" {
		t.Errorf("stderr=%q", got)
	}

	// Positional arguments are rejected
	out.Reset()
	errb.Reset()
	if code := run([]string{"contacts.vcf"}, strings.NewReader(in), &out, &errb); code != 2 {
		t.Errorf("expected exit code 2 with positional argument, got %d", code)
	}

	// Empty input
	out.Reset()
	errb.Reset()
	if code := run(nil, strings.NewReader(""), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	if out.String() != "" {
		t.Errorf("expected empty stdout, got %q", out.String())
	}
	if got := errb.String(); got != "Total processed contacts: 0\n" {
		t.Errorf("stderr=%q", got)
	}
}
