var net = require('net');
var http = require('http');
var sockjs = require('sockjs');
var carrier = require('carrier');

var settings = {
    sockjs_server: {sockjs_url: "http://cdn.sockjs.org/sockjs-0.3.min.js"},
    backend_server: {
        host: 'localhost',
        port: 9090
    },
    listen_port: 8080
};

// Base Http Server
var server = http.createServer();

// SockJS server
var sockjs_tcp = sockjs.createServer(settings.sockjs_server);

sockjs_tcp.on('connection', function(socket) {
    var connected = false;
    var buffer = [];

    var backend = net.createConnection(settings.backend_server.port, settings.backend_server.host);

    backend.setKeepAlive(true);
    backend.on('connect', function(){
        connected = true;

        for (var i = 0; i < buffer.length; i++){
            backend.write(buffer[i]);
        }
    });

    carrier.carry(backend, function(message){
        socket.write(message);
    });

    socket.on('data', function(message) {
        if (connected){
            backend.write(message);
        } else {
            buffer.push(message);
        }
    });

    // On End
    socket.on('close', function(){
        backend.end();
    });
    backend.on('close', function(){
        socket.end();
    });
    backend.on('error', function(){
        socket.end();
        backend.end();
    });

});


sockjs_tcp.installHandlers(server);
console.log(' [*] Listening on 0.0.0.0:8080' );
server.listen(settings.listen_port, '0.0.0.0');
