package tokenizer

import (
	"bufio"
	"io"
)

type RuneReader struct {
	reader   bufio.Reader
	current  rune
	previous rune
	line     int
	column   int
	done     bool
	position textPositionCounter
}

func NewRuneReader(reader io.Reader) RuneReader {
	bufReader := bufio.NewReader(reader)

	return RuneReader{
		reader:   *bufReader,
		position: newTextPositionCounter(),
		done:     false,
	}
}

func (r *RuneReader) move() error {
	rn, _, err := r.reader.ReadRune()

	if err != nil {
		if err == io.EOF {
			r.done = true
			return nil
		}
		return err
	}

	r.position.increase(rn)

	r.line = r.position.line
	r.column = r.position.column

	r.previous = r.current
	r.current = rn

	return nil
}

func (r *RuneReader) rewind() error {
	err := r.reader.UnreadRune()
	if err != nil {
		return err
	}

	r.position.decrease(r.current)

	r.current = r.previous

	r.line = r.position.line
	r.column = r.position.column

	return nil
}
