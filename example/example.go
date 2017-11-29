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
	"fmt"
	"github.com/ctessum/gobra"
	"github.com/ctessum/gobra/example/cmd"
	"html/template"
	"os"
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
	<footer>
		Add more content below. But here's a footer. 2017.
	</footer>
</div>
</body>
</html>
`
	output := template.Must(template.New("outputPage").Parse(tmpl))

	fmt.Println("Generating front-end html.")

	c := &gobra.CommandFromCobra{cmd.Root, ""}

	val, err := c.Render()
	if err != nil {
		panic(err)
	}

	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	output.Execute(f, template.HTML(val))

	fmt.Println("Successfully generated front-end html.")
	fmt.Println("Starting server, with front-end html/js.")
	fmt.Println("------------------")

	server := gobra.Server{cmd.Root, 8080, false, false}
	server.Start()
}
