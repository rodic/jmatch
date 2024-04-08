package jmatch

import "fmt"

type ParsingType int

const (
	Object ParsingType = iota
	Array
)

type ParsingContext struct {
	path        string
	elemsCount  int
	parsingType ParsingType
}

func (p ParsingContext) inArray() bool {
	return p.parsingType == Array
}

func (p ParsingContext) inObject() bool {
	return p.parsingType == Object
}

type Stack struct {
	stack []ParsingContext
	cnt   int
}

func (s *Stack) pop() ParsingContext {
	s.cnt--
	top := s.stack[s.cnt]
	s.stack = s.stack[:s.cnt]
	return top
}

func (s *Stack) push(stackFame ParsingContext) {
	s.stack = append(s.stack, stackFame)
	s.cnt++
}

func NewStack() Stack {
	return Stack{
		stack: []ParsingContext{},
		cnt:   0,
	}
}

type ParsingResult map[string]Token

type Parser struct {
	tokens         Tokens
	parsingContext ParsingContext
	stack          Stack
	lastKey        string
	result         ParsingResult
	err            error
}

func NewParser(tokens []Token) Parser {
	return Parser{
		tokens: NewTokens(tokens),
		stack:  NewStack(),
		result: make(ParsingResult),
	}
}

func (p *Parser) IsValue(t Token) bool {
	return t.IsString() || t.IsNumber() || t.IsBoolean() || t.IsNull()
}

func (p *Parser) parseObject() {
	currentToken, nextToken := p.tokens.current(), p.tokens.next()
	if currentToken.IsRightBrace() {
		p.parsingContext = p.stack.pop()
	} else if currentToken.IsString() && nextToken.IsColon() { // key found
		p.lastKey = fmt.Sprintf("%s.%s", p.lastKey, currentToken.Value)
	} else if currentToken.IsLeftBrace() && nextToken.IsString() {
		// pass, will catch later new object start with ': {' match
	} else if p.IsValue(currentToken) && (nextToken.IsComma() || nextToken.IsRightBrace()) {
		// pass, will catch later values with ': v
	} else if currentToken.IsComma() && nextToken.IsString() {
		// pass, it's only valid to have string after comma
	} else if currentToken.IsColon() {
		if p.IsValue(nextToken) {
			p.result[p.lastKey] = nextToken
			p.lastKey = p.parsingContext.path
		} else if nextToken.IsLeftBrace() {
			p.stack.push(p.parsingContext)
			p.parsingContext = ParsingContext{
				path:        p.lastKey,
				parsingType: Object,
			}
		} else if nextToken.IsLeftBracket() {
			p.stack.push(p.parsingContext)
			p.parsingContext = ParsingContext{
				path:        p.lastKey,
				parsingType: Array,
				elemsCount:  0,
			}
		} else {
			p.err = fmt.Errorf("invalid JSON %s -> %s", currentToken.Value, nextToken.Value)
		}
	} else {
		p.err = fmt.Errorf("invalid JSON %s -> %s", currentToken.Value, nextToken.Value)
	}
}

func (p *Parser) parseArray() {
	currentToken, nextToken := p.tokens.current(), p.tokens.next()
	if currentToken.IsRightBracket() {
		p.parsingContext = p.stack.pop()
		p.lastKey = p.parsingContext.path
	} else if p.IsValue(currentToken) && (nextToken.IsComma() || nextToken.IsRightBracket()) {
		newPath := fmt.Sprintf("%s.[%d]", p.parsingContext.path, p.parsingContext.elemsCount)
		p.parsingContext.elemsCount++
		p.result[newPath] = currentToken
	} else if nextToken.IsLeftBracket() {
		newPath := fmt.Sprintf("%s.[%d]", p.parsingContext.path, p.parsingContext.elemsCount)
		p.parsingContext.elemsCount++
		p.stack.push(p.parsingContext)
		p.parsingContext = ParsingContext{
			path:        newPath,
			parsingType: Array,
			elemsCount:  0,
		}
	} else if currentToken.IsLeftBracket() {
		// skipping the first '[', all else covered in the previous case.
	} else if currentToken.IsComma() &&
		!(nextToken.IsComma() || nextToken.IsRightBrace() || nextToken.IsRightBracket()) {
		//
	} else if currentToken.IsLeftBrace() {
		p.lastKey = fmt.Sprintf("%s.[%d]", p.parsingContext.path, p.parsingContext.elemsCount)
		p.parsingContext.elemsCount++
		p.stack.push(p.parsingContext)
		p.parsingContext = ParsingContext{
			path:        p.lastKey,
			parsingType: Object,
		}
	} else {
		p.err = fmt.Errorf("invalid JSON %s -> %s", currentToken.Value, nextToken.Value)
	}
}

func (p *Parser) parse() (ParsingResult, error) {
	p.parsingContext = ParsingContext{path: "", parsingType: Object}
	p.lastKey = ""

	if p.tokens.tokensCount < 2 || !p.tokens.current().IsLeftBrace() {
		p.err = fmt.Errorf("invalid JSON. must start with { and end with }")
	}

	for p.tokens.hasNext() {
		if p.parsingContext.inObject() {
			p.parseObject()
		} else if p.parsingContext.inArray() {
			p.parseArray()
		} else {
			p.err = fmt.Errorf("invalid JSON")
		}

		if p.err != nil {
			return p.result, p.err
		}

		p.tokens.move()
	}
	if !p.tokens.next().IsRightBrace() {
		p.err = fmt.Errorf("invalid JSON. the last token must be }")
	}

	return p.result, p.err
}
