# Step by Step Tutorial
* Step by Step tutorial using this code is available as three part blog series
* [Part I](https://medium.com/@software_factotum/cross-origin-resource-sharing-a-hands-on-tutorial-fb19748cb3b7)
* [Part II](https://medium.com/@software_factotum/cross-origin-resource-sharing-a-hands-on-tutorial-part-ii-complex-requests-ed5be46fadcf)
* [Part III](https://medium.com/@software_factotum/cross-origin-resource-sharing-a-hands-on-tutorial-part-iii-cookies-a60ecbee0983)

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
