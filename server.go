package main

import (
	"bufio"

	"bytes"
	"fmt"
	"github.com/willf/bloom"
	"net"
	"os"
	"strconv"
)

const (
	CONN_HOST           = "localhost"
	CONN_PORT           = 3333
	CONN_TYPE           = "tcp"
	NUMBER_OF_ELEMENTS  = 1000000
	MULT                = 20
	NUMBER_OF_HASH_FUNC = 5
)

type (
	Server struct {
		conn net.Conn
		f    *bloom.BloomFilter
	}
)

func main() {
	server, err := net.Listen(CONN_TYPE, CONN_HOST+":"+strconv.Itoa(CONN_PORT))
	if server == nil {
		panic("couldn't start listening: " + err.Error())
	}
	f := bloom.New(MULT*NUMBER_OF_ELEMENTS, NUMBER_OF_HASH_FUNC)
	conns := clientConns(server)
	for {
		go handleConn(&Server{conn: <-conns, f: f})
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Printf("couldn't accept: " + err.Error())
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(),
				client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(s *Server) {
	reader := bufio.NewReader(s.conn)
	for {
		request, err := reader.ReadBytes('\n')
		if err != nil { // EOF, or worse
			break
		}
		s.prepareResponse(request)
	}
	s.conn.Close()
}

func (s *Server) prepareResponse(request []byte) {
	message := bytes.Split(request, []byte{'|'})
	var response []byte
	if len(message) != 2 {
		response = []byte("400|Malformed request\n" + string(request))
	}

	verb := string(message[0])
	item := bytes.TrimSuffix(message[1], []byte("\n"))

	if verb == "TEST" {
		if s.f.Test(item) {
			response = []byte("200|true\n")
		} else {
			response = []byte("200|false\n")
		}
	} else if verb == "ADD" {
		s.f.Add(item)
		response = []byte("201\n")
	} else {
		response = []byte("400|Unknown Verb\n")
	}
	s.conn.Write(response)
	fmt.Println("Request:", string(request), "Response: ", string(response))
}

func checkError(err error) {
	if err != nil {
		fmt.Println("hi")
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
