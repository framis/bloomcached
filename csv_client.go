package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net"
	"os"
)

func mainT() {
	conn, err := net.Dial("tcp", "localhost:3333")
	defer conn.Close()
	checkErrorC(err)

	//go handleClientRequest(conn)

	readCsv("/Users/francois/projects/go/src/revinate.com/bloom/tripadvisorUrls.csv", conn)
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
