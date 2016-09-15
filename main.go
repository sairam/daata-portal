package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World")
}
/*
  Phase 1
1. SSH into a host
2. execute a command
3. send the command output to stdout
4. Close the connection

Phase 2
0. Connect to hosts
1. Add a web server on bootstrap css
2. Open a default page with host information
3. Text Area for input to send (like date and TZ, uptime, user)
4. Execute the commands and display stdout and stderr
5. Repeat (3)
6. Exit/Close connection and web server

Next Phases (Simplify this):
1. CLI Configuration
2. Multiple hosts on CLI
3. Hosts file
4. Display hosts which we are unable to connect to w/ counts
5. Favourites of commands typed to 1-click execute or `uptime` etc., - Handle commands like `top`
6. Aggregated / Group output by hosts or output spit
7. Option to download zip file of output-hostname as well as script executed
8. Use a screen to execute/connect (not sure how to do this)
9. Option to merge data by simple operations like concat, regexp match, pipe output
10. Save workflow (the complete execution that happened like ssh, commands run, piping output to localhost and passing it back to different set of hosts) - this actually looks like a shell script
11. Button to 'Upload' data to server
12. Run command / script on X of Y hosts at a time - Case for deployments to maintain service up
00. Metrics via Javascript client to GA/other location via CLI server w/o voilating privacy
13. MVP - no workflows, no fancy stuff, let something work!
*/
