/*
Package redirect allows you to add short endpoints for long URLs under /r/

This is your typical bitly/is.gd/goo.gl/tinyurl without the analytics

Any request without a "short_url" will auto generate a "short_url" with alphanumeric of 6 character length which is present in the response

Restrictions

"long_url" length can no longer be more than 1k characters

"long_url" is limited to http and https protocols (data and ssh protcols are not allowed)

"short_url" can no longer be more than 256 characters

"override" only accepts true to override. (setting it to 1 does not set override flag to true)


Try out Examples

Typical input

  curl -D - -X POST -H "Content-Type: multipart/form-data" -F "short_url=yelo" -F "long_url=https://www.google.com" https://example.daata.in/r/


Visit https://example.daata.in/r/yelo to redirect to https://www.google.com

Visit https://example.daata.in/r/yelo+ will display a link to https://www.google.com

=====================================================================================

Works with utf8 character set

  curl -D - -X POST -H "Content-Type: multipart/form-data" -F "short_url=😉" -F "long_url=https://www.google.com" https://example.daata.in/r/

Visit "https://example.daata.in/r/😉+"

=====================================================================================

Use "override=true" to overwrite any existing redirect already set.

  curl -D - -X POST -H "Content-Type: multipart/form-data" -F "short_url=yelo" -F "long_url=https://www.google.com" -F "override=true" https://example.daata.in/r/

=====================================================================================

Ignore "short_url" to auto generate

  curl -D - -X POST -H "Content-Type: multipart/form-data" -F "long_url=https://www.google.com" https://example.daata.in/r/

*/
package redirect
