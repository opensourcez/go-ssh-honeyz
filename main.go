package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func main() {
	sshConfig := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			remoteAddr := c.RemoteAddr().String()
			ip := remoteAddr[0:strings.Index(remoteAddr, ":")]
			log.Printf("ip=[%s], username=[%s], password=[%s], version=[%s]", ip, c.User(), pass, c.ClientVersion())
			return nil, fmt.Errorf("invalid credentials")
		},
	}

	privateBytes, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatalf("Failed to load private key.  Run make gen_ssh_key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}
	sshConfig.AddHostKey(private)

	portComplete := fmt.Sprintf(":%s", os.Args[1])
	listener, err := net.Listen("tcp4", portComplete)
	if err != nil {
		log.Fatalf("failed to listen on *:%s", os.Args[1])
	}

	log.Printf("listening on %s", os.Args[1])

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			continue
		}
		log.Println("got connection !!")
		go handleConnection(sshConfig, tcpConn)
	}
}

func handleConnection(sshConfig *ssh.ServerConfig, tcpConn net.Conn) {
	defer tcpConn.Close()
	sshConn, _, _, err := ssh.NewServerConn(tcpConn, sshConfig)
	if err != nil {
		log.Println("Connection error..", err)
	} else {
		sshConn.Close()
	}
}
