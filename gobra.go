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

// Package gobra is an HTML-based graphical user interface (GUI) for the cobra
// command line interface (CLI; github.com/spf13/cobra).
package gobra

import (
	"bytes"
	"fmt"
	"html/template"

	"honnef.co/go/js/dom"

	"github.com/go-humble/view"
)

// Command holds information about a cobra command
type Command struct {
	// Name is the name of the command
	Name string

	// Use is the documentation for the command.
	Doc string

	// Flags flags for the command.
	Flags []Flag

	// Children are subcommands of this command.
	Children []*Command

	ParentElement dom.Element
	view.DefaultView
	listener *view.EventListener
}

var tCmd, tFlag *template.Template

func init() {
	const commandTpl = `
  <div>
  <h3>{{.Name}}</h3>
  <div>{{range .Flags}}<p>{{.Name}} {{.Value}}</p>{{end}}</div>
  {{if .Children}}<select id="{{.Name}}select">
    <option value="" disabled selected>Select</option>
    {{range .Children}}<option>{{ .Name }}</option>{{end}}
    </select>{{end}}
  </div>`

	var err error
	tCmd, err = template.New("cmd").Parse(commandTpl)
	if err != nil {
		panic(err)
	}
}

// Render renders the view of the command.
func (c *Command) Render() error {
	b := new(bytes.Buffer)
	if err := tCmd.Execute(b, c); err != nil {
		return err
	}
	c.Element().SetInnerHTML(b.String())
	view.AppendToEl(c.ParentElement, c)
	if c.listener != nil {
		c.listener.Remove()
	}
	c.listener = view.AddEventListener(c, "input", fmt.Sprintf("#%sselect", c.Name), c.renderChild)
	return nil
}

func (c *Command) renderChild(ev dom.Event) {
	sel := ev.Target().(*dom.HTMLSelectElement)
	for i, child := range c.Children {
		child.ParentElement = c.Element()
		child.Render()
		if i == sel.SelectedIndex-1 {
			view.Append(c, child)
		} else {
			view.Remove(child)
		}
	}
}

// Flag holds informaation about a flag for a command.
type Flag struct {
	// Name is the name of the argument.
	Name string

	// Use is usage information.
	Use string

	// Value is the argument value.
	Value string
}
