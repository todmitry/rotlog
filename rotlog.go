// Package rotlog daily rotated log implementation.
// Example:
// 	import rotlog
// 	file, err := rotlog.Rotate("/var/log", "rotlog.test", 9)
// 	log.SetOutput(io.MultiWriter(file, os.Stdout))
package rotlog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Rotate opens the output file, trims older files
func Rotate(
	folder string, // folder where log files go
	pattern string, // common prefix of log files, .2006-01-02.log will be appended to the filename
	numkeep int, // number of files to keep, oldest (by pattern name) files will be deleted
) (*os.File, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	// list of matching files
	old := make([]string, 0)
	for _, fileinfo := range files {
		name := fileinfo.Name()
		if strings.HasPrefix(name, pattern) {
			old = append(old, name)
		}
	}
	sort.Strings(old)
	numold := len(old)

	// today's file name
	day := time.Now().Format("2006-01-02")
	logName := fmt.Sprintf("%s.%s.log", pattern, day)
	if numold == 0 || old[numold-1] != logName {
		old = append(old, logName)
	}

	// remove older files
	if numkeep > 0 && numold > numkeep {
		for i, name := range old {
			if i < numold-numkeep {
				os.Remove(filepath.Join(folder, name))
			}
		}

	}

	// open and return today's file
	logPath := filepath.Join(folder, logName)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

// Set create a file with today's date and add it to root logger's output streams
func Set(
	folder string, // folder where log files go
	pattern string, // common prefix for log file names (.2006-01-02.log will be appended to the name)
	numkeep int, // keep only these many files
	streams ...io.Writer, // os.Stdout, os.Stderr - streams to use in addition to the log file
) error {
	// create and open the log file
	logFile, err := Rotate(folder, pattern, numkeep)
	if err != nil {
		return err
	}

	streams = append(streams, logFile)
	log.SetOutput(io.MultiWriter(streams...))
	return nil
}
