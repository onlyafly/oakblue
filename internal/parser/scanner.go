/*
See "Lexical Scanning in Go" by Rob Pike for the basic theory behind this
module: http://www.youtube.com/watch?v=HxaD_trXwRE
*/

// TODO: Remove the non-ASM code from the scanner (there's still some left from Lisp related stuff)

package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/onlyafly/oakblue/internal/syntax"
)

////////// Token

type Token struct {
	Location *syntax.Location
	Code     TokenCode
	Value    string
}

func (t Token) String() string {
	switch t.Code {
	case TcEOF:
		return "EOF"
	case TcNewline:
		return "NEWLINE"
	case TcError:
		return t.Value
	}

	return fmt.Sprintf("%v", t.Value)
}

func (t Token) Loc() *syntax.Location {
	return t.Location
}

////////// Token code

type TokenCode int

// TODO remove the unused token codes
const (
	TcError TokenCode = iota // important that TcError has value 0

	TcCaret
	TcChar
	TcColon
	TcDecimalNumber
	TcHexNumber
	TcLeftParen
	TcNewline
	TcRightParen
	TcSingleQuote
	TcString
	TcSymbol
	TcRegister

	TcEOF
)

const eof = -1

type ErrorHandler func(t Token, message string)

////////// Scanner struct

type Scanner struct {
	name   string     // used only for error reports
	input  string     // the string being scanned
	start  int        // start position of this item
	pos    int        // current position in the input
	line   int        // current line number in the input
	width  int        // width of last rune read from input
	Tokens chan Token // channel of scanned items

	// Error handling
	errorCount   int
	errorHandler ErrorHandler
}

func (s *Scanner) String() string {
	return fmt.Sprintf("<scanner next=%#v>", s.input[s.start:s.pos])
}

func (s *Scanner) run() {
	for state := scanBegin; state != nil; {
		state = state(s)
	}
	close(s.Tokens)
}

func (s *Scanner) emit(code TokenCode) {
	s.Tokens <- Token{
		Location: &syntax.Location{Pos: s.start, Line: s.line, Filename: s.name},
		Code:     code,
		Value:    s.input[s.start:s.pos],
	}
	s.start = s.pos
}

func (s *Scanner) next() (r rune) {
	if s.pos >= len(s.input) {
		s.width = 0
		r = eof
		return
	}
	r, s.width = utf8.DecodeRuneInString(s.input[s.pos:])
	s.pos += s.width
	return
}

