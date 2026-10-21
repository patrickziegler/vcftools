package vcf

import (
	"io"
	"mime/quotedprintable"
	"strings"
	"unicode"
	"unicode/utf8"
)

// DecodeValue converts a raw property value to UTF-8. Quoted-printable
// values are decoded; bytes that are not valid UTF-8 are treated as Latin-1.
func DecodeValue(params, val string) string {
	if strings.Contains(strings.ToUpper(params), "QUOTED-PRINTABLE") {
		b, err := io.ReadAll(quotedprintable.NewReader(strings.NewReader(val)))
		if err == nil || len(b) > 0 {
			val = string(b)
		}
	}
	if !utf8.ValidString(val) {
		r := make([]rune, len(val))
		for i := 0; i < len(val); i++ {
			r[i] = rune(val[i])
		}
		return string(r)
	}
	return val
}

// NormalizePhone strips formatting characters and converts a leading 0 into
// the given country prefix (e.g. "+49"). A "00" international prefix becomes "+".
func NormalizePhone(s, prefix string) string {
	var b strings.Builder
	for i, r := range strings.TrimSpace(s) {
		switch {
		case unicode.IsDigit(r):
			b.WriteRune(r)
		case r == '+' && i == 0:
			b.WriteRune(r)
		}
	}
	n := b.String()
	switch {
	case strings.HasPrefix(n, "00"):
		n = "+" + n[2:]
	case strings.HasPrefix(n, "0"):
		n = prefix + n[1:]
	}
	return n
}
