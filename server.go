// +build ignore

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "rootpost"
	dbname   = "GoCompiler"
)

type Code struct {
	id     int
	date   string
	code   string
	result string
	name   string
}

//var addr = flag.String("addr", "localhost:8080", "http service address")
var addr = flag.String("addr", "localhost:8080", "http service address")

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

		fileName = "data" + r.Host

		//write to file
		f, err := os.Create(fileName + "go")
		n2, err := f.Write(message)
		f.Sync()
		f.Close()
		log.Printf("wrote %d bytes\n", n2)

		//execute file:
		execStatus, execResult = FileExecute(fileName)

		if !fileStatus {
			//send failure message to cleint
			err = c.WriteMessage(mt, "NOK")
			if err != nil {
				log.Println(" write:", err)
				break
			}
		} else {
			//handle request
			//send success message to cleint
			err = c.WriteMessage(mt, "OK")
			if err != nil {
				log.Println(" write:", err)
				break
			}
			//send program output to client
			err = c.WriteMessage(mt, execResult)
			if err != nil {
				log.Println(" write:", err)
				break
			}
		}

	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func FileExecute(fileName string) (bool, string) {
	cmd := exec.Command("cmd", "/C", "go run", fileName+".go")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, ""
	}
	log.Printf("combined out:\n%s\n", string(out))

	f, err := os.Create(fileName + ".txt")
	n2, err := f.Write(out)
	f.Sync()
	f.Close()
	log.Printf("wrote %d bytes\n", n2)

	return true, out
}

func insertToDatabase(codee string, resultt string, namee string, db *sql.DB, err error) {

	sqlInsert := `
	INSERT INTO public."Code" (date,code,result,name)
	VALUES (current_timestamp, $1 ,  $2, $3)`
	_, err = db.Exec(sqlInsert, codee, resultt, namee)
	if err != nil {
		panic(err)
	}
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	codee := "cobbbbbde"
	resultt := "result"
	namee := "name"
	insertToDatabase(codee, resultt, namee, db, err)

	//sqlSelect := `SELECT * FROM public."Code" WHERE id=14;`
	//var code Code
	//row := db.QueryRow(sqlSelect)
	//errs := row.Scan(&code.id, &code.date, &code.code,
	//	&code.result,&code.name)
	//switch errs {
	//case sql.ErrNoRows:
	//	fmt.Println("No rows were returned!")
	//	return
	//case nil:
	//	fmt.Println(code.id)
	//default:
	//	panic(errs)
	//}

	fmt.Println("Successfully connected!")
	//start-up:
	dupa()

	log.Printf("Still going")
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
