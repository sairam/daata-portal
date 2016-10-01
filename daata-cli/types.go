package main

// Host is identified by the user connected and the hostname provided which resolves to the IP address
type Host struct {
	Username  string
	IPAddress string
	Hostname  string
}

// CommandExecution is a host with input and output
type CommandExecution struct {
	Host *Host
	PWD  string
	// TODO - Output struct should of stdout and stderr
	// Stdin  string
	// Stdout string
	// Stderr string
	Input        string
	Output       string
	TerminalAttr TerminalAttr
	// TODO - Take TerminalSize input from the UI so that you can fit in the output screen
}

// TerminalAttr should already be a struct
type TerminalAttr struct {
	Rows    int
	Columns int
	// Color string . Color is mono or 256. Output would be colored similarly
}

// UserInterfaceData is used by the UI
type UserInterfaceData struct {
	Details []CommandExecution
	// History [][]CommandExecution
	History []string
	Current string
}
