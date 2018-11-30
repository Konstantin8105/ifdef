package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var osStdout *os.File = os.Stdout

const (
	beginToken  = "#ifdef"
	finishToken = "#endif"
)

type arrayStrings []string

func (a *arrayStrings) String() string {
	return fmt.Sprintf("%v", []string(*a))
}

func (a *arrayStrings) Set(value string) error {
	v := []string(*a)
	v = append(v, value)
	*a = arrayStrings(v)
	return nil
}

// flags
var (
	listFlag   *bool
	gofmtFlag  *bool
	inputFile  *string
	outputFile *string
	pres       arrayStrings
)

func init() {
	listFlag = flag.Bool("l", false, "show list of preprocessor names")
	gofmtFlag = flag.Bool("f", false, "gofmt output file")
	inputFile = flag.String("i", "", "name of input Go source")
	outputFile = flag.String("o", "", "name of output Go source")
	flag.Var(&pres, "p", "allowable preprocessors #ifdef...#endif")
}

func main() {
	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Name of input file is empty")
		return
	}

	// show list of preprocessors names
	if *listFlag {
		if err := list(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
		}
		return
	}

	if *outputFile == "" {
		fmt.Fprintf(os.Stderr, "Name of output file is empty")
		return
	}

	if len(([]string)(pres)) == 0 {
		fmt.Fprintf(os.Stderr, "List of allowable preprocessor names is empty")
		return
	}

	err := change()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return
	}

	if *gofmtFlag {
		cmd := exec.Command("gofmt", "-w", *outputFile)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			return
		}
		fmt.Fprintf(osStdout, string(out))
	}
}

func list() error {
	b, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return fmt.Errorf("cannot read file `%s` : %v", *inputFile, err)
	}

	lines := bytes.Split(b, []byte("\n"))
	var names []string
	for i := range lines {
		// filter
		if !bytes.Contains(lines[i], []byte(beginToken)) {
			continue
		}

		// find name
		index := bytes.Index(lines[i], []byte(beginToken))
		if index < 0 {
			continue
		}

		// get name
		name := string(lines[i][index+len(beginToken):])
		names = append(names, strings.TrimSpace(name))
	}

	// show names
	fmt.Fprintf(osStdout, "Preprocessor names :\n")
	for i := range names {
		fmt.Fprintf(osStdout, "* %s\n", names[i])
	}

	return nil
}

func change() error {
	b, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return fmt.Errorf("cannot read file `%s` : %v", *inputFile, err)
	}

	ps := []string(pres)

	lines := bytes.Split(b, []byte("\n"))
	var addComment bool
	for i := range lines {
		// beginToken
		index := bytes.Index(lines[i], []byte(beginToken))
		if index >= 0 {
			// get name
			name := strings.TrimSpace(string(lines[i][index+len(beginToken):]))
			addComment = true
			for j := range ps {
				if name == ps[j] {
					addComment = false
					break
				}
			}
			continue
		}

		// finishToken
		index = bytes.Index(lines[i], []byte(finishToken))
		if index >= 0 {
			// get name
			addComment = false
			continue
		}

		tmp := string(lines[i])
		tmp = strings.TrimSpace(tmp)
		hasComment := strings.HasPrefix(tmp, "//")

		if addComment && hasComment {
			continue
		}
		if addComment {
			lines[i] = append([]byte("/"+"/"), lines[i]...)
			continue
		}

		// remove begin comment
		if hasComment {
			lines[i] = bytes.Replace(lines[i], []byte("// "), []byte(""), 1)
		}
	}

	return ioutil.WriteFile(*outputFile, bytes.Join(lines, []byte("\n")), 0644)
}
