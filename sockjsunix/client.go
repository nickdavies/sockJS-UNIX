package sockjsunix

import "io"
import "net"
import "encoding/json"

import "crypto/rand"
import "encoding/base64"

func getId() (string) {
    b := make([]byte, 30)
    rand.Read(b)
    en := base64.StdEncoding
    d := make([]byte, en.EncodedLen(len(b)))
    en.Encode(d, b)

    return string(d)
}

func UnixSockJSClient(path string, direct bool, log Logger) (inbound chan Packet, outbound chan Packet, err error) {
    log.Debugf("SockJSUnix: Dialing")
    sock, err := net.Dial("unix", path)
    if err != nil {
        return nil, nil, err
    }
    log.Debugf("SockJSUnix: Connection established")
    dec := json.NewDecoder(sock)
    dec.UseNumber()

    if direct {
        // Handshake is require to make sure that no
        // packets are lost before sockJS gets the 
        // onopen event
        log.Debugf("SockJSUnix: waiting for handshake")
        var handshake = make(map[string]interface{})
        err = dec.Decode(&handshake)
        if err != nil {
            log.Debugf("SockJSUnix: handshake BAD")
            return nil, nil, err
        }
        log.Debugf("SockJSUnix: handshake OK")
    } else {
        log.Debugf("SockJSUnix: sending header")

        enc := json.NewEncoder(sock)
        err := enc.Encode(Header{getId()})
        if err != nil {
            log.Errorf("SockJSUnix: failed to send header: %s", err)
            return nil, nil, err
        }
    }

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
            log.Debugf("SockJSUnix: writing - %s", output)
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
            log.Debugf("SockJSUnix: received - %s", message)
            inbound <- message
        }
    }()

    return inbound, outbound, nil
}
