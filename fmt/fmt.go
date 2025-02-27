package fmt

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
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
					var (
						v     string
						names []string
					)
					switch l := field.Value.(type) {
					case *ast.BasicLit:
						v, _, err = stringValue(field.Value)
						if err != nil {
							return
						}
					case *ast.Interpolation:
						for _, e := range l.Elts {
							str, n, err := stringValue(e)
							names = append(names, n...)
							if err != nil {
								return
							}
							v += str
						}
					default:
						return
					}

					// multiline string literal only
					if !strings.Contains(v, "\n") || !strings.Contains(v, `"""`) {
						return
					}

					var (
						repPairs    []string
						revertPairs []string
					)
					for i, n := range names {
						repPairs = append(revertPairs, fmt.Sprintf(`\(%s)`, n), fmt.Sprintf("cuestr%dRep", i))
						revertPairs = append(repPairs, fmt.Sprintf("cuestr%dRep", i), fmt.Sprintf(`\(%s)`, n))
					}
					rep := strings.NewReplacer(repPairs...)
					revert := strings.NewReplacer(revertPairs...)

					v = rep.Replace(stripIndent(strings.Trim(v, `"`)))
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
					formatted := revert.Replace(`"""` + "\n" + strings.Trim(string(b), "\n") + "\n" + `"""`)
					field.Value = &ast.BasicLit{
						ValuePos: field.Value.Pos(),
						Kind:     token.STRING,
						Value:    formatted,
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
	var indents []int
	for _, l := range lines {
		if strings.Trim(l, " \t") == "" {
			continue
		}
		indents = append(indents, (len(l) - len(strings.TrimLeft(l, " \t"))))
	}
	slices.Sort(indents)
	indent := indents[0]
	if indent == 0 {
		return v
	}
	for i, l := range lines {
		if strings.Trim(l, " \t") == "" {
			lines[i] = ""
			continue
		}
		lines[i] = l[indent:]
	}

	return strings.Join(lines, "\n")
}

func stringValue(v ast.Expr) (string, []string, error) {
	var names []string
	switch v := v.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			return v.Value, nil, nil
		}
		return "", nil, fmt.Errorf("not a string literal: %v", v)
	case *ast.SelectorExpr:
		vx, names, err := stringValue(v.X)
		if err != nil {
			return "", nil, err
		}
		i, ok := v.Sel.(*ast.Ident)
		if !ok {
			return "", nil, fmt.Errorf("not an identifier: %v", v.Sel)
		}
		name := vx + "." + i.Name
		names = append(names, name)
		return name, names, nil
	case *ast.Ident:
		names = append(names, v.Name)
		return v.Name, names, nil
	default:
		return "", nil, fmt.Errorf("not a string: %v", v)
	}
}
