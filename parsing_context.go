package jmatch

import "fmt"

type context interface {
	getPath() string
	setValue()
	isObject() bool
	isArray() bool

	// only object use, refactor...
	setKey(string)
	isKeySet() bool
}

type objectParsingContext struct {
	path string
	key  string
}

func (o *objectParsingContext) isKeySet() bool {
	return o.path != o.key
}

func (o *objectParsingContext) setKey(key string) {
	o.key = fmt.Sprintf("%s.%s", o.path, key)
}

func (o *objectParsingContext) getPath() string {
	return o.key
}

func (o *objectParsingContext) isArray() bool {
	return false
}

func (o *objectParsingContext) isObject() bool {
	return true
}

func (o *objectParsingContext) setValue() {
	o.key = o.path
}

func newObjectContext(path string) *objectParsingContext {
	return &objectParsingContext{
		path: path,
		key:  path,
	}
}

type arrayParsingContext struct {
	path       string
	elemsCount int
}

func (a *arrayParsingContext) isKeySet() bool {
	panic("unimplemented")
}

func (a *arrayParsingContext) setKey(string) {
	panic("unimplemented")
}

func (a *arrayParsingContext) getPath() string {
	return fmt.Sprintf("%s.[%d]", a.path, a.elemsCount)
}

func (a *arrayParsingContext) setValue() {
	a.elemsCount++
}

func (a *arrayParsingContext) isObject() bool {
	return false
}

func (a *arrayParsingContext) isArray() bool {
	return true
}

func newArrayContext(path string) *arrayParsingContext {
	return &arrayParsingContext{
		path:       path,
		elemsCount: 0,
	}
}

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
