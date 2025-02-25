# cuestr

`cuestr` is a utility tool for string literals in CUE files.

## Usage

### `cuestr fmt [FILE ...]`

Format CUE files and string literals in CUE files.

For each string literal field, a different formatter can be specified.

```console
find . -type f -name '*.cue' | xargs -I{} cuestr fmt {} --field 'Expr:deno fmt ${FILE} --ext js' --field 'Query:prettier ${FILE} --parser graphql'
```

Any formatter can be specified for each field with the `--field` option ( `field:format command` ).

By formatting the file of the `FILE` environment variable, it can format string literals.

## Install

**homebrew tap:**

```console
$ brew install k1LoW/tap/cuestr
```

**go install:**

```console
$ go install github.com/k1LoW/cuestr@latest
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/cuestr/releases)
