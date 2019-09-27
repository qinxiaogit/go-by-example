package main

import (
	"bufio"
	"crypto/tls"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"socket5/common"
	"sync"
	"time"
)

type Relayer interface {
	Begin(addr string, src net.Conn) error
	Relay() *RelayResult
	ToRemote(b []byte) error
	ToLocal(b []byte) error
}

type RelayResult struct {
	errTx error
	errRx error
	nTx   int64
	nRx   int64
}

type LocalRelayer struct {
	src  net.Conn
	dst  net.Conn
	addr string
}

// Begin implements Relayer.Begin.
func (r *LocalRelayer) Begin(addr string, src net.Conn) error {
	r.src = src
	r.addr = addr

	dst, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		tslog.Red("? Dial host:%s - %v", addr, err)
		return err
	}
	r.dst = dst
	return nil
}

// Begin implements Relayer.ToLocal.
func (r *LocalRelayer) ToLocal(b []byte) error {
	r.src.Write(b)
	return nil
}

func (r *LocalRelayer) ToRemote(b []byte) error {
	r.dst.Write(b)
	return nil
}
func (r *LocalRelayer) Relay() *RelayResult {
	tslog.Log("> [Direct] %s", r.addr)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	var tx int64
	var rx int64

	var errTx, errRx error

	go func() {
		tx, errTx = io.Copy(r.dst, r.src)
		wg.Done()
		r.src.Close()
		r.dst.Close()
	}()

	go func() {
		rx, errRx = io.Copy(r.src, r.dst)
		wg.Done()
		r.src.Close()
		r.dst.Close()
	}()

	wg.Wait()
	r.src.Close()
	r.dst.Close()

	tslog.Gray("< [Direct] %s [TX:%d,RX:%d]", r.addr, tx, rx)

	return &RelayResult{
		errTx: errTx,
		errRx: errRx,
		nTx:   tx,
		nRx:   rx,
	}
}

const gVersion string = "taosocks/20190722"

var (
	ErrCannotDialRemoteServer        = errors.New("cannot dial remote server")
	ErrRemoteServerCannotConnectHost = errors.New("remote proxy server cannot connect to the specified host")
)

type RemoteRelayer struct {
	src  net.Conn
	dst  net.Conn
	addr string
}

func (r *RemoteRelayer) dialServer() (net.Conn, error) {
	tlscfg := &tls.Config{
		InsecureSkipVerify: config.InSecure,
	}
	conn, err := tls.Dial("tcp4", config.Server, tlscfg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", config.Path, nil)
	if err != nil {
		return nil, err
	}
	req.Host = config.Server
	req.Header.Add("Connection", "upgrade")
	req.Header.Add("Upgrade", gVersion)
	req.Header.Add("Authorization", "taosocks"+config.Key)

	err = req.Write(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	bio := bufio.NewReader(conn)
	resp, err := http.ReadResponse(bio, nil)
	if err != nil {
		conn.Close()
		return nil, err
	}
	resp.Body.Close()
	if resp.StatusCode != 101 {
		conn.Close()
		return nil, errors.New("server upgrade protocol error")
	}
	return conn, nil
}

func (r *RemoteRelayer) Begin(addr string, src net.Conn) error {
	r.src = src
	r.addr = addr

	dst, err := r.dialServer()
	if err != nil {
		tslog.Red("%v", err)
		return ErrCannotDialRemoteServer
	}
	r.dst = dst

	enc := gob.NewEncoder(r.dst)
	dec := gob.NewDecoder(r.dst)

	err = enc.Encode(common.OpenMessage{Addr: r.addr})
	if err != nil {
		return err
	}
	var oamsg common.OpenAckMessage
	err = dec.Decode(&oamsg)
	if err != nil {
		return err
	}
	if !oamsg.Status {
		return ErrRemoteServerCannotConnectHost
	}
	return nil
}

func (r *RemoteRelayer) ToLocal(b []byte) error {
	r.src.Write(b)
	return nil
}
func (r *RemoteRelayer) ToRemote(b []byte) error {
	var msg common.RelayMessage
	msg.Data = b

	enc := gob.NewEncoder(r.dst)
	enc.Encode(msg)
	return nil
}

func (r *RemoteRelayer) Relay() *RelayResult {
	tslog.Log("> [Proxy ]%s", r.addr)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	var tx, rx int64
	var errTx, errRx error

	go func() {
		tx, errTx = r.src2dst()
		wg.Done()
		if errRx != nil {
			r.src.Close()
			r.dst.Close()
		}
	}()
	wg.Wait()

	r.src.Close()
	r.dst.Close()

	tslog.Gray("< [Proxy ] %s [TX:%d , RX:%d]", r.addr, tx, rx)
	return &RelayResult{
		errTx: errTx,
		errRx: errRx,
		nTx:   tx,
		nRx:   rx,
	}
}

func (r *RemoteRelayer) src2dst() (int64, error) {
	enc := gob.NewEncoder(r.dst)

	buf := make([]byte, common.ReadBuffSize)

	var all int64
	var err error
	var n int

	for {
		n, err = r.src.Read(buf)
		if err != nil {
			break
		}
		var msg common.RelayMessage
		msg.Data = buf[:n]

		err = enc.Encode(msg)
		if err != nil {
			break
		}
		all += int64(n)
	}
	return all, err
}

func (r *RemoteRelayer) dst2src() (int64, error) {
	dec := gob.NewDecoder(r.dst)

	var all int64
	var err error

	for {
		var msg common.RelayMessage

		err = dec.Decode(&msg)
		if err != nil {
			break
		}
		_, err = r.src.Write(msg.Data)
		if err != nil {
			break
		}
		all += int64(len(msg.Data))
	}
	return all, err
}

type SmartRelayer struct {
}

func (o *SmartRelayer) Relay(host string, conn net.Conn, beforeRelay func(r Relayer) error) error {
	hostname, portstr, _ := net.SplitHostPort(host)
	proxyType := filter.Test(hostname, portstr)

	var r Relayer

	switch proxyType {
	case proxyTypeDirect, proxyTypeAutoDirect:
		r = &LocalRelayer{}
	case proxyTypeProxy, proxyTypeAutoProxy:
		r = &RemoteRelayer{}
	case proxyTypeReject:
		return fmt.Errorf("x host is rejected: %s", hostname)
	}

	var beginErr error
	useRemote := false

	if err := r.Begin(host, conn); err != nil {
		beginErr = err
		switch r.(type) {
		case *LocalRelayer:
			if proxyType != proxyTypeDirect {
				r = &RemoteRelayer{}
				if err := r.Begin(host, conn); err == nil {
					beginErr = nil
					useRemote = true
				}
			}
		}
	}
	if beginErr != nil {
		conn.Close()
		if beginErr != ErrCannotDialRemoteServer {
			switch proxyType {
			case proxyTypeAutoDirect, proxyTypeAutoProxy:
				filter.DeleteHost(hostname)
			}
		}
		return errors.New("no relayer can relay such host")
	}
	if useRemote {
		filter.AddHost(hostname, proxyTypeAutoProxy)
	}
	if beforeRelay != nil {
		if beforeRelay(r) != nil {
			conn.Close()
			return errors.New("before Relay returns an error")
		}
	}
	rr := r.Relay()
	_ = rr
	return nil
}
