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

package main

import (
	"os" // for now
	"github.com/ctessum/gobra"
	"html/template"
)

func main() {
	const tmpl = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Some page</title>
	<style>
		html, body {padding: 0; margin: 2% 0; font-family: sans-serif;}
		.container { max-width: 700px; margin: 0 auto; padding: 10px; }
	</style>
</head>
<body>
	<div class="container">
		<h1>Some valid HTML page here</h1>
		<span>Below is a div where Gobra is added.</span>
		<div>
			{{.}}
		</div>
	</div>
</body>
</html>
`
	output := template.Must(template.New("outputPage").Parse(tmpl))

	cmd := &gobra.Command{
		TopLevel: true,
		Name:          "go",
		Flags: []gobra.Flag{
			{
				Name:  "param",
				Value: "value",
			},
			{
				Name: "descr",
				Value: "testi test",
				Use: "Flag used for description",
			},
		},
		Children: []*gobra.Command{
			{
				Name: "run",
				Flags: []gobra.Flag{
					{
						Name:  "background",
						Value: "true",
					},
				},
				Children: []*gobra.Command{
					{
						Name: "example.go",
					},
					{
						Name: "somefile.go",
					},
				},
			},
			{
				Name: "test",
			},
		},
	}

	val, err := cmd.Render();
	if err != nil {
		panic(err)
	}

	f, err := os.Create("index.html")

	output.Execute(f, template.HTML(val))
}
