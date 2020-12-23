package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	u "github.com/mishazawa/steve/utils"
)

var (
	datafolder string
	port       string
)

func init() {
	flag.StringVar(&datafolder, "datafolder", "./temp", "Data folder")
	flag.StringVar(&port, "port", "8080", "Server port")
	flag.Parse()
}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(datafolder))

	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.HandleFunc("/api/ls", ls)
	mux.HandleFunc("/api/upload", upload)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	fmt.Println("Running...", port)
	server.ListenAndServe()
}

func ls(w http.ResponseWriter, req *http.Request) {
	files, err := u.GetPaths(datafolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range files {
		rel, _ := filepath.Rel(datafolder, files[i])
		files[i] = rel
	}

	resp, err := json.Marshal(files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func upload(w http.ResponseWriter, req *http.Request) {
	if len(req.Header["X-Filepath"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filePath := fmt.Sprintf("%s/%s", datafolder, req.Header["X-Filepath"][0])

	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(file, req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
