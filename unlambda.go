package main

import (
	"fmt"
	"log"
	"regexp"
	//	"log"
)

//////////////////////////////////////////////
type TokenReader struct {
	tokens []string
	ptr    int
}

func tokenize(program string) []string {
	result := []string{}
	re := regexp.MustCompile(`s|k|i|r|\..|` + "`")
	for _, t := range re.FindAllStringSubmatch(program, -1) {
		result = append(result, t[0])
	}
	return result
}

func (tr *TokenReader) next() string {
	tr.ptr++
	return tr.tokens[tr.ptr]
}
func (tr *TokenReader) eos() bool {
	return (tr.ptr+1 >= len(tr.tokens))
}

///////////////////////////////////////////////
type Node interface {
	eval() Obj
	String() string
	String2() string
}

func make_node(tr *TokenReader) Node {
	if tr.eos() {
		log.Fatal("[[make_node() -> EOS]]")
	}
	s := tr.next()
	if s == "`" {
		f := make_node(tr)
		a := make_node(tr)
		return Apply{f, a}
	} else {
		return SimpleFunction{s}
	}
}

///////////////////////////////////////////////
type Obj interface {
	call(a Obj) Obj
	String() string
}

///////////////////////////////////////////////////
type SimpleFunction struct {
	name string
}

func (self SimpleFunction) eval() Obj {
	return self
}
func (self SimpleFunction) String() string {
	return self.name
}
func (self SimpleFunction) String2() string {
	return self.String()
}

func (self SimpleFunction) call(a Obj) Obj {
	switch self.name {
	case "i":
		return a
	case "r":
		fmt.Println()
		return a
	case "k":
		return ComplexFunction{"K2", a, nil}
	case "s":
		return ComplexFunction{"S2", a, nil}
	default:
		if len(self.name) == 2 && self.name[0] == '.' {
			fmt.Printf("%c", self.name[1])
			return a
		}
	}
	log.Fatalf("[[SimpleFunction::call(%s,%s) Error]]", self, a)
	return nil
}

/////////////////////////////////////////////////
type Apply struct {
	f Node
	a Node
}

func (self Apply) eval() Obj {
	ef := self.f.eval()
	ea := self.a.eval()
	return ef.call(ea)
}

func (self Apply) String() string {
	return "`" + self.f.String() + self.a.String()
}
func (self Apply) String2() string {
	return "(" + self.f.String2() + " " + self.a.String2() + ")"
}

////////////////////////////////////////////
type ComplexFunction struct {
	name string
	val1 Obj
	val2 Obj
}

func (self ComplexFunction) String() string {
	switch self.name {
	case "K2":
		return "(k " + self.val1.String() + ")"
	case "S2":
		return "(s " + self.val1.String() + ")"
	case "S3":
		return "(s " + self.val1.String() + " " + self.val2.String() + ")"
	}
	log.Fatalf("[[ComplexFunction::String(%s) Error]]", self)
	return ""
}
func (self ComplexFunction) call(a Obj) Obj {
	switch self.name {
	case "K2":
		return self.val1
	case "S2":
		return ComplexFunction{"S3", self.val1, a}
	case "S3":
		x := self.val1
		y := self.val2
		z := a
		return x.call(z).call(y.call(z))
	}
	log.Fatalf("[[ComplexFunction::call(%s) Error]]", self, a)
	return nil
}

///////////////////////////////////////////////
func main() {
	//	flag.Parse()
	//	fmt.Printf("%v\n", flag.Args())

	var src string
	//src = "```skk.a"
	//src = ".a"
	//src = "`i.a"
	//src = "`.a`.bi"
	//src = "```.a.b.ci"
	//src = "`sk"
	src = "`r```````````.H.e.l.l.o. .w.o.r.l.di"

	ts := tokenize(src)
	tr := &TokenReader{ts, -1}
	tree := make_node(tr)
	fmt.Println("DEBUG=>", ts)
	fmt.Println("DEBUG=>", tr)
	fmt.Println("DEBUG=>", tree)
	fmt.Println("DEBUG=>", tree.String2())
	fmt.Println("====== EVAL ======")
	o := tree.eval()
	fmt.Println("\n====== END of EVAL ======")
	fmt.Println("RESULT=>", o)
}
