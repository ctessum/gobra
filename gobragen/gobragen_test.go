/*
MIT License

Copyright (c) 2017 Chris Tessum

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gobragen

import (
	"testing"

	"github.com/ctessum/gobra"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

func TestGenerate(t *testing.T) {
	root := &cobra.Command{
		Use:  "Root",
		Long: "Documentation for the root command.",
	}
	root.Flags().String("abc", "123", "you and me")
	root.PersistentFlags().Float64("=", 42, "the answer")
	sub1 := &cobra.Command{
		Use:  "sub1",
		Long: "sub1 documentation",
	}
	sub2 := &cobra.Command{
		Use:  "sub2",
		Long: "sub2 documentation",
	}
	subsub1 := &cobra.Command{
		Use:  "subsub1",
		Long: "subsub1 documentation",
	}
	subsub1.Flags().String("sub", "sub", "1")
	root.AddCommand(sub1, sub2)
	sub1.AddCommand(subsub1)

	cmd := Generate(root)
	want := &gobra.Command{
		Name: "Root",
		Doc:  "Documentation for the root command.",
		Flags: []gobra.Flag{
			gobra.Flag{Name: "=", Use: "the answer", Value: "42"},
			gobra.Flag{Name: "abc", Use: "you and me", Value: "123"},
		},
		Children: []*gobra.Command{
			&gobra.Command{
				Name: "sub1",
				Doc:  "sub1 documentation",
				Children: []*gobra.Command{
					&gobra.Command{
						Name: "subsub1",
						Doc:  "subsub1 documentation",
						Flags: []gobra.Flag{
							gobra.Flag{Name: "sub", Use: "1", Value: "sub"},
						},
					},
				},
			},
			&gobra.Command{
				Name: "sub2",
				Doc:  "sub2 documentation",
			},
		},
	}
	if diff := pretty.Diff(cmd, want); len(diff) != 0 {
		t.Errorf("have != want; diff=%+v", diff)
	}
}
