# vcf tools

```
cmd/mergevcf    merge several .vcf files into unique contacts (fields stay untouched)
cmd/cleanvcf    rewrite a .vcf as plain UTF-8 vCard 3.0 with name + normalized phone numbers
internal/vcf    shared library: parsing, writing, merging, sorting, cleaning
```

Build with `make build` (binaries in `bin/`), test with `make test`.

```
bin/cleanvcf [--prefix +49] input.vcf > clean.vcf
bin/mergevcf a.vcf b.vcf ... > merged.vcf
```

Tip: run `cleanvcf` on each file first, then `mergevcf` on the results.

See [REQUIREMENTS.md](REQUIREMENTS.md) for the full specification.

## Development

```sh
podman compose run --rm dev
```
