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
		fileString   string
		err          error
		expected     []RespType
	}{
		{
			testcaseName: "Simple String",
			fileString:   "+OK",
			err:          nil,
			expected:     []RespType{SimpleString{Value: "OK"}},
		},
		{
			testcaseName: "Simple Error",
			fileString:   "-ERR unknown command",
			err:          nil,
			expected:     []RespType{SimpleError{Value: "ERR unknown command"}},
		},
		{
			testcaseName: "Integer",
			fileString:   ":12345",
			err:          nil,
			expected:     []RespType{Integer{Value: 12345}},
		},
		{
			testcaseName: "Bulk String",
			fileString:   "$5\r\nhello\r\n",
			err:          nil,
			expected:     []RespType{BulkStrings{Value: "hello", Length: 5}},
		},
		{
			testcaseName: "Empty Bulk String",
			fileString:   "$0\r\n\r\n",
			err:          nil,
			expected:     []RespType{BulkStrings{Value: "", Length: 0}},
		},
		{
			testcaseName: "Null Bulk String",
			fileString:   "$-1\r\n",
			err:          nil,
			expected:     []RespType{BulkStrings{Value: "", Length: -1}},
		},
		{
			testcaseName: "Empty Array",
			fileString:   "*0\r\n",
			err:          nil,
			expected:     []RespType{Array{Values: make([]RespType, 0), Length: 0}},
		},
		{
			testcaseName: "Array",
			fileString:   "*3\r\n:1\r\n:2\r\n:3\r\n",
			err:          nil,
			expected: []RespType{Array{Values: []RespType{
				Integer{Value: 1},
				Integer{Value: 2},
				Integer{Value: 3},
			}, Length: 3}},
		},
		{
			testcaseName: "All different types",
			fileString:   "*5\r\n:1\r\n+OK\r\n$6\r\nfoobar\r\n-Error message\r\n*2\r\n:2\r\n+World\r\n",
			err:          nil,
			expected: []RespType{
				Array{
					Values: []RespType{
						Integer{Value: 1},
						SimpleString{Value: "OK"},
						BulkStrings{Value: "foobar", Length: 6},
						SimpleError{Value: "Error message"},
						Array{
							Values: []RespType{
								Integer{Value: 2},
								SimpleString{Value: "World"},
							},
							Length: 2,
						},
					},
					Length: 5,
				}},
		},
		{
			testcaseName: "Unknown Type",
			fileString:   "?Unknown",
			err:          ErrUnkType,
			expected:     make([]RespType, 0),
		},
	}

	for _, testcase := range testcases {

		t.Run(testcase.testcaseName, func(t *testing.T) {
			p := NewParser()
			cleanedLines := cleanLines(testcase.fileString)
			err := p.Parse(cleanedLines)
			if err != testcase.err {
				t.Errorf("expected error: %v, got: %v", testcase.err, err)
			}

			if !reflect.DeepEqual(p.RespList, testcase.expected) {
				t.Fatalf("mismatch:\nexpected: %#v\ngot: %#v", testcase.expected, p.RespList)
			}

		})
	}
}
