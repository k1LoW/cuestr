/*
Copyright © 2025 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	cuefmt "github.com/k1LoW/cuestr/fmt"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var cmds []string

var fmtCmd = &cobra.Command{
	Use:   "fmt [...FILES]",
	Short: "Format CUE files and string literals in CUE files",
	Long: `Format CUE files and string literals in CUE files.
For each string literal format, a different formatter can be specified.
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmdcmds := map[string]string{}
		for _, c := range cmds {
			splitted := strings.Split(c, ":")
			if len(splitted) != 2 {
				return fmt.Errorf("invalid format command: %s", c)
			}
			fmdcmds[splitted[0]] = splitted[1]
		}
		cf := cuefmt.New(fmdcmds)
		eg := new(errgroup.Group)
		for _, fp := range args {
			if _, err := os.Stat(fp); err != nil {
				return err
			}
			b, err := os.ReadFile(fp)
			if err != nil {
				return err
			}
			eg.Go(func() error {
				formatted, err := cf.Format(b)
				if err != nil {
					return err
				}
				if err := os.WriteFile(fp, formatted, 0600); err != nil {
					return err
				}
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fmtCmd)
	fmtCmd.Flags().StringSliceVarP(&cmds, "field", "f", []string{}, "format command for string literal field in CUE files. e.g. 'Expr:deno fmt ${FILE} --ext js'")
}
