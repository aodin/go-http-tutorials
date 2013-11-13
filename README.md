Go HTTP Tutorials
=================

Various code snippets for writing http services in Go using the `net/http` package, including how to:


### Handlers

* Create a Handler that can wrap another Handler


#### Reponse Headers

* Return a 404 or custom status code
* Setting content (MIME) type
* Set and retrieve a Cookie


#### Request Data

* Upload and read POST form data


#### Static Files

* Serve static files using `http.FileServer`


#### Templates

* Nest templates using Go's `html/template` package
* Change the default variable delimiters
* Inject safely escaped CSS and JS into template variables
* Write the correct template syntax for iterators


In Progress
-----------

#### Config

* Use flags to accept a custom server address


#### Multiplexer

* Issues with Go's default route multiplexer
* Write your own route multiplexer


#### Session

* How to use cookies to maintain sessions data (insecure)


#### Upload Data

* How to upload a file to a server
* How to read an uploaded file
* How to save an uploaded image