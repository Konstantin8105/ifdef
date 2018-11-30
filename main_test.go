package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		*inputFile = "./testdata/test.gotmpl"
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
		s2 := "Preprocessor names :\n* pre1\n* pre2\n"
		if strings.Compare(s1, s2) != 0 {
			t.Fatalf("Not equal:\n%s\n%s", s1, s2)
		}
	})
}
