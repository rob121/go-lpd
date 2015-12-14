package lpd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// PrintJob describes a print job for a queue in an LPD printer with a
// control file and data file.
type PrintJob struct {
	QueueName   string
	ControlFile *ControlFile
	DataFile    *os.File
}

func (p *PrintJob) String() string {
	return fmt.Sprintf("<Queue=\"%s\" ControlFile=\"%v\" DataFile=\"%v\">", p.QueueName, p.ControlFile, p.DataFile)
}

// GetFile return the data file
func (p *PrintJob) GetFile() *os.File {
	return p.DataFile
}

// NewPrintJob returns a PrintJob configured for a queue with a data file.
// NewPrintJob returns an error if and only if the data file fails to copy
// to a temporary file.
func NewPrintJob(queue string, dataFile io.Reader) (*PrintJob, error) {
	job := new(PrintJob)
	job.QueueName = queue
	job.ControlFile = NewControlFile()

	tempFile, err := ioutil.TempFile(os.TempDir(), "go-lpd")

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(tempFile, dataFile)

	if err != nil {
		return nil, err
	}

	job.DataFile = tempFile

	return job, nil
}

// Handle series of subcommands that describe a print job
func receiveJob(reader io.Reader, writer io.Writer) (*PrintJob, error) {
	job := new(PrintJob)
	var minibuff = []byte{0x0}

	bufReader := bufio.NewReader(reader)

	for {
		select {
		default:
			// fmt.Print("Yo")
			rawSubCommand, err := bufReader.ReadBytes(0x0a)

			if err != nil {
				if job.ControlFile != nil && job.DataFile != nil {
					return job, nil
				}

				return nil, err
			}

			subCmd, err := unmarshalSubCommand(rawSubCommand)
			if err != nil {
				return nil, err
			}
			if subCmd.Code == 0x0 {
				return nil, errors.New("Code Invalid\n")
			}
			// fmt.Printf("\nsub : %v\n", subCmd.Code)

			if err != nil {
				return nil, err
			}

			switch subCmd.Code {
			case AbortJob:
				ackSubCommand(writer)
				return nil, errors.New("job aborted")
			case ReceiveControlFile:
				// fmt.Println("Start RecCon")
				ackSubCommand(writer)
				cFile, err := ReadControlFile(reader, subCmd.NumBytes)
				if err != nil {
					return nil, err
				}

				_, err = reader.Read(minibuff)
				if err != nil || minibuff[0] != 0x0 {
					log.Fatal(err)
				}

				ackSubCommand(writer)
				// fmt.Print("Yo3")

				job.ControlFile = cFile
				// fmt.Println("End RecCon")
			case ReceiveDataFile:
				// fmt.Println("Start RecDat")
				ackSubCommand(writer)
				// fmt.Println("Start RecDat.")

				dataFile, err := readDataFile(reader, subCmd.NumBytes)
				// fmt.Println("Start RecDat..")

				if err != nil {
					return nil, err
				}
				// fmt.Println("Start RecDat...")
				// fmt.Println(subCmd.NumBytes)

				_, err = reader.Read(minibuff)
				if err != nil || minibuff[0] != 0x0 {
					log.Fatal(err)
				}
				// fmt.Println("Start RecDat....")

				ackSubCommand(writer)
				// fmt.Println("Start RecDat.....")

				job.DataFile = dataFile

				// fmt.Println("End RecCon")

				return job, nil
			}
		}
	}
}

// Send an octect of 0 to acknowledge a subcommand
func ackSubCommand(writer io.Writer) error {
	_, err := writer.Write([]byte{0x0})

	if err != nil {
		return err
	}

	return nil
}

// Send an octect of 1 to refuse a subcommand
func nackSubCommand(writer io.Writer) error {
	_, err := writer.Write([]byte{0x1})

	if err != nil {
		return err
	}

	return nil
}
