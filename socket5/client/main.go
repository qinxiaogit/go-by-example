package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"socket5/common"
	"syscall"
)

type xConfig struct {
	Listen   string
	Server   string
	InSecure bool
	Key      string
	Path     string
}

var config xConfig
var filter HostFilter
var tcpChecker = NewTCPChecker()
var tslog common.TSLog

const (
	rulePath     = `../config/rules.txt`
	autoRulePath = `../config/auto-rules.txt`
)

type Server struct {
}

func (s *Server) Run(network, add string) error {
	l, err := net.Listen(network, add)

	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()

		if err != nil {
			return err
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) error {
	defer conn.Close()

	var bior = bufio.NewReader(conn)
	var biow = bufio.NewWriter(conn)
	var biorw = bufio.NewReadWriter(bior, biow)
	firsts, err := bior.Peek(1)
	if err != nil {
		return err
	}
	switch firsts[0] {
	case '\x04':
	case '\x05':
		var sp SocksProxy
		sp.Handle(conn, biorw)
	default:
		var hp HTTPProxy
		hp.Handle(conn, biorw)
	}
	return nil
}
func parseConfig() {
	flag.StringVar(&config.Listen, "listen", "0.0.0.0:1082", "listen address(host:port)")
	flag.StringVar(&config.Server, "server", "127.0.0.1:1081", "server address(host:port)")
	flag.BoolVar(&config.InSecure, "insecure", true, "don't verify server certificate")
	flag.StringVar(&config.Key, "key", "", "login key")
	flag.StringVar(&config.Path, "path", "/", "/your/path")
	flag.Parse()
}

func handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGKILL)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		filter.SaveAuto(autoRulePath)
		fmt.Println()
		os.Exit(0)
	}()
}

func main() {
	handleInterrupt()
	parseConfig()

	filter.Init(rulePath)
	filter.LoadAuto(autoRulePath)
	s := Server{}

	if err := s.Run("tcp4", config.Listen); err != nil {
		filter.SaveAuto(autoRulePath)
		panic(err)
	}
}
