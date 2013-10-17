var net = require('net');
var SockJS = require('sockjs-client-node');

var settings = {
    listen_socket: '/tmp/sockjs-unix-client.sock',
    backend: process.argv[2] || 'http://localhost:1334/',
    log_level: "debug"
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
    socket_count ++;

    var data_buffer = [];
    var sockjs = new SockJS(settings.backend);

    client.on('close', function(){
        socket_count--;
        sockjs.close();
    });

    client.on('data', function(packet) {
        reply_buffer += packet.toString();
        if (reply_buffer.indexOf('\n') != -1) {

            var lines = reply_buffer.split('\n');
            for (var i = 0; i < lines.length - 1; i++){
                sockjs.send(lines[i]);
            }
            reply_buffer = lines[lines.length - 1];
        }
    });

    sockjs.onopen = function() {
        sockjs_count++;
        sockjs.onmessage = function(e) {
            client.write(e.data);
        }

        client.write(JSON.stringify({"connect": true}));
    }

    sockjs.onclose = function() {
        sockjs_count--;
        client.end()
    }

});

socket_server.listen(settings.listen_socket);

setInterval(function(){
    log_message("info", "Connections:", socket_count, sockjs_count)
}, 2000);
