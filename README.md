# cuestr

`cuestr` is a utility tool for string literals in CUE files.

## `cuestr fmt [...FILES]`

Format CUE files and string literals in CUE files.

For each string literal format, a different formatter can be specified.

```console
find . -type f -name '*.cue' | xargs -I{} cuestr fmt {} --cmd 'Expr:deno fmt ${FILE} --ext js' --cmd 'Query:prettier ${FILE} --parser graphql'
```
