package test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onlyafly/oakblue/internal/analyzer"
	"github.com/onlyafly/oakblue/internal/emitter"
	"github.com/onlyafly/oakblue/internal/parser"
	"github.com/onlyafly/oakblue/internal/syntax"
	"github.com/onlyafly/oakblue/internal/util"
	"github.com/onlyafly/oakblue/internal/vm"
	"github.com/stretchr/testify/assert"
)

const (
	assemblerSuiteTestDataDir = "test/testdata_assembler"
	vmSuiteTestDataDir        = "test/testdata_vm"
	fileExtPattern            = "*.asm"
	objFileExtension          = ".obj"
	errFileExtension          = ".err"
	regFileExtension          = ".reg"
)

func TestMain(m *testing.M) {
	// We set the base directory so that the test cases can use paths that make
	// sense. If this is not set, the current working directory while the tests
	// run will be "test"
	err := os.Chdir("..")
	if err != nil {
		panic("Error changing directory to the base directory: " + err.Error())
	}

	// NOTE: If this is not here, the tests will not run and they will not exit
	os.Exit(m.Run())
}

// TestAssemblerSuite runs the entire language test suite
func TestAssemblerSuite(t *testing.T) {

	err := filepath.Walk(assemblerSuiteTestDataDir, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil // Can't visit this node, but continue walking elsewhere
		}
		if fi.IsDir() {
			return nil // Not a file, ignore.
		}

		name := fi.Name()
		matched, err := filepath.Match(fileExtPattern, name)
		if err != nil {
			return err // malformed pattern
		}

		if matched {
			testAssemblingFile(fp, t)
		}

		return nil
	})

	if err != nil {
		t.Errorf("Error walking test suite directory <" + assemblerSuiteTestDataDir + ">: " + err.Error())
	}
}

func TestExecutingSuite(t *testing.T) {

	err := filepath.Walk(vmSuiteTestDataDir, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil // Can't visit this node, but continue walking elsewhere
		}
		if fi.IsDir() {
			return nil // Not a file, ignore.
		}

		name := fi.Name()
		matched, err := filepath.Match(fileExtPattern, name)
		if err != nil {
			return err // malformed pattern
		}

		if matched {
			testExecutingFile(fp, t)
		}

		return nil
	})

	if err != nil {
		t.Errorf("Error walking test suite directory <" + vmSuiteTestDataDir + ">: " + err.Error())
	}
}

func testAssemblingFile(sourceFilePath string, t *testing.T) {
	sourceDirPart, sourceFileNamePart := filepath.Split(sourceFilePath)
	parts := strings.Split(sourceFileNamePart, ".")
	testName := parts[0]

	input, errIn := util.ReadTextFile(sourceFilePath)
	if errIn != nil {
		t.Errorf("Error reading file <" + sourceFilePath + ">: " + errIn.Error())
		return
	}

	errorList := syntax.NewErrorList("Syntax")
	listing, _ := parser.Parse(input, sourceFilePath, errorList) // the error return is ignored because it will be combined with the analyzer's errors
	program, err := analyzer.Analyze(listing, errorList)

	if err != nil {
		outputFilePath := sourceDirPart + testName + errFileExtension
		expectedRaw, errOut := util.ReadTextFile(outputFilePath)
		if errOut != nil {
			expectedRaw = "SUITE_TEST FOUND NO .ERR FILE AT <" + outputFilePath + ">"
		}

		// Remove any carriage return line endings from .out file
		expectedWithUntrimmed := strings.Replace(expectedRaw, "\r", "", -1)
		expected := strings.TrimSpace(expectedWithUntrimmed)

		verify(t, sourceFilePath, input, expected, err.Error())
	}

	actual, emitError := emitter.Emit(program, syntax.NewErrorList("Emit"))
	if emitError != nil {
		return
	}

	outputFilePath := sourceDirPart + testName + objFileExtension

	expected, errOut := util.ReadBinaryFile(outputFilePath)
	if errOut != nil {
		t.Errorf("Error reading file <" + outputFilePath + ">: " + errOut.Error())
		return
	}

	verifyBinary(t, sourceFilePath, input, expected, actual)
}

func testExecutingFile(sourceFilePath string, t *testing.T) {
	sourceDirPart, sourceFileNamePart := filepath.Split(sourceFilePath)
	parts := strings.Split(sourceFileNamePart, ".")
	testName := parts[0]

	input, errIn := util.ReadTextFile(sourceFilePath)
	if errIn != nil {
		t.Errorf("Error reading file <" + sourceFilePath + ">: " + errIn.Error())
		return
	}

	errorList := syntax.NewErrorList("Syntax")
	listing, _ := parser.Parse(input, sourceFilePath, errorList) // the error return is ignored because it will be combined with the analyzer's errors
	program, err := analyzer.Analyze(listing, errorList)

	if err != nil {
		outputFilePath := sourceDirPart + testName + errFileExtension
		expectedRaw, errOut := util.ReadTextFile(outputFilePath)
		if errOut != nil {
			expectedRaw = "SUITE_TEST FOUND NO .ERR FILE AT <" + outputFilePath + ">"
		}

		// Remove any carriage return line endings from .out file
		expectedWithUntrimmed := strings.Replace(expectedRaw, "\r", "", -1)
		expected := strings.TrimSpace(expectedWithUntrimmed)

		verify(t, sourceFilePath, input, expected, err.Error())
		return
	}

	machineCode, emitError := emitter.Emit(program, syntax.NewErrorList("Emit"))
	if !assert.NoError(t, emitError) {
		return
	}

	m := vm.NewMachine()
	m.LoadMemory(machineCode, 0x3000)
	executeError := m.Execute()
	if executeError != nil {
		t.Errorf("Error during execution of test <%s>: %s", sourceFilePath, executeError.Error())
		return
	}

	registerDump := m.RegisterDump()

	regFilePath := sourceDirPart + testName + regFileExtension
	expectedRegisterDump, errOut := util.ReadTextFile(regFilePath)
	if errOut != nil {
		t.Errorf("Error reading file <" + regFilePath + ">: " + errOut.Error())
		return
	}
	verify(t, sourceFilePath, input, expectedRegisterDump, registerDump)

	/*
		outPath := sourceDirPart + testName + ".out"
		expected, errOut := util.ReadTextFile(outPath)
		if errOut != nil {
			t.Errorf("Error reading file <" + outPath + ">: " + errOut.Error())
			return
		}
		verify(t, sourceFilePath, input, expected, actual)
	*/
}

func verify(t *testing.T, testCaseName, input, expected, actual string) {
	if expected != actual {
		t.Errorf(
			"\n===== TEST SUITE CASE FAILED: %s\n"+
				"===== INPUT\n%v\n"+
				"===== EXPECTED\n%v\n"+
				"===== ACTUAL\n%v\n"+
				"===== END\n",
			testCaseName,
			strings.TrimSpace(input),
			expected,
			actual)
	}
}

func verifyBinary(t *testing.T, testCaseName, input string, expected, actual []byte) {
	result := bytes.Compare(expected, actual)

	if result == 0 {
		// All good
	} else {
		t.Errorf(
			"\n===== TEST SUITE CASE FAILED: %s\n"+
				"===== INPUT\n%v\n"+
				"===== EXPECTED\n%v\n"+
				"===== ACTUAL\n%v\n"+
				"===== END\n",
			testCaseName,
			input,
			expected,
			actual)
	}
}
