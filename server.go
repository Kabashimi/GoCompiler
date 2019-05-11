// +build ignore

package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
)

//var addr = flag.String("addr", "localhost:8080", "http service address")
var addr = flag.String("addr", "192.168.0.150:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		//write to file
		f, err := os.Create("data.go")
		n2, err := f.Write(message)
		f.Sync()
		f.Close()
		log.Printf("wrote %d bytes\n", n2)

		//execute commands:
		//cmd := exec.Command("dir")
		//
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		//
		//cmd_err := cmd.Run()
		//if cmd_err != nil {
		//	log.Fatalf("cmd.Run() failed with %s\n", cmd_err)
		//}

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println(" write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	//start-up:
	cmd := exec.Command("cmd", "/C", "dir", "")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd_err := cmd.Run()
	if cmd_err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", cmd_err)
	}

	//main code
	flag.Parse()
	log.SetFlags(0)
	log.Printf("Server Listening")
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))

}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))

// go run server.go
