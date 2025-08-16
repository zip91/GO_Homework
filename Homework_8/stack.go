package main

import (
	"errors"
	"fmt"
)

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, errors.New("stack is empty")
	}
	lastIndex := len(s.items) - 1
	element := s.items[lastIndex]
	s.items = s.items[:lastIndex]
	return element, nil
}

func (s *Stack[T]) Peek() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, errors.New("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func main() {
	var stack Stack[int]

	stack.Push(10)
	stack.Push(20)
	stack.Push(30)

	top, _ := stack.Peek()
	fmt.Println("Peek:", top)

	for !stack.IsEmpty() {
		val, _ := stack.Pop()
		fmt.Println("Pop:", val)
	}

	_, err := stack.Pop()
	if err != nil {
		fmt.Println("Ошибка:", err)
	}
}
