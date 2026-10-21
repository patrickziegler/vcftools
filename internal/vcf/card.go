// Package vcf handles vCards. Parsing, writing, merging and sorting never
// change property values; cleaning (clean.go, decode.go) decodes to UTF-8,
// normalizes phone numbers and rebuilds minimal vCard 3.0 contacts.
package vcf

import (
	"io"
	"strings"
)

// Prop is a single (unfolded) vCard property line. Raw keeps the original
// physical line(s) so that cards can be written back byte-for-byte.
type Prop struct {
	Group  string // e.g. "item1" for "item1.TEL"
	Name   string // upper-cased property name, e.g. "TEL"
	Params string // raw parameter string without leading ";"
	Value  string // raw (still encoded) value
	Raw    []string
}

// Card is a vCard without its BEGIN/END lines.
type Card struct {
	Props []Prop
}

// Get returns the first property with the given name.
func (c *Card) Get(name string) (Prop, bool) {
	for _, p := range c.Props {
		if p.Name == name {
			return p, true
		}
	}
	return Prop{}, false
}

// All returns all properties with the given name.
func (c *Card) All(name string) []Prop {
	var out []Prop
	for _, p := range c.Props {
		if p.Name == name {
			out = append(out, p)
		}
	}
	return out
}

// Names returns the full, family and given name using the raw property values.
func (c *Card) Names() (full, family, given string) {
	return c.NamesWith(func(p Prop) string { return p.Value })
}

// NamesWith is like Names but reads property values through text, which
// allows callers to supply decoded values.
// Missing parts are derived from the other name property where possible.
func (c *Card) NamesWith(text func(Prop) string) (full, family, given string) {
	if p, ok := c.Get("N"); ok {
		parts := strings.Split(text(p), ";")
		family = strings.TrimSpace(parts[0])
		if len(parts) > 1 {
			given = strings.TrimSpace(parts[1])
		}
	}
	if p, ok := c.Get("FN"); ok {
		full = strings.TrimSpace(text(p))
	}
	if full == "" {
		full = strings.TrimSpace(given + " " + family)
	}
	if family == "" && given == "" && full != "" {
		f := strings.Fields(full)
		family = f[len(f)-1]
		given = strings.Join(f[:len(f)-1], " ")
	}
	return full, family, given
}

// Write renders the card using the original raw lines of each property.
func (c *Card) Write(w io.Writer) error {
	var b strings.Builder
	b.WriteString("BEGIN:VCARD\n")
	for _, p := range c.Props {
		for _, l := range p.Raw {
			b.WriteString(l + "\n")
		}
	}
	b.WriteString("END:VCARD\n")
	_, err := io.WriteString(w, b.String())
	return err
}
