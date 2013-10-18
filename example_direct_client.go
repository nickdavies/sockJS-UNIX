package main

import "fmt"
import "time"
import "github.com/vaughan0/go-logging"
import "./sockjsunix"

func main() {
    logging.DefaultSetup()
    log := logging.Get("example.logger")
    log.Threshold = logging.Debug

    log.Debug("Starting example!")
    inbound, outbound, err := sockjsunix.UnixSockJSClient("/tmp/sockjs-unix.sock", false, log)
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

    // this is here to give the client.go time to shutdown
    // only needs to be here to show that it does close cleanly
    <-time.After(100 * time.Millisecond)
}

