package pkg

import (
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
			err: nil,
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
		
	}{
		// {
		// 	testcaseName: "Normal case: ini file is present",
		// },
		// {
		// 	testcaseName: "corner case: empty file",
		// },
		{
			testcaseName: "corner case: file not found",
		},
	}

	for _, testcase := range testcases {

		t.Run(testcase.testcaseName, func(t *testing.T) {
			p := NewParser()
			p.LoadFromFile(filePath)

			// for _, resp := range p.RespList {
			// 	switch resp.Type() {
			// 	case "SimpleString":
			// 		t.Logf("Simple String: %s", resp)
			// 	case "SimpleError":
			// 		t.Logf("Simple Error: %s", resp)
			// 	case "Integer":
			// 		t.Logf("Integer: %s", resp)
			// 	default:
			// 		t.Errorf("Unknown RESP type: %s", resp.Type())
			// 	}
			// }
			
		})
	}
}
