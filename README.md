# gobra
Gobra is a graphical user interface (GUI) for the cobra command line interface (cli).

# Usage

Gobra generates an HTML snippet.

```go
import (
	"github.com/ctessum/gobra"
	"html/template"
)

func main () {
	const wrapper = `
		[Your webpage's HTML here]
			{{.}}
		[...]
	`
	output := template.Must(template.New("outputPage").Parse(wrapper))

	cmd := &gobra.Command{
		TopLevel: true,
		Name:          "fakewget",
		Flags: []gobra.Flag{
			{
				Name:  "method",
				Value: "GET",
				Use: "HTTP request method"
			},
			{
				Name: "header",
				Value: "Content-Length:test/test",
				Use: "Header",
			},
		},
		Children: []*gobra.Command{
			{
				Name: "fetch",
				Flags: []gobra.Flag{
					{
						Name:  "in-background",
						Value: "true",
					},
				},
				Children: []*gobra.Command{
					{
						Name: "example.tld",
					},
					{
						Name: "website.site",
					},
				},
			},
			{
				Name: "ping",
			},
		},
	}

	val, err := cmd.Render();
	if err != nil {
		panic(err)
	}

	f, err := os.Create("out.html")

	output.Execute(f, template.HTML(val))
}

```