// ignore skips over the pending input before this point.
func (s *Scanner) ignore() {
	s.start = s.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (s *Scanner) backup() {
	s.pos -= s.width
}

// peek returns but does not consume
// the next rune in the input.
func (s *Scanner) peek() rune {
	r := s.next()
	s.backup()
	return r
}

// accept consumes the next rune
// if it's from the valid set.
func (s *Scanner) accept(valid string) bool {
	if strings.ContainsRune(valid, s.next()) {
		return true
	}
	s.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (s *Scanner) acceptRun(valid string) {
	for strings.ContainsRune(valid, s.next()) {
	}
	s.backup()
}

func (s *Scanner) emitErrorf(format string, args ...interface{}) {
	t := Token{
		Location: &syntax.Location{Pos: s.start, Line: s.line, Filename: s.name},
		Code:     TcError,
		Value:    s.input[s.start:s.pos],
	}

	if s.errorHandler != nil {
		message := fmt.Sprintf(format, args...)
		s.errorHandler(t, message)
	}

	s.Tokens <- t
	s.start = s.pos

	s.errorCount++
}

////////// Scanning

type stateFn func(*Scanner) stateFn

func Scan(name, input string) (*Scanner, chan Token) {
	s := &Scanner{
		name:   name,
		input:  input,
		line:   1,
		Tokens: make(chan Token),
	}
	go s.run()
	return s, s.Tokens
}

func scanBegin(s *Scanner) stateFn {
Outer:
	for {
		switch r := s.next(); {
		case isWhiteSpace(r):
			s.ignore()
		case isNewLine(r):
			s.line++
			s.emit(TcNewline)
		case r == '(':
			s.emit(TcLeftParen)
		case r == ')':
			s.emit(TcRightParen)
		case r == 'x':
			return scanHexNumber
		case r == '0':
			rnext := s.next()
			if rnext == 'x' {
				s.backup()
				s.backup()
				return scanHexNumber
			}
			s.backup()
			return scanDecimalNumber
		case ('0' <= r && r <= '9') || r == '+' || r == '-' || r == '#':
			s.backup()
			return scanDecimalNumber
		case r == '^':
			s.emit(TcCaret)
		case r == '\'':
			s.emit(TcSingleQuote)
		case r == '\\':
			return scanChar
		case r == ';':
			return scanSingleLineComment
		case r == ':':
			s.emit(TcColon)
		case r == 'r' || r == 'R':
			rnext := s.peek()
			if '0' <= rnext && rnext <= '9' {
				s.backup()
				return scanRegister
			}
			return scanSymbol
		case isSymbolic(r):
			s.backup()
			return scanSymbol
		case r == '"':
			return scanString
		/* TODO change multiline comment to something else
		case r == '#':
			rnext := s.next()
			if rnext == '|' {
				return scanMultiLineComment
			}
			s.emitErrorf("unrecognized character sequence: '%c%c' = %v,%v", r, rnext, r, rnext)
		*/
		case r == eof:
			break Outer
		default:
			s.emitErrorf("unrecognized character: '%c' = %v", r, r)
		}
	}

	s.emit(TcEOF)
	return nil
}

func scanSingleLineComment(s *Scanner) stateFn {
	for isSingleLineCommentContent(s.next()) {
	}
	s.backup()
	s.ignore()
	return scanBegin
}

/* TODO
func scanMultiLineComment(s *Scanner) stateFn {
	r := s.next()
	for r != eof {
		switch {
		case isNewLine(r):
			s.line++
		case r == '|':
			rnext := s.next()
			if rnext == '#' {
				s.ignore()
				return scanBegin
			}
			s.backup() // in case there is a '|' following the '|'
		}
		r = s.next()
	}

	s.emitErrorf("non-terminated multiline comment")
	return scanBegin
}
*/

func scanString(s *Scanner) stateFn {
	for isStringContent(s.next()) {
	}
	s.emit(TcString)
	return scanBegin
}

func scanChar(s *Scanner) stateFn {
	s.next() // first character in literal
	if isAlpha(s.peek()) {
		for isAlpha(s.next()) {
		}
		s.backup()
	}
	s.emit(TcChar)
	return scanBegin
}

func scanSymbol(s *Scanner) stateFn {
	for isSymbolic(s.next()) {
	}
	s.backup()
	s.emit(TcSymbol)
	return scanBegin
}

func scanRegister(s *Scanner) stateFn {
	for isSymbolic(s.next()) {
	}
	s.backup()
	s.emit(TcRegister)
	return scanBegin
}

func scanHexNumber(s *Scanner) stateFn {

	s.accept("0")

	s.accept("xX")

	digits := "0123456789abcdefABCDEF"
	s.acceptRun(digits)

	// Next thing must not be alphanumeric
	if isAlphaNumeric(s.peek()) {
		s.next()
		s.emitErrorf("bad hex number syntax: %q", s.input[s.start:s.pos])
	} else {
		s.emit(TcHexNumber)
	}

	return scanBegin
}

func scanDecimalNumber(s *Scanner) stateFn {

	// Accept leading decimal flag
	s.accept("#")

	// Optional leading sign
	s.accept("+-")

	digits := "0123456789"

	s.acceptRun(digits)
	if s.accept(".") {
		s.acceptRun(digits)
	}
	if s.accept("eE") {
		s.accept("+-")
		s.acceptRun("0123456789")
	}

	// Is it imaginary?
	s.accept("i")

	// Next thing must not be alphanumeric
	if isAlphaNumeric(s.peek()) {
		s.next()
		s.emitErrorf("bad decimal number syntax: %q", s.input[s.start:s.pos])
	} else {
		s.emit(TcDecimalNumber)
	}

	return scanBegin
}

////////// Helpers

func isAlphaNumeric(r rune) bool {
	switch {
	case '0' <= r && r <= '9':
		return true
	case 'a' <= r && r <= 'z':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	}

	return false
}

func isAlpha(r rune) bool {
	switch {
	case 'a' <= r && r <= 'z':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	}

	return false
}

func isSymbolic(r rune) bool {
	switch {
	case '0' <= r && r <= '9':
		return true
	case 'a' <= r && r <= 'z':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	// NOTE: Don't ever allow these characters: [ ] { } ( ) " , ' ` : ; # | \ ~
	case r == '?' ||
		r == '+' || r == '-' || r == '*' || r == '/' ||
		r == '=' || r == '<' || r == '>' || r == '!' ||
		r == '&' || r == '_' || r == '.':
		return true
	}

	return false
}

func isWhiteSpace(r rune) bool {
	switch r {
	case ' ':
		return true
	case '\t':
		return true
	}

	return false
}

func isNewLine(r rune) bool {
	switch r {
	case '\r':
		return true
	case '\n':
		return true
	}

	return false
}

func isSingleLineCommentContent(r rune) bool {
	switch r {
	case '\r':
		return false
	case '\n':
		return false
	case eof:
		return false
	}

	return true
}

func isStringContent(r rune) bool {
	switch r {
	case '"':
		return false
	case eof:
		return false
	}

	return true
}
