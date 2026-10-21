package vcf

import (
	"bufio"
	"io"
	"strings"
)

type logicalLine struct {
	text string
	raw  []string
	qp   bool // quoted-printable: a trailing "=" is a soft line break
}

func (l *logicalLine) open() bool { return l.qp && strings.HasSuffix(l.text, "=") }

// Parse reads all vCards from r. Folded lines and quoted-printable soft
// line breaks are unfolded for parsing, but the original lines are retained.
func Parse(r io.Reader) ([]*Card, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1024*1024), 64*1024*1024)

	var lines []*logicalLine
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		if n := len(lines); n > 0 {
			last := lines[n-1]
			switch {
			case last.open():
				last.text = strings.TrimSuffix(last.text, "=") + line
				last.raw = append(last.raw, line)
				continue
			case line != "" && (line[0] == ' ' || line[0] == '\t'):
				last.text += line[1:]
				last.raw = append(last.raw, line)
				continue
			}
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, &logicalLine{text: line, raw: []string{line}, qp: isQP(line)})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	var cards []*Card
	var cur *Card
	for _, l := range lines {
		switch strings.ToUpper(strings.TrimSpace(l.text)) {
		case "BEGIN:VCARD":
			cur = &Card{}
		case "END:VCARD":
			if cur != nil {
				cards = append(cards, cur)
			}
			cur = nil
		default:
			if cur != nil {
				cur.Props = append(cur.Props, parseProp(l))
			}
		}
	}
	return cards, nil
}

func isQP(line string) bool {
	key, _, ok := strings.Cut(line, ":")
	return ok && strings.Contains(strings.ToUpper(key), "QUOTED-PRINTABLE")
}

func parseProp(l *logicalLine) Prop {
	p := Prop{Raw: l.raw}
	key, val, ok := strings.Cut(l.text, ":")
	if !ok {
		p.Value = l.text
		return p
	}
	name, params, _ := strings.Cut(key, ";")
	if g, n, found := strings.Cut(name, "."); found {
		p.Group, name = g, n
	}
	p.Name = strings.ToUpper(strings.TrimSpace(name))
	p.Params = params
	p.Value = val
	return p
}
