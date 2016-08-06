### Reproduce browsers failur to follow redirects in case of Pre-flighted CORS requests

# Running the repro servers
* Download and install [golang[ (https://golang.org/dl/)
* Clone the repo
    `git clone git@github.com:monmohan/cors-experiment.git`
* Running the API and page server
    + Change to the `issue` directory
    + Type `go run apiserver_303.go`
    + This will run the api server
    + Type `go run pageserver.go`
    + This will start the page server
    + Go to `http://localhost:12345/redirectfails.html`
    + This will show a page with two buttons, one with successfull pre-flighted request
    + The other with a failed 303 response (same URI)

