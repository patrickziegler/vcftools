# vcf tools

```
cmd/mergevcf    merge .vcf contacts from stdin into unique contacts (fields stay untouched)
cmd/cleanvcf    rewrite a .vcf stream from stdin as plain UTF-8 vCard 3.0 with name + normalized phone numbers
cmd/showvcf     print a CSV version of all contacts from stdin sorted by last name
internal/vcf    shared library: parsing, writing, merging, sorting, cleaning
```

Build with `make build` (binaries in `bin/`), test with `make test`.

Pipeline usage:

```sh
cat contacts1.vcf contacts2.vcf ... | bin/cleanvcf [--prefix +49] | bin/mergevcf > merged.vcf
cat contacts1.vcf contacts2.vcf ... | bin/cleanvcf [--prefix +49] | bin/mergevcf | bin/showvcf > contacts.csv
```

Individual tool usage:

```sh
bin/cleanvcf [--prefix +49] < input.vcf > clean.vcf
bin/mergevcf < input.vcf > merged.vcf
bin/showvcf < input.vcf > contacts.csv
```

See [REQUIREMENTS.md](REQUIREMENTS.md) for the full specification.

## Development

```sh
podman compose run --rm dev
```
