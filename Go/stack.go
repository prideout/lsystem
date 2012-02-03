package main

import "vmath"

// pkg/container/list might make more sense; just an interesting exercise!

type StackNode struct {
    RuleIndex int
    Depth     int
    Transform *vmath.M4
}

type Stack struct {
    s []StackNode
}

func (self *Stack) Clone() *Stack {
    x := new(Stack)
    x.s = make([]StackNode, len(self.s))
    copy(x.s, self.s)
    return x
}

func (self *Stack) Push(stackNode StackNode) {
    self.s = append(self.s, stackNode)
}

func (self *Stack) Pop() StackNode {
    var e StackNode
    n := len(self.s) - 1
    e, self.s = self.s[n], self.s[:n]
    return e
}

func (self *Stack) Len() int {
    return len(self.s)
}
