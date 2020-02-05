package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Version export
var Version = "1.0.1"

func downloadHandler(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	_, fileName := filepath.Split(req.URL.Path)
	fmt.Println("Downloading [", fileName, "]")

	resp.Header().Set("Connection", "close")

	if len(fileName) == 0 {
		fmt.Println("Invalid file [", fileName, "]")
		resp.Header().Set("Content-Type", "text/plain")
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte("No file name provided"))
		return
	}

	file, err := os.Open("files/" + fileName)
	if err != nil {
		fmt.Println("Failed to open file [", fileName, "]", err)
		resp.Header().Set("Content-Type", "text/plain")
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte("File not found " + fileName))
		return
	}
	defer file.Close()

	fileStat, _ := file.Stat()
	resp.Header().Set("Content-Type", "application/octet-stream")
	resp.Header().Add("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	resp.WriteHeader(http.StatusOK)
	io.Copy(resp, file)
}

func main() {
	port := "82"
	fmt.Println("Starting GoDL v", Version, " on port", port)
	http.HandleFunc("/", downloadHandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
