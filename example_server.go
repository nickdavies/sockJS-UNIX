package main

import "github.com/vaughan0/go-logging"
import "./sockjsunix"

func echoHandler(header sockjsunix.Header, inbound chan sockjsunix.Packet, outbound chan interface{}) {
    for packet := range inbound {
        outbound <- packet
    }
}

func main() {
    logging.DefaultSetup()
    log := logging.Get("example.logger")
    log.Threshold = logging.Debug

    err := sockjsunix.UnixSockJSServer("/tmp/sockjs-unix.sock", echoHandler, log)
    if err != nil {
        panic(err)
    }
}
