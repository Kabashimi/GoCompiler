package main

import (
	"fmt"
	"testing"
)

func TestCompare(t *testing.T) {

	firstCode := "package main import " + "fmt" + "func main() { fmt.Println(" + "hello world" + ") }"
	total, percentage := Diff(firstCode, firstCode)
	if percentage != 100 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", percentage, 100)
		fmt.Print(total)
	}

}

//func TestCompare2(t *testing.T) {
//
//
//	firstCode := "package main import " +"fmt" +"func main() { fmt.Println("+"hello world"+") }"
//	total, percentage:= Diff(firstCode, firstCode)
//	if total != "" {
//		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, "kk")
//		fmt.Print(percentage)
//	}
//}
