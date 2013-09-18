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

    go run example.go backend.go

Which will listen on the unix socket `/tmp/sockjs-unix.sock` 

Finally running:

    python2 -m SimpleHTTPServer

then visit `localhost:8000/test.html` and in the console run `sock.send(JSON.stringify({"body": "woo", "channel": "test"}));`

## Things to note

In my usage of SockJS I enforce that you must have a `body` and a `channel` (which is not normally something SockJS mandates)
