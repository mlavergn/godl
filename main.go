package main

import (
	"io"
	"io/ioutil"
	"log"
	oslog "log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Version export
const Version = "1.0.3"

// logger stand-in
var dlog *oslog.Logger

// DEBUG export
var DEBUG = false

func init() {
	if DEBUG {
		dlog = oslog.New(os.Stderr, "GoDL ", oslog.Ltime|oslog.Lshortfile)
	} else {
		dlog = oslog.New(ioutil.Discard, "", 0)
	}
}

func downloadHandler(resp http.ResponseWriter, req *http.Request) {
	dlog.Println("godl.downloadHandler")
	defer req.Body.Close()

	_, fileName := filepath.Split(req.URL.Path)
	dlog.Println("Downloading [", fileName, "]")

	resp.Header().Set("Connection", "close")

	if len(fileName) == 0 {
		log.Println("Invalid file [", fileName, "]")
		resp.Header().Set("Content-Type", "text/plain")
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte("No file name provided"))
		return
	}

	file, err := os.Open("files/" + fileName)
	if err != nil {
		log.Println("Failed to open file [", fileName, "]", err)
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
	log.Println("Starting GoDL v", Version, " on port", port)
	http.HandleFunc("/", downloadHandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("failed to start listener", err)
	}
}
