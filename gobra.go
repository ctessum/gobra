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
	"html/template"
)

// Command holds information about a cobra command
type Command struct {
	// TopLevel specificies whether the command is a top-level one or not.
	TopLevel bool 

	// Name is the name of the command
	Name string

	// Use is the documentation for the command.
	Doc string

	// Flags flags for the command.
	Flags []Flag

	// Children are subcommands of this command.
	Children []*Command
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

var tCmd, tFlag *template.Template

func init() {
	const commandTpl = `
<div>
{{ define "command" }}
	<div data-gobra-name={{.Name}} style="{{if not .TopLevel}}display:none; {{end}}">
		<h3>{{.Name}}</h3>
		<ul>
			{{ range .Flags }}
				<li><code>--{{ .Name }}="{{ .Value }}"</code> {{ .Use }}</li>
			{{ end }}
		</ul>
		{{ if .Children}}
			<select data-gobra-select>
				<option selected disabled>Select</option>
				{{ range .Children }}
				<option value="{{.Name}}">{{ .Name }}</option>
				{{ end }}
			</select>
			{{range .Children}}
				{{ template "command" .}}
			{{ end }}
		{{ end }}
	</div>
{{ end }}
{{ template "command" .}}
<script>
	document.querySelectorAll("[data-gobra-select]").forEach( option => {
		option.onchange = e => {
			[...e.target.parentElement.children].forEach( el => {
				if (el.tagName == "DIV")
					el.style.display = (el.dataset.gobraName == e.target.value)? "" : "none";
			})
		}
	})
</script>
</div>
`

	tCmd = template.Must(template.New("commands").Parse(commandTpl))
}

// Render renders the view of the command.
func (c *Command) Render() ([]byte, error) {
	b := new(bytes.Buffer)
	if err := tCmd.Execute(b, c); err != nil {
		return b.Bytes(), err
	}
	
	return b.Bytes(), nil
}

