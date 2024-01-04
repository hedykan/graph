package graph

import "errors"

type Stack[T any] struct {
	Size int
	Data []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		Size: 0,
		Data: make([]T, 0),
	}
}

func (stack *Stack[T]) Push(data T) {
	stack.Size += 1
	stack.Data = append(stack.Data, data)
}

func (stack *Stack[T]) Pop() (T, error) {
	if stack.Size == 0 {
		var zero T
		return zero, errors.New("stack: size is zero")
	}
	stack.Size -= 1
	res := stack.Data[stack.Size]
	stack.Data = stack.Data[:stack.Size]
	return res, nil
}

func FindInt(idArr []int, id int) bool {
	for i := 0; i < len(idArr); i++ {
		if idArr[i] == id {
			return true
		}
	}
	return false
}

func FindIntIndex(idArr []int, id int) int {
	for i := 0; i < len(idArr); i++ {
		if idArr[i] == id {
			return i
		}
	}
	return -1
}

func DeleteInt(idArr []int, index int) []int {
	idArr = append(idArr[:index], idArr[index+1:]...)
	return idArr
}
