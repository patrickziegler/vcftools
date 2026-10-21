package vcf

import (
	"regexp"
	"strings"
)

// standardTelTypes are the vCard 3.0 TEL types that are kept; everything
// else (PREF, X-... labels) is dropped.
var standardTelTypes = map[string]bool{
	"CELL": true, "HOME": true, "WORK": true, "VOICE": true, "FAX": true,
	"PAGER": true, "MSG": true, "CAR": true, "VIDEO": true, "TEXT": true,
	"BBS": true, "MODEM": true, "ISDN": true, "PCS": true,
}

// appleLabelTypes maps (lower-cased) Apple X-ABLabel texts to TEL types.
var appleLabelTypes = map[string]string{
	"mobile": "CELL", "iphone": "CELL", "home": "HOME", "work": "WORK",
	"homefax": "FAX", "workfax": "FAX", "otherfax": "FAX", "pager": "PAGER",
}

var (
	newlines    = strings.NewReplacer("\r\n", `\n`, "\n", `\n`, "\r", `\n`)
	phoneLike   = regexp.MustCompile(`^\+?[\d\s()./-]{6,}$`)
	appleMarker = regexp.MustCompile(`^_\$!<(.*)>!\$_$`)
)

func text(p Prop) string { return DecodeValue(p.Params, p.Value) }

func newProp(name, params, value string) Prop {
	line := name
	if params != "" {
		line += ";" + params
	}
	line += ":" + value
	return Prop{Name: name, Params: params, Value: value, Raw: []string{line}}
}

// telTypes extracts the standard TEL types from a parameter string, handling
// 2.1 style bare parameters (TEL;CELL), TYPE=A,B and repeated TYPE= parameters.
func telTypes(params string) []string {
	var types []string
	for _, p := range strings.Split(params, ";") {
		k, v, hasV := strings.Cut(p, "=")
		if hasV && !strings.EqualFold(k, "TYPE") {
			continue
		}
		if !hasV {
			v = k
		}
		for _, t := range strings.Split(v, ",") {
			if t = strings.ToUpper(strings.TrimSpace(t)); standardTelTypes[t] {
				types = append(types, t)
			}
		}
	}
	return types
}

// appleLabelType returns the TEL type for an Apple label such as
// "_$!<Mobile>!$_", or "" if unknown.
func appleLabelType(label string) string {
	label = strings.TrimSpace(label)
	if m := appleMarker.FindStringSubmatch(label); m != nil {
		label = m[1]
	}
	return appleLabelTypes[strings.ToLower(strings.ReplaceAll(label, " ", ""))]
}

// Clean returns a minimal, standardized vCard 3.0 in UTF-8 containing only
// the name (FN, N) and the phone numbers (TEL) of c. Phone numbers found in
// Apple style grouped properties (item1.TEL) are kept, and a phone number
// stored as an Apple label (item1.X-ABLabel:+49 ...) becomes a normal TEL.
// ok is false if the card has neither a name nor a phone number.
func Clean(c *Card, prefix string) (clean *Card, ok bool) {
	full, family, given := c.NamesWith(text)

	// Apple labels per group, used to derive a type for grouped TEL entries.
	labels := map[string]string{}
	for _, p := range c.All("X-ABLABEL") {
		labels[p.Group] = text(p)
	}

	var order []string
	types := map[string][]string{}
	add := func(raw string, ts []string) {
		n := NormalizePhone(raw, prefix)
		if n == "" {
			return
		}
		if _, seen := types[n]; !seen {
			order = append(order, n)
			types[n] = nil
		}
		for _, t := range ts {
			dup := false
			for _, e := range types[n] {
				dup = dup || e == t
			}
			if !dup {
				types[n] = append(types[n], t)
			}
		}
	}

	for _, p := range c.Props {
		switch p.Name {
		case "TEL":
			ts := telTypes(p.Params)
			if p.Group != "" {
				if t := appleLabelType(labels[p.Group]); t != "" {
					ts = append(ts, t)
				}
			}
			add(text(p), ts)
		case "X-ABLABEL":
			if v := text(p); phoneLike.MatchString(strings.TrimSpace(v)) {
				add(v, nil)
			}
		}
	}
	if full == "" && len(order) == 0 {
		return nil, false
	}

	out := &Card{Props: []Prop{newProp("VERSION", "", "3.0")}}
	out.Props = append(out.Props,
		newProp("FN", "", newlines.Replace(full)),
		newProp("N", "", newlines.Replace(family)+";"+newlines.Replace(given)+";;;"))
	for _, n := range order {
		params := ""
		if len(types[n]) > 0 {
			params = "TYPE=" + strings.Join(types[n], ",")
		}
		out.Props = append(out.Props, newProp("TEL", params, n))
	}
	return out, true
}
