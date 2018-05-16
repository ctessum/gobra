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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/net/websocket"
)

type flagType struct {
	*pflag.Flag
	Type string
}

// flagSetToSlice converts pflag.FlagSet to slices for iteration
// It also combines two flagsets.
func flagSetToSlice(fl *pflag.FlagSet, fl2 *pflag.FlagSet) []flagType {
	var out []flagType

	fl.VisitAll(func(f *pflag.Flag) {
		if f.Name != "help" {
			out = append(out, flagType{
				Flag: f,
				Type: f.Value.Type(),
			})
		}
	})

	fl2.VisitAll(func(f *pflag.Flag) {
		if f.Name != "help" {
			out = append(out, flagType{
				Flag: f,
				Type: f.Value.Type(),
			})
		}
	})
	return out
}

func notHelpCommand(use string) bool {
	return use != "help [command]"
}

func (s *Server) canUploadFile(name string) bool {
	_, ok := s.uploadableFlags[name]
	return ok
}

const commandTpl = `
<div id="gobra-{{.Root.Use}}">
{{ define "command" }}
	<div data-gobra-name={{.Use}} style="{{if .HasParent }}display:none;{{end}}">
		<h3>{{.Use}}</h3>
		<p>{{.Long}}</p>
		<ul class="flags">
			{{ range (flagSetToSlice .PersistentFlags .LocalNonPersistentFlags) }}
				<li><code data-name={{ .Name }} data-type={{.Type}}>--{{ .Name }}=<input type="text" value={{ .Value.String }}></input>
					{{ if (canUploadFile .Name) }}
						<input type="file" name="{{ .Name }}" {{ if (isStringSlice .Type) }}multiple{{ end }}>
					{{ end }}
					</code><br>
					<blockquote>{{ .Usage }}</blockquote>
				</li>
			{{ end }}
		</ul>
		{{ if .HasSubCommands }}
			<select data-gobra-select>
				<option selected disabled>Select</option>
				{{ range .Commands }}
				{{ if (notHelpCommand .Use) }}<option value="{{.Use}}">{{ .Use }}</option>{{ end }}
				{{ end }}
			</select>
			{{range .Commands }}
				{{ template "command" .}}
			{{ end }}
		{{ end }}
	</div>
{{ end }}

{{ template "command" .Root }}
<br/>
<button>Execute</button>

<pre class="gobraStatus" style="padding:10px; background:lightgray; height:30em; overflow-y:scroll; white-space: pre-wrap; word-break: break-all;">
</pre>

<script>
const serverAddress = {{ if .ServerAddress}} "{{ .ServerAddress }}" {{ else }} "" {{ end }};

{{ with .Root }}
const logger = document.querySelector("#gobra-{{.Use}} .gobraStatus");
const execBtn = document.querySelector("#gobra-{{.Use}}>button");

// printData prints appends data to destination and scrolls to bottom.
const printData = (dest, str) => {
	dest.textContent += str;
	dest.scrollTop = dest.scrollHeight;
}

// clearLogger clears content of an output logger.
const clearLogger = (dest) => {
	dest.textContent = "";
}

// serverSend sends a request to the server and returns a Promise.
// It takes in the commands and flags as an array
// where each flag are of the format "name=value".
const serverSend = (cmds, flags) => {
	return fetch("http://"+serverAddress+"/"+cmds.join("/")+"?"+flags.join("&"));
}

// When an option is chosen, display the correct sub-command.
document.querySelectorAll("#gobra-{{.Use}} [data-gobra-select]").forEach( option =>
	option.onchange = e =>
		[...e.target.parentElement.children].forEach( el =>
			el.tagName != "DIV"? 1 :
				el.style.display = el.dataset.gobraName == e.target.value? "" : "none"
		)
);

let files = document.querySelectorAll("#gobra-{{.Use}} input[type^=f]");
for (const file of files) {
	file.addEventListener("change", e => {
		file.previousElementSibling.value = "";
		file.previousElementSibling.disabled = true;
	})
}

// Compile query when Execute is clicked
execBtn.onclick = e => {
	execBtn.setAttribute("disabled", "disabled");
	clearLogger(logger);

	// find file inputs and upload them
	let promisesOfFiles = [];
	let files = document.querySelectorAll("#gobra-{{.Use}} input[type^=f]");


	for (const file of files) {

		if (file.files.length === 0) continue;

		let formData = new FormData();
		formData.set("name", file.parentElement.dataset.name);
		formData.set("type", file.parentElement.dataset.type);
		for (const fileData of file.files) {
			formData.append("data", fileData);
		}

		let request = fetch("http://" + serverAddress + "/upload", {
			method: "POST",
			body: formData
		})
		.catch(err => {
			return Promise.reject("Failed uploading: " + err + "\n");
		})
		.then(res => res.json())
		.then(res => {
			file.previousElementSibling.value = res.path;
			file.previousElementSibling.disabled = false;
			file.value= '';
		})
		.catch(err => {
			return Promise.reject("Failed processing file: " + err + "\n");
		})

		promisesOfFiles.push(request);
	}

	// wait for all files to finish uploaded before moving executing command
	Promise.all(promisesOfFiles)
	.then( () => {
		// recurse through page to populate commands and flags
		let recurse = el => {
			let cmds = [],
				flags = [];
			if (el.tagName = "DIV" && el.style.display !== "none") {
				if (el.dataset.gobraName) {
					cmds.push(el.dataset.gobraName);
					[...el.querySelector("ul.flags").querySelectorAll("code")].forEach(f => {
						if(f.children[0]) flags.push(f.dataset.name + "=" + encodeURIComponent(f.children[0].value));
					})
				}
				[...el.children].forEach( child => {
					if (child.style.display !== "none") {
						let childRes = recurse(child);
						Array.prototype.push.apply(cmds, childRes[0]);
						Array.prototype.push.apply(flags, childRes[1]);
						return;
					}
				})
			}
			return [cmds, flags];
		}
		let resultCmd = recurse(document.getElementById("gobra-{{.Use}}"));


		printData(logger, "→ "+resultCmd.reduce((x,y) => {
				return x.join(" ") + " "
					+ y.map(z =>
						"--" + z.split("=")[0] + "=\"" + decodeURIComponent(z.split("=")[1]) + "\""
					).join(" ")
			})+ "\n");

		serverSend(resultCmd[0], resultCmd[1])
			.then(res => res.text()).then( d => {
				printData(logger,"← " + d + "\n");
				execBtn.removeAttribute("disabled");
			})
			.catch(e => {
				printData(logger,"⤬ Failed communicating with server: " + e + "\n");
				execBtn.removeAttribute("disabled");
			});
	})
	.catch( err => {
		printData(logger, "⤬ Failed data uploading, command not executed. " + err);

	})

	// re-enable
	execBtn.removeAttribute("disabled");
}
{{ end }}

window.onload = () => {
	let sock = new WebSocket("ws://" + serverAddress + "/ws");

	sock.onopen = () => {
		printData(logger, "* Connected.");
	}

	sock.onclose = (e) => {
		printData(logger, "* Connection Closed. ", e.reason);
		alert("Lost connection with server");
	}

	sock.onmessage = (e) => {
		printData(logger, e.data)
	}
}
</script>
</div>
`

