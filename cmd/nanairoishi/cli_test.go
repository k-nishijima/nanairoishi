package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRun_usage(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := CLI{outStream: outStream, errStream: errStream}

	status := cli.Run([]string{"nanairoish"})
	if status != ExitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOK)
	}

	expected := Usage
	if !strings.Contains(outStream.String(), expected) {
		t.Errorf("Output=%q, want %q", outStream.String(), expected)
	}
}

func TestRun_versionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("nanairoish -version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("nanairoish version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("Output=%q, want %q", errStream.String(), expected)
	}
}

func TestRun_updateFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := CLI{outStream: outStream, errStream: errStream}

	// sg名無し
	args := strings.Split("nanairoish -update", " ")

	status := cli.Run(args)
	if status != ExitCodeParseFlagError {
		t.Errorf("ExitStatus=%d, want %d", status, ExitCodeParseFlagError)
	}

	expected := fmt.Sprintf("please set security group name")
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("Output=%q, want %q", errStream.String(), expected)
	}

	// sg名あり
	args = strings.Split("nanairoish -update sgName", " ")

	status = cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOK)
	}

	expected = fmt.Sprintf("you specified the security group name called '%v'\n", "sgName")
	if !strings.Contains(outStream.String(), expected) {
		t.Errorf("Output=%q, want %q", outStream.String(), expected)
	}
}
