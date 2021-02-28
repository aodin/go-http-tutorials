A Simple HTTP Server
====

A simple HTTP web server can be created in Go by writing a function that accepts two parameters with the types `ResponseWriter` and `Request`, both from the `net/http` package in the standard library.

```go
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Go!")
}
```

This function is added to the default Go request multiplexer by passing it to the package function `HandleFunc` along with a given path, in this case the root directory `/`.

```go
http.HandleFunc("/", handler)
```

The server is then started by calling the package function `ListenAndServe` with a listening address and `nil`. Since this function may return an error, for instance if the listening address is already in use, it should be handled. In this case, it is logged before exiting the program:

```go
log.Fatal(http.ListenAndServe(":8080", nil))
```

Running [this short program](handler.go) produces no output, but will continue to run without exiting. It is waiting for a request, which we can send using a command like [`curl`](https://curl.se) or by visiting the listening address in a browser.

```console
user:~$ curl localhost:8080
Hello, Go!
```

This server will continue to run and listen for requests until we tell it to exit.


## More Information Than You Require

It may only takes a few lines of Go to create a functioning web server, but we've skipped over an incredible amount of detail about how this program actually works.

Let's start with the `handler` function. Its first parameter type, `ResponseWriter`, is actually an interface with three required method signatures:

```go
type ResponseWriter interface {
	Header() Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}
```

In the standard library `net/http` server, this interface is implemented by an unexported `response` type.

We'll discuss `Header` and `WriteHeader` in later sections and focus on `Write` for now.

The signature of `Write` matches the sole required method signature of `io.Writer`. Therefore, any type that implements the `ResponseWriter` interface also implements `io.Writer`, and any function that uses `io.Writer` can also use the `ResponseWriter`. This cross-compatibility allows us to write the response content with the `fmt` package's [`Fprintf`](https://golang.org/pkg/fmt/#Fprintf) function, which takes an `io.Writer` and a formatting string along with any number of format arguments.

Instead of using `Fprintf`, we could have written the string using the response's `Write` method after converting the string to a byte slice:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Go!"))
}
```

The second parameter of the handler is a `Request`, specifically a pointer to a struct type. It is a pointer to improve performance, since a request, especially its content body, can be quite large. The use of a pointer has nothing to do with mutating the request, as the [standard library specifically mentions](https://golang.org/pkg/net/http/#Handler) that the given `Request` should not be modified. We'll discuss the `Request` struct and its fields in later sections.

The signature of the handler function itself is also important, and is actually declared as an exported type in the `net/http` package.

```go
type HandlerFunc func(ResponseWriter, *Request)
```

Interestingly, this type has a [single method](https://golang.org/pkg/net/http/#HandlerFunc.ServeHTTP) wherein it calls itself.

```go
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```

This roundabout function call allows arbitrary functions matching the `HandlerFunc` signature to implement a generic `Handler` interface.

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

In our example web server, the conversion from an arbitrary function to a `Handler` occurs when we add it to default multiplexer with `http.HandleFunc`. Attempting to add a handler function of a different signature will result in a compile-time error.

We should note that adding a handler function isn't required. The web server will compile and run without adding anything to the default multiplexer. The server will be essentially useless, however, since all requests will return a `404 page not found`.

Finally, to run our server, we call `ListenAndServe` with a listening address. The function also takes a second parameter that allows us to specify a specific `Handler` interface to use for the server. By passing `nil`, we tell the server to use the default Go multiplexer with our now attached function.

Since `ListenAndServe` accepts any `Handler` as its second parameter, we could also do the type conversion and interface implementation ourselves using our existing handler function, completely skipping the default multiplexer.

```go
func main() {
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(handler)))
}
```

Since no multiplexer is used, a path is not required, and all requests will use our handler function, regardless of specified path.

There is little functional difference between the above and our example server, however, since the behavior of the default multiplexer is to use the handler added to the root directory `/` unless a handler has been attached to a more specific path.

Next: [Requests](/02_request/)
