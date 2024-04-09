package jmatch

import "fmt"

type parsingType int

const (
	Object parsingType = iota
	Array
)

type parsingContext struct {
	path       string
	elemsCount int
	lastKey    string
	_type      parsingType
}

func (p *parsingContext) setLastKey(key string) {
	p.lastKey = fmt.Sprintf("%s.%s", p.path, key)
}

func (p *parsingContext) resetLastKey() {
	p.lastKey = p.path
}

func (p parsingContext) inArray() bool {
	return p._type == Array
}

func (p parsingContext) inObject() bool {
	return p._type == Object
}

func (p parsingContext) arrayPath() string {
	path := fmt.Sprintf("%s.[%d]", p.path, p.elemsCount)
	return path
}

func (p *parsingContext) increaseElemsCount() {
	p.elemsCount++
}

func newParsingContext(path string, pt parsingType) parsingContext {
	return parsingContext{
		path:       path,
		lastKey:    path,
		elemsCount: 0,
		_type:      pt,
	}
}

type contextStack struct {
	stack []parsingContext
	cnt   int
}

func (s *contextStack) pop() parsingContext {
	s.cnt--
	top := s.stack[s.cnt]
	s.stack = s.stack[:s.cnt]
	return top
}

func (s *contextStack) push(stackFame parsingContext) {
	s.stack = append(s.stack, stackFame)
	s.cnt++
}

func newContextStack() contextStack {
	return contextStack{
		stack: []parsingContext{},
		cnt:   0,
	}
}
