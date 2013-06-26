Go HTTP Tutorials
=================

Various code snippets for writing http services in Go using the `net/http` package, including examples on how to:


#### Static Files

* Serve static files using `http.FileServer`


#### Reponse Headers

* Return a 404 or custom status code
* Setting content (MIME) type


#### Templates

* Nest templates using Go's `html/template` package
* Change the default variable delimiters.
* Inject safely escaped CSS and JS into template variables


In Progress
-----------

#### Cookies

* How to set cookies on a response
* How to use cookies to maintain sessions data (insecure)


#### Multiplexer

* Issues with Go's default route multiplexer 


#### Upload Data

* How to POST form data to a server 
* How to read POST form data
* How to upload a file to a server
* How to read an uploaded file
* How to save an uploaded image