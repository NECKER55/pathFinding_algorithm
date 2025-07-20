package main

type Stack struct {
	head *StackNode
}

type StackNode struct {
	value          Point
	actions        []Point
	counter_action int
	next           *StackNode
}

func (q *Stack) push(value Point, actions []Point) {
	if q.head == nil {
		q.head = &StackNode{value, actions, 0, nil}
		return
	}

	newNode := &StackNode{value, actions, 0, q.head}
	q.head = newNode
}

func (q *Stack) pop() (Point, []Point, int) {
	if q.head != nil {
		q.head = q.head.next
	}
	return q.head.value, q.head.actions, q.head.counter_action

}

func (q *Stack) top() Point {
	if q.head != nil {
		return q.head.value
	}
	return Point{}
}

func (q *Stack) isEmpty() bool {
	return q.head == nil
}
