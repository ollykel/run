package main

/**
 * @author Oliver Kelton, oakelton@gmail.com
 * @date Apr 7, 2019
 * Simple command-line utility to format printing of a process'
 * stdout and stderr using different colors.
 */

import (
	"errors"
	"os"
	"fmt"
	"io"
	"bufio"
	"log"
	"strings"
)

type color string

const (
	Black color = "30"
	Red color = "31"
	Green color = "32"
	Yellow color = "33"
	Blue color = "34"
	Magenta color = "35"
	Cyan color = "36"
	White color = "37"
)//-- end color enums

func PrintColor (w io.Writer, r io.Reader, col color) {
	colorFmt := fmt.Sprint("\x1b[", col, "m")
	format := fmt.Sprint(colorFmt, "%s\x1b[0m")
	reader := bufio.NewReader(r)
	var (
		input string
		err error
	)
	for {
		input, err = reader.ReadString('\n')
		if err != nil {
			if err != io.EOF { log.Print(err.Error()) }
			return
		}
		fmt.Fprintf(w, format, input)
	}//-- end infinite for
}//-- end func PrintColor

var ErrExecutableNotFound = errors.New("executable not found")

func GetExecutableName (name string) (string, error) {
	path, exists := os.LookupEnv("PATH")
	if !exists { log.Fatal("PATH environment var not found") }
	paths := strings.Split(path, ":")
	paths = append([]string{ "." }, paths...)
	var (
		filename string
		err error
	)
	for i := range paths {
		filename = paths[i] + "/" + name
		_, err = os.Stat(filename)
		if err == nil { return filename, nil }
	}//-- end for range paths
	return "", ErrExecutableNotFound
}//-- end func GetExecutableName

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr,
			"\x1b[8mrun requires at least one arg\n\x1b[0m")
		return
	}
	stderrReader, stderrWriter, err := os.Pipe()
	defer stderrReader.Close()
	if err != nil { log.Fatal(err.Error()) }
	attribs := os.ProcAttr{
		Files: []*os.File { os.Stdin, os.Stdout, stderrWriter } }
	go PrintColor(os.Stderr, stderrReader, Red)
	procName, err := GetExecutableName(os.Args[1])
	if err == ErrExecutableNotFound {
		fmt.Fprintf(os.Stderr, "\x1b[31m\"%s\" not found\n\x1b[0m",
			os.Args[1])
		return
	}
	if err != nil { log.Fatal(err) }
	proc, err := os.StartProcess(procName, os.Args[1:], &attribs)
	if err != nil { log.Fatal(err.Error()) }
	_, err = proc.Wait()
	if err != nil { log.Fatal(err.Error()) }
	os.Exit(0)
}//-- end func main

