package pkg

import (
	"errors"
	"fmt"
	"strings"
)

type RespType interface {
	Type() string
	String() string
}

type SimpleString struct {
	Value string
}

func (s SimpleString) Type() string   { return "SimpleString" }
func (s SimpleString) String() string { return fmt.Sprintf("+%s\r\n", s.Value) }

type SimpleError struct {
	Value string
}

func (e SimpleError) Type() string   { return "SimpleError" }
func (e SimpleError) String() string { return fmt.Sprintf("-%s\r\n", e.Value) }

type Integer struct {
	Value int
}

func (i Integer) Type() string   { return "Integer" }
func (i Integer) String() string { return fmt.Sprintf(":%d\r\n", i.Value) }

type BulkStrings struct {
	Length int
	Value  string
}

func (b BulkStrings) Type() string   { return "BulkStrings" }
func (b BulkStrings) String() string { return fmt.Sprintf("$%d\r\n%s\r\n", len(b.Value), b.Value) }

type Array struct {
	Length int
	Values []RespType
}

func (b Array) Type() string { return "Array" }
func (a Array) String() string {

	s := fmt.Sprintf("*%d\r\n", len(a.Values))
	for _, v := range a.Values {
		s += v.String()
	}
	return s
}

type Parser struct {
	RespList []RespType
}

func NewParser() *Parser {
	return &Parser{
		RespList: make([]RespType, 0),
	}
}

var ErrFileNotExist = errors.New("file doesn't exist")
var ErrEmptyFile = errors.New("file is empty")

var ErrInvalidRESP = errors.New("invalid RESP syntax")
var ErrUnkType = errors.New("unknown RESP type")

func (p *Parser) Parse(data []byte) error {
	for i := 0; i < len(data); {
		element, consumed, err := p.parseOne(data[i:])
		if err != nil {
			return err
		}
		p.RespList = append(p.RespList, element)
		i += consumed
	}
	return nil
}

func (p *Parser) parseOne(data []byte) (RespType, int, error) {

	idx := strings.Index(string(data), "\r\n")
	if idx == -1 {
		return nil, 0, ErrInvalidRESP
	}
	switch data[0] {
	case '+': //simple String	+OK\r\n
		return SimpleString{Value: string(data[1:idx])}, idx + 2, nil

	case '-': //simple error	-Error message\r\n
		return SimpleError{Value: string(data[1:idx])}, idx + 2, nil

	case ':': //integer		:1000\r\n
		var value int
		if _, err := fmt.Sscanf(string(data[1:idx]), "%d", &value); err != nil {
			return nil, 0, ErrInvalidRESP
		}
		return Integer{Value: value}, idx + 2, nil

	case '$': //Bulk String		$<length>\r\n<data>\r\n
		var length int
		if _, err := fmt.Sscanf(string(data[1:idx]), "%d", &length); err != nil {
			return nil, 0, ErrInvalidRESP
		}

		if length < 0 { // null bulk string		$-1\r\n
			return BulkStrings{Value: "", Length: -1}, idx + 2, nil
		}

		if length == 0 { // empty bulk string 		$0\r\n\r\n
			return BulkStrings{Value: "", Length: 0}, idx + 4, nil
		}

		start := idx + 2
		end := start + length
		value := string(data[start:end])

		return BulkStrings{Value: value, Length: length}, end + 2, nil

	case '*': //Array *<count>\r\n<value1>\r\n<value2>\r\n...
		var length int
		if _, err := fmt.Sscanf(string(data[1:idx]), "%d", &length); err != nil {
			return nil, 0, ErrInvalidRESP
		}

		if length == 0 { // empty array
			return Array{Values: make([]RespType, 0), Length: 0}, idx + 2, nil
		}

		arr := Array{Values: make([]RespType, 0, length), Length: length}

		consumed := idx + 2 // already consumed the count line

		for j := 0; j < length; j++ {
			element, parsed, err := p.parseOne(data[consumed:])
			if err != nil {
				return nil, 0, err
			}
			arr.Values = append(arr.Values, element)
			consumed += parsed
		}
		return arr, consumed, nil
	}
	return nil, 0, ErrUnkType
}
