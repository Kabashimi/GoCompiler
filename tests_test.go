package main

import (
	"testing"
)

func TestCompareSameCodes(t *testing.T) {

	code := "package main import " + "fmt" + "func main() { fmt.Println(" + "hello world" + ") }"
	_, percentage := Diff(code, code)

	if percentage != 100 {
		t.Errorf("Percentage was incorrect, got: %d, want: %d.", percentage, 100)
	}

}

func TestCompareDiffCodes(t *testing.T) {
	firstCode := "package main \nimport " + "fmt" + "\nfunc main() {\n fmt.Println(" + "hello" + ") \n fmt.Println(" + "hello1" + ")}"
	secondCode := "package main \nimport " + "fmt" + "\nfunc main() {\n fmt.Println(" + "hello" + ")\n + for j := 7; j <= 9; j++ {fmt.Println(j)\n} }"
	_, percentage := Diff(firstCode, secondCode)
	if percentage != 44 {
		t.Errorf("Percentage was incorrect, got: %d, want: %d.", percentage, 44)
	}

}

func TestCompareByRaportSameCodes(t *testing.T) {

	code := "package main import " + "fmt" + "func main() { fmt.Println(" + "hello world" + ") }"
	total, _ := Diff(code, code)
	if total != "" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", total, "")
	}

}
