package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestCompareSameCodes(t *testing.T) {

	code := "package main import " + "fmt" + "func main() { fmt.Println(" + "hello world" + ") }"
	_, percentage := Diff(code, code)

	if percentage != 100 {
		t.Errorf("Percentage was incorrect, got: %d, want: %d.", percentage, 100)
	}

}

func TestDiffDiffCodes(t *testing.T) {
	firstCode := "package main \nimport " + "fmt" + "\nfunc main() {\n fmt.Println(" + "hello" + ") \n fmt.Println(" + "hello1" + ")}"
	secondCode := "package main \nimport " + "fmt" + "\nfunc main() {\n fmt.Println(" + "hello" + ")\n + for j := 7; j <= 9; j++ {fmt.Println(j)\n} }"
	_, percentage := Diff(firstCode, secondCode)
	if percentage != 44 {
		t.Errorf("Percentage was incorrect, got: %d, want: %d.", percentage, 44)
	}

}

func TestDiffByRaportSameCodes(t *testing.T) {

	code := "package main import " + "fmt" + "func main() { fmt.Println(" + "hello world" + ") }"
	total, _ := Diff(code, code)
	if total != "" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", total, "")
	}

}

func TestCompare(t *testing.T) {
	var codes []Code

	firstCodeObject := Code{name: "name",
		date:   "time",
		result: "result",
		code:   "package main\nimport " + "fmt" + "\nfunc main(){\nfmt.Println(" + "hello" + ")\nfmt.Println(" + "hello1" + ")}"}

	codeObject := Code{name: "name2",
		date:   "time2",
		result: "result2",
		code:   "package main\nimport " + "fmt" + "\nfunc main(){\nfmt.Println(" + "hello" + ")\n+ for j := 7; j <= 9; j++ {fmt.Println(j)\n} }"}

	codes = append(codes, firstCodeObject)
	result := compareData(codes, codeObject)

	expectedResult, err := ioutil.ReadFile("result.txt")

	if result != string(expectedResult) {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, expectedResult)
		fmt.Print(err)
	}

}

func TestFileIfExecuted(t *testing.T) {

	fileName := "testingData/testFile"
	isDone, _ := FileExecute(fileName)
	if isDone != true {
		t.Errorf("Result was incorrect, got: %s, want: %s", (strconv.FormatBool(isDone)), (strconv.FormatBool(true)))
	}

}

func TestFileReturnResultExecute(t *testing.T) {

	fileName := "testingData/testFile"
	expectedResult := "Hello World"
	_, result := FileExecute(fileName)
	if result != expectedResult {
		t.Errorf("Result was incorrect, got: %s, want: %s", result, expectedResult)
	}

}
