# Step by Step Tutorial
* The code in this repository is used to demonstrate CORS concepts like - Simple GET with Allow-Origin, requests needing pre-flight and handing credentials
* Step by Step tutorial using this code is available as three part blog series
* [Part I](https://medium.com/@software_factotum/cross-origin-resource-sharing-a-hands-on-tutorial-fb19748cb3b7)
* [Part II](https://medium.com/@software_factotum/cross-origin-resource-sharing-a-hands-on-tutorial-part-ii-complex-requests-ed5be46fadcf)
* [Part III](https://medium.com/@software_factotum/cross-origin-resource-sharing-a-hands-on-tutorial-part-iii-cookies-a60ecbee0983)


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
- Follow along the Blogs or play around with the code .

Have Fun !