// Render renders the view of the command. If the HTML field of the receiver
// is not nil, it will render the whole page, otherwise it will just render
// the gobra section.
func (s *Server) Render(w io.Writer) error {
	if s.HTML != nil {
		b := new(bytes.Buffer)
		if err := s.tCmd.Execute(b, s); err != nil {
			return err
		}
		return s.HTML.Execute(w, template.HTML(b.Bytes()))
	}
	return s.tCmd.Execute(w, s)
}

// Server struct/class that holds configuration for a Cobra back-end instance
type Server struct {
	// Root is the Cobra command tree root
	Root *cobra.Command

	// ServerAddress is the address that the front-end will communicate with.
	ServerAddress string

	// Allow Cross-Origin. If set to true, everyone can use the Gobra instance on client-side
	// Set this to true if you're planning to expose the API to public.
	AllowCORS bool

	// HTML is an HTML template.
	// If this is not nil, it will be served as an HTML front end.
	HTML *template.Template

	// socketChannel is a channel to the websocket handler.
	// Whatever gets sent to this channel will be sent by the websocket.
	socketChannel chan string

	// FileUploadFunc is a function that stores uploaded files and returns the
	// stored location. The default FileUploadFunc saves files in a temporary
	// directory.
	FileUploadFunc func(data io.Reader, name string) (filename string, err error)

	tCmd *template.Template

	// uploadableFlags is a set of flag names that can accept file uploads.
	uploadableFlags map[string]struct{}
}

// MakeFlagUploadable registers the given flag name(s) as allowing file uploads.
// Once a flag is registered, a file input will appear next to it in the
// user interface, allowing the user to upload a file for the flag argument.
func (s *Server) MakeFlagUploadable(names ...string) {
	if s.uploadableFlags == nil {
		s.uploadableFlags = make(map[string]struct{})
	}
	for _, name := range names {
		s.uploadableFlags[name] = struct{}{}
	}
}

