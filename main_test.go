package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	codestyle "github.com/Konstantin8105/cs"
)

func Test(t *testing.T) {
	*inputFile = "./testdata/test.gotmpl"

	t.Run("list", func(t *testing.T) {
		*listFlag = true
		defer func() {
			*listFlag = false
		}()

		f, err := ioutil.TempFile("", "list")
		if err != nil {
			t.Fatal(err)
		}
		tmp := osStdout
		osStdout = f
		defer func() {
			osStdout = tmp
		}()

		err = list()
		if err != nil {
			t.Fatal(err)
		}

		filename := f.Name()
		err = f.Close()
		if err != nil {
			t.Fatal(err)
		}

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}

		s1 := string(b)
		s2 := "Preprocessor names :\n* Float32\n* Float64\n* pre1\n* pre2\n"
		if strings.Compare(s1, s2) != 0 {
			t.Fatalf("Not equal:\n%s\n%s", s1, s2)
		}
	})

	t.Run("preprocessors", func(t *testing.T) {
		tcs := []struct {
			outputFile string
			pres       string
		}{
			{
				outputFile: "./testdata/pre1.go",
				pres:       "pre1",
			},
			{
				outputFile: "./testdata/pre2.go",
				pres:       "pre2",
			},
			{
				outputFile: "./testdata/f32.go",
				pres:       "Float32",
			},
			{
				outputFile: "./testdata/f64.go",
				pres:       "Float64",
			},
		}
		defer func() {
			*outputFile = ""
			*pres = ""
		}()

		for _, tc := range tcs {
			t.Run(tc.pres, func(t *testing.T) {
				*outputFile = tc.outputFile
				*pres = tc.pres
				err := change()
				if err != nil {
					t.Fatal(err)
				}

				// compare results
				out, err := ioutil.ReadFile(tc.outputFile)
				if err != nil {
					t.Fatal(err)
				}
				exp, err := ioutil.ReadFile(tc.outputFile + ".expect")
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(out, exp) {
					t.Fatalf("result is not same")
				}
			})
		}
	})
}

func TestIfDef(t *testing.T) {
	tcs := []struct {
		line, ps   string
		has, found bool
	}{
		{
			line:  "// #ifdef FL",
			ps:    "FL",
			has:   true,
			found: true,
		},
		{
			line:  "// #ifdef FL FL2 DL",
			ps:    "FL",
			has:   true,
			found: true,
		},
		{
			line:  "// #ifdef FL3 FL2 DL1",
			ps:    "FL",
			has:   true,
			found: false,
		},
		{
			line:  "// #ifdef FL3 FL DL",
			ps:    "FL",
			has:   true,
			found: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.line, func(t *testing.T) {
			has, _ := hasIfdef(tc.line, tc.ps)
			if has != tc.has {
				t.Errorf("not same")
			}
		})
	}
}

func TestCodeStyle(t *testing.T) {
	codestyle.All(t)
}
