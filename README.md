# gobra
Gobra is a graphical user interface (GUI) for the cobra command line interface (cli).

# Usage

Gobra generates an HTML snippet.

```go
import (
	"github.com/ctessum/gobra"
	"html/template"
	"github.com/ctessum/gobra/example/cmd"
)

func main () {
	const wrapper = `
		[Your webpage's HTML here]
			{{.}}
		[...]
	`
	output := template.Must(template.New("outputPage").Parse(wrapper))

	cmd := &gobra.CommandFromCobra{ cmd.Root }

	val, err := cmd.Render();
	if err != nil {
		panic(err)
	}

	f, err := os.Create("out.html")

	output.Execute(f, template.HTML(val))
}

```