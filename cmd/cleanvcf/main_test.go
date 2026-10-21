package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	in := "BEGIN:VCARD\nVERSION:3.0\nFN:Thea Kummeth\nN:Kummeth;Thea;;;\n" +
		"item1.TEL:+49 171 5400000\nitem1.X-ABLabel:Мобільний\n" +
		"TEL;TYPE=X-Мобільний;TYPE=PREF:+491715400000\nCATEGORIES:myContacts\nEND:VCARD\n" +
		"BEGIN:VCARD\nVERSION:3.0\nORG:nothing\nEND:VCARD\n"
	path := filepath.Join(t.TempDir(), "in.vcf")
	os.WriteFile(path, []byte(in), 0o644)

	var out, errb bytes.Buffer
	if code := run([]string{"--prefix", "+43", path}, &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	want := "BEGIN:VCARD\nVERSION:3.0\nFN:Thea Kummeth\nN:Kummeth;Thea;;;\nTEL:+491715400000\nEND:VCARD\n"
	if out.String() != want {
		t.Errorf("stdout:\n%s\nwant:\n%s", out.String(), want)
	}
	if got := errb.String(); got != "skipping contact without name and phone number\nTotal processed contacts: 2\n" {
		t.Errorf("stderr=%q", got)
	}
	if run(nil, &out, &errb) == 0 {
		t.Error("expected failure without file")
	}
}
