package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	in := "BEGIN:VCARD\nVERSION:3.0\nFN:Jane Smith\nN:Smith;Jane;;;\n" +
		"item1.TEL:+49 171 1234567\nitem1.X-ABLabel:Мобільний\n" +
		"TEL;TYPE=X-Мобільний;TYPE=PREF:+491711234567\nCATEGORIES:myContacts\nEND:VCARD\n" +
		"BEGIN:VCARD\nVERSION:3.0\nORG:nothing\nEND:VCARD\n"

	var out, errb bytes.Buffer
	if code := run([]string{"--prefix", "+43"}, strings.NewReader(in), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	want := "BEGIN:VCARD\nVERSION:3.0\nFN:Jane Smith\nN:Smith;Jane;;;\nTEL:+491711234567\nEND:VCARD\n"
	if out.String() != want {
		t.Errorf("stdout:\n%s\nwant:\n%s", out.String(), want)
	}
	if got := errb.String(); got != "skipping contact without name and phone number\nTotal processed contacts: 2\n" {
		t.Errorf("stderr=%q", got)
	}

	// Unexpected positional arguments should fail
	out.Reset()
	errb.Reset()
	if code := run([]string{"file.vcf"}, strings.NewReader(in), &out, &errb); code != 2 {
		t.Errorf("expected exit code 2 with positional argument, got %d", code)
	}

	// Default run with no flags should succeed
	out.Reset()
	errb.Reset()
	if code := run(nil, strings.NewReader(in), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
}
