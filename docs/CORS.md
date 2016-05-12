###Cross-Origin Resource Sharing (CORS)
 is a W3C spec that allows cross-domain communication from the browser. CORS is becoming increasingly more important as we use multiple API's and services to create a mashup/stitched user experience

But, in order to understand cross origin resoure sharing, first we need to understand the concept of an "origin".

###What is an Origin?
Two pages have the same origin if the protocol, port (if one is specified), and host are the same for both pages. 
So 
- http://api.autodesk.com/resource.html
has same origin as
- http://api.autodesk.com/somepath/resource2.html
but different from 
- http://api.autodesk.com:99/resource.html (different port)
or
- https://api.autodesk.com:99/resource.html (different protocol)
There are some exeptions to the above rule (mostly by, suprise surprise IE !) but they are non-standard.

###Same Origin Policy
By default, Browsers enforce "Same Origin Policy" for HTTP requests initiated from within scripts. A web application using XMLHttpRequest could only make HTTP requests to its own domain. One important thing to be aware of is that cross orgin "embedding" is allowed. Browsers can load scripts(source),images, media files embedded within the page even if they are from a different origin.

In this blog we will focus on the main restriction, cross origin requests using XMLHttpRequest

### Enter CORS
The Cross-Origin Resource Sharing standard (https://www.w3.org/TR/cors/) works by adding new HTTP headers that allow servers to describe the set of origins that are permitted to read that information using a web browser. Important thing to note is that its the Servers which are in control, not the client.

### Code Samples 
All the code shown in the blog is available at [LINK]
The server code is written in golang and the client samples use XMLHttpRequest Object (javascript/html). Although the code is in golang, the reader doesn't require knowledge of the language to understand what's going on. Its fairly obvious.
You can either build the code from source (See Readme) or download the binaries from here(OS X only) [LINK]. If you are interested to give it a spin, and would like to get binaries for any other platform feel free to reach out to me 

==Example 1
So lets first see what happens when we do a cross origin XMLHttpRequest. For this example, we will be running two servers.

PageServer : A simple server which serves requested page. 
	This server runs on a given port and serves an HTML file

    var port = flag.Int("port", 10001, "help message for flagname")

    func fileHandler(w http.ResponseWriter, r *http.Request) {
      fmt.Printf("Requested URL %v\n", r.URL.Path)
      http.ServeFile(w, r, r.URL.Path[1:])
    }

    func main() {
      flag.Parse()
      http.HandleFunc("/", fileHandler)
      log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), nil))
    }

- Start the pageserver on port 12345

> $pagesever -port 12345
 
 Apiserver : This is a simple server that looks up a user sent in the URL request and returns JSON data 

    var userData = map[string]User{
      "john": User{"jdoe", "John", "Doe", "France"},
    }
    var port = flag.Int("port", 10001, "help message for flagname")

    func userHandler(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "application/json")
      b, _ := json.Marshal(userData[r.URL.Path[len("/users/"):]])
      io.WriteString(w, string(b))

    }

    func main() {
      flag.Parse()
      http.HandleFunc("/users/", userHandler)
      log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), nil))
    }  
  
Run the apiserver

>$ apiserver -port 12346

- Open the browser and load the html http://localhost:12345/showuser.html
  Here is how this looks
![ShowUser](https://raw.githubusercontent.com/monmohan/cors-experiment/master/docs/showuser.png)

if you click "show", it is supposed to go to http://localhost:12346/users/john and get the user json to display but instead you see this error in console :

>showuser.html:1 XMLHttpRequest cannot load http://localhost:12346/users/john. No 'Access-Control-Allow-Origin' header is present on the requested resource. Origin 'http://localhost:12345' is therefore not allowed access.

This is called a "Simple" Cross origin GET request
Simple requests are requests that meet the following criteria:
HTTP Method matches one of
- HEAD, GET or POST
and 
HTTP Headers matches one or more of these
- Accept
- Accept-Language
- Content-Language
- Content-Type, but only if the value is one of:
- application/x-www-form-urlencoded, multipart/form-data, text/plain

Lets see what we can do to succeed in serving a simple cross origin request

- Stop the apiserver 
- Start the apiserver\_allow\_origin server. 

> $ apiserver\_allow\_origin -port 12346

What we have done here is added the _Access-Control-Allow-Origin_ header for any incoming GET request. The value of the header is same as the value sent by browser for the Origin header in the request. This is equivalent to allowing requests that come from any origin (*)

    func corsWrapper(fn func(http.ResponseWriter, *http.Request)) httpHandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            fmt.Printf("Request Origin header %s\n", origin)
            if origin != "" {
                w.Header().Set("Access-Control-Allow-Origin", origin)
            }
            fn(w, r)
        }
    }


Lets attempt clicking the "show" button again and Voila we see the data returned by the server:

> {"UserName":"jdoe","FirstName":"John","LastName":"Doe","Country":"France"}

Its all good until we realize that just adding Access-Control-Allow-Origin isn't sufficient for certain "complex" requests (or anything which isn't covered in the Simple request). 
An example of such a request is POST request with Content-Type as application/json.

Lets understand this with another example.
Point your browser to http://localhost:12345/createUser.html
This is a simple form which looks like below. Entering the data and clicking "create" send a POST request to the ApiServer in-memory store
Here is how this looks
![CreateUser](/path/to/img.jpg)

So, lets add some string data in the form fields and click "create" button. This should convert the data to JSON and do a POST to http://localhost:12346/users with the json data as the body of the request 
Here is the relevant code in createUser.html:
 > function sendRequest(url) {

        var oReq = new XMLHttpRequest();
        oReq.addEventListener("load", reqListener);
        oReq.open("POST", url);
        oReq.setRequestHeader("Content-Type", "application/json")
        var data = serializeUser($('#fcreate').serializeArray());
        console.log(data)
        oReq.send(JSON.stringify(data));
    }

