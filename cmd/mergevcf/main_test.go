package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	card := "BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nN:Doe;John;;;\nTEL;TYPE=CELL:0176 1234 5678\nCATEGORIES:myContacts\nEND:VCARD\n"
	in := card + card

	var out, errb bytes.Buffer
	if code := run(nil, strings.NewReader(in), &out, &errb); code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	if out.String() != card {
		t.Errorf("stdout should be unmodified:\n%s", out.String())
	}
	if errb.String() != "Total unique contacts: 1\n" {
		t.Errorf("stderr=%q", errb.String())
	}

	// Positional arguments are no longer accepted
	out.Reset()
	errb.Reset()
	if code := run([]string{"file.vcf"}, strings.NewReader(in), &out, &errb); code != 2 {
		t.Errorf("expected exit code 2 with positional argument, got %d", code)
	}
}
