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
	"strconv"
	"time"

	"bytes"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"strings"
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
var addr = flag.String("addr", "192.168.0.17:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	fmt.Println("REQUEST OCCURED: " + c.RemoteAddr().String())

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		//currentTime := time.Now()
		//time := currentTime.String()
		fileName := "data" + c.RemoteAddr().String()

		//write to file
		f, err := os.Create(fileName + ".go")
		n2, err := f.Write(message)
		f.Sync()
		f.Close()
		log.Printf("wrote %d bytes\n", n2)

		//execute file:
		execStatus, execResult := FileExecute(fileName)

		if !execStatus {
			//send failure message to cleint
			err = c.WriteMessage(mt, []byte("NOK"))
			if err != nil {
				log.Println(" write:", err)
				break
			}
		} else {
			//handle request
			//send success message to cleint
			err = c.WriteMessage(mt, []byte("OK"))
			if err != nil {
				log.Println(" write:", err)
				break
			}
			//send program output to client
			err = c.WriteMessage(mt, []byte(execResult))
			if err != nil {
				log.Println(" write:", err)
				break
			}
			raport := insertToDatabase(string(message), string(execResult), string(fileName))
			//fmt.Println("Raport:\n"+ raport)
			err = c.WriteMessage(mt, []byte(raport))
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

	return true, string(out)
}

func insertToDatabase(codee string, resultt string, namee string) string {

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

	currentTime := time.Now()
	time := currentTime.String()

	fmt.Println("Successfully connected!")

	sqlInsert := `
	INSERT INTO public."Code" (date,code,result,name)
	VALUES ($1, $2 , $3, $4)`
	_, err = db.Exec(sqlInsert, time, codee, resultt, namee)
	if err != nil {
		panic(err)
	}
	codeObject := Code{name: namee,
		date:   time,
		result: resultt,
		code:   codee}

	return selectFromDatabase(db, err, codeObject)

}

func selectFromDatabase(db *sql.DB, err error, codeObject Code) string {
	var codes []Code
	sqlSelect := `SELECT * FROM public."Code"`
	var code Code

	rows, row := db.Query(sqlSelect)
	for rows.Next() {
		errs := rows.Scan(&code.id, &code.date, &code.code,
			&code.result, &code.name)
		if err != nil {
			panic(err)
			fmt.Println(row)
			fmt.Println(errs)

		}
		codes = append(codes, Code{id: code.id,
			date:   code.date,
			code:   code.code,
			result: code.result,
			name:   code.name})
	}

	return compareData(codes, codeObject)
}

func compareData(codes []Code, codeObject Code) string {

	raport := " "
	raports := " "
	percentage := 0
	for index, element := range codes {
		if element.date != codeObject.date {
			raport, percentage = Diff(string(element.code), string(codeObject.code))
			raports += "\n Procentowa zgodnosc aktualnego pliku z plikiem z : " + element.date + " wynosi: " + strconv.Itoa(percentage) + "% \n" + "Roznice plikow: \n" + raport
		}
		fmt.Printf(string(index))
	}

	return raports
}

func Diff(A, B string) (string, int) {

	aLines := strings.Split(A, "\n")
	bLines := strings.Split(B, "\n")

	chunks := DiffChunks(aLines, bLines)
	equalSum := 1
	addedSum := 0
	deletedSum := 0

	buf := new(bytes.Buffer)
	for _, c := range chunks {
		for _, line := range c.Added {
			addedSum += 1
			fmt.Fprintf(buf, "+%s\n", line)
		}
		for _, line := range c.Deleted {
			deletedSum += 1
			fmt.Fprintf(buf, "-%s\n", line)
		}
		for _, line := range c.Equal {
			equalSum += 1
			fmt.Fprintf(buf, " %s\n", line)
		}
	}

	equalPercent := (equalSum * 100) / (equalSum + addedSum + deletedSum)
	return strings.TrimRight(buf.String(), "\n"), equalPercent
}

func DiffChunks(a, b []string) []Chunk {
	// algorithm: http://www.xmailserver.org/diff2.pdf

	// We'll need these quantities a lot.
	alen, blen := len(a), len(b) // M, N

	// At most, it will require len(a) deletions and len(b) additions
	// to transform a into b.
	maxPath := alen + blen // MAX
	if maxPath == 0 {
		// degenerate case: two empty lists are the same
		return nil
	}

	// Store the endpoint of the path for diagonals.
	// We store only the a index, because the b index on any diagonal
	// (which we know during the loop below) is aidx-diag.
	// endpoint[maxPath] represents the 0 diagonal.
	//
	// Stated differently:
	// endpoint[d] contains the aidx of a furthest reaching path in diagonal d
	endpoint := make([]int, 2*maxPath+1) // V

	saved := make([][]int, 0, 8) // Vs
	save := func() {
		dup := make([]int, len(endpoint))
		copy(dup, endpoint)
		saved = append(saved, dup)
	}

	var editDistance int // D
dLoop:
	for editDistance = 0; editDistance <= maxPath; editDistance++ {
		// The 0 diag(onal) represents equality of a and b.  Each diagonal to
		// the left is numbered one lower, to the right is one higher, from
		// -alen to +blen.  Negative diagonals favor differences from a,
		// positive diagonals favor differences from b.  The edit distance to a
		// diagonal d cannot be shorter than d itself.
		//
		// The iterations of this loop cover either odds or evens, but not both,
		// If odd indices are inputs, even indices are outputs and vice versa.
		for diag := -editDistance; diag <= editDistance; diag += 2 { // k
			var aidx int // x
			switch {
			case diag == -editDistance:
				// This is a new diagonal; copy from previous iter
				aidx = endpoint[maxPath-editDistance+1] + 0
			case diag == editDistance:
				// This is a new diagonal; copy from previous iter
				aidx = endpoint[maxPath+editDistance-1] + 1
			case endpoint[maxPath+diag+1] > endpoint[maxPath+diag-1]:
				// diagonal d+1 was farther along, so use that
				aidx = endpoint[maxPath+diag+1] + 0
			default:
				// diagonal d-1 was farther (or the same), so use that
				aidx = endpoint[maxPath+diag-1] + 1
			}
			// On diagonal d, we can compute bidx from aidx.
			bidx := aidx - diag // y
			// See how far we can go on this diagonal before we find a difference.
			for aidx < alen && bidx < blen && a[aidx] == b[bidx] {
				aidx++
				bidx++
			}
			// Store the end of the current edit chain.
			endpoint[maxPath+diag] = aidx
			// If we've found the end of both inputs, we're done!
			if aidx >= alen && bidx >= blen {
				save() // save the final path
				break dLoop
			}
		}
		save() // save the current path
	}
	if editDistance == 0 {
		return nil
	}
	chunks := make([]Chunk, editDistance+1)

	x, y := alen, blen
	for d := editDistance; d > 0; d-- {
		endpoint := saved[d]
		diag := x - y
		insert := diag == -d || (diag != d && endpoint[maxPath+diag-1] < endpoint[maxPath+diag+1])

		x1 := endpoint[maxPath+diag]
		var x0, xM, kk int
		if insert {
			kk = diag + 1
			x0 = endpoint[maxPath+kk]
			xM = x0
		} else {
			kk = diag - 1
			x0 = endpoint[maxPath+kk]
			xM = x0 + 1
		}
		y0 := x0 - kk

		var c Chunk
		if insert {
			c.Added = b[y0:][:1]
		} else {
			c.Deleted = a[x0:][:1]
		}
		if xM < x1 {
			c.Equal = a[xM:][:x1-xM]
		}

		x, y = x0, y0
		chunks[d] = c
	}
	if x > 0 {
		chunks[0].Equal = a[:x]
	}
	if chunks[0].empty() {
		chunks = chunks[1:]
	}
	if len(chunks) == 0 {
		return nil
	}
	return chunks
}

type Chunk struct {
	Added   []string
	Deleted []string
	Equal   []string
}

func (c *Chunk) empty() bool {
	return len(c.Added) == 0 && len(c.Deleted) == 0 && len(c.Equal) == 0
}

func main() {
	//start-up:

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
