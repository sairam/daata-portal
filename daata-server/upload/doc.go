/*
Package upload determines based on the parameters or existing directory structure
which functions to call.

1. Static files which are versioned and/or aliased identified by Type: 'VD'
2. Static files which are just uploaded
3. Data sent in form of key/value for graphing
4. Static files edited in UI (via markdown etc., to be updated in place)

Static:
curl -i -X POST -H "Content-Type: application/zip" --data-binary "@data.zip" https://my.daata.xyz/u/docs/spokes-platform

Versioned:
curl -X POST -H 'X-Version: "2.1.3"' -H 'X-Alias: release-20160707, master, stable' -H 'Content-Type: application/zip' --file-binary="@filename.zip" https://my.daata.xyz/docs/spokes-platform

fetching the data, attributes based on headers, file extension type, parsing, etc.,

*/
package upload
