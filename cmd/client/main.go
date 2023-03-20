package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type serverResponse struct {
	bid string
}

const EndpointPath = "cotacao"
const ServerApiUrl = "http://localhost:8080/" + EndpointPath
const FileType = ".txt"
const FileName = EndpointPath + FileType
const FilePath = "./web/static/" + FileName

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	res, err := http.NewRequestWithContext(ctx, "GET", ServerApiUrl, nil)
	if err != nil {
		panic(err)
	}
	defer func(b io.ReadCloser) {
		err := b.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	var decodedResponse = new(serverResponse)
	err = json.NewDecoder(res.Body).Decode(decodedResponse)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	_, err = file.WriteString("Dolar: " + decodedResponse.bid)
	if err != nil {
		panic(err)
	}
}
