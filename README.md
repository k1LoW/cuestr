# cuestr

`cuestr` is a utility tool for string literals in CUE files.

## Usage

### `cuestr fmt [...FILES]`

Format CUE files and string literals in CUE files.

For each string literal format, a different formatter can be specified.

```console
find . -type f -name '*.cue' | xargs -I{} cuestr fmt {} --cmd 'Expr:deno fmt ${FILE} --ext js' --cmd 'Query:prettier ${FILE} --parser graphql'
```

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
