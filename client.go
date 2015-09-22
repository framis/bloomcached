package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

// NOT WORKING YET. THE MESSAGES ARE NOT SEPARATED
func main() {
	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	defer conn.Close()
	checkErrorC(err)

	//go handleClientRequest(conn)

	go readCsv("/Users/francois/projects/go/src/revinate.com/bloom/tripadvisorUrlsSmall.csv", conn)
}

func readCsv(filePath string, conn net.Conn) {

	csvfile, err := os.Open(filePath)
	defer csvfile.Close()

	checkErrorC(err)

	reader := csv.NewReader(csvfile)
	rawCSVdata, err := reader.ReadAll()
	checkErrorC(err)

	writer := bufio.NewWriter(conn)
	for n, each := range rawCSVdata {
		writer.WriteString(each[0] + "\n")
		if n%100 == 0 {
			writer.Flush()
		}
	}
	writer.Flush()
}

//func handleClientRequest(conn net.Conn) {
//	buf := make([]byte, 4096)
//	n, err := conn.Read(buf)
//	fmt.Println(string(buf[0:n]))
//	checkErrorC(err)
//}

func checkErrorC(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
