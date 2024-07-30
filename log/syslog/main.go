package main

import (
	"fmt"
	"log"
	"log/syslog"

	syslog2 "gopkg.in/mcuadros/go-syslog.v2"
)

func server() {
	channel := make(syslog2.LogPartsChannel)
	handler := syslog2.NewChannelHandler(channel)

	server := syslog2.NewServer()
	server.SetFormat(syslog2.RFC5424)
	server.SetHandler(handler)
	server.ListenUDP("0.0.0.0:514")
	server.Boot()

	go func(channel syslog2.LogPartsChannel) {
		for logParts := range channel {
			fmt.Println(logParts)
		}
	}(channel)

	server.Wait()
}
func main() {
	network := "udp"
	raddr := "10.10.66.248:8200"

	s, err := syslog.Dial(network, raddr, syslog.LOG_WARNING|syslog.LOG_LOCAL7, "tool")
	if err != nil {
		log.Fatal(err)
	}

	s.Emerg("log test")
	s.Info("Thu nghiem he thong")
	fmt.Printf("Done.\n")
}
