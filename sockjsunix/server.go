package sockjsunix

import "io"
import "os"
import "net"
import "os/signal"
import "encoding/json"

func packetStreamer(fd net.Conn, handler HandlerFunc, log Logger) {
    defer func(){
        if recovery := recover(); recovery != nil {
            log.Errorf("Request Paniced: %v", recovery)
        }
    }()
    // Without the buffer its possible that the call inbound <- p
    // would block causing the decoder to never register the error
    // and the go routine would never end
    var inbound = make(chan Packet, 100)
    var outbound = make(chan interface{}, 5)

    defer fd.Close()
    defer close(outbound)

    go func(){
        for reply := range outbound {
            output, err := json.Marshal(reply)
            if err != nil {
                log.Errorf("SockJSUnix: Failed to json encode reply - %s when encoding %s", err, reply)
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
        log.Errorf("SockJSUnix: Did not receive header properly - %s", err)
        return
    }

    go func () {
        defer close(inbound)
        for {
            var p Packet
            err := dec.Decode(&p)
            if err == io.EOF {
                return
            } else if err != nil {
                if err.Error() != "read unix @: use of closed network connection" {
                    log.Warnf("SockJSUnix: Error decoding inbound json - %s", err)
                }
                return
            }
            log.Debugf("SockJSUnix: received packet - %s", p)
            inbound <- p
        }
    }()

    log.Debugf("SockJSUnix: Handler Enter")
    handler(h, inbound, outbound)
    log.Debugf("SockJSUnix: Handler Leave")
}

func UnixSockJSServer(path string, handler HandlerFunc, log Logger) error {
    l, err := net.Listen("unix", path)
    if err != nil {
        log.Fatalf("Listen error: %s", err)
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
                log.Errorf("SockJSUnix: Server accept error - %s", err)
                die_ch <-err
                return
            }

            log.Infof("SockJsUnix: Accepted connection")
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

