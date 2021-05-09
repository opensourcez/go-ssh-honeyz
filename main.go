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
	Listen("5555")
}
func Listen(port string) {
	sshConfig, listener := setupSSHListener(os.Args[1])
	processConnections(&sshConfig, listener)
}

func setupSSHListener(port string) (ssh.ServerConfig, net.Listener) {
	sshConfig := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			remoteAddr := c.RemoteAddr().String()
			ip := remoteAddr[0:strings.Index(remoteAddr, ":")]
			log.Printf("SSH connection from ip=[%s], username=[%s], password=[%s], version=[%s]", ip, c.User(), pass, c.ClientVersion())
			return nil, fmt.Errorf("invalid credentials")
		},
	}

	privateBytes, err := ioutil.ReadFile("sshkey")
	if err != nil {
		log.Fatalf("Failed to load private key.  Run make gen_ssh_key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}
	sshConfig.AddHostKey(private)

	portComplete := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp4", portComplete)
	if err != nil {
		log.Fatalf("failed to listen on *:%s", port)
	}

	log.Printf("listening on %s", port)

	return *sshConfig, listener
}

func processConnections(sshConfig *ssh.ServerConfig, listener net.Listener) {
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
	} else {
		sshConn.Close()
	}
}
