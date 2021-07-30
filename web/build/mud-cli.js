// TODO: make these URLs configurable when served
var mudURL = "http://localhost:8080/mud"
var heartbeatURL = "http://localhost:8080/heartbeat"
var shouldReload = true;

default_mudfiles = {};

var es = new EventSource(heartbeatURL);
es.addEventListener("message", function(e){
    console.log(e.data);
});

fetch(mudURL).then(function(response) {
    return response.json();
}).then(function(data) {
    console.log(data);
    network.ready_to_draw = false;
    network.add_mudfile(data);
    network.create_network();
    var interval = setInterval(function () {
        if (network.ready_to_draw == false) {
            return;
        }
        clearInterval(interval);
        network_data = network.get_nodes_links_json();
        mud_drawer(network_data);
    }, 100);
}).catch(function() {
    console.log("error when retrieving the MUD");
});