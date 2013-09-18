package main

import "github.com/vaughan0/go-logging"
import "bitbucket.org/ndavies/sockjs-unix/sockjsunix"

func echoHandler(header sockjsunix.Header, inbound chan sockjsunix.Packet, outbound chan interface{}) {
    for packet := range inbound {
        outbound <- packet
    }
}

func main() {
    log := logging.Get("example.logger")

    err := sockjsunix.UnixSockJSServer("/tmp/sockjs-unix.sock", echoHandler, log)
    if err != nil {
        panic(err)
    }
}
