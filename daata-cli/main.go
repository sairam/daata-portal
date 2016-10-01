package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var connection *ssh.Client
var channel = make(chan []byte)

func main() {
	fmt.Println("Started server")
	user := "deploy"
	host := "198.211.126.170:22"
	port := "8002"
	connection = makeConnection(user, host)
	http.HandleFunc("/", hello)
	http.ListenAndServe(":"+port, nil)
}

const testOutput = `## Done so far
1. Connect to Host
2. Get a connection
3. Execute a command
4. Start HTTP endpoint
5. Send output to HTTP endpoint

## Known issues
1. Parallel connections wont work because of channels global variable
2. Fetch all data from channels
3. Close SSH connections properly

## TODO
1. SSH into a host
2. execute a command
3. send the command output to stdout,stderr
4. Close the connection
5. Add a web server
`

func hello(w http.ResponseWriter, r *http.Request) {

	var host *Host
	var commandExecution *CommandExecution
	var currentCommand = ""

	host = &Host{"sairam", "127.0.0.2", "example.com"}
	// 2 methods
	switch r.Method {

	// 1. GET
	// GET will just lists the hosts that we are connected to
	case http.MethodGet:
		commandExecution = &CommandExecution{Host: host, PWD: "~"}

	// 2. POST
	// POST will execute the command display the same page with the output
	case http.MethodPost:
		r.ParseForm()
		command := r.Form.Get("command")
		command = strings.TrimSpace(command)
		output := ""
		if command != "" {
			// command := "ls -l /home/deploy"
			session := makeSession(connection)
			defer session.Close()

			setSession(session)
			err := session.Run(command)

			if err != nil {
				fmt.Println(err)
			}

			data := <-channel
			output = string(data) // testOutput
		} else {
			output = ""
		}
		// fmt.Println(string(data))
		commandExecution = &CommandExecution{Host: host, PWD: "~/test/", Input: command, Output: output}
		currentCommand = command

	default:
		http.NotFound(w, r)
		return
	}

	history := strReverse(inputHistory, 10)

	info := []CommandExecution{*commandExecution, *commandExecution, *commandExecution, *commandExecution, *commandExecution, *commandExecution}
	uiData := UserInterfaceData{
		Details: info,
		History: history,
		Current: currentCommand,
	}

	if currentCommand != "" {
		addToHistory(currentCommand)
	}

	renderTemplate(w, &uiData)
}

func (h *Host) String() string {
	return fmt.Sprintf("%s@%s", h.Username, h.Hostname)
}

var inputHistory []string

// TODO - should have a mutex
func addToHistory(input string) {
	if strings.TrimSpace(input) == "" {
		return
	}
	inputHistory = append(inputHistory, input)
}

func renderTemplate(w http.ResponseWriter, uiData *UserInterfaceData) {
	w.Header().Set("Content-Type", "text/html;utf8")
	filename := "ui.tmpl"
	t, err := template.New(filename).ParseFiles("tmpl/" + filename)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	t.Execute(w, uiData)
}

// max = 4
// strs[] = 10
// we need 10..7
// max = 10
// strs[] = 4
// we need all 4..1
// if its only 1, we get 0,0, we need to process that as well
func strReverse(strs []string, max int) []string {
	var start, end int
	if len(strs) >= max {
		start, end = len(strs)-1, len(strs)-max
	} else {
		start, end = len(strs)-1, 0
	}

	var count = start - end + 1

	if start < 0 || count <= 0 {
		return make([]string, 0)
	}

	newstrs := make([]string, count)
	for i := 0; start-i+1 != end; i++ {
		newstrs[i] = strs[start-i]
	}

	return newstrs
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

	// TODO - add timeout after a command is sent
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
