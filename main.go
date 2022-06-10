package main

import (
	"net/http"
	"path/filepath"
)

type neuteredFilesystem struct {
	fs http.FileSystem
}

func (nfs neuteredFilesystem) Open(name string) (http.File, error) {
	f, err := nfs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	// check if f is a directory and serve index.html if it exists
	s, _ := f.Stat()
	if s.IsDir() {
		_, err = nfs.fs.Open(filepath.Join(name, "index.html"))
		if err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	// return file otherwise
	return f, nil
}

func main() {
	http.Handle("/", router())

	static := http.FileServer(neuteredFilesystem{http.Dir("assets")})
	http.Handle("/assets/", http.StripPrefix("/assets/", static))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
