package stacks

import "golang.org/x/tools/go/ssa"

// Used by CFG to traverse the graph. It uses both as a stack for traversal by order, and as a map to count occurrences
// and fast retrieval of items.
type BasicBlockStack struct {
	stack     blocksStack
	blocksMap blocksMap
}

func NewBasicBlockStack() *BasicBlockStack {
	stack := make([]*ssa.BasicBlock, 0)
	blocksMap := make(blocksMap)
	basicBlockStack := &BasicBlockStack{stack: stack, blocksMap: blocksMap}
	return basicBlockStack
}

func (s *BasicBlockStack) Push(v *ssa.BasicBlock) {
	s.stack.Push(v)
	if _, ok := s.blocksMap[v.Index]; !ok {
		s.blocksMap[v.Index] = &basicBlockWithCount{block: v, count: 1}
	} else {
		s.blocksMap[v.Index].count += 1
	}
}

func (s *BasicBlockStack) Pop() *ssa.BasicBlock {
	v := s.stack.Pop()
	if v == nil {
		return nil
	}
	s.blocksMap[v.Index].count -= 1
	if s.blocksMap[v.Index].count == 0 {
		delete(s.blocksMap, v.Index)
	}
	return v
}

func (s *BasicBlockStack) GetItems() []*ssa.BasicBlock {
	return s.stack
}

type basicBlockWithCount struct {
	block *ssa.BasicBlock
	count int
}

type blocksMap map[int]*basicBlockWithCount

type blocksStack []*ssa.BasicBlock

func (s *blocksStack) Push(v *ssa.BasicBlock) {
	*s = append(*s, v)
}

func (s *blocksStack) Pop() *ssa.BasicBlock {
	if len(*s) == 0 {
		return nil
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}


// The function check if the block already exists in the stack and if it does, return true only if it appear twice or more.
// This is used to handle cycles where a path might visit a node more then once. This way a cycle is allowed, but only once
// to avoid infinite recursion. Finding all valid flows is hard so this is chosen as an approximation. There are edge
//cases such as: A->B->D->C->D->E(end) where the graph is not circular, but the traversal won't reach to E since D
//appear a second time and it will exit at that point. This flow is rare and can be achieved using gotos
func (s *BasicBlockStack) Contains(v *ssa.BasicBlock) bool {
	block, ok := s.blocksMap[v.Index]
	if !ok {
		return false
	}
	if block.count >= 2 {
		return true
	}
	return false
}