// Write method makes Server implements io.Writer.
// It also transforms the input bytes to string then sends it to the websocket.
func (s *Server) Write(p []byte) (n int, err error) {
	s.socketChannel <- fmt.Sprintf("%s", p)
	return len(p), nil
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		// Serves front-end if root is requested
		if s.HTML != nil {
			if err := s.Render(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	} else if strings.HasPrefix(r.URL.Path, "/"+s.Root.Name()) {
		// Serves API if path starts with root command name

		if s.AllowCORS {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cmds := strings.Split(r.URL.Path[1:], "/")
		flags := r.Form

		// Set arguments to run.
		// Set cobra output to send to server instead.
		s.Root.SetArgs(cmds[1:])
		s.Root.SetOutput(s)

		// Getting the command we need to set flags
		c, _, _ := s.Root.Find(cmds[1:])
		for key, values := range flags {
			if err := c.Flags().Set(key, values[0]); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		fmt.Println("Executing: ", cmds, flags)
		_, err := s.Root.ExecuteC()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Finished. ")

	} else if strings.HasPrefix(r.URL.Path, "/upload") {
		// API end-point for file uploading
		// Store uploaded files to temporary folder.

		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB is held in memory.
			http.Error(w, fmt.Sprintf("while parsing upload form: %v", err), http.StatusInternalServerError)
			return
		}
		// flagName := r.MultipartForm.Value["name"][0]
		flagType := r.MultipartForm.Value["type"][0]
		fhs := r.MultipartForm.File["data"]
		paths := make([]string, len(fhs))
		for i, fh := range fhs {
			file, err := fh.Open()
			if err != nil {
				http.Error(w, fmt.Sprintf("failed retrieving uploaded file: %v", err), http.StatusInternalServerError)
				return
			}
			localPath, err := s.FileUploadFunc(file, fh.Filename)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed opening/copying uploaded file: %v", err), http.StatusInternalServerError)
				return
			}
			paths[i] = localPath
		}

		var pathResponse []byte
		if flagType == "stringSlice" {
			var err error
			pathResponse, err = writeAsCSV(paths)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			pathResponse = []byte(paths[0])
		}
		response, err := json.Marshal(map[string]string{
			"path": string(pathResponse),
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed opening/copying uploaded file: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(response))

	} else {
		// Everything else gets a 404
		http.Error(w, "404 Page not Found", http.StatusNotFound)
	}
}

// This is from github.com/spf13/pflag for string slice flags.
func writeAsCSV(vals []string) ([]byte, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err := w.Write(vals)
	if err != nil {
		return nil, err
	}
	w.Flush()
	return bytes.TrimSuffix(b.Bytes(), []byte("\n")), nil
}

func saveTempFileFunc() (func(io.Reader, string) (string, error), error) {
	tempDir, err := ioutil.TempDir("", "gobra")
	if err != nil {
		return nil, fmt.Errorf("gobra: creating temporary directory: %v", err)
	}
	return func(data io.Reader, name string) (string, error) {
		localPath := filepath.Join(tempDir, name)
		f, err := os.OpenFile(localPath, os.O_WRONLY|os.O_CREATE, 0660) // wr for owner/group only
		if err != nil {
			return "", fmt.Errorf("gobra: opening uploaded file: %v", err)
		}
		defer f.Close()
		_, err = io.Copy(f, data)
		if err != nil {
			return "", fmt.Errorf("gobra: saving uploaded file: %v", err)
		}
		return localPath, nil
	}, nil
}

func (s *Server) wsHandler(ws *websocket.Conn) {
	for {
		// Receiving data from channel
		data := <-s.socketChannel
		if err := websocket.Message.Send(ws, data); err != nil {
			fmt.Println("Error sending data.")
			break
		}
	}
}

// Start starts the server.
func (s *Server) Start() error {
	s.socketChannel = make(chan string)
	if s.FileUploadFunc == nil {
		var err error
		s.FileUploadFunc, err = saveTempFileFunc()
		if err != nil {
			return err
		}
	}

	if s.uploadableFlags == nil {
		s.uploadableFlags = make(map[string]struct{})
	}
	var funcMaps = template.FuncMap{
		"flagSetToSlice": flagSetToSlice,
		"notHelpCommand": notHelpCommand,
		"canUploadFile":  s.canUploadFile,
		"isStringSlice":  func(s string) bool { return s == "stringSlice" },
	}
	s.tCmd = template.Must(template.New("commands").Funcs(funcMaps).Parse(commandTpl))

	http.HandleFunc("/", s.handler)
	http.Handle("/ws", websocket.Handler(s.wsHandler))
	return http.ListenAndServe(s.ServerAddress, nil)
}
