# cors-experiment
* Simple GET with Allow-Origin
* POST with application/json content type, support for OPTIONS
* Share cookie, xhrCredentials

# Running the demo
* Download and install [golang](https://golang.org/dl/)
* Clone the repo
    `git clone git@github.com:monmohan/cors-tutorial-practical.git`
* Running the Servers
    + Map couple of loopback interfaces to page and apiserver. For example here is my etc/hosts file
    ```
       127.0.0.2 pageserver.cors.com
       127.0.0.3 apiserver.cors.com
    ```
    + The servers can be run directly from the source
    + Change to the relevent directory
    + Type `go run <filename>`
- For example to run the apiserver_preflight, 
    + cd to apiserver directory
    + `go run apiserver_preflight.go`

Follow CORS.MD to try out the examples in the blog article
