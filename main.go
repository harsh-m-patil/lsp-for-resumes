package main

import (
	"bufio"
	"encoding/json"
	"log"
	"lspfromscratch/lsp"
	"lspfromscratch/rpc"
	"os"
)

func main() {
	logger := getLogger("/home/harsh/repos/projects/lsp-go/log.txt")
	logger.Println("Lsp started")

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Println("Error decoding message:", err)
			continue
		}

		handleMessage(logger, method, content)
	}
}

func handleMessage(logger *log.Logger, method string, contents []byte) {
	logger.Println("Received message with method:", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Println("Error unmarshaling initialize request:", err)
			return
		}

		logger.Printf("Connected to: %s %s\n",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)

		// reply
		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)
		writer := os.Stdout
		writer.Write([]byte(reply))

		logger.Println("Sent initialize response")
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("error opening file")
	}

	return log.New(logfile, "[lspfromscratch] ", log.Ldate|log.Ltime|log.Lshortfile)
}
