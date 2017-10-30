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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CommandFromCobra is a wrapper for cobra.Command to work with Gobra
type CommandFromCobra struct {
	// CobraCmd holds the pointer to a Cobra command
	CobraCmd *cobra.Command 
}

var tCmd *template.Template

// convert pflag.FlagSet to slices for range to iterate
func FlagSetToSlice(fl *pflag.FlagSet) []*pflag.Flag {
	var out []*pflag.Flag

	fl.VisitAll(func (f *pflag.Flag) {
		out = append(out, f)
	})
	return out
}

func init() {

	var funcMaps = template.FuncMap{
		"flagSetToSlice" : FlagSetToSlice,
	}

	const commandTpl = `
<div>
{{ define "command" }}
	<div data-gobra-name={{.Use}} style="{{if .HasParent }}display:none; {{end}}">
		<h3>{{.Use}}</h3>
		<ul>
			{{ range (flagSetToSlice .PersistentFlags) }}
				<li><code>--{{ .Name }}="{{ .Value.String }}"</code> {{ .Usage }} </li>
			{{ end }}
			{{ range (flagSetToSlice .Flags) }}
				<li><code>--{{ .Name }}="{{ .Value.String }}"</code> {{ .Usage }} </li>
			{{ end }}
		</ul>
		{{ if .HasSubCommands }}
			<select data-gobra-select>
				<option selected disabled>Select</option>
				{{ range .Commands }}
				<option value="{{.Use}}">{{ .Use }}</option>
				{{ end }}
			</select>
			{{range .Commands }}
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

	tCmd = template.Must(template.New("commands").Funcs(funcMaps).Parse(commandTpl))
}

// Render renders the view of the command.
func (c *CommandFromCobra) Render() ([]byte, error) {
	b := new(bytes.Buffer)
	if err := tCmd.Execute(b, c.CobraCmd); err != nil {
		return b.Bytes(), err
	}
	
	return b.Bytes(), nil
}
