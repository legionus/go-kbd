// To build it:
// go tool yacc -p "parser" parser.y (produces y.go)
// go build -o expr y.go

%{

package parser

%}

%token EOL NUMBER LITERAL CHARSET KEYMAPS KEYCODE EQUALS

%%

top:

%%

func Foo() bool {
	return true
}
