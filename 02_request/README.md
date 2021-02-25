Requests
====

Our simple HTTP server doesn't do much. It has only one handler, and simply returns the same pre-determined string at every path for any type of request. Let's add some functionality, with a focus on the [`Request` type](https://golang.org/pkg/net/http/#Request).

There's a number of interesting properties available on the `Request`. Here's a few of those properties with example values that you might see when interacting with the server on localhost:

| Property       | Accessor                    | Example                           |
| -------------- | --------------------------- | --------------------------------- |
| URL            | .URL                        | /path/item                        |
| Method         | .Method                     | POST                              |
| Remote Address | .RemoteAddr                 | [::1]:56310                       |
| Host           | .Host                       | :8080                             |
| User Agent     | .Header.Get("User-Agent")   | curl/7.64.1                       |
| Content Type   | .Header.Get("Content-Type") | application/x-www-form-urlencoded |
| Form Data      | .Form                       | map[name:[Go]]                    |

Many of these properties are also Go types with further functionality, or have additional accessors that simplify their usage. We'll examine many of these in this and future tutorials.

One way of looking at the `Request` type and its fields is a structured type that the Go standard library has parsed from incoming unstructured bytes.


## Well... how did I get here?

Let's examine that process of parsing further with the help of a [more advanced server](request.go). We'll continue using `curl`, but this time in a verbose mode.

```console
user:~$ curl -v :8080
> GET / HTTP/1.1
> Host: :8080
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Fri, 05 Feb 2021 01:54:36 GMT
< Content-Length: 10
< Content-Type: text/plain; charset=utf-8
<
What is your name?
```

All lines beginning with `>` are header data sent by `curl`, with `<` preceding received header data. For clarity, I've omitted some lines of additional data provided by `curl`.

The server parses this incoming request by reading a buffer of bytes from the current connection. Most importantly, the `Request` type is constructed by the connection using an unexported function `readRequest` that will populate the type's fields with a series of parsers. Once the request has been constructed, the connection itself will set a few additional properties, such as `RemoteAddr`.

The parsing of the request data starts with an unexported function `parseRequestLine` that sets the method, URI, and protocol on the `Request` type from the first line of our HTTP request header. Subsequent lines of the header will all be parsed by the [`ReadMIMEHeader` method](https://golang.org/pkg/net/textproto/#Reader.ReadMIMEHeader) of `textproto.Reader`.

These parsers can error when given poorly constructed data. For example, when we remove the request target from the first line of the header.

```console
user:~$ curl -v --request-target "" :8080
> GET  HTTP/1.1
> Host: :8080
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 400 Bad Request
< Content-Type: text/plain; charset=utf-8
< Connection: close
<
400 Bad Request
```

The header lines parsed by `ReadMIMEHeader` will have their names standardized into a canonical format that capitalizes the first character and any character following a hyphen. Values are then aggregated within the header map using these canonical names. Additionally, accessing a name via the header's `Get` or `Values` method will also apply the canonical format to the name. `Get` will return only the first value, while `Values` will return all values as a slice.

We can see how Go formats header names into the internal map representation using `curl`.

```console
user:~$ curl -v -H "Special: A" -H "special: B" :8080
> GET / HTTP/1.1
> Host: :8080
> User-Agent: curl/7.64.1
> Accept: */*
> Special: A
> special: B
```

    map[Accept:[*/*] Special:[A B] User-Agent:[curl/7.64.1]]

And how the header values vary depending on how we access them.

| Accessor                   | Value | Type     |
| -------------------------- | ----- | -------- |
| r.Header.Get("Special")    | A     | string   |
| r.Header.Get("special")    | A     | string   |
| r.Header.Values("special") | [A B] | []string |
| r.Header.Values("special") | [A B] | []string |
| r.Header["Special"]        | [A B] | []string |
| r.Header["special"]        | []    | []string |

It is possible to send too much header data to the Go server. By default, the server will return a `431 Request Header Fields Too Large` error when the headers exceed the server's `MaxHeaderBytes`, which defaults to 1 MB, plus 4 KB of buffer.

In addition to headers, we can send information to the server using query strings and the request body. Query strings are parsed from the URL and placed into a [`Values` type](https://golang.org/pkg/net/url/#Values) that shares the same underlying Go builtin as our headers: `map[string][]string`. Although the rules governing query strings are complex, once parsed, there is no canonical format for the keys and are, therefore, case sensitive as the following example shows.

```console
user:~$ curl -v ":8080/?name=A&name=B"
> GET /?name=A&name=B HTTP/1.1
> Host: :8080
> User-Agent: curl/7.64.1
> Accept: */*
```

| Accessor                  | Value | Type     |
| ------------------------- | ----- | -------- |
| r.URL.Query()["name"]     | [A B] | []string |
| r.URL.Query()["Name"]     | []    | []string |
| r.URL.Query().Get("name") | A     | string   |
| r.URL.Query().Get("Name") |       | string   |

Note that the query string is parsed every time `Query()` is called, so these parsed values should be assigned to a variable if they will be repeatedly accessed.

Data sent via the request body is even more complex, so we'll focus on the common use case of sending form data via `POST`.

```console
user:~$ curl -v -d "name=Go" :8080
> POST / HTTP/1.1
> Host: :8080
> User-Agent: curl/7.64.1
> Accept: */*
> Content-Length: 7
> Content-Type: application/x-www-form-urlencoded
```

The above command will send a request to our Go server with the body `name=Go`. This body, however, will not be parsed automatically by the server. We must either access a value using the request's `PostFormValue` method or by manually calling either of the `ParseMultipartForm` or `ParseForm` methods. Note that this parsing is only done when the media type is `application/x-www-form-urlencoded` and a valid method is used, which `curl` has done for us.

As with headers, Go sets a sane default limit on the amount of data in the request body, in this case 10 MB. And just as with query strings, the data will be stored in a `Values` type after being parsed using [`ParseQuery`](https://golang.org/pkg/net/url/#ParseQuery).

There are countless other ways to send data to a server via the request body. We'll examine some of the more common encodings in further tutorials.

Next: [Response](/03_response/)
