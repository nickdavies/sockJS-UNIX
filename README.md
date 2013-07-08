# SockJS-UNIX

SockJS-Unix terminates a sockJS connection and proxies all traffic to and from a UNIX socket.

It is written in node.js and depends only on the sockJS library itself.

## Why would you want to do that?

The reason why I wrote this was because I wanted to use SockJS with a backend server that did not have a SockJS library (for me it was Golang).
This allows any language that can talk via a unix socket (everything) to use SockJS.

