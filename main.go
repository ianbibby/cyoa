package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Options struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Arc struct {
	Title   string    `json:"title"`
	Story   []string  `json:"story"`
	Options []Options `json:"options"`
}

func main() {
	const (
		defaultFilePath = "gopher.json"
	)
	var (
		filePath string
	)

	flag.StringVar(&filePath, "file", defaultFilePath, "Path to json story file.")
	flag.Parse()

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	adventure := make(map[string]Arc)
	err = json.NewDecoder(bufio.NewReader(f)).Decode(&adventure)
	if err != nil {
		log.Fatal(err)
	}

	data := `
		<html>
			<head>
			</head>
			<body>
				<h1>{{ .Title }}</h1>
				<div>
					{{ range .Story}}
						<p>
							{{ . }}
						</p>
					{{ end }}
				</div>
				<div>
					{{ if .Options }}
					{{ range .Options }}
						<p>
							<a href="/{{ .Arc }}">{{ .Text }}</a>
						</p>
					{{ end }}
					{{ else }}
						<p>
							<a href="/intro">Start again</a>
						</p>
					{{ end }}
				</div>
			</body>
		</html>
`

	tpl, err := template.New("story").Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")[1]
		if path == "" {
			path = "intro"
		}

		if _, ok := adventure[path]; !ok {
			http.Error(w, "arc does not exist", http.StatusNotFound)
			return
		}
		tpl.Execute(w, adventure[path])
	})

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
