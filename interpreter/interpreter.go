package interpreter

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jbert/gof/stack"
)

type Interpreter struct {
	stack stack.Stack[int64]
	funcs map[string]func() error
}

func New() *Interpreter {
	i := &Interpreter{
		stack: stack.New[int64](),
		funcs: make(map[string]func() error),
	}
	i.installStdLib()
	return i
}

func (i *Interpreter) installStdLib() {
	i.funcs["+"] = i.intAdd
	i.funcs["-"] = i.intSub
	i.funcs["."] = i.print
	i.funcs["dup"] = i.dup

	var keys []string
	for k, _ := range i.funcs {
		keys = append(keys, k)
	}
	for _, key := range keys {
		ucKey := strings.ToUpper(key)
		i.funcs[ucKey] = i.funcs[key]
	}
}

func (i *Interpreter) dup() error {
	return i.runFunc(1, func(xs []int64) error {
		i.Push(xs[0])
		i.Push(xs[0])
		return nil
	})
}

func (i *Interpreter) print() error {
	return i.runFunc(1, func(xs []int64) error {
		fmt.Printf("%d\n", xs[0])
		return nil
	})
}

func (i *Interpreter) intAdd() error {
	return i.runFunc(2, func(xs []int64) error {
		i.Push(xs[0] + xs[1])
		return nil
	})
}

func (i *Interpreter) intSub() error {
	return i.runFunc(2, func(xs []int64) error {
		i.Push(xs[1] - xs[0])
		return nil
	})
}

func (i *Interpreter) runFunc(nArgs int, f func(xs []int64) error) error {
	xs := make([]int64, nArgs)
	for ii := range xs {
		var err error
		xs[ii], err = i.Pop()
		if err != nil {
			return fmt.Errorf("Error on arg [%d/%d]: %s", ii, nArgs, err)
		}
	}
	return f(xs)
}

func (i *Interpreter) Push(v int64) {
	(&(i.stack)).Push(v)
}

func (i *Interpreter) Pop() (int64, error) {
	return (&(i.stack)).Pop()
}

func (i *Interpreter) MustPop() int64 {
	return (&(i.stack)).MustPop()
}

func (i *Interpreter) DumpStack(w io.Writer) {
	first := true
	fmt.Fprintf(w, "[")
	i.stack.ForEach(func(n int64) {
		if first {
			first = false
		} else {
			fmt.Fprintf(w, ", ")
		}
		fmt.Fprintf(w, "%d", n)
	})
	fmt.Fprintf(w, "]\n")
}

func (i *Interpreter) Run(prog string) error {
	tokens := strings.Split(prog, " ")
	for _, t := range tokens {
		err := i.RunToken(t)
		if err != nil {
			return fmt.Errorf("Error running token [%s]: %s", t, err)
		}
	}
	return nil
}

func (i *Interpreter) RunToken(token string) error {
	f, ok := i.funcs[token]
	if ok {
		err := f()
		if err != nil {
			return fmt.Errorf("Error running [%s]: %s", token, err)
		}
		return nil
	}

	n, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid numeric literal [%s]: %s", token, err)
	}
	i.Push(n)
	return nil
}
