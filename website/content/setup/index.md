+++
title = ""
description = ""
+++

Should contain configurations like public/private setups.

Check options present in the toml file etc., so that we can spit out the toml file here.
some may want to use uploads, some redirects, some dashboards or based on operating system etc.,

Workflow:
1. User comes in. Enters email id or verifies user via robot authentication.
2. User copies the command from the browser, runs the bash script which takes in the uname required to generate the go binary
3. This will be of the format https://setup.daata.in/machines/random-api-key.txt
4. We pull the same data off in the browser after the user clicks 'Check'
5. We generate the next command so that he can copy/paste once he has configured the setup.

2. We get him an API key and a url https://setup.daata.in/email-name/file-name.txt
3. API key can be used to ping our server back data like uname etc., for the build to be pulled.
4. We use the above data/file which we can use to spit out the executable to be pulled.

TODO - nice to host an API key based version
