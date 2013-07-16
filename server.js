var net = require('net');
var http = require('http');
var sockjs = require('sockjs');

var settings = {
    sockjs_server: {sockjs_url: "http://cdn.sockjs.org/sockjs-0.3.min.js"},
    listen_port: 8080,
    backend_socket: '/tmp/sockjs-unix.sock'
};

// Base Http Server
var server = http.createServer();

// SockJS server
var sockjs_tcp = sockjs.createServer(settings.sockjs_server);

var socket_count = 0;
var backend_count = 0;

sockjs_tcp.on('connection', function(socket) {
    socket_count ++;
    var reply_buffer = [];

    var backend = net.createConnection({path: settings.backend_socket}, function(){backend_count++;})
    backend.setKeepAlive(true);
    backend.write(JSON.stringify({"id": socket.id}));

    socket.on('data', function(message) {
        try {
            var data = JSON.parse(message);

            if (data.body === undefined) {
                throw "Missing Body";
            }
            if (data.channel === undefined) {
                throw "Missing Channel";
            }
        } catch (err) {
            console.log("bad packet:", message, err);
            return;
        }
        backend.write(message);
    });

    backend.on('data', function(packet){

        reply_buffer += packet.toString();
        if (reply_buffer.indexOf('\n') != -1) {

            var lines = reply_buffer.split('\n');
            for (var i = 0; i < lines.length - 1; i++){
                socket.write(lines[i]);
            }
            reply_buffer = lines[lines.length - 1];
        }
    });

    var socket_closed = false;
    var backend_closed = false;

    // On End
    socket.on('close', function(){
        if (socket_closed) {
            return;
        }
        socket_closed = true;
        socket_count--;

        console.log("socket closed killing backend", arguments);
        backend.end();
    });

    backend.on('close', function(){
        if (backend_closed) {
            return;
        }
        backend_closed = true;
        backend_count--;

        console.log("backend closed killing socket");
        socket.end();
    });

    backend.on('error', function(){
        console.log("an error!", arguments);
        socket.end();
        backend.end();
    });

});

sockjs_tcp.installHandlers(server);
console.log(' [*] Listening on 0.0.0.0:8080' );
server.listen(settings.listen_port, '0.0.0.0');

setInterval(function(){
    console.log("Connections:", socket_count, backend_count)
}, 2000);
