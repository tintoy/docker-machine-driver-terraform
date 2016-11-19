package terraform

import (
	"bufio"
	"io"

	"github.com/docker/machine/libmachine/log"
)

// Scan STDOUT and STDERR pipes for a process.
//
// Calls the supplied PipeHandler once for each line encountered.
func scanProcessPipes(stdioPipe io.ReadCloser, stderrPipe io.ReadCloser, pipeOutput PipeHandler) {
	go scanPipe(stdioPipe, pipeOutput, "STDOUT")
	go scanPipe(stdioPipe, pipeOutput, "STDERR")
}

// Scan a process output pipe, and call the supplied PipeHandler once for each line encountered.
func scanPipe(pipe io.ReadCloser, pipeOutput PipeHandler, pipeName string) {
	lineScanner := bufio.NewScanner(pipe)
	for lineScanner.Scan() {
		line := lineScanner.Text()
		pipeOutput(line)
	}

	scanError := lineScanner.Err()
	if scanError != nil {
		log.Errorf("Error scanning pipe %s: %s",
			pipeName,
			scanError.Error(),
		)
	}

	pipe.Close()
}
