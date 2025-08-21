package pkg

import (
	"errors"
	"fmt"
	"os"
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
	Values []RespType
}

func (b Array) Type() string { return "Array" }
func (a Array) String() string {

	s := fmt.Sprintf("*%d\r\n", len(a.Values))
	for _, v := range a.Values {
		s += v.String() + "\r\n"
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

func (p *Parser) parse(lines []string) error {
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		switch line[0] {
		case '+': //simple String
			p.RespList = append(p.RespList, SimpleString{Value: line[1:]})

		case '-': //simple error
			p.RespList = append(p.RespList, SimpleError{Value: line[1:]})

		case ':': //integer
			var value int
			if _, err := fmt.Sscanf(line[1:], "%d", &value); err != nil {
				return ErrInvalidRESP
			}
			p.RespList = append(p.RespList, Integer{Value: value})
		case '$': //Bulk String $<length>\r\n<data>\r\n
			var length int
			if _, err := fmt.Sscanf(line[1:], "%d", &length); err != nil {
				return ErrInvalidRESP
			}

			if length < 0 { // null bulk string
				p.RespList = append(p.RespList, BulkStrings{Value: "", Length: -1})
				continue
			}
			value := lines[i+1]

			if len(value) != length {
				return ErrInvalidRESP
			}

			p.RespList = append(p.RespList, BulkStrings{Value: value, Length: length})
			i++

		default:
			return ErrUnkType
		}
	}
	return nil
}

func (p *Parser) cleanLines(data string) []string {
	lines := strings.Split(string(data), "\n")

	cleanedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) == 0 {
			continue
		}
		cleanedLines = append(cleanedLines, trimmedLine)
	}
	return cleanedLines
}

func (p *Parser) LoadFromFile(path string) error {

	data, err := os.ReadFile(path)
	if err != nil {
		return ErrFileNotExist
	}

	if len(data) == 0 {
		return ErrEmptyFile
	}

	cleanedLines := p.cleanLines(string(data))
	p.parse(cleanedLines)
	return nil
}
