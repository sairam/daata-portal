+++
title = "Paste Bin"
description = "A regular Pastebin"
see_also = ["append-data", "send-json-on-rest"]
+++

## Regular Paste Bin
* `any` data in any format can be sent along with a filename.

### Input
```
curl -X POST --data=demo.txt  $portalURL/test/filename.txt
```
### Output
