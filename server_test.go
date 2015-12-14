package lpd

import (
	"encoding/hex"
	"errors"
	"io"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	// timeoutChan := time.After(time.Duration(1) * time.Second)
	clientState := "Nope"
	// ServerState := new(string)

	server, err := NewServer(5515)
	if err != nil {
		t.Fatal(err)
	}

	var ch = make(chan int, 1)

	testHandle := func(p PrintJob) {
		// server.Stop()
		ch <- 1
	}
	server.HandleFunc("LIVRET_A4", testHandle)

	go func() {
		server.Serve()
	}()

	go func() {
		// time.Sleep(100 * time.Millisecond)
		startClient(t, &clientState)
	}()

	select {
	case <-ch:
	case <-time.After(time.Duration(1) * time.Second):
		t.Fatalf("Server timeout : %v", clientState)
	}

}

func startClient(t *testing.T, s *string) {

	// *s = "Connecting..."
	var conn net.Conn
	var err error

	for conn == nil {
		conn, err = net.Dial("tcp", "127.0.0.1:5515")
	}
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	err = clientSendTCmd(conn)
	if err != nil {
		t.Error(err)
	}
	// *s = "1 Command sent..."

	err = waitForACK(conn)
	if err != nil {
		t.Error(err)
	}
	// *s = "1 Command sent...ACK"

	err = clientSendTSCmdCtrl(conn)
	if err != nil {
		t.Error(err)
	}
	*s = "2 SubCtrl sent..."

	err = waitForACK(conn)
	if err != nil {
		t.Error(err)
	}
	*s = "2 SubCtrl sent...ACK"

	err = clientSendTCtrlFile(conn)
	if err != nil {
		t.Error(err)
	}
	*s = "3 CtrlFile sent..."

	err = waitForACK(conn)
	if err != nil {
		t.Error(err)
	}
	*s = "3 CtrlFile sent...ACK"

	err = clientSendTSCmdData(conn)
	if err != nil {
		t.Error(err)
	}
	*s = "4 SubData sent..."

	err = waitForACK(conn)
	if err != nil {
		t.Error(err)
	}
	*s = "4 SubData sent...ACK"

	err = clientSendTDataFile(conn)
	if err != nil {
		t.Error(err)
	}
	*s = "5 DataFile sent..."

	err = waitForACK(conn)
	if err != nil {
		t.Error(err)
	}
	*s = "5 DataFile sent...ACK"

	return
}

////
// Send Recieve job command for queue LIVRET_A4
func clientSendTCmd(w io.Writer) error {
	rawCommand, _ := hex.DecodeString("024c49565245545f41340a")
	_, err := w.Write(rawCommand)
	return err
}

// Send Recieve job Subcommand for controlFile
func clientSendTSCmdCtrl(w io.Writer) error {
	rawCommand, _ := hex.DecodeString("023631206366413130317072696e74636c69656e7431320a")
	_, err := w.Write(rawCommand)
	return err
}

// Send ControlFile
func clientSendTCtrlFile(w io.Writer) error {
	rawCommand, _ := hex.DecodeString("487072696e74636c69656e7431320a50746f746f0a4a626f6f6b6c65742e7064660a4c746f746f0a666466413130317072696e74636c69656e7431320a00")
	_, err := w.Write(rawCommand)
	return err
}

// Send Recieve job Subcommand for DataFile
func clientSendTSCmdData(w io.Writer) error {
	rawCommand := []byte{0x3}
	rawCommand = append(rawCommand, strconv.Itoa(8+1)...)
	dec, _ := hex.DecodeString("206466413130317072696e74636c69656e7431320a")

	rawCommand = append(rawCommand, dec...)
	_, err := w.Write(rawCommand)
	return err
}

// Send DataFile
func clientSendTDataFile(w io.Writer) error {
	rawCommand, _ := hex.DecodeString("497420776f726b730a00")
	_, err := w.Write(rawCommand)
	return err
}

func waitForACK(r io.Reader) error {
	buff := make([]byte, 1)
	_, err := r.Read(buff)
	if err != nil || buff[0] != 0x0 {
		err = errors.New("Error waiting for ACK")
	}
	return err
}
