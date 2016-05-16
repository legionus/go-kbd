//go:generate -command yacc go tool yacc
//go:generate yacc -o expr.go -p "parser" parser.y

package parser
