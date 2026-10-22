// Command showvcf reads vCards from stdin and outputs a CSV representation
// sorted by last name.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"vcftools/internal/vcf"
)

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("showvcf", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: showvcf < input.vcf > contacts.csv")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return 2
	}

	cards, err := vcf.Parse(stdin)
	if err != nil {
		fmt.Fprintf(stderr, "error reading input: %v\n", err)
		return 1
	}

	vcf.SortByLastName(cards)

	w := csv.NewWriter(stdout)
	w.Comma = ';'
	for _, c := range cards {
		_, family, given := c.Names()
		record := []string{family, given}
		for _, p := range c.All("TEL") {
			if v := strings.TrimSpace(p.Value); v != "" {
				record = append(record, v)
			}
		}
		if err := w.Write(record); err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	fmt.Fprintf(stderr, "Total processed contacts: %d\n", len(cards))
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
