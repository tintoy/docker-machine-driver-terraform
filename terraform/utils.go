package terraform

import (
	"io"
	"os/exec"
)

type multiCloser struct {
	innerClosers []io.Closer
}

func (closer *multiCloser) Close() error {
	for _, innerCloser := range closer.innerClosers {
		err := innerCloser.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func newMultiCloser(innerClosers ...io.Closer) *multiCloser {
	return &multiCloser{
		innerClosers: innerClosers,
	}
}

func pipeCombinedOutput(command *exec.Cmd) (combinedOutput io.Reader, outputCloser io.Closer, err error) {
	var stdoutPipe, stderrPipe io.ReadCloser

	stdoutPipe, err = command.StdoutPipe()
	if err != nil {
		return
	}

	stderrPipe, err = command.StderrPipe()
	if err != nil {
		return
	}

	combinedOutput = io.MultiReader(stderrPipe, stdoutPipe)
	outputCloser = newMultiCloser(stdoutPipe, stderrPipe)

	return
}
