//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------
package scan

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

//-----------------------------------------------------------------------------

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

//-----------------------------------------------------------------------------
// scanner token items

// item represents a token returned from the scanner.
type item struct {
	typ itemType // Type, such as itemNumber.
	val string   // Value, such as "23.2".
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError      itemType = iota // error occurred; value is text of error
	itemDot                        // the cursor, spelled '.'
	itemEOF                        // end of file
	itemElse                       // else keyword
	itemEnd                        // end keyword
	itemField                      // identifier, starting with '.'
	itemIdentifier                 // identifier
	itemIf                         // if keyword
	itemLeftMeta                   // left meta-string
	itemNumber                     // number
	itemPipe                       // pipe symbol
	itemRange                      // range keyword
	itemRawString                  // raw quoted string (includes quotes)
	itemRightMeta                  // right meta-string
	itemString                     // quoted string (includes quotes)
	itemText                       // plain text
)

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

//-----------------------------------------------------------------------------

const leftMeta = "{{"
const rightMeta = "}}"
const eof = 0

//-----------------------------------------------------------------------------
// lexer and helper functions

// lexer holds the state of the scanner.
type lexer struct {
	name  string    // used only for error reports.
	input string    // the string being scanned.
	start int       // start position of this item.
	pos   int       // current position in the input.
	width int       // width of last rune read from input.
	items chan item // channel of scanned items.
}

// return a new lexer and its item channel
func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run() // Concurrently run state machine.
	return l, l.items
}

// run lexes the input by executing state functions until the state is nil.
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, n := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = n
	l.pos += l.width
	return r
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// accept consumes the next rune
// if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

//-----------------------------------------------------------------------------
// lexer state functions

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

func lexText(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], leftMeta) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeftMeta // Next state.
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF) // Useful to make EOF a token.
	return nil      // Stop the run loop.
}

func lexLeftMeta(l *lexer) stateFn {
	l.pos += len(leftMeta)
	l.emit(itemLeftMeta)
	return lexInsideAction // Now inside {{ }}.
}

func lexInsideAction(l *lexer) stateFn {
	// Either number, quoted string, or identifier.
	// Spaces separate and are ignored.
	// Pipe symbols separate and are emitted.
	for {
		if strings.HasPrefix(l.input[l.pos:], rightMeta) {
			return lexRightMeta
		}
		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed action")
		case isSpace(r):
			l.ignore()
		case r == '|':
			l.emit(itemPipe)
		case r == '"':
			return lexQuote
		case r == '`':
			return lexRawQuote
		case r == '+' || r == '-' || '0' <= r && r <= '9':
			l.backup()
			return lexNumber
		case isAlphaNumeric(r):
			l.backup()
			return lexIdentifier
		}
	}
}

func lexNumber(l *lexer) stateFn {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}

	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q",
			l.input[l.start:l.pos])
	}
	l.emit(itemNumber)
	return lexInsideAction
}

func lexRightMeta(l *lexer) stateFn {
	return nil
}
func lexQuote(l *lexer) stateFn {
	return nil
}

func lexRawQuote(l *lexer) stateFn {
	return nil
}

func lexIdentifier(l *lexer) stateFn {
	return nil
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

//-----------------------------------------------------------------------------
