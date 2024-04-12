package parser

import (
	"fmt"
	"strings"
)

type context interface {
	getPath() string
	setValue()
	isObject() bool
	isArray() bool

	// only object use, refactor...
	setKey(string)
	isKeySet() bool
}

type objectContext struct {
	path string
	key  string
}

func (o *objectContext) isKeySet() bool {
	return o.path != o.key
}

func (o *objectContext) setKey(key string) {
	if strings.ContainsAny(key, " .") { // keys with . or space.
		o.key = fmt.Sprintf("%s.\"%s\"", o.path, key)
	} else {
		o.key = fmt.Sprintf("%s.%s", o.path, key)
	}
}

func (o *objectContext) getPath() string {
	return o.key
}

func (o *objectContext) isArray() bool {
	return false
}

func (o *objectContext) isObject() bool {
	return true
}

func (o *objectContext) setValue() {
	o.key = o.path
}

func newObjectContext(path string) *objectContext {
	return &objectContext{
		path: path,
		key:  path,
	}
}

type arrayContext struct {
	path       string
	elemsCount int
}

func (a *arrayContext) isKeySet() bool {
	panic("unimplemented")
}

func (a *arrayContext) setKey(string) {
	panic("unimplemented")
}

func (a *arrayContext) getPath() string {
	return fmt.Sprintf("%s.[%d]", a.path, a.elemsCount)
}

func (a *arrayContext) setValue() {
	a.elemsCount++
}

func (a *arrayContext) isObject() bool {
	return false
}

func (a *arrayContext) isArray() bool {
	return true
}

func newArrayContext(path string) *arrayContext {
	return &arrayContext{
		path:       path,
		elemsCount: 0,
	}
}
