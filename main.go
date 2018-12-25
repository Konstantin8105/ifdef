package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var osStdout *os.File = os.Stdout

const (
	ifdefToken = "#ifdef"
	kvToken    = "#kv"
	endToken   = "#endif"
)

var (
	listFlag   *bool
	inputFile  *string
	outputFile *string
	pres       *string
)

func init() {
	{
		var b bool
		listFlag = &b
	}
	{
		var s string
		inputFile = &s
	}
	{
		var s string
		outputFile = &s
	}
	{
		var s string
		pres = &s
	}
}

func main() {
	// flags
	listFlag = flag.Bool("l", false, "show list of preprocessor names")
	inputFile = flag.String("i", "", "name of input Go source")
	outputFile = flag.String("o", "", "name of output Go source")
	pres = flag.String("p", "", "allowable preprocessors #ifdef...#endif")

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

	if *pres == "" {
		fmt.Fprintf(os.Stderr, "List of allowable preprocessor names is empty")
		return
	}

	err := change()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return
	}
}

// list of preprocessors
func list() error {
	b, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return fmt.Errorf("cannot read file `%s` : %v", *inputFile, err)
	}

	lines := bytes.Split(b, []byte("\n"))
	var names []string
	for i := range lines {
		// filter
		if !bytes.Contains(lines[i], []byte(ifdefToken)) {
			continue
		}

		// find name
		index := bytes.Index(lines[i], []byte(ifdefToken))
		if index < 0 {
			continue
		}

		// get name
		name := string(lines[i][index+len(ifdefToken):])
		ls := strings.Split(name, " ")
		for j := range ls {
			ls[j] = strings.TrimSpace(ls[j])
			if ls[j] == "" {
				continue
			}
			names = append(names, ls[j])
		}
	}

	// show names
	fmt.Fprintf(osStdout, "Preprocessor names :\n")
	sort.Strings(names)
	for i := 0; i < len(names); i++ {
		isUniq := true
		for j := i + 1; j < len(names); j++ {
			if names[i] == names[j] {
				isUniq = false
			}
		}
		if !isUniq {
			continue
		}
		fmt.Fprintf(osStdout, "* %s\n", names[i])
	}

	return nil
}

func change() error {
	b, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return fmt.Errorf("cannot read file `%s` : %v", *inputFile, err)
	}

	ps := string(*pres)

	lines := bytes.Split(b, []byte("\n"))
	var buf bytes.Buffer
	addLine := true

	kv := map[string]string{}

	for i := range lines {
		line := string(lines[i])

		// ifdefToken
		has, found := hasIfdef(line, ps)
		if has {
			addLine = false
			if found {
				addLine = true
			}
			continue
		}

		// kvToken
		has, found, key, value := hasKv(line, ps)
		if has {
			if found {
				kv[key] = value
			}
			continue
		}

		// endToken
		index := strings.Index(line, endToken)
		if index >= 0 {
			// get name
			addLine = true
			continue
		}

		if !addLine {
			continue
		}

		// key-value changing
		for k, v := range kv {
			line = strings.Replace(line, "#"+k, v, -1)
		}

		buf.WriteString(line)
		buf.Write([]byte("\n"))
	}

	err = ioutil.WriteFile(*outputFile, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	// gofmt
	_, _ = exec.Command("gofmt", "-s", "-w", *outputFile).CombinedOutput()
	// goimports
	_, _ = exec.Command("goimports", "-w", *outputFile).CombinedOutput()

	return nil
}

func hasIfdef(line, ps string) (has, found bool) {
	index := strings.Index(line, ifdefToken)
	if index >= 0 {
		has = true
		preprocessors := strings.Split(line[index+len(ifdefToken):], " ")
		for i := range preprocessors {
			preprocessors[i] = strings.TrimSpace(preprocessors[i])
			if preprocessors[i] == "" {
				continue
			}
			if preprocessors[i] == ps {
				found = true
				break
			}
		}
	}
	return
}

func hasKv(line, ps string) (has, found bool, key, value string) {
	index := strings.Index(line, kvToken)
	if index >= 0 {
		has = true
		// example : preprocessors = "Float32 keys:value"
		preprocessors := strings.TrimSpace(line[index+len(kvToken):])
		index = strings.Index(preprocessors, " ")
		if index < 0 {
			return
		}
		if preprocessors[:index] != ps {
			return
		}
		found = true
		preprocessors = preprocessors[index:]

		// example : preprocessors = "keys:value"
		preprocessors = strings.TrimSpace(preprocessors)
		index = strings.Index(preprocessors, ":")
		if index < 0 {
			return
		}
		key = preprocessors[:index]
		value = preprocessors[index+1:]
	}
	return
}
