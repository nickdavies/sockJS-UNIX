package main

import "fmt"
import "github.com/vaughan0/go-logging"
import "./sockjsunix"

func main() {
    logging.DefaultSetup()
    log := logging.Get("example.logger")
    log.Threshold = logging.Debug

    log.Debug("Starting example!")
    inbound, outbound, err := sockjsunix.UnixSockJSClient("/tmp/sockjs-unix-client.sock", log)
    if err != nil {
        panic(err)
    }

    outbound <- sockjsunix.Packet{"hello!", "meh"}
    fmt.Println(<-inbound)

    outbound <- sockjsunix.Packet{"hello 1!", "meh"}
    outbound <- sockjsunix.Packet{"hello 2!", "meh"}

    fmt.Println(<-inbound)
    fmt.Println(<-inbound)

    close(outbound)
}

