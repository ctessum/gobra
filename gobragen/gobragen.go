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

// Package gobragen generates inputs for gobra from cobra.
package gobragen

import (
	"github.com/ctessum/gobra"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Generate creates a gobra command from a cobra command.
func Generate(cmd *cobra.Command) *gobra.Command {
	children := cmd.Commands()
	c := &gobra.Command{
		Name:     cmd.Name(),
		Doc:      cmd.Long,
		Children: make([]*gobra.Command, len(children)),
	}
	cmd.PersistentFlags().VisitAll(func(arg1 *pflag.Flag) {
		c.PersistentFlags = append(c.PersistentFlags, gobra.Flag{
			Name:  arg1.Name,
			Use:   arg1.Usage,
			Value: arg1.Value.String(),
		})
	})
	cmd.LocalNonPersistentFlags().VisitAll(func(arg1 *pflag.Flag) {
		c.LocalFlags = append(c.LocalFlags, gobra.Flag{
			Name:  arg1.Name,
			Use:   arg1.Usage,
			Value: arg1.Value.String(),
		})
	})
	for i, child := range children {
		c.Children[i] = Generate(child)
	}
	return c
}
