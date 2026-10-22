# vcf tools

```
cmd/mergevcf    merge .vcf contacts from stdin into unique contacts (fields stay untouched)
cmd/cleanvcf    rewrite a .vcf stream from stdin as plain UTF-8 vCard 3.0 with name + normalized phone numbers
internal/vcf    shared library: parsing, writing, merging, sorting, cleaning
```

Build with `make build` (binaries in `bin/`), test with `make test`.

Pipeline usage:

```sh
cat contacts1.vcf contacts2.vcf ... | bin/cleanvcf [--prefix +49] | bin/mergevcf > merged.vcf
```

Individual tool usage:

```sh
bin/cleanvcf [--prefix +49] < input.vcf > clean.vcf
bin/mergevcf < input.vcf > merged.vcf
```

See [REQUIREMENTS.md](REQUIREMENTS.md) for the full specification.

## Development

```sh
podman compose run --rm dev
```
