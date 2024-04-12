package parser

type contextStack struct {
	stack []context
	cnt   int
}

func (s *contextStack) pop() context {
	s.cnt--
	top := s.stack[s.cnt]
	s.stack = s.stack[:s.cnt]
	return top
}

func (s *contextStack) push(stackFame context) {
	s.stack = append(s.stack, stackFame)
	s.cnt++
}

func (s *contextStack) isEmpty() bool {
	return s.cnt == 0
}

func newContextStack() contextStack {
	return contextStack{
		stack: []context{},
		cnt:   0,
	}
}