But once you hit, "create", the browser reports the following error :-

>XMLHttpRequest cannot load http://localhost:12346/users. Response to preflight request doesn't pass access control check: No 'Access-Control-Allow-Origin' header is present on the requested resource. Origin 'http://localhost:12345' is therefore not allowed access.

Pre-Flight
Non-Simple requests actually cause two http requests under the covers. The browser first issues a preflight or an OPTIONS request first, which is basically asking the server for permission to make the actual request. Once permissions have been granted, the browser makes the actual request. 

So, in this case, the pre-flight request is something like below
*OPTIONS /users HTTP/1.1*    
Host: localhost:12346    
Connection: keep-alive   
*Access-Control-Request-Method: POST*   
*Origin: http://localhost:12345*    
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36   
*Access-Control-Request-Headers: content-type*   
Accept: */*   
Referer: http://localhost:12345/createUser.html    

The preflight request contains a few additional headers:
Access-Control-Request-Method - The HTTP method of the actual request. 
Access-Control-Request-Headers - A comma-delimited list of non-simple headers that are included in the request. Notice that all CORS related headers are prefixed with "Access-Control-". 

In order for the POST to succeed, the server should support this request, "granting" permission based on the above request headers.
Lets do that.
- Stop apiserver_allow_origin
- Start apiserver_preflight

> $ ./apiserver_preflight -port 12346

What we have done here is added some code in the apiserver to respond to OPTIONS request, granting the permission for GET, POST and OPTIONS calls with Content-Type as application/json.

> func optionsWrapper(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "OPTIONS" {
            //set other headers
            w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
            return
        }

       fn(w, r)
   }
}

Enter data and hit "create" button again. You will see that the request succeeeded. Using chrome tools or similar debugger, the response to OPTIONS request can be examined as well. 

>HTTP/1.1 200 OK
*Access-Control-Allow-Headers: Content-Type*
*Access-Control-Allow-Methods: POST, GET, OPTIONS*
Access-Control-Allow-Origin: http://localhost:12345
Date: Thu, 12 May 2016 10:10:13 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8 

The response headers from the server grant permission to the different cross origin request methods (comma separated list of GET, POST, OPTIONS) and also the allowed headers (in this case Content-Type header).
In addtion, the server can also return a header called Access-Control-Max-Age. The value of the header indicates how long the pre-flight response can be cached by the browser and hence browsers can skip the check for that duration

### Handling credentials
By default, cookies are not included in CORS requests. This means that a cookie set by one origin will not sent as part of the HTTP request sent to the different origin. 
Lets see an example 
- Stop apiserver_preflight
- Start apiserver_creds_fail 
- Stop pageserver
- Start pageserver_cookie

Point your browser to http://localhost:12345/showusermore.html. 

The UI is same as showuser.html but the pageserver_cookie server now adds a cookie (name="token", value="secret_token")to the page when its served. 

>func fileHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Requested URL %v\n", r.URL.Path)
    http.SetCookie(w, &http.Cookie{Name: "token", Value: "secret_token"})
    http.ServeFile(w, r, r.URL.Path[1:])
}

Also, the apiserver will attempt to read this cookie, and respond with additional secret data.

>func userHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    b, _ := json.Marshal(userData[r.URL.Path[len("/users/"):]])
    io.WriteString(w, string(b))
    //respond with secret data when cookie secret token is recieved
    if c, err := r.Cookie("token"); err == nil && c.Value == "secret_token" {
        io.WriteString(w, "<br/>Show Secret Data !!")
    }

}

Enter "john" in the text box and hit "show". The request doesn't succeed!
You will see following error in the console

>XMLHttpRequest cannot load http://localhost:12346/users/john. Credentials flag is 'true', but the 'Access-Control-Allow-Credentials' header is ''. It must be 'true' to allow credentials. Origin 'http://localhost:12345' is therefore not allowed access.

What happened here is that page tried to send the cookie to the different origin API server. Here is the sendRequest method from page

 >function sendRequest(url) {
        var oReq = new XMLHttpRequest();
        oReq.addEventListener("load", reqListener);
        oReq.withCredentials = true;
        oReq.open("GET", url);
        oReq.send();
    }

Notice the "oReq.withCredentials = true;" statement. The XMLHttpRequest object needs to set a property called "withCredentials" in order to share the cookie to the different origin server. However that's not enough. 
Remember the OPTIONS request? The server should have responded with a header called _Access-Control-Allow-Credentials_ with value as true in order for this cookie to be accepted. This request header works in conjunction with the XMLHttpRequest property. If .withCredentials is true, but there is no Access-Control-Allow-Credentials header, the request will fail (and vice versa).

Lets try again
- Stop apiserver_creds_fail
- Start apiserver_allow_creds

> $ ./apiserver_allow_creds -port 12346

What we done now is added support for _Access-Control-Allow-Credentials_  header
>func credsWrapper(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        fn(w, r)
    }
}

Again enter "john" in the text box and hit "show".
You will see the following response with the secret data text, showing the 
>{"UserName":"jdoe","FirstName":"John","LastName":"Doe","Country":"France"}
Show Secret Data !!

Hopefully this has given a hands on experience with supporting CORS.
To read more on this subject, please take a look at the links in the reference section.

References
- [Browser Security Handbook]
(https://code.google.com/archive/p/browsersec/wikis/Part2.wiki#Same-origin_policy)
- [W3C Cross-Origin Resource Sharing]
(https://www.w3.org/TR/cors/)
- [HTML5 Rocks]
(http://www.html5rocks.com/en/tutorials/cors/#toc-introduction)
- [CORS on MDN]
(https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS)