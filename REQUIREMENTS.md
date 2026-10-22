# Requirements

Go project (module `mergevcf`, Go 1.22) with two CLI tools sharing the `internal/vcf` library.
Input: `.vcf` stream with many vCards (versions 2.1 and 3.0 occur) via stdin. Results go to
stdout; all diagnostics go to stderr so stdout can be redirected to a file.

## Layout

- `cmd/mergevcf` – merge tool; `cmd/cleanvcf` – cleaning tool
- `internal/vcf` – the single shared library. Generic part (`card.go`,
  `parse.go`, `merge.go`): parse, write, merge, sort; never interprets or
  rewrites values. Cleaning part (`decode.go`, `clean.go`): UTF-8 /
  quoted-printable decoding, phone normalization, Apple-label handling,
  building clean cards; used by cleanvcf only, never by mergevcf.
- `make build` -> `bin/mergevcf`, `bin/cleanvcf`; `make test` runs all tests
- Each tool has a testable `run(args, stdin, stdout, stderr) int`

## Parsing (shared)

- Handles CRLF, folded lines (leading space/tab) and quoted-printable soft
  line breaks (line ending in `=`, needed to keep a card's lines intact).
  Property groups (`item1.TEL`) are recognized. Original raw lines of each
  property are retained and written back verbatim.
- Names come from the raw `FN`/`N` values; missing parts are derived from the
  other (last word of FN = family name).
- Never lose contacts.

## mergevcf

`mergevcf < input.vcf > merged.vcf` (reads from stdin, no flags, no positional arguments)

- Reads contacts from stdin.
- Pure merge: cards are written back unmodified (raw lines, original
  `VERSION`, encodings, categories, ...). It does no decoding or phone
  normalization; for reliable duplicate matching run cleanvcf on the inputs first.
- Duplicates = identical (trimmed) raw `TEL` value, or same case-insensitive
  raw name. First card is kept; `TEL` lines of duplicates with values the kept
  card lacks are appended unchanged.
- Output sorted by last name, then first name, then full name (case-insensitive).
- stderr: only `Total unique contacts: N` (no per-contact lines).

## cleanvcf

`cleanvcf [--prefix +49] < input.vcf > output.vcf` (reads from stdin, no positional arguments)

- Reads contacts from stdin.
- Output is plain UTF-8 vCard 3.0 containing only `VERSION`, `FN`, `N`
  (`Family;Given;;;`) and `TEL`. Everything else (CATEGORIES, ORG, ...) is dropped.
- Decoding: quoted-printable values (`ENCODING=QUOTED-PRINTABLE`, UTF-8
  bytes) are decoded; invalid UTF-8 is treated as Latin-1. Newlines in values
  are written as `\n`.
- Phone normalization: drop everything except digits and a leading `+`;
  leading `00` -> `+`; leading `0` -> `--prefix` (default `+49`).
- Apple style: phone numbers in grouped props (`item1.TEL`) are kept as normal
  `TEL`. An `X-ABLabel` whose value is itself a phone number becomes a `TEL`;
  other label texts are dropped. A known label (Mobile/iPhone/Home/Work/fax/
  pager, also in `_$!<Mobile>!$_` form) on a grouped TEL provides its type.
- Same number within a card is collapsed. Only standard TEL types are kept
  (CELL, HOME, WORK, VOICE, FAX, PAGER, ...), written as `TYPE=A,B`; `PREF`
  and `X-` types are dropped; 2.1 bare params (`TEL;CELL`) become types.
- Input order preserved. Cards with neither name nor phone are skipped
  (message on stderr).
- stderr: `Total processed contacts: N` (counts all input cards).

## Known gaps / ideas

- Output lines are not folded at 75 octets.
- Non-UTF-8 charsets are only handled as Latin-1 fallback.
- Name-only duplicate matching may merge two different people with the same name.
- A phone-like `X-ABLabel` is detected heuristically (6+ digits/separators).
