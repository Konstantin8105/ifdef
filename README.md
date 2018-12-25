# ifdef

#ifdef golang code generation

![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)


### Preprocessor pattern

```golang
// #ifdef LIST_OF_MACRO_NAMES

controlled text

// #endif
```

### Preprocessor key-value pattern

```golang
// #kv MACRO_NAME key:value
```
That pattern are removed at the result source.

Example:
```golang
// #kv Float64 short:F64
// #kv Float64 type:float64
// #kv Float32 short:F32
// #kv Float32 type:float32

// Min#short return minimal of 2 values with type #type
func Min#short(a,b #type) #type{
	if a > b {
		return b
	}
	return a
}
```

### Installation

```
go install
```

### CLI commands

```cmd
./ifdef -h
```

```
Usage of ./ifdef:
  -i string
    	name of input Go source
  -l	show list of preprocessor names
  -o string
    	name of output Go source
  -p string
    	allowable preprocessors #ifdef...#endif
```

### Example on test file

Go template (see file `./testdata/test.gotmpl`):
```go
package test

// some comment

// #ifdef pre1
// a1 function return zero value of int
 func a1() int {
// #endif
// #ifdef pre2
// a2 function return zero value of float64
func a2() float64{
// #endif
// templorary variable
	b := 0
	return b
}
```

```cmd
# show list of preprocessor names
ifdef -l -i=./testdata/test.gotmpl

# Preprocessor names :
# * pre1
# * pre2
```

Example of generate `pre1`:

```cmd
# generate file `pre1.go` with preprocessor flag `pre1`
ifdef -p=pre1 -i=./testdata/test.gotmpl -o=./testdata/pre1.go
```

Result of file pre1.go:
```golang
package test

// some comment

// a2 function return zero value of float64
func a2() float64{
// templorary variable
	b := 0
	return b
}
```

Example of generate `pre2`:

```cmd
# generate file `pre2.go` with preprocessor flag `pre2`
ifdef -p=pre2 -i=./testdata/test.gotmpl -o=./testdata/pre2.go
```


Result of file pre2.go:
```golang
package test

// some comment

// a1 function return zero value of int
func a1() int {
	// templorary variable
	b := 0
	return b
}
```
