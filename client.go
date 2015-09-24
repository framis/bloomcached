package main

import (
	"bufio"
	"net"
)

type (
	Client struct {
		conn net.Conn
	}
)

// Creates a new TCP Client
func NewClient(dsn string) (*Client, error) {
	conn, err := net.Dial("tcp", dsn)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

// If response is 1, the item might be in the filter
// If response is 0, the item is definitely not in the filter
func (c *Client) Test(item string) (bool, error) {
	writer := bufio.NewWriter(c.conn)
	_, err := writer.WriteString(item + "\n")
	if err != nil {
		return false, err
	}
	writer.Flush()
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	return (buf[0:n] == "200"), err
}

func (c *Client) Add(item string) (bool, error) {
	writer := bufio.NewWriter(c.conn)
	_, err := writer.WriteString(item + "\n")
	writer.Flush()
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	return (buf[0:n] == "201"), err
}
