package fmt

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/cue/token"
	"github.com/cli/safeexec"
)

const (
	envFileKey         = "FILE"
	defaultShell       = "sh"
	defaultPlaceholder = "?"
)

type Cue struct {
	shell       string
	placeholder string
	fmtcmds     map[string]string
}

func New(fmtcmds map[string]string) *Cue {
	return &Cue{
		shell:       defaultShell,
		placeholder: defaultPlaceholder,
		fmtcmds:     fmtcmds,
	}
}

func (c *Cue) Format(in []byte) ([]byte, error) {
	bin, err := safeexec.LookPath(c.shell)
	if err != nil {
		return nil, fmt.Errorf("failed to find %s: %w", c.shell, err)
	}
	f, err := parser.ParseFile("", in, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CUE: %w", err)
	}

	var errr error
	ast.Walk(f, nil, func(n ast.Node) {
		if errr != nil {
			return
		}
		if field, ok := n.(*ast.Field); ok {
			if ident, ok := field.Label.(*ast.Ident); ok {
				if fmtcmd, ok := c.fmtcmds[ident.Name]; ok {
					if strLit, ok := field.Value.(*ast.BasicLit); ok && strLit.Kind == token.STRING {
						// multiline string literal only
						if !strings.Contains(strLit.Value, "\n") || !strings.Contains(strLit.Value, `"""`) {
							return
						}
						v := stripIndent(strings.Trim(strLit.Value, `"`))
						f, err := os.CreateTemp("", "cuestr")
						if err != nil {
							errr = fmt.Errorf("failed to create temp file: %w", err)
							return
						}
						if _, err := f.WriteString(v); err != nil {
							errr = fmt.Errorf("failed to write to temp file: %w", err)
							return
						}
						fmtcmd = strings.ReplaceAll(fmtcmd, c.placeholder, f.Name())
						cmd := exec.Command(bin, "-c", fmtcmd)
						cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", envFileKey, f.Name()))
						if out, err := cmd.CombinedOutput(); err != nil {
							errr = fmt.Errorf("failed to run %q: %s", fmtcmd, string(out))
							return
						}
						b, err := os.ReadFile(f.Name())
						if err != nil {
							errr = fmt.Errorf("failed to read temp file: %w", err)
							return
						}
						formatted := `"""` + "\n" + strings.Trim(string(b), "\n") + "\n" + `"""`
						strLit.Value = formatted
					}
				}
			}
		}
	})
	if errr != nil {
		return nil, errr
	}

	formatted, err := format.Node(f)
	if err != nil {
		return nil, fmt.Errorf("failed to format CUE: %w", err)
	}

	return formatted, nil
}

func stripIndent(v string) string {
	lines := strings.Split(v, "\n")
	indent := -1
	for _, l := range lines {
		if l == "" {
			continue
		}
		tmp := len(l) - len(strings.TrimLeft(l, " \t"))
		if tmp > 0 && (indent == -1 || tmp < indent) {
			indent = tmp
		}
	}
	if indent == -1 {
		return v
	}
	for i, l := range lines {
		if l == "" {
			continue
		}
		lines[i] = l[indent:]
	}
	return strings.Join(lines, "\n")
}
