+++
title = "URL Shortner"
description = "Shorten and Redirect URLs with a named shortner like bitly or tinyurl"
warning = "Target/Long URL is limited to only <kbd>http/s</kbd> protocols for reasons of security"
+++

### Consider this a URL Shortener or Redirector

<kbd>portalURL={{% demo %}}</kbd>
```
curl -X POST --data='{shortURL:"google", longURL:"https://www.google.co.uk"}'  $portalURL/r/
```

#### Sample Output
<samp>fdask fdakl fsa </samp>
