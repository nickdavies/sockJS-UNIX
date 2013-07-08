package main

import "fmt"
import "os"
import "net"
import "os/signal"
import "encoding/json"

type Packet struct {
    Body interface{}
    Channel string
}

func echoHandler(fd net.Conn) {
    defer fd.Close()

    dec := json.NewDecoder(fd)
    for {
        var p Packet
        err := dec.Decode(&p)
        if err != nil {
            fmt.Println("error:", err)
            return
        }
        fmt.Println("got:", p)

        output, err := json.Marshal(p)
        if err != nil {
            fmt.Println("error:", err)
            return
        }
        fd.Write(output)
        fd.Write([]byte("\n"))
        fd.Write(output)
        fd.Write([]byte("\n"))
        fd.Write(output)
        fd.Write([]byte("\n"))
    }
}

func echoServer(l net.Listener, die_ch chan bool){
    defer func(){
        die_ch <- true
    }()

    for {
        fd, err := l.Accept()
        if err != nil {
            fmt.Println("accept error", err)
            return
        }

        fmt.Println("accepted connection")
        go echoHandler(fd)
    }
}

func main() {
    l, err := net.Listen("unix", "/tmp/sockjs-unix.sock")
    if err != nil {
        fmt.Println("listen error", err)
        return
    }
    defer l.Close()

    var sig_ch = make(chan os.Signal, 1)
    var die_ch = make(chan bool, 1)

    signal.Notify(sig_ch, os.Interrupt)

    go echoServer(l, die_ch)

    select {
    case <-sig_ch:
        return
    case <-die_ch:
        return
    }
}
