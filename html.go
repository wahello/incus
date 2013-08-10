package main

import "html/template"
import "net/http"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, listenAddr)
}

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<script>

function MQ(url, UID) {
    this.MAXRETRIES   = 6;
    
    this.retries      = 0;
    this.url          = url;
    this.UID          = UID;
    this.onMessageCbs = {};
    
    this.connect();
}

MQ.prototype.connect = function() {
    this.socket = new WebSocket(this.url);
    
    var self = this;
    this.socket.onopen    = function() { self.authenticate() };
    this.socket.onmessage = function(e) { self.onMessage(e) };
    this.socket.onclose   = function() { self.onClose() };
}

MQ.prototype.newMessage = function(Event, Body) {
    var obj = {
        "Event": Event,
        "Body": Body,
        "Time": Math.round(new Date().getTime() / 1000)
    };
    
    return JSON.stringify(obj);
}

MQ.prototype.authenticate = function() {
    this.retries = 0;
    var message = this.newMessage("Authenticate", {"UID": this.UID});
    
    this.socket.send(message);
    console.log("Authenticated");
}

MQ.prototype.on = function(name, func) {
    this.onMessageCbs[name] = func;
}

MQ.prototype.onMessage = function(e) {
    var msg = JSON.parse(e.data);

    if ("Event" in msg && msg.Event in this.onMessageCbs) {
        if(typeof this.onMessageCbs[msg.Event] == "function") {
            this.onMessageCbs[msg.Event].call(null, msg.Body);
	}
    }
}

MQ.prototype.onClose = function() {
    if (this.retries > this.MAXRETRIES) {
        return;
    }
    
    var self = this;
    this.retries++;
    window.setTimeout(function() {
        console.log("Connection closed, retrying");
        
        self.connect();
    }, 1000);
}

MQ.prototype.MessageUser = function(event, UID, message) { // need to send sender's UID
    var body = {"Event": event, "UID": UID, "Message": message};
    
    var msg = this.newMessage("MessageUser", body);
    return this.socket.send(msg);
}

MQ.prototype.MessageAll = function(event, message) {
    var body = {"Event": event, "Message": message};
    
    var msg = this.newMessage("MessageAll", body);
    return this.socket.send(msg);
}

MQ.prototype.setPage = function(page) {
    var body = {"Page": page}
    
    var msg = this.newMessage("SetPage", body);
    return this.socket.send(msg);
}

function init() {
    var socket = new MQ("ws://{{.}}/socket", "USER1");
    
    socket.on("TEST", function(e) {
        alert(e);
        console.log(e);
    });
	
}

window.addEventListener("load", init, false);

</script>
</head>
<body>
<input id="input" type="text">
<div id="output"></div>
</body>
</html>
`))