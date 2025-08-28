package rpc_test

import (
	"lspfromscratch/rpc"
	"testing"
)

type EncodingExample struct {
	Testing bool `json:"testing"`
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing: true})

	if expected != actual {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMsg := "Content-Length: 15\r\n\r\n{\"method\":\"hi\"}"
	method, content, err := rpc.DecodeMessage([]byte(incomingMsg))
	contentLength := len(content)

	if err != nil {
		t.Fatal(err)
	}

	if contentLength != 15 {
		t.Fatalf("Expected 15 got %d", contentLength)
	}

	if method != "hi" {
		t.Fatalf("Expected hi got %s", method)
	}
}
