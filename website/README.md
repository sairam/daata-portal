
## Validation
* Start the server in development mode `hugo server`
* Test if all links are valid
```
cd /tmp/
wget --recursive http://localhost:1313/
echo $?
rm -rf localhost:1313
```

## Deploying
* `hugo`
* `rsync -avrz --delete public/ deploy@daata.in:/var/websites/daata.in/`

### TODO
* Upload to private `daata server` to a directory (aka Eat your own Dog Food)
