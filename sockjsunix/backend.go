package sockjsunix

import "fmt"
import "os"
import "net"
import "os/signal"
import "encoding/json"

type HandlerFunc func (chan Packet, chan interface{})

type Packet struct {
    Body interface{}    `json:"body"`
    Channel string      `json:"channel"`
}

func packetStreamer(fd net.Conn, handler HandlerFunc) {
    // Without the buffer its possible that the call inbound <- p
    // would block causing the decoder to never register the error
    // and the go routine would never end
    var inbound = make(chan Packet, 100)
    var outbound = make(chan interface{})

    defer fd.Close()
    defer close(outbound)

    go func(){
        for reply := range outbound {
            output, err := json.Marshal(reply)
            if err != nil {
                fmt.Println("error:", err)
            }
            fd.Write(output)
            fd.Write([]byte("\n"))
        }
    }()

    go func () {
        dec := json.NewDecoder(fd)
        for {
            var p Packet
            err := dec.Decode(&p)
            if err != nil {
                fmt.Println("lib error:", err)
                close(inbound)
                return
            }
            fmt.Println("got:", p)
            inbound <- p
        }
    }()

    fmt.Println("Handler: Enter")
    handler(inbound, outbound)
    fmt.Println("Handler: Leave")
}

func UnixSockJSServer(path string, handler HandlerFunc) error {
    l, err := net.Listen("unix", path)
    if err != nil {
        fmt.Println("listen error", err)
        return err
    }
    defer l.Close()

    var sig_ch = make(chan os.Signal, 1)
    var die_ch = make(chan error, 1)

    signal.Notify(sig_ch, os.Interrupt)

    go func (){
        for {
            fd, err := l.Accept()
            if err != nil {
                fmt.Println("accept error", err)
                die_ch <-err
                return
            }

            fmt.Println("accepted connection")
            go packetStreamer(fd, handler)
        }
    }()

    select {
    case <-sig_ch:
        return nil
    case err = <-die_ch:
        return err
    }
}

