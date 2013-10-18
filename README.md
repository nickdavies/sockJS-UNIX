# SockJS-UNIX

SockJS-Unix terminates a sockJS connection and proxies all traffic to and from a UNIX socket.

It is written in node.js and depends only on the sockJS library itself.

## Why would you want to do that?

The reason why I wrote this was because I wanted to use SockJS with a backend server that did not have a SockJS library (for me it was Golang).
This allows any language that can talk via a unix socket (everything) to use SockJS.

## How do I use it?

You can run the server with:

    node server.js

Which will start a sockJS server listening on port `1334`

Then you can run the backend server with:

    go run example_server.go

Which will listen on the unix socket `/tmp/sockjs-unix.sock` 

Finally running:

    python2 -m SimpleHTTPServer

then visit `localhost:8000/test.html` and in the console run `sock.send(JSON.stringify({"body": "woo", "channel": "test"}));`

## Things to note

In my usage of SockJS I enforce that you must have a `body` and a `channel` (which is not normally something SockJS mandates)


# SockJS-UNIX Client

SockJS-Unix also provides a basic client that can do the oposite of the server translation.

The client listens on a UNIX socket (`/tmp/sockjs-unix-client.sock` by default) and 
will make connections to a given SockJS endpoint (`http://localhost:1334/' by default) whenever a client connects

You can use this to write a client that can communicate with a SockJS server without having to write the client in node.

The client can also talk directly to the server side UNIX socket (`/tmp/sockjs-unix.sock` by default) without having to use SockJS at all but provides the same interface
as the client that talks over SockJS. You can use this for unit testing, healthchecks or reduce the amount of things you have to run when developing the server.

## How To use the client

First step is to run the server as in the above description

Next run the client javascript with:

    node client.js

Which will start a sockJS client listening `/tmp/sockjs-unix-client.sock`

Then you can simply run the example client:

    go run example_client.go

If you wish to run without the SockJS layer you can run the `example_server.go` and then run `example_direct_server.go`
