package lpd

import (
	"bytes"
	"fmt"
	"io"
)

// ControlFile get some options
type ControlFile struct {
	File         []byte
	Hostname     string
	User         string
	JobName      string
	BannerPage   string
	PrintOptions []PrintOption

	// get rid of options for now
}

func (c *ControlFile) String() string {
	return fmt.Sprintf("<Hostname=\"%s\" User=\"%v\" JobName=\"%v\">", c.Hostname, c.User, c.JobName)
}

// AddOption ..
func (c *ControlFile) AddOption(opt PrintOption) []PrintOption {
	c.PrintOptions = append(c.PrintOptions, opt)
	return c.PrintOptions
}

// PrintOption ..
type PrintOption struct {
	Option   rune
	Filename string
}

// NewControlFile Create a new ControlFile
func NewControlFile() *ControlFile {
	cFile := new(ControlFile)

	return cFile
}

// ReadControlFile ...
func ReadControlFile(reader io.Reader, numBytes uint64) (*ControlFile, error) {
	file := make([]byte, numBytes)

	_, err := io.ReadFull(reader, file)

	if err != nil {
		return nil, err
	}

	// For now, absorb control file so we can play with data first
	cFile := new(ControlFile)
	cFile.File = file

	buf := bytes.NewBuffer(cFile.File)

	for {
		line, err := buf.ReadBytes(0x0a)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		cmd := rune(line[0])
		value := string(line[1 : len(line)-1])

		switch cmd {
		case 'H':
			cFile.Hostname = value
			break
		case 'P':
			cFile.User = value
			break
		case 'J':
			cFile.JobName = value
			break
		case 'c', 'd', 'f', 'g', 'l', 'n', 'o', 'p', 'r', 't', 'v':
			cFile.AddOption(PrintOption{cmd, value})
			break
		default:
		}
	}

	return cFile, nil
}
