## Start Coding HERE in this order and `write tests`
1. Create new (upload) entities to store information
1. Unzip / Host files at location based on url endpoint and version
1. Store data points in a file by time
1. Make a graph with data points
1. Append to file based on name (per directory from multiple hosts)
1. Create / Edit (`txt`/`md`) files (<10KB) and save back from UI
1. UI features for `grep`, `split`, `prepend`, `append`, `clean` data etc.,
1. Create one Organisation and User entity
1. User Auth with `Google`/`Github` Auth Login only
1. Authorization token by user for all entities
1. Public website and Dashboard development in parallel
1. Webhooks

## Workflows
### Upload documentation / Release Notes
1. Create a permalink with name "Spokes Platform" (Note: A zip or gz or bz2 file is sent, we automatically de-compress it)
1. Headers are sent as versions as array
  `curl -D - -X POST -H "X-Version: 0.10" -H "X-Alias: master,default" -H "Content-Type: application/zip" --data-binary "@data.zip" localhost:8001/u/snapdeal/mvp/`
  `curl -X POST -H 'X-Version: 2.1.3' -H 'X-Alias: release-20160707,master,stable' -H 'Authorization: "abcdefghijklmnopqrstuvwxyz"' --file-binary="@filename.zip" -H 'Content-Type: application/zip' https://my.daata.xyz/docs/spokes-platform`
1. Main page of https://my.daata.xyz/docs/project-name will contain all the list of recently uploaded data (index page)
1. Triggers webhook for actions

### Code Coverage
1. Create a 'Data Point' project (enable History) with name "Spokes Platform" under Code Coverage (code-coverage)
1. Create data point name as "coverage" and select "number". (Use this template and modify the name)
1. Headers are sent as versions as array
  curl -X POST -H 'Authorization: "abcdefghijklmnopqrstuvwxyz"' -H 'Content-Type: application/vnd.datapoint' --data="coverage,89.9,time" https://my.daata.xyz/code-coverage/spokes-platform
1. If additional data points check is not present, extra data will be ignored.

### Data/Information from database/url through a Cron/Deployment Status
1. Create a 'Data Point' project with name "Shipments Delivered" under Ekart (ekart)
1. Allow any data to be sent. Useful when new keys come in or structure is not well defined (Billing should be based on no. of data points - like stathat.com)
1. Data is sent as file with json content (as data points)
  curl -X POST -H 'Authorization: "abcdefghijklmnopqrstuvwxyz"' -H 'Content-Type: application/json' --data="@file.json" https://my.daata.xyz/ekart/shipments-delivered
1. Visualization on Table based on type of Hash sent !!!! this is not clear
1. (Will require versions or latest information in case of delays in publishing or backfilling)

### Data Points by Host
1. During a deployment, only few have Visualization tools to see the status of the deployments.
1. New Deployment Project .
1. Hosted on `hostname`. If docker, find unique name.
1. Different if blue/green deployment vs rolling deploy.
1. tool deploy --progress 0 [| start] --hostname=hostname project_name --type rolling --version="1.234" # links to documentation
1. tool deploy --progress 10  ... # if you have steps
1. tool deploy --progress 100  ...
1. On start/end, triggers can be set for email/slack/webhooks to teams.

### Sending a static file(s) of text/html/image/directory/log files
1. If its under a project name, file(s) are just pushed into the master. No versioning allowed like docs
1. Multiple files can be added to the directory and files are overridden. All names should be uniq or will be made uniq if conflicts occur.
1. Multiple files sending the same filename will be appended with "/a" flag to avoid overriding
1. `curl -D - -X POST -H "X-Append: true" -H "X-File-Name: index.txt" --data-binary "@index.txt" localhost:8001/u/example/sample/`
1. Your standard S3 file with an append mode on locks which will be merged as soon as file uploads start/end !!!!!! what?
1. Use the UI to uniq/sort/filter/append/prepend data on the processed file using Javascript

### Create snippets
1. Create snippets with a markdown/WSIWYG editor and link snippets etc.,
1. Embed graphs or dashboards !!!!
1. Publish as a page with information so that you can showcase tech metrics of your company to the rest of the world!
1. Send Daily snapshots of dashboards to your customers/stakeholders

### Creating a Dashboard
1. 'New Dashboard'
1. 'New Graph' or 'New Data Point' or 'New Text' (from local data or remote data API like Graphite or Graphana (client side))
1. Select from the UI to make a graph with the data point using auto complete and time range.
1. Resize/Move the widgets to arrange on the screen resolution.
1. A text widget has header, text from the data point.
1. External widgets like images/text can be refreshed by frequency
1. UI may need data like MIN, MAX, UNIQUE, COUNT, SUM, AVERAGE, P99, P95 etc.,

## Main Website
* www.daata.xyz will contain the static website  (see if you can generate with hugo with a theme?)
* /blog/...
* /about
* /usage
* /use-cases/ - another blog / detail with data like how ifttt does with recipes
* /how-it-is-being-used/ - testimonials from customers
* /innovative-dashboard/
* /pricing
* /{feature}
(think about SEO keywords to use here)

## Types of URLs
1. Company accounts - has subdomain like google.daata.xyz/repo/url
1. Personal accounts - at my.daata.xyz/repo/data.txt
1. Repositories are usernames at my.daata.xyz. anyone can create a new data point. The owners are defined based on the repo.

1. /temp/ is a different repo based on files which go away in few hours/days
1. /a/ or /data is a permanent default url.
1. /r/ for redirects

## AuthS & AuthZ
1. Auth S is based on domain and you can whitelist users by domain and add others and ensure they have Authorization per project/endpoint.
1. Auth Z Token authorization based on -
  1. Upload anything into your organisation
  1. Upload to Data Point
  1. Upload to Project?
  1. Upload allowed to Group (useful for build systems to push anything to specific repos)
  1. Upload means 'w'/'a', but cannot read the data
1. Reads are same, everything, data, project/group, history.

## By User
* Everyone will login through subdomain.daata.xyz/auth or my.daata.xyz/auth. If a domain is authenticated
* Every user in every organisation has an account. Data can never move across organisations.
* An employee is part of multiple teams/projects and a Company has many teams
* This is required for authorization by project/team for sensitive information
* Freeze / Snapshot capability
* Revision history for overriding values (as a pro feature)
* For hosting repos/data, look at gogs UI(MIT)
* Hooks/Webhooks can be created to trigger website updates


## TODO
### Upload
1. Access restriction to write/read
1. If unzip file does not contain index.html, generate one with the tree (already being done when displaying the directory)
1. Paginate if there are too many files/directories
1. Generate template with bootstrap
1. Streaming can be a Enterprise feature since it requires authentication to encrypt traffic
1. Graphs - line, column, area, bar, table, json -- options to group by (keen.io)
1. Graphs - add timeframe - relative/ absolute with TZ (keen.io)
