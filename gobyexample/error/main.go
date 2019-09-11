package main

import (
	"errors"
	"fmt"
)

func f1(args int) (int, error) {
	if args == 43 {
		return -1, errors.New("can't work 43")
	}
	return args + 3, nil
}

type AragError struct {
	arg  int
	prob string
}

func (e *AragError) Error() string {
	return fmt.Sprintf("%d-%s", e.arg, e.prob)
}
func f2(arg int) (int, error) {
	if arg == 43 {
		return -1, &AragError{1, "can't network"}
	}
	return arg + 3, nil
}

func main() {

	for _, index := range []int{1, 43, 50} {
		if r, e := f1(index); e != nil {
			fmt.Println("f1 failed:", e)
		} else {
			fmt.Println("f1 worked", r)
		}
	}

	for _, index := range []int{1, 43, 50} {
		if r, e := f2(index); e != nil {
			fmt.Println("f2 failed:", e)
		} else {
			fmt.Println("f2 worked", r)
		}
	}

	_, e := f2(42)
	if ae, ok := e.(*AragError); ok {
		fmt.Println(ae.arg)
		fmt.Println(ae.prob)
	}

}
