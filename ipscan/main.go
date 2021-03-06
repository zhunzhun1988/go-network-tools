package main

import (
	"flag"
	"fmt"
	"go-network-tools/utils"
	"go-network-tools/utils/log"
	"net"
	"strings"
	"sync"
	"time"
	"golang.org/x/crypto/ssh"
)

var (
	startip   = flag.String("startip", "", "scan start ip")
	stopip    = flag.String("endip", "", "scan end ip")
	startport = flag.Int("startport", -1, "scan start port")
	endport   = flag.Int("endport", -1, "scan end port")
	threadnum = flag.Int("threadnum", 100, "number goroutine to run")
	timeout   = flag.Int("timeout", 10, "connect port timeout second")
	checkSSH  = flag.Bool("checkssh", true, "check whether the port is ssh port")
	sshuser = flag.String("sshuser", "root", "scan ssh user")
	sshpasswd = flag.String("sshpasswd", "123456", "scan ssh password")
)

func checkArgs() error {
	if *startip == "" {
		return fmt.Errorf("please input startip")
	}
	if *startport <= 0 {
		return fmt.Errorf("please input starport")
	}
	return nil
}

type scanTask struct {
	ip   string
	port int
}

type Scanner struct {
	taskCh      chan scanTask
	numGorotine int
}

func isSSHConn(conn net.Conn) bool {
	buf := make([]byte, 30)
	n, err := conn.Read(buf)
	if n > 0 && err == nil && strings.Contains(string(buf), "SSH") {
		return true
	}
	return false
}

func checkSSHPasswd(ip string , port int,user, password string) bool {
	PassWd := []ssh.AuthMethod{ssh.Password(password)}
	Conf := ssh.ClientConfig{
	User: user, 
	Auth: PassWd,
	HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
	}
	Client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), &Conf)
	if err != nil { 
		return false
	}
	defer Client.Close()
	return true
}

func worker(id, timeout int, ch chan scanTask, wg *sync.WaitGroup) {
exitFor:
	for {
		select {
		case task := <-ch:
			if task.port == 0 {
				break exitFor
			}
			addr := fmt.Sprintf("%s:%d", task.ip, task.port)

			conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)

			if err != nil {
				log.MyLogS("[%s] tcp port [%d] is not open", task.ip, task.port)
			} else {
				defer conn.Close()
				if *checkSSH == true && isSSHConn(conn) {
					if checkSSHPasswd(task.ip, task.port ,*sshuser, *sshpasswd) {
						log.MyLogI("[%s] tcp port [%d] is open is SSH port and password is %s:%s", task.ip, task.port ,*sshuser, *sshpasswd)
					} else {
						log.MyLogI("[%s] tcp port [%d] is open is SSH port", task.ip, task.port)
					}
					
				} else {
					log.MyLogI("[%s] tcp port [%d] is open", task.ip, task.port)
				}
			}
			/*_, err = net.Dial("udp", addr)
			if err != nil {
				//log.MyLogI("[%s] udp port [%d] is not open", task.ip, task.port)
			} else {
				log.MyLogI("[%s] udp port [%d] is open", task.ip, task.port)
			}*/
		}
	}
	wg.Done()
}

func (s *Scanner) producer() {
	sip := *startip
	eip := sip
	sport := *startport
	eport := sport

	if *stopip != "" {
		eip = *stopip
	}
	if *endport > 0 {
		eport = *endport
	}
	cmp := utils.IPCmp(sip, eip)
	if cmp == -2 {
		close(s.taskCh)
		return
	} else if cmp == 1 {
		sip, eip = *stopip, *startip
	}
	if sport > eport {
		sport, eport = eport, sport
	}
	log.MyLogI("start scan ip[%s-%s], port[%d-%d]", sip, eip, sport, eport)
	for ip := sip; ip != ""; ip = utils.NextIP(ip) {
		for port := sport; port <= eport; port++ {
			s.taskCh <- scanTask{
				ip:   ip,
				port: port,
			}
		}
		if ip == eip {
			break
		}
	}
	close(s.taskCh)
}
func (s *Scanner) Run() {
	go s.producer()
	var wg sync.WaitGroup
	for i := 0; i < s.numGorotine; i++ {
		wg.Add(1)
		go worker(i+1, *timeout, s.taskCh, &wg)
	}
	wg.Wait()
}

func main() {
	defer log.Clean()
	flag.Parse()
	utils.CheckError(checkArgs())
	scanner := &Scanner{
		taskCh:      make(chan scanTask, *threadnum),
		numGorotine: *threadnum,
	}
	scanner.Run()
}
