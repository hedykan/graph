package graph

import "errors"

type IntStack struct {
	Size int
	Data []int
}

func NewStack() IntStack {
	return IntStack{
		Size: 0,
		Data: make([]int, 0),
	}
}

func (stack *IntStack) Push(data int) {
	stack.Size += 1
	stack.Data = append(stack.Data, data)
}

func (stack *IntStack) Pop() (int, error) {
	if stack.Size == 0 {
		return -1, errors.New("stack: size is zero")
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
