# ifdef

#ifdef golang code generation

```
// #ifdef MACRO

controlled text

// #endif
```

### Installation

```
go install
```

### Example on test file

```cmd
# show list of preprocessor names
ifdef -l -i=./testdata/test.gotmpl

# Preprocessor names :
# * pre1
# * pre2

# generate file `pre1.go` with preprocessor flag `pre1`
ifdef -p=pre1 -i=./testdata/test.gotmpl -o=./testdata/pre1.go
```

**result of file pre1.go**
```golang
package test

// #ifdef pre1
func a1() int {
	// #endif
	// #ifdef pre2
	// 		func a2() float64{
	// #endif
	b := 0
	return b
}
```

```cmd
# generate file `pre2.go` with preprocessor flag `pre2` with `gofmt` result Go source
ifdef -p=pre2 -i=./testdata/test.gotmpl -o=./testdata/pre2.go -f
```


**result of file pre2.go**
```golang
package test

// #ifdef pre1
// func a1() int {
// #endif
// #ifdef pre2
func a2() float64 {
	// #endif
	b := 0
	return b
}
```
