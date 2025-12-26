package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

const DEFAULT_PORT = 8000

func main() {
	sPort := flag.Int("p", DEFAULT_PORT, "Specify the port to serve on.")
	flag.Parse()

	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get working directory")
	}
	fileSystem := os.DirFS(workdir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		target := strings.TrimSuffix(strings.TrimPrefix(p, "/"), "/")
		dirp := "."

		if target != "" {
			dirp = target
		}

		log.Println("Requested", dirp)

		info, err := os.Stat(dirp)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Could not find the requested resource %v", dirp)
			return
		}

		switch {
		// override special behaviour for http.ServeFile to redirect to / if index.html is requested
		case info.Name() == "index.html":
			fptr, err := os.Open(dirp) // os.Open implements io.ReadSeeker needed by http.ServeContent
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Could not open the file", target)
				return
			}
			http.ServeContent(w, r, info.Name(), info.ModTime(), fptr)

		case info.IsDir():
			http.ServeFileFS(w, r, fileSystem, dirp)
			entries, err := fs.ReadDir(fileSystem, dirp)
			if err != nil {
				return
			}
			fmt.Fprintf(w, "<pre>Total files: %d</pre>", len(entries))

		case !info.IsDir():
			http.ServeFile(w, r, dirp)
		}
	})

	log.Printf("Listening on http://localhost:%v", *sPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *sPort), nil))
}
