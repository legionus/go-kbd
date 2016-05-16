package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
)

type State int

const (
	StateNormal State = iota
	StateValue
	StateInclude
)

const MaxString int = 512

var AltIsMetaRe = regexp.MustCompile("^[aA][lL][tT][-_][iI][sS][-_][mM][eE][tT][aA]$")
var HexRe = regexp.MustCompile("^0[xX][0-9a-fA-F]+$")
var UnicodeRe = regexp.MustCompile("^U[+]([0-9a-fA-F]){4}$")
var LiteralRe = regexp.MustCompile("^[a-zA-Z][a-zA-Z_0-9]*$")

func isOcta(c byte) bool {
	return '0' <= c && c <= '7'
}

func isHex(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

func NewLexer(rd io.Reader) *Lexer {
	return &Lexer{
		input: bufio.NewReader(rd),
	}
}

type Lexer struct {
	state State
	input *bufio.Reader
	buf   []byte
}

func (l *Lexer) setState(s State) {
	//fmt.Printf("State = %d\n", s)
	l.state = s
}

func (l *Lexer) read() error {
	c, err := l.input.ReadByte()
	if err != nil {
		return err
	}
	l.buf = append(l.buf, c)
	return nil
}

func (l *Lexer) unread() error {
	err := l.input.UnreadByte()
	if err != nil {
		return err
	}
	l.buf = l.buf[:len(l.buf)-1]
	return nil
}

func (l *Lexer) consume(kind Kind) Leaf {
	leaf := Leaf{
		Kind: kind,
		Data: l.buf,
	}
	l.buf = []byte{}
	return leaf
}

func (l *Lexer) consumeFunc(f func(c byte) bool, kind Kind) Leaf {
	if !f(l.buf[len(l.buf)-1]) {
		panic("nothing consumed")
	}

	for {
		if err := l.read(); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if !f(l.buf[len(l.buf)-1]) {
			l.unread()
			break
		}
	}
	return l.consume(kind)
}

func (l *Lexer) consumeWhile(b []byte, kind Kind) Leaf {
	return l.consumeFunc(func(c byte) bool {
		return bytes.IndexByte(b, c) != -1
	}, kind)
}

func (l *Lexer) consumeUntil(b []byte, kind Kind) Leaf {
	return l.consumeFunc(func(c byte) bool {
		return bytes.IndexByte(b, c) == -1
	}, kind)
}

func (l *Lexer) getString() (Node, error) {
	if l.buf[0] != '"' {
		panic("expected \"")
	}
Loop:
	for {
		if err := l.read(); err != nil {
			panic(err)
		}
		switch l.buf[len(l.buf)-1] {
		case '\\':
			if err := l.read(); err != nil {
				panic(err)
			}
			// TODO(legion) check Octa
			if l.buf[len(l.buf)-1] == 'n' {
				l.buf[len(l.buf)-2] = '\n'
			} else {
				l.buf[len(l.buf)-2] = l.buf[len(l.buf)-1]
			}
			l.buf = l.buf[:len(l.buf)-1]
		case '"':
			break Loop
		}
	}
	l.buf = l.buf[1 : len(l.buf)-1]

	if len(l.buf) >= MaxString {
		panic("string too long: " + string(l.buf))
	}
	return l.consume(String), nil
}

func (l *Lexer) Get() (Node, error) {
	for {
		if err := l.read(); err != nil {
			if err == io.EOF {
				if len(l.buf) != 0 {
					break
				}
				return nil, nil
			}
			return nil, err
		}

		if len(l.buf) > 1 {
			if l.buf[len(l.buf)-1] == '\n' && l.buf[len(l.buf)-2] == '\\' {
				l.buf = l.buf[:len(l.buf)-2]
				continue
			}

			if data := UnicodeRe.Find(l.buf); data != nil {
				return l.consume(Unicode), nil
			}

			if data := HexRe.Find(l.buf); data != nil {
				more, err := l.input.Peek(1)
				if err != nil {
					panic(err)
				}
				if isHex(more[0]) {
					continue
				}
				return l.consume(Number), nil
			}

			if l.state == StateValue {
				if bytes.IndexByte([]byte(" \t\r\n"), l.buf[len(l.buf)-1]) != -1 {
					l.unread()
					if data := LiteralRe.Find(l.buf); data != nil {
						return l.consume(Literal), nil
					}
					panic("unexpected spacing: <" + string(l.buf) + ">")
				}
				continue
			}

			switch string(l.buf) {
			case "include":
				l.setState(StateInclude)
				return l.consume(Include), nil
			case "alt", "Alt", "ALT":
				more, err := l.input.Peek(1)
				if err != nil {
					panic(err)
				}
				if bytes.IndexByte([]byte("gG"), more[0]) != -1 {
					break
				}
				return l.consume(Alt), nil
			case "altgr", "Altgr", "AltGr", "ALTGR":
				return l.consume(AltGr), nil
			case "string", "String", "STRING":
				more, err := l.input.Peek(1)
				if err != nil {
					panic(err)
				}
				if bytes.IndexByte([]byte("sS"), more[0]) != -1 {
					break
				}
				l.setState(StateValue)
				return l.consume(String), nil
			case "strings", "Strings", "STRINGS":
				return l.consume(Strings), nil
			case "shift", "Shift", "SHIFT":
				more, err := l.input.Peek(1)
				if err != nil {
					panic(err)
				}
				if bytes.IndexByte([]byte("rRlL"), more[0]) != -1 {
					break
				}
				return l.consume(Shift), nil
			case "shiftl", "ShiftL", "SHIFTL":
				return l.consume(ShiftL), nil
			case "shiftr", "ShiftR", "SHIFTR":
				return l.consume(ShiftR), nil
			case "keycode", "Keycode", "KeyCode", "KEYCODE":
				return l.consume(Keycode), nil
			case "charset", "Charset", "CharSet", "CHARSET":
				return l.consume(Charset), nil
			case "keymaps", "Keymaps", "KeyMaps", "KEYMAPS":
				return l.consume(Keymaps), nil
			case "plain", "Plain", "PLAIN":
				return l.consume(Plain), nil
			case "control", "Control", "CONTROL":
				if l.state == StateValue {
					break
				}
				return l.consume(Control), nil
			case "ctrll", "CtrlL", "CTRLL":
				return l.consume(CtrlL), nil
			case "ctrlr", "CtrlR", "CTRLR":
				return l.consume(CtrlR), nil
			case "capsshift", "Capsshift", "CapsShift", "CAPSSHIFT":
				return l.consume(CapsShift), nil
			case "compose", "Compose", "COMPOSE":
				return l.consume(Compose), nil
			case "usual", "Usual", "USUAL":
				return l.consume(Usual), nil
			case "for", "For", "FOR":
				return l.consume(For), nil
			case "as", "As", "AS":
				return l.consume(As), nil
			case "on", "On", "ON":
				return l.consume(On), nil
			case "to", "To", "TO":
				l.setState(StateValue)
				return l.consume(To), nil
			}

			if data := AltIsMetaRe.Find(l.buf); data != nil {
				return l.consume(AltIsMeta), nil
			}

			continue
		}

		switch l.buf[0] {
		case ' ', '\t', '\r':
			return l.consumeWhile([]byte(" \t\r"), Space), nil
		case '#', '!':
			return l.consumeUntil([]byte("\n"), Comment), nil
		case '-':
			return l.consume(Dash), nil
		case '+':
			return l.consume(Plus), nil
		case ',':
			return l.consume(Comma), nil
		case '\n':
			l.setState(StateNormal)
			return l.consume(NewLine), nil
		case '=':
			l.setState(StateValue)
			return l.consume(Equals), nil
		case '"':
			l.setState(StateNormal)
			return l.getString()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if l.buf[0] == '0' {
				more, err := l.input.Peek(1)
				if err != nil {
					panic(err)
				}
				if bytes.IndexByte([]byte("xX"), more[0]) != -1 {
					continue
				}
			}
			return l.consumeWhile([]byte("0123456789"), Number), nil
		case '\\':
			if err := l.read(); err != nil {
				panic(err)
			}

			l.buf[0] = l.buf[1]
			l.buf = l.buf[:1]

			i := 0
			for isOcta(l.buf[i]) && i < 3 {
				if err := l.read(); err != nil {
					panic(err)
				}
				i += 1
			}
			if i != 0 {
				l.unread()
			}
			return l.consume(Char), nil
		}
	}
	return nil, fmt.Errorf("ERROR: state=%d: %q", l.state, string(l.buf))
}
