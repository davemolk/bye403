# bye403
bypass 403 responses (or try to)

## Overview
By default, bye403 uses a combination of path, header, and method manipulation in the attempt to bypass a 403 response. Any non-403 response is considered a success, although you can filter by status code to determine what gets printed to stdout.

## Examples
bye403 has a default url of https://www.example.com/secret. Here's a sample of what's printed to stdout without adding any additional flags:
```
404: GET request
http://www.example.com/secret/.

404: GET request
https://www.example.com/secret
headers: [X-Remote-IP localhost:443]

501: REBIND request
https://www.example.com/secret

404: GET request
http://www.example.com/%2e/secret

404: GET request
https://www.example.com/secret
headers: [X-Remote-IP localhost:80]

404: POST request
https://www.example.com/secret

404: GET request
https://www.example.com/secret
headers: [X-Remote-IP 127.0.0.1]

404: GET request
https://www.example.com/secret
headers: [X-Remote-IP 127.0.0.1:80]

404: GET request
https://www.example.com/secret;/
```

## Install
First, you'll need to [install go](https://golang.org/doc/install). Then, run the following command:

```
go install github.com/davemolk/bye403/cmd/bye403@latest
```

## Flags
```
Usage of bye403:
  -c int
    	max number of goroutines to use at any given time
  -h bool
    	manipulate headers
  -i bool
    	read url off stdin
  -ignore string
    	status code responses to ignore (403 is ignored by default)
  -insecure bool
    	accept any certificate and host name presented by server
  -m bool
    	manipulate http method
  -os string
    	operating system (used in request header creation)
  -p bool
    	manipulate path
  -proxy string
    	proxy to use
  -r bool
    	allow redirects
  -rh bool
    	include response headers in output
  -s bool
    	silent error reporting
  -sc string
    	filter output by status code(s)
  -t int
    	request timeout (in ms)
  -u string
    	target url
  -v bool
    	validate url before running program
```

## Note
Each request gets a randomly assigned user agent corresponding to your os as well as appropriate headers (50/50 chance of chrome or firefox). Go unfortunately doesn't preserve header order, so if that's important to you and what you're up to, you'll have to look elsewhere.

## Thanks
I looked at these, more or less ported them to Go, and added a bunch of new features. Enjoy!
* [byp4xx](https://github.com/lobuhi/byp4xx)
* [403bypasser](https://github.com/yunemse48/403bypasser)
