package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var connection *ssh.Client
var channel = make(chan []byte)

func main() {
	fmt.Println("Hello World")
	user := "deploy"
	host := "198.211.126.170:22"
	connection = makeConnection(user, host)
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	command := "ls -l /home/deploy"

	// extract command from input

	session := makeSession(connection)
	defer session.Close()

	setSession(session)
	err := session.Run(command)

	if err != nil {
		fmt.Println(err)
	}

	data := <-channel
	// fmt.Println(string(data))
	io.WriteString(w, string(data)+"\n")

}
func setSession(session *ssh.Session) {
	stdin, err := session.StdinPipe()
	if err != nil {
		fmt.Println(fmt.Errorf("Unable to setup stdin for session: %v", err))
		return
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		fmt.Println(fmt.Errorf("Unable to setup stdout for session: %v", err))
		return
	}

	bufSize := 1024
	go func() {
		for {
			b := make([]byte, bufSize)
			n, err1 := stdout.Read(b)
			if err1 != nil {
				return
			}
			if n > 0 {
				channel <- b[0:n]
			}
		}
	}()

	// go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		fmt.Println(fmt.Errorf("Unable to setup stderr for session: %v", err))
		return
	}
	go io.Copy(os.Stderr, stderr)

}

func makeConnection(user, host string) *ssh.Client {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			sshAgent(),
			// publicKeyFile("/Users/ram/.ssh/id_rsa"),
		},
	}

	connection, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		fmt.Println(fmt.Errorf("Failed to dial: %s", err))
		return connection
	}
	return connection

}

func makeSession(connection *ssh.Client) *ssh.Session {
	s := &ssh.Session{}

	session, err := connection.NewSession()
	if err != nil {
		fmt.Println(fmt.Errorf("Failed to create session: %s", err))
		return s
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = session.RequestPty("xterm", 150, 70, modes)
	if err != nil {
		session.Close()
		fmt.Println(fmt.Errorf("request for pseudo terminal failed: %s", err))
		return s
	}
	return session
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func sshAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}
