## Done so far
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
6. Connect to multiple hosts
7. Add UI components
  1. bootstrap css
  2. Add sections for hosts
  3. Add UI for command (text area)
  4. Add UI for stdout/stderr
  5. Page with hosts information and connection in green
  6. Execute command button to run on multiple hosts
8. Screen to add new hosts and save configuration
9. Pass same UI configs via CLI
10. Option to download zip file of output-hostname as well as script executed
11. Option to merge data by simple operations like concat, regexp match via UI

## Wishlist
1. Save complete session output
1. Display hosts which we are unable to connect grouped by counts
1. Favorites of commands typed to 1-click execute or `uptime` etc., - Handle commands like `top`
1. Aggregated / Group output by hosts or output split
1. Open a console for one of the hosts
1. Save workflow (the complete execution that happened like ssh, commands run, piping output to localhost and passing it back to different set of hosts) - this actually looks like a shell script
1. Button to 'Upload' data to server
1. Run command / script on X of Y hosts at a time - Case for deployments to maintain service up
1. Metrics via Javascript client to GA/other location via CLI server w/o violating privacy
MVP - no workflows, no fancy stuff, let something work!
