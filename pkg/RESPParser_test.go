package pkg

import (
	"reflect"
	"testing"
)

func TestLoadFromFile(t *testing.T) {

	testcases := []struct {
		testcaseName string
		filePath     string
		err          error
	}{
		{
			testcaseName: "Normal case: ini file is present",
			filePath:     "testdata/resp.txt",
			err:          nil,
		},
		{
			testcaseName: "corner case: empty file",
			filePath:     "testdata/empty.txt",
			err:          ErrEmptyFile,
		},
		{
			testcaseName: "corner case: file not found",
			filePath:     "testdata/filex.txt",
			err:          ErrFileNotExist,
		},
	}

	for _, testcase := range testcases {

		t.Run(testcase.testcaseName, func(t *testing.T) {
			p := NewParser()
			err := p.LoadFromFile(testcase.filePath)
			if err != testcase.err {
				t.Errorf("expected error: %v, got: %v", testcase.err, err)
			}
		})
	}
}

func TestParse(t *testing.T) {

	const filePath = "testdata/resp.txt"

	testcases := []struct {
		testcaseName string
		fileLines    []string
		err          error
		expected     []RespType
	}{
		{
			testcaseName: "Simple String",
			fileLines:    []string{"+OK"},
			err:          nil,
			expected:     []RespType{SimpleString{Value: "OK"}},
		},
		{
			testcaseName: "Simple Error",
			fileLines:    []string{"-ERR unknown command"},
			err:          nil,
			expected:     []RespType{SimpleError{Value: "ERR unknown command"}},
		},
		{
			testcaseName: "Integer",
			fileLines:    []string{":12345"},
			err:          nil,
			expected:     []RespType{Integer{Value: 12345}},
		},
		{
			testcaseName: "Unknown Type",
			fileLines:    []string{"?Unknown"},
			err:          ErrUnkType,
			expected:     make([]RespType, 0),
		},
	}

	for _, testcase := range testcases {

		t.Run(testcase.testcaseName, func(t *testing.T) {
			p := NewParser()
			err := p.parse(testcase.fileLines)
			if err != testcase.err {
				t.Errorf("expected error: %v, got: %v", testcase.err, err)
			}

			if !reflect.DeepEqual(p.RespList, testcase.expected) {
				t.Errorf("expected %v, got %v", testcase.expected, p.RespList)
			}

		})
	}
}
