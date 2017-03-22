package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"log/syslog"
	"net"
	"net/textproto"
	"os/exec"
)

func main() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "av_server")
	if e == nil {
		log.SetOutput(logwriter)
	}
	fmt.Printf("starting server...  ")
	socket, err := net.Listen("tcp", ":55111")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("OK\n")
	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal("cannot accept connection:", err.Error())
		}
		go handleConection(conn)
	}
}

func handleConection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	textReader := textproto.NewReader(reader)
	filename, _ := textReader.ReadLine()
	fmt.Println(filename)
	cmd := exec.Command("/opt/kaspersky/kav4fs/bin/kav4fs-control", "--action=Remove", "--scan-file="+filename)
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		fmt.Println("error on exec: ", err.Error())
		log.Println(err.Error(), stderr.String())
		conn.Write([]byte(stderr.String()))
	}

	log.Printf("%s\n%s", filename, stdout.String())
	conn.Write([]byte(stdout.String()))
	conn.Close()
}
