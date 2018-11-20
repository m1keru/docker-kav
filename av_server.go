package main

import (
	"bufio"
	"bytes"
	"fmt"
	. "github.com/blacked/go-zabbix"
	"io/ioutil"
	"log"
	"log/syslog"
	"net"
	"net/textproto"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHost = `system-001`
	defaultPort = 10051
)

func sendToZabbix(checkDuration int64) {
	fqdn, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
	}
	strDuration := strconv.FormatInt(checkDuration, 10)
	hostname := strings.Split(fqdn, ".")[0]
	log.Printf("sending metrics to zabbix server %s : %s", defaultHost, string(strDuration))
	var metrics []*Metric
	metrics = append(metrics, NewMetric(hostname, "avcheck.duration", string(strDuration), time.Now().Unix()))
	packet := NewPacket(metrics)
	z := NewSender(defaultHost, defaultPort)
	resp, err := z.Send(packet)
	if err != nil {
		log.Println("unable to send data to zabbix server" + err.Error())
	}
	log.Printf("resp: %s", string(resp))
}

func main() {
	pid := fmt.Sprintf("%d", os.Getpid())
	err := ioutil.WriteFile("/var/run/av_server.pid", []byte(pid), 0644)
	if err != nil {
		log.Fatal(err.Error())
	}

	logwriter, e := syslog.New(syslog.LOG_NOTICE, "av_server")
	if e == nil {
		log.SetOutput(logwriter)
	}
	file, err := os.Open("/opt/kaspersky/kesl/bin/kesl-control")
	if err != nil {
		fmt.Println("kav binary not installed")
		log.Fatalln("kav binary not installed")
	}
	file.Close()
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
	cmd := exec.Command("/opt/kaspersky/kesl/bin/kesl-control", "--action", "Remove", "--scan-file", filename)
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	// Меряем время исполнения
	start := time.Now()
	if err := cmd.Run(); err != nil {
		fmt.Println("error on exec: ", err.Error())
		log.Println(err.Error(), stderr.String())
		conn.Write([]byte(stderr.String()))
	}
	milliseconds := int64(time.Since(start) / time.Millisecond)
	sendToZabbix(milliseconds)
	log.Printf("%s time:%d \n%s\n", filename, milliseconds, stdout.String())
	conn.Write([]byte(stdout.String()))
	conn.Close()
}
