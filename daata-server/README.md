# Server

## `redirect`
This is your typical bitly/is.gd/goo.gl/tinyurl service without the analytics

The service resides at `/r/`

Any request without a `short_url` will auto generate a `short_url` with alphanumeric of 6 character length which is present in the response

### Restrictions

* `long_url` length can no longer be more than 1k characters
* `long_url` is limited to http and https protocols (data and ssh protcols are not allowed)
* `short_url` can no longer be more than 256 characters
* `override` only accepts true to override. (setting it to 1 does not set override flag to true)


### Usage

#### Example 1

```
curl -D - -X POST -H "Content-Type: multipart/form-data" -F "short_url=yelo" -F "long_url=https://www.google.com" https://example.daata.in/r/
```

Visit https://example.daata.in/r/yelo to redirect to https://www.google.com

##### Add `+` at the end of the url to display it

Visit https://example.daata.in/r/yelo+ will display a link to https://www.google.com

#### Example 2

Works with utf8 character set

  curl -D - -X POST -H "Content-Type: multipart/form-data" -F "short_url=he😉lo" -F "long_url=https://www.google.com" https://example.daata.in/r/

Visit "https://example.daata.in/r/he😉lo+"

#### Example 3

Use `override=true` to overwrite any existing redirect already set.

  curl -D - -X POST -H "Content-Type: multipart/form-data" -F "short_url=yelo" -F "long_url=https://www.google.com" -F "override=true" https://example.daata.in/r/

#### Example 4

Ignore `short_url` to auto generate

  curl -D - -X POST -H "Content-Type: multipart/form-data" -F "long_url=https://www.google.com" https://example.daata.in/r/


## `upload/display`

### Usage

#### Example 1
When you have a `tar`||`zip` file to unarchive and view.

**Usecase:** When someone shares a zipped version of content like html pages with fonts and css

```
curl -D - -X POST -H "Content-Type: application/zip" --data-binary "@new-react-prototype.zip" https://example.daata.in/u/static/react-prototype/
```
Your unarchived and uncompressed content will be hosted at https://example.daata.in/u/static/react-prototype/

#### Example 2

Upload zipped content in versions and create versions and link to latest aliases

**Usecase:** When you want to host documentation in a versioned manner and maintain multiple aliases by `branch` or `tag`

```
curl -D - -X POST -H "X-Version: 0.10" -H "X-Alias: master,default" -H "Content-Type: application/zip" --data-binary "@data.zip" https://example.daata.in/u/your/project/
```

#### Example 3

Upload any readable file to share with your team

```
curl -D - -X POST -H "Content-Type: application/json" --data-binary "@stackexchange.json" https://example.daata.in/u/static/dir/stackoverflow.json
```

#### Example 4

Upload a single file to view on a hosted service

**Usecase:** When you generate a heatmap profile from your code

```
curl -D - -X POST -H "X-File-Name: heatmap-20161010.svg" --data-binary "@heatmap.svg" https://example.daata.in/u/code/project/
```

#### Example 5

Send files from multiple hosts to append with each other

**Usecase:** When querying for logs/static data and needs to be passed merged. Can be run on all hosts in parallel

curl -D - -X POST -H "X-Append: true" -H "X-File-Name: index.txt" --data-binary "@index.txt" https://example.daata.in/u/company/logfile/

#### Example 6
Send output from `MySQL` or `Postgres` to visualize as smart sortable tables

#### Example 7
Send metrics to be plotted alongside time

If your usecase does not allow you to create new metrics in your Graphite/statsd/collectd system, use this as a metrics system to send information about cron usage/runtime.

Also, aggregate all metrics inside a project path to view it as a dashboard on your monitor/TV dedicated to metrics

**Usecase:** This can be your code coverage metrics in a CI box which does not natively have information about Code Coverage by branch

```
curl -D - -X POST -H 'Content-Type: application/vnd.datapoint+value' --data "twitter-stream,78,`date +%s`" https://example.daata.in/u/code-coverage/mycode/
```

## `display`
