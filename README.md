# cuestr

`cuestr` is a utility tool for string literals in CUE files.

## Usage

### `cuestr fmt [FILE ...]`

Format CUE files and string literals in CUE files.

For each string literal field, a different formatter can be specified.

```console
find . -type f -name '*.cue' | xargs -I{} cuestr fmt {} --field 'Expr:deno fmt ? --ext js' --field 'Query:prettier ? --write --parser graphql'
```

Any formatter can be specified for each field with the `--field` option ( `field:format command` ).

By formatting the file of placeholder `?` or the `FILE` environment variable, it can format string literals.

```console
--field 'Expr:deno fmt ? --ext js'
```

or

```console
--field 'Expr:deno fmt ${FILE} --ext js'
```


## Install

**homebrew tap:**

```console
$ brew install k1LoW/tap/cuestr
```

**[aqua](https://aquaproj.github.io/):**

```console
$ aqua g -i k1LoW/cuestr
```

**go install:**

```console
$ go install github.com/k1LoW/cuestr@latest
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/cuestr/releases)
