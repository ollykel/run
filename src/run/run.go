package main

/**
 * @author Oliver Kelton, oakelton@gmail.com
 * @date Apr 7, 2019
 * Simple command-line utility to format printing of a process'
 * stdout and stderr using different colors.
 */

import (
	"os"
	"fmt"
	"io"
	"bufio"
	"log"
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

func main() {
	if len(os.Args) < 2 { log.Fatal("requires at least one arg") }
	stderrReader, stderrWriter, err := os.Pipe()
	defer stderrReader.Close()
	if err != nil { log.Fatal(err.Error()) }
	attribs := os.ProcAttr{
		Files: []*os.File { os.Stdin, os.Stdout, stderrWriter } }
	go PrintColor(os.Stderr, stderrReader, Red)
	procName := os.Args[1]
	proc, err := os.StartProcess(procName, os.Args[1:], &attribs)
	if err != nil { log.Fatal(err.Error()) }
	_, err = proc.Wait()
	if err != nil { log.Fatal(err.Error()) }
	os.Exit(0)
}//-- end func main

