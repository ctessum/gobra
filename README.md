# Gobra
Gobra generates web interfaces to interact with [cobra](https://github.com/spf13/cobra)

## Usage

There are two parts to Gobra: generating client-side interface and serving a HTTP API. 

### Generating client-side HTML

Gobra has a `CommandFromCobra` struct which takes in two arguments: a `cobra.Command` and a string which is the address for the server the client-side will connect to. If the address is left blank, the client will connect to the same server it resides.

Calling `CommandFromCobra.Render()` will return an HTML string, which you can use to insert into your already existing webpage.

### Running the API

Gobra also has a `Server` struct which takes in 4 arguments: `cobra.Command`, port number, AllowCORS and Frontless.

The `cobra.Command` argument should be the same one that you give to `CommandFromGobra`. AllowCORS is a bool determines whether the API will have `Access-Allow-Control-Origin: *` or not. If `Frontless` is true, Gobra will not serve the `index.html` file from the folder it's run.

## Example

Here is an example in the case where you would run both the client-side and API on the same server:

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

	cmd := &gobra.CommandFromCobra{ cmd.Root, "" }

	val, err := cmd.Render();
	if err != nil {
		panic(err)
	}

	// We dump this HTML file into the current folder
	f, err := os.Create("index.html")
	output.Execute(f, template.HTML(val))

	// The server will start serving the index.html file we've just made
	// on port 8080, and the API will be served from /gobra as well.
	server := gobra.Server { cmd.Root, 8080, false, false }
	server.Start()
}

```

## API endpoints

If you decide to use Gobra only as a server, the API endpoint works like so:

If the command you want to run is: `app math add --num1=3 --num2=6`

You would want to make a GET request to: `//<serverAddress>/app/math/add?num1=3&num2=6`

