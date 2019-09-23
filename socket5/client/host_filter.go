package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
)

type ProxyType byte

const (
	proxyTypeNone ProxyType = iota
	proxyTypeDirect
	proxyTypeProxy
	proxyTypeReject
	proxyTypeAutoDirect
	proxyTypeAutoProxy
)

func (t ProxyType) String()string{
	switch t {
	case proxyTypeDirect:
		return "direct"
	case proxyTypeProxy:
		return "proxy"
	case proxyTypeReject:
		return "reject"
	case proxyTypeAutoDirect:
		return "auto-direct"
	case proxyTypeAutoProxy:
		return "auto-proxy"
	default:
		return "unknown"
	}
}
// IsAuto returns true if the rule is an auto-generated rule.
func (t ProxyType)IsAuto()bool{
	switch t{
	case proxyTypeAutoDirect,proxyTypeAutoProxy:
			return true
	default:
		return false
	}
}

func proxyTypeFromString(name string)ProxyType{
	switch name {
	case "direct":
		return proxyTypeDirect
	case "proxy":
		return proxyTypeProxy
	case "reject":
		return proxyTypeReject
	case "auto-direct":
		return proxyTypeAutoDirect
	case "auto-proxy":
		return proxyTypeAutoProxy
	default:
		return proxyTypeNone
	}
}

type AddrType uint

const (
	_ AddrType = iota
	IPv4
	Domain
)

var reIsComment = regexp.MustCompile(`^[ \t]*#`)

func isComment(line string)bool{
	return reIsComment.MatchString(line)
}
// HostFilter returns the proxy type on specified host.
type HostFilter struct {
	mu sync.RWMutex
	hosts	map[string]ProxyType
	cidrs map[*net.IPNet]ProxyType
}
// auto save
func (f *HostFilter)SaveAuto(path string){
	f.mu.Lock()
	defer f.mu.Unlock()

	file,err:= os.Create(path)
	if err!=nil{
		return
	}
	defer file.Close()
	w:=bufio.NewWriter(file)
	for k,t:=range f.hosts{
		switch t {
		case proxyTypeAutoProxy,proxyTypeAutoDirect:
			fmt.Fprintf(w,"%s,%s\n",k,t)
		}
	}
	w.Flush()
}

func (f *HostFilter)LoadAuto(path string){
	file,err:=os.Open(path)
	if err!=nil{
		return
	}
	defer file.Close()
	f.scanFile(file)
}
//init
func (f *HostFilter)Init(path string){
	f.hosts = make(map[string]ProxyType)
	f.cidrs = make(map[*net.IPNet]ProxyType)

	if file,err := os.Open(path);err!=nil{
		tslog.Red("rule file not found :%s",path)
	}else{
		f.scanFile(file)
		file.Close()
	}
}

func (f *HostFilter)scanFile(reader io.Reader){
	scanner := bufio.NewScanner(reader)

	for scanner.Scan(){
		rule:=strings.Trim(scanner.Text(),"\t")
		if isComment(rule)||rule == ""{
			continue
		}
		toks := strings.Split(rule,",")
		if len(toks) == 2{
			ptype := proxyTypeFromString(toks[1])
			if ptype == 0{
				tslog.Red("invalid proxy type:%s",toks[1])
				continue
			}
			if strings.IndexByte(toks[0],'/') == -1{
				f.hosts[toks[0]] = ptype
			}else{
				_,ipnet,err := net.ParseCIDR(toks[0])
				if err!=nil{
					f.cidrs[ipnet] = ptype
				}else{
					tslog.Red("bad cidr :%s",toks[0])
				}
			}
		}else{
			tslog.Red("invalid rule :%s",rule)
		}
	}
}

func (f *HostFilter)AddHost(host string,ptype ProxyType){
	f.mu.Lock()
	defer f.mu.Unlock()

	ty,ok := f.hosts[host]
	f.hosts[host] = ptype
	if !ok{
		tslog.Green("+ Add rule[%s] %s ",ptype,host)
	}else{
		if ty!=ptype{
			tslog.Green("Change Rule [%s-%s] %s",ty,ptype,host)
		}
	}
}

func (f *HostFilter)DeleteHost(host string){
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.hosts,host)
	tslog.Red("- Delete Rule %s",host)
}

func (f *HostFilter) Test(host,port string)(proxyType ProxyType){
	defer func() {
		if proxyType == proxyTypeNone{
			pty := proxyTypeAutoDirect
			if !tcpChecker.Check(host,port){
				pty = proxyTypeAutoProxy
			}
			f.AddHost(host,pty)
			proxyType = pty
		}
	}()

	f.mu.RLock()
	defer f.mu.RUnlock()

	host = strings.ToLower(host)

	if !strings.Contains(host,"."){
		return proxyTypeDirect
	}
	aty := Domain
	if net.ParseIP(host).To4() !=nil{
		aty =IPv4
	}
	if aty == IPv4{
		if ty,ok := f.hosts[host];ok{
			return ty
		}
		ip := net.ParseIP(host)
		for ipnet,ty :=range f.cidrs{
			if ipnet.Contains(ip){
				return ty
			}
		}
	}else if aty == Domain{
		part := host
		for {
			if ty,ok := f.hosts[part];ok{
				if !ty.IsAuto(){
					return ty
				}
			}
			index :=strings.IndexByte(part,'.')
			if index == -1{
				break
			}
			part = part[index+1:]
		}
	}
	return proxyTypeNone
}


