package repl

import (
	"bicep-go/lexer"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		l.Lex()

		for _, tok := range l.GetTokens() {
			fmt.Printf("%+v \n", tok.ToString())
		}
	}
}
