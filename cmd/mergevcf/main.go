// Command mergevcf merges .vcf contacts from stdin into one list of unique contacts.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"vcftools/internal/vcf"
)

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("mergevcf", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: mergevcf < input.vcf > merged.vcf")
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

	merged := vcf.Merge(cards)
	vcf.SortByLastName(merged)

	w := bufio.NewWriter(stdout)
	for _, c := range merged {
		if err := c.Write(w); err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	fmt.Fprintf(stderr, "Total unique contacts: %d\n", len(merged))
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
