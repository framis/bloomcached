package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/willf/bloom"
	"net"
	"strconv"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = 3333
	SIZE      = 1000000 // Size of the Bloom filter
	HASH      = 5       // Number of Hash functions
)

type (
	ClientConn struct {
		conn net.Conn
		f    *bloom.BloomFilter
	}

	Server struct {
		listener net.Listener
		f        *bloom.BloomFilter
	}
)

func main() {
	server, err := NewServer(CONN_HOST+":"+strconv.Itoa(CONN_PORT), SIZE, HASH)
	if server == nil {
		panic("couldn't start listening: " + err.Error())
	}
	conns := clientConns(server.listener)
	for {
		go handleConn(&ClientConn{conn: <-conns, f: server.f})
	}
}

// Creates a new TCP Server with a bloom filter (m, k)
func NewServer(dsn string, m uint, k uint) (*Server, error) {
	listener, err := net.Listen("tcp", dsn)
	if listener == nil {
		return nil, err
	}
	f := bloom.New(m, k)
	return &Server{listener: listener, f: f}, nil
}

// Accepts connections
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

// Handles connections: from a request, send a response
func handleConn(c *ClientConn) {
	reader := bufio.NewReader(c.conn)
	for {
		request, err := reader.ReadBytes('\n')
		if err != nil { // EOF, or worse
			break
		}
		c.writeResponse(request)
	}
	c.conn.Close()
}

// Write the response from the request
func (c *ClientConn) writeResponse(request []byte) {
	message := bytes.Split(request, []byte{'|'})
	var response []byte
	if len(message) != 2 {
		response = []byte("400|Malformed request\n" + string(request))
	}

	verb := string(message[0])
	item := bytes.TrimSuffix(message[1], []byte("\n"))

	if verb == "TEST" {
		if c.f.Test(item) {
			response = []byte("200|true\n")
		} else {
			response = []byte("200|false\n")
		}
	} else if verb == "ADD" {
		c.f.Add(item)
		response = []byte("201\n")
	} else {
		response = []byte("400|Unknown Verb\n")
	}
	c.conn.Write(response)
	fmt.Println("Request:", string(request), "Response: ", string(response))
}
