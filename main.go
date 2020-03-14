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
	"strings"
)

// Version export
const Version = "1.1.0"

// logger stand-in
var dlog *oslog.Logger

// DEBUG export
var DEBUG = true

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

	resp.Header().Set("Connection", "close")

	_, fileName := filepath.Split(req.URL.Path)

	file, err := os.Open("files/" + fileName)
	if err != nil {
		log.Println("Failed to open file [", fileName, "]", err)
		resp.Header().Set("Content-Type", "text/plain")
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte("File not found " + fileName))
		return
	}

	fileStat, _ := file.Stat()
	fileLenStr := strconv.FormatInt(fileStat.Size(), 10)
	file.Close()

	byteRange := req.Header.Get("Range")
	if len(byteRange) > 0 {
		log.Println(byteRange)
		byteFromTo := strings.Split(byteRange[6:], "-")
		byteFrom, _ := strconv.ParseInt(byteFromTo[0], 10, 64)
		byteTo, toErr := strconv.ParseInt(byteFromTo[1], 10, 64)
		if toErr != nil {
			byteFromTo[1] = "1"
			byteTo = 1
		}
		byteLen := int64(byteTo - byteFrom + 1)
		buffer := make([]byte, byteLen)
		file.ReadAt(buffer, int64(byteFrom))

		// TODO apply mimetype
		resp.Header().Set("Content-Type", "application/octet-stream")
		resp.Header().Set("Content-Range", "bytes "+byteFromTo[0]+"-"+byteFromTo[1]+"/"+fileLenStr)
		resp.Header().Set("Content-Length", strconv.FormatInt(byteLen, 10))

		resp.WriteHeader(http.StatusOK)
		resp.Write(buffer)
	} else {
		// TODO apply mimetype
		resp.Header().Set("Content-Type", "application/octet-stream")
		resp.Header().Set("Content-Length", fileLenStr)
		io.Copy(resp, file)
	}
}

func fileList(resp http.ResponseWriter) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("failed to walk files", err)
		}
		if info.IsDir() {
			return nil
		}
		_, fileName := filepath.Split(path)
		resp.Write([]byte("<a href='/files/" + fileName + "'>" + fileName + "</a><br/>"))
		return nil
	}
}

func listHandler(resp http.ResponseWriter, req *http.Request) {
	dlog.Println("godl.listHandler")
	defer req.Body.Close()

	_, fileName := filepath.Split(req.URL.Path)
	dlog.Println("Downloading [", fileName, "]")

	resp.Header().Set("Connection", "close")

	if len(fileName) == 0 {
		resp.Header().Set("Content-Type", "text/html")
		err := filepath.Walk("files", fileList(resp))
		if err != nil {
			log.Fatal("failed to list files", err)
			resp.Write([]byte("Failed to walk files"))
			return
		}
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
	resp.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	resp.WriteHeader(http.StatusOK)
	io.Copy(resp, file)
}

func main() {
	log.Println("Starting GoDL v", Version)
	http.HandleFunc("/files/", downloadHandler)
	http.HandleFunc("/", listHandler)
	err := http.ListenAndServe(":80", nil)
	// err := http.ListenAndServeTLS(":443", "demo.crt", "demo.key", nil)
	if err != nil {
		log.Fatal("failed to start listener", err)
	}
}
