package sockjsunix

import "os"
import "net"
import "os/signal"
import "encoding/json"
import "github.com/vaughan0/go-logging"

type HandlerFunc func (Header, chan Packet, chan interface{})

type Packet struct {
    Body interface{}    `json:"body"`
    Channel string      `json:"channel"`
}

type Header struct {
    Id string
}

func packetStreamer(fd net.Conn, handler HandlerFunc, log *logging.Logger) {
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
                log.Error("SockJSUnix: Failed to json encode reply - ", err, " when encoding ", reply)
                continue
            }
            fd.Write(output)
            fd.Write([]byte("\n"))
        }
    }()

    dec := json.NewDecoder(fd)
    dec.UseNumber()

    var h Header
    err := dec.Decode(&h)
    if err != nil {
        log.Error("SockJSUnix: Did not receive header properly - ", err)
        return
    }

    go func () {
        for {
            var p Packet
            err := dec.Decode(&p)
            if err != nil {
                log.Warn("SockJSUnix: Error decoding inbound json - ", err)
                close(inbound)
                return
            }
            log.Debug("SockJSUnix: received packet - ", p)
            inbound <- p
        }
    }()

    log.Debug("SockJSUnix: Handler Enter")
    handler(h, inbound, outbound)
    log.Debug("SockJSUnix: Handler Leave")
}

func UnixSockJSServer(path string, handler HandlerFunc, log *logging.Logger) error {
    l, err := net.Listen("unix", path)
    if err != nil {
        log.Fatal("Listen error: ", err)
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
                log.Error("SockJSUnix: Server accept error - ", err)
                die_ch <-err
                return
            }

            log.Info("SockJsUnix: Accepted connection")
            go packetStreamer(fd, handler, log)
        }
    }()

    select {
    case <-sig_ch:
        return nil
    case err = <-die_ch:
        return err
    }
}

