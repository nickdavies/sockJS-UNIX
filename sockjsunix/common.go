package sockjsunix

type Logger interface {
    Debugf(format string, args ...interface{})
    Noticef(format string, args ...interface{})
    Infof(format string, args ...interface{})
    Warnf(format string, args ...interface{})
    Errorf(format string, args ...interface{})
    Fatalf(format string, args ...interface{})
}

type HandlerFunc func (Header, chan Packet, chan interface{})

type Packet struct {
    Body interface{}    `json:"body"`
    Channel string      `json:"channel"`
}

type Header struct {
    Id string
}

