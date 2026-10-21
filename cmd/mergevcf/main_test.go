package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()
	card := "BEGIN:VCARD\nVERSION:3.0\nFN:Ahmed Kashef\nN:Kashef;Ahmed;;;\nTEL;TYPE=CELL:0176 811 23123\nCATEGORIES:myContacts\nEND:VCARD\n"
	f1, f2 := filepath.Join(dir, "1.vcf"), filepath.Join(dir, "2.vcf")
	os.WriteFile(f1, []byte(card), 0o644)
	os.WriteFile(f2, []byte(card), 0o644)

	var out, errb bytes.Buffer
	if code := run([]string{f1, f2}, &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	if out.String() != card {
		t.Errorf("stdout should be unmodified:\n%s", out.String())
	}
	if errb.String() != "Total unique contacts: 1\n" {
		t.Errorf("stderr=%q", errb.String())
	}
	if run(nil, &out, &errb) == 0 {
		t.Error("expected failure without files")
	}
}
