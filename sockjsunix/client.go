package sockjsunix

import "io"
import "net"
import "encoding/json"

func UnixSockJSClient(path string, log Logger) (inbound chan Packet, outbound chan Packet, err error) {
    log.Debugf("SockJSUnix: Dialing")
    sock, err := net.Dial("unix", path)
    if err != nil {
        return nil, nil, err
    }
    log.Debugf("SockJSUnix: Connection established")
    dec := json.NewDecoder(sock)

    log.Debugf("SockJSUnix: waiting for handshake")
    var handshake = make(map[string]interface{})
    err = dec.Decode(&handshake)
    if err != nil {
        log.Debugf("SockJSUnix: handshake BAD")
        return nil, nil, err
    }
    log.Debugf("SockJSUnix: handshake OK")

    inbound = make(chan Packet, 5)
    outbound = make(chan Packet, 5)

    var closing = false
    go func () {
        for reply := range outbound {
            output, err := json.Marshal(reply)
            if err != nil {
                log.Errorf("SockJSUnix: Failed to json encode packet - %s when encoding %s", err, reply)
                continue
            }
            sock.Write(output)
            sock.Write([]byte("\n"))
        }
        log.Debugf("SockJSUnix: outbound close")
        closing = true
        sock.Close()
    }()

    go func () {
        for {
            var message Packet
            err := dec.Decode(&message)
            if err != nil {
                if err != io.EOF && !closing {
                    log.Errorf("SockJSUnix: Error reading inbound data: %s", err)
                }
                log.Debugf("SockJSUnix: inbound close")
                close(inbound)
                return
            }
            inbound <- message
        }
    }()

    return inbound, outbound, nil
}
