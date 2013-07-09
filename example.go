package main

func echoHandler(inbound chan Packet, outbound chan interface{}) {
    for packet := range inbound {
        outbound <- packet
    }
}

func main() {
    err := UnixSockJSServer("/tmp/sockjs-unix.sock", echoHandler)
    if err != nil {
        panic(err)
    }
}
