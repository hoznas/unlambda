package main

import (
	"fmt"
	"regexp"
)

type F interface {
	to_s() string
	call(arg F) F
}

type I struct{}

func (i I) to_s() string {
	return "I"
}
func (i I) call(arg F) F {
	return arg
}

type K struct{}

func (k K) to_s() string {
	return "K"
}
func (k K) call(arg F) F {
	return K2{arg}
}

type K2 struct {
	arg F
}

func (k2 K2) to_s() string {
	return fmt.Sprintf("K2{%s}", k2.arg.to_s())
}
func (k2 K2) call(arg F) F {
	// eval(arg)
	return k2.arg
}

type S struct{}

func (s S) to_s() string { return "S" }
func (s S) call(arg1 F) F {
	return S2{arg1}
}

type S2 struct{ arg1 F }

func (s2 S2) to_s() string {
	return fmt.Sprintf("S2{%s}", s2.arg1.to_s())
}
func (s2 S2) call(arg2 F) F {
	return S3{s2.arg1, arg2}
}

type S3 struct {
	arg1 F
	arg2 F
}

func (s3 S3) to_s() string {
	return fmt.Sprintf("S3{%s, %s}", s3.arg1.to_s(), s3.arg2.to_s())
}
func (s3 S3) call(arg3 F) F {
	//```Sxyz == ``xz`yz
	//(((S x) y) z) == ((x z) (y z)
	x := s3.arg1
	y := s3.arg2
	z := arg3
	return x.call(z).call(y.call(z))
}

func tokenize(program string) []string {
	result := []string{}
	re := regexp.MustCompile("s|k|i|`")
	for _, t := range re.FindAllStringSubmatch(program, -1) {
		result = append(result, t[0])
	}

	//re.ReplaceAllStringFunc(program, func(s string) string {
	//	fmt.Println("AAAA ", s)
	//	result = append(result, s)
	//	return ""
	//})
	return result
}

type TokenReader struct {
	tokens []string
	ptr    int
}

func (tr *TokenReader) next() string {
	tr.ptr++
	return tr.tokens[tr.ptr]
}
func (tr *TokenReader) eos() bool {
	return (tr.ptr+1 >= len(tr.tokens))
}

func main() {
	var f F

	i := I{}
	f = i
	fmt.Println("hello", f.to_s())

	k := K{}
	k2 := k.call(i)
	f = k
	fmt.Println("hello", f.to_s())
	f = k2
	fmt.Println("hello", f.to_s())
	kii := k.call(i).call(i)
	f = kii
	fmt.Println("hello", f.to_s())

	s := S{}
	s2 := s.call(k)
	s3 := s2.call(k)
	f = s
	f = s2
	f = s3
	f = s3.call(s)
	fmt.Println(f.to_s())

	fmt.Println("====================")
	code := "```skkk"
	fmt.Println(eval_string(code).to_s())
}
func eval(tr *TokenReader) F {
	if tr.eos() {
		fmt.Println("ERROR ::: ")
	}
	t := tr.next()
	switch t {
	case "`":
		f := eval(tr)
		a := eval(tr)
		return f.call(a)
	case "s":
		return S{}
	case "k":
		return K{}
	case "i":
		return I{}
	default:
		fmt.Printf("ERROR:eval()  >> [%s]", t)
		return nil
	}
}
func eval_string(program string) F {
	tokens := tokenize(program)
	tr := &TokenReader{tokens, -1}
	return eval(tr)
}
