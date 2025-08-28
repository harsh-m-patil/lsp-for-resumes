package main

import (
	"bufio"
	"log"
	"lspfromscratch/rpc"
	"os"
)

func main() {
	logger := getLogger("/home/harsh/repos/projects/lsp-go/log.txt")
	logger.Println("Lsp started")

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Text()
		handleMessage(logger, msg)
	}
}

func handleMessage(logger *log.Logger, msg any) {
	logger.Println(msg)
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("error opening file")
	}

	return log.New(logfile, "[lspfromscratch] ", log.Ldate|log.Ltime|log.Lshortfile)
}
