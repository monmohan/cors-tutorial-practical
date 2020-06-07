## EDIT : This has been fixed and now both the buttons will show same behavior :)
See [Chrome Bug Report] (https://bugs.chromium.org/p/chromium/issues/detail?id=580796)
and https://github.com/whatwg/fetch/issues/204

### Reproduce browsers failure to follow redirects in case of Pre-flighted CORS requests

# Running the repro servers
* Download and install [golang[ (https://golang.org/dl/)
* Clone the repo
* Running the API and page server
    + Map couple of loopback interfaces to page and apiserver. For example here is my etc/hosts file
    ```
       127.0.0.2 pageserver.cors.com
       127.0.0.3 apiserver.cors.com
    ```
    + Change to the `issue` directory
    + Type `go run apiserver_303.go`
    + This will run the api server
    + Type `go run pageserver.go`
    + This will start the page server
    + Go to `http://pageserver.cors.com:12345/redirectfails.html`
    + This will show a page with two buttons, one with successfull pre-flighted request
    + The other with a failed 303 response (same URI)

