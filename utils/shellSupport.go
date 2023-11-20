package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func ExecShell(name string, arg ...string) error {
	// var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	co := exec.Command(name, arg...)
	stdoutIn, _ := co.StdoutPipe()
	stderrIn, _ := co.StderrPipe()
	// var errStdout error
	// var errStderr error
	// stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	_ = co.Start()
	go func() {
		_, _ = io.Copy(ioutil.Discard, stdoutIn)
	}()
	go func() {
		_, _ = io.Copy(stderr, stderrIn)
	}()
	// if errStderr != nil {
	// 	log.Printf("%v", errStderr)
	// }
	// if errStdout != nil {
	// 	log.Printf("%v", errStdout)
	// }
	err := co.Wait()
	// outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	//println(outStr + errStr)
	return err
}
