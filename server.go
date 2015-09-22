package main

import (
	"bufio"

	"bytes"
	"fmt"
	"github.com/willf/bloom"
	"net"
	"os"
)

const (
	CONN_HOST           = "localhost"
	CONN_PORT           = "3333"
	CONN_TYPE           = "tcp"
	NUMBER_OF_ELEMENTS  = 1000000
	MULT                = 20
	NUMBER_OF_HASH_FUNC = 5
)

func main() {

	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	//defer l.Close()

	checkError(err)

	f := bloom.New(MULT*NUMBER_OF_ELEMENTS, NUMBER_OF_HASH_FUNC)

	for {
		conn, err := l.Accept()
		checkError(err)
		go handleRequest(conn, f)
	}
}

func handleRequest(conn net.Conn, f *bloom.BloomFilter) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil { // EOF, or worse
			break
		}
		s := bytes.TrimSuffix(line, []byte("\n"))
		if f.Test(s) {
			fmt.Println(string(s), "is in the filter")
			conn.Write([]byte("{\"url\":\"" + string(s) + "\", \"present\": true}"))
		} else {
			f.Add(s)
			fmt.Println("Added", string(s), "to the filter")
			conn.Write([]byte("{\"url\":\"" + string(s) + "\", \"present\": false}"))
		}
		checkError(err)
	}
	conn.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("hi")
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
