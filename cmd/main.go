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
			fmt.Fprintf(w, "Cannot find the requested resource %v", dirp)
			return
		}

		switch {
		case info.IsDir():
			serveDir(w, fileSystem, dirp)
		case !info.IsDir():
			serveFileContents(w, fileSystem, dirp)
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *sPort), nil))
}

func serveDir(w http.ResponseWriter, fileSystem fs.FS, target string) {
	entries, err := fs.ReadDir(fileSystem, target)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Could not read requested dir, %q", target)
		return
	}
	n := len(entries)
	if n == 0 {
		fmt.Sprintln(w, "Empty Dir")
	}

	s := ""
	s += fmt.Sprintln("Total files:", n)
	for _, entry := range entries {
		s += fmt.Sprintln(fs.FormatDirEntry(entry))
	}

	fmt.Fprintf(w, s)
}

func serveFileContents(w http.ResponseWriter, fileSystem fs.FS, target string) {
	contents, err := fs.ReadFile(fileSystem, target)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Could not open the file", target)
		return
	}

	fmt.Fprint(w, string(contents))
}
