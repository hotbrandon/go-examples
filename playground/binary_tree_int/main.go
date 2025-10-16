package main

import "fmt"

type IntTree struct {
	val         int
	left, right *IntTree
}

func (it *IntTree) Insert(val int) *IntTree {
	if it == nil {
		return &IntTree{val: val}
	}
	if val < it.val {
		it.left = it.left.Insert(val)
	} else if val > it.val {
		it.right = it.right.Insert(val)
	}
	return it
}

func (it *IntTree) Contains(val int) bool {
	switch {
	case it == nil:
		return false
	case val < it.val:
		return it.left.Contains(val)
	case val > it.val:
		return it.right.Contains(val)
	default:
		return true
	}
}

func (it *IntTree) Traverse() []int {
	if it == nil {
		return nil
	}

	result := []int{}
	result = append(result, it.left.Traverse()...)
	result = append(result, it.val)
	result = append(result, it.right.Traverse()...)
	return result

}

func main() {
	var it *IntTree
	it = it.Insert(5)
	it = it.Insert(3)
	it = it.Insert(10)
	it = it.Insert(2)
	it = it.Insert(8)
	it = it.Insert(6)

	fmt.Println(it.Contains(2))  // true
	fmt.Println(it.Contains(12)) // false

	fmt.Println(it.Traverse())
}
