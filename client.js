var net = require('net');
var SockJS = require('sockjs-client-node');

var settings = {
    listen_socket: '/tmp/sockjs-unix-client.sock',
    backend: process.argv[2] || 'http://localhost:1334/',
    log_level: "warn"
};

var log_levels = {
    "none": 7,
    "fatal": 6,
    "error": 5,
    "warn": 4,
    "info": 3,
    "debug": 2,
    "verbose": 1
};

settings.log_level = log_levels[settings.log_level];
var log_message = function(severity, message) {
    if (log_levels[severity] >= settings.log_level) {
        arguments[0] = arguments[0].toUpperCase();
        console.log.apply(this, arguments);
    }
}

var socket_count = 0;
var sockjs_count = 0;

var socket_server = net.createServer(function(client) {
    log_message("info", "Accepted Connection");
    socket_count ++;

    var data_buffer = [];
    var sockjs = new SockJS(settings.backend);

    client.on('close', function(){
        log_message("debug", "Client closed connection");
        socket_count--;
        sockjs.close();
    });

    client.on('data', function(packet) {
        log_message("debug", "Client data: " + packet);
        data_buffer += packet.toString();
        if (data_buffer.indexOf('\n') != -1) {

            var lines = data_buffer.split('\n');
            for (var i = 0; i < lines.length - 1; i++){
                sockjs.send(lines[i]);
            }
            data_buffer = lines[lines.length - 1];
        }
    });

    sockjs.onmessage = function(e) {
        log_message("debug", "SockJS message");
        log_message("verbose", "SockJS message: " + e.data);
        client.write(e.data);
    }

    sockjs.onopen = function() {
        log_message("debug", "SockJS open");
        sockjs_count++;
        client.write(JSON.stringify({"connect": true}));
        log_message("debug", "SockJS sent handshake");
    }

    sockjs.onclose = function() {
        log_message("debug", "SockJS close");
        sockjs_count--;
        client.end()
    }

});

socket_server.listen(settings.listen_socket);

setInterval(function(){
    log_message("info", "Connections:", socket_count, sockjs_count)
}, 2000);
