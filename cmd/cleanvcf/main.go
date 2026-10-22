// Command cleanvcf rewrites a .vcf stream from stdin as plain UTF-8 vCard 3.0 contacts
// that only contain a name and normalized phone numbers.
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
	fs := flag.NewFlagSet("cleanvcf", flag.ContinueOnError)
	fs.SetOutput(stderr)
	prefix := fs.String("prefix", "+49", "country code prefix replacing a leading 0 in phone numbers")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: cleanvcf [--prefix +49] < input.vcf > output.vcf")
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

	w := bufio.NewWriter(stdout)
	for _, c := range cards {
		cc, ok := vcf.Clean(c, *prefix)
		if !ok {
			fmt.Fprintln(stderr, "skipping contact without name and phone number")
			continue
		}
		if err := cc.Write(w); err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	fmt.Fprintf(stderr, "Total processed contacts: %d\n", len(cards))
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
