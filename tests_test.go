package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestDiffShouldReturnFullSimilarityForSameCodes(t *testing.T) {
	//Code initialization
	code := "package main import " + "fmt" + "func main() { fmt.Println(" + "hello world" + ") }"
	// Diff method execution
	_, percentage := Diff(code, code)
	// Is percentage isn't 100
	if percentage != 100 {
		// If test failed
		t.Errorf("Percentage was incorrect, got: %d, want: %d.", percentage, 100)
	}
}

func TestDiffShouldReturnPercentageSimilarityForDifferentCodes(t *testing.T) {
	//Codes initialization
	firstCode := "package main \nimport " + "fmt" + "\nfunc main() {\n fmt.Println(" + "hello" + ") \n fmt.Println(" + "hello1" + ")}"
	secondCode := "package main \nimport " + "fmt" + "\nfunc main() {\n fmt.Println(" + "hello" + ")\n + for j := 7; j <= 9; j++ {fmt.Println(j)\n} }"
	// Diff method execution
	_, percentage := Diff(firstCode, secondCode)
	// Is percentage isn't 44
	if percentage != 44 {
		// If test failed
		t.Errorf("Percentage was incorrect, got: %d, want: %d.", percentage, 44)
	}
}

func TestDiffShouldReturnEmptyRaportIfTheSameCodes(t *testing.T) {
	//Code initialization
	code := "package main import " + "fmt" + "func main() { fmt.Println(" + "hello world" + ") }"
	// Diff method execution
	total, _ := Diff(code, code)
	// If raport isn't empty
	if total != "" {
		// If test failed
		t.Errorf("Result was incorrect, got: %s, want: %s.", total, "")
	}

}

func TestCompareDataShouldCompareRaportForAllFiles(t *testing.T) {
	// Array of code table objects
	var codes []Code
	// Create object
	codeToArray := Code{name: "name",
		date:   "time",
		result: "result",
		code:   "package main\nimport " + "fmt" + "\nfunc main(){\nfmt.Println(" + "hello" + ")\nfmt.Println(" + "hello1" + ")}"}

	// Create mock object from Client
	codeFromClient := Code{name: "name2",
		date:   "time2",
		result: "result2",
		code:   "package main\nimport " + "fmt" + "\nfunc main(){\nfmt.Println(" + "hello" + ")\n+ for j := 7; j <= 9; j++ {fmt.Println(j)\n} }"}

	// Add object to ArrayCodes
	codes = append(codes, codeToArray)
	// CompareData method execution
	result := compareData(codes, codeFromClient)

	// reading expected result form file
	expectedResult, err := ioutil.ReadFile("result.txt")

	// If raports are different
	if result != string(expectedResult) {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, expectedResult)
		fmt.Print(err)
	}

}

func TestFileExecuteShouldReturnTrueIfExecute(t *testing.T) {

	// define file name
	fileName := "testingData/testFile"
	// FileExecute method execution
	isDone, _ := FileExecute(fileName)
	// If file isn't execute
	if isDone != true {
		t.Errorf("Result was incorrect, got: %s, want: %s", (strconv.FormatBool(isDone)), (strconv.FormatBool(true)))
	}

}

func TestFileExecuteShouldReturnResultOfExecution(t *testing.T) {
	// define file name
	fileName := "testingData/testFile"
	// define expected result
	expectedResult := "Hello World"
	// FileExecute method execution
	_, result := FileExecute(fileName)
	// If reaults are different
	if result != expectedResult {
		t.Errorf("Result was incorrect, got: %s, want: %s", result, expectedResult)
	}

}
