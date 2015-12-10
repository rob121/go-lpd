package lpd

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// Test ReceivePrintJob Command
func TestBinUnmarshalCommand(t *testing.T) {
	rawCommand, _ := hex.DecodeString("024c49565245545f4134")

	cmd, err := unmarshalCommand(rawCommand)

	if err != nil {
		t.Fatal(err)
	}

	if cmd.Code != ReceivePrintJob {
		t.Fatal("Code not decoded correctly: ", cmd.Code)
	}

	if cmd.Queue != "LIVRET_A4" {
		t.Fatalf("Queue not decoded correctly: '%x'", cmd.Queue)
	}

}

// Test Receive control file subcommand
func TestBinUnmarshalSubCommand(t *testing.T) {
	// fmt.Printf("Dec : %s", 0x20)
	rawCommand, _ := hex.DecodeString("023631206366413130317072696e74636c69656e743132")

	cmd, err := unmarshalSubCommand(rawCommand)

	if err != nil {
		t.Fatal(err)
	}

	if cmd.Code != 0x2 {
		t.Fatal("Code not decoded correctly: ", cmd.Code)
	}

	if cmd.NumBytes != 61 {
		t.Fatalf("NumBytes not decoded correctly: '%v'", cmd.NumBytes)
	}

	if cmd.FileName != "cfA101printclient12" {
		t.Fatalf("FileName not decoded correctly: '%x'", cmd.FileName)
	}

}

// Test control file decode
func TestBinUnmarshalControlFile(t *testing.T) {
	rawFile, _ := hex.DecodeString("487072696e74636c69656e7431320a50746f746f0a4a626f6f6b6c65742e7064660a4c746f746f0a666466413130317072696e74636c69656e7431320a00")

	rd := bytes.NewBuffer(rawFile)

	cFile, err := ReadControlFile(rd, 61)
	if err != nil {
		t.Fatal(err)
	}

	if cFile.Hostname != "printclient12" {
		t.Fatalf("Hostname not decoded correctly: '%v'", cFile.Hostname)
	}

	if cFile.User != "toto" {
		t.Fatalf("User not decoded correctly: '%v'", cFile.User)
	}

	if cFile.JobName != "booklet.pdf" {
		t.Fatalf("JobName not decoded correctly: '%v'", cFile.JobName)
	}

	if cFile.PrintOptions == nil {
		t.Fatal("PrintOptions is empty")
	}

	if cFile.PrintOptions[0].Option != 'f' {
		t.Fatalf("PrintOptions Option not decoded correctly: '%v'", cFile.PrintOptions[0].Option)
	}

	if cFile.PrintOptions[0].Filename != "dfA101printclient12" {
		t.Fatalf("PrintOptions Filename not decoded correctly: '%x'", cFile.PrintOptions[0].Filename)
	}
}

func TestBinUnmarshalSubCommand2(t *testing.T) {
	rawCommand, _ := hex.DecodeString("03313139333038206466413130317072696e74636c69656e743132")

	cmd, err := unmarshalSubCommand(rawCommand)

	if err != nil {
		t.Fatal(err)
	}

	if cmd.Code != 0x3 {
		t.Fatal("Code not decoded correctly: ", cmd.Code)
	}

	if cmd.NumBytes != 119308 {
		t.Fatalf("NumBytes not decoded correctly: '%v'", cmd.NumBytes)
	}

	if cmd.FileName != "dfA101printclient12" {
		t.Fatalf("FileName not decoded correctly: '%x'", cmd.FileName)
	}
}
