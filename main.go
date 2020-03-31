package main

import (
	"errors"
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/hemu/antlr-calc/parser"
	"strconv"
)

type calcListener struct {
	// embedded type
	*parser.BaseCalcListener

	stack []int
}

func (cl *calcListener) Pop() (int, error) {
	if len(cl.stack) < 1 {
		return 0, errors.New("no elements to pop")
	}
	num := cl.stack[len(cl.stack)-1]
	cl.stack = cl.stack[:len(cl.stack)-1]
	return num, nil
}

func (cl *calcListener) getLastTwoOperands() (int, int, error) {
	num1, err := cl.Pop()
	if err != nil {
		return 0, 0, err
	}
	num2, err := cl.Pop()
	if err != nil {
		return 0, 0, err
	}
	return num1, num2, nil
}

func (cl *calcListener) ExitMulDiv(c *parser.MulDivContext) {
	ttype := c.GetOp().GetTokenType()

	switch ttype {
	case parser.CalcParserMUL:
		num1, num2, err := cl.getLastTwoOperands()
		if err != nil {
			return
		}
		cl.stack = append(cl.stack, num1*num2)
	case parser.CalcParserDIV:
		num1, num2, err := cl.getLastTwoOperands()
		if err != nil {
			return
		}
		cl.stack = append(cl.stack, num1/num2)
	}
}

func (cl *calcListener) ExitAddSub(c *parser.AddSubContext) {
	ttype := c.GetOp().GetTokenType()

	switch ttype {
	case parser.CalcParserADD:
		num1, num2, err := cl.getLastTwoOperands()
		if err != nil {
			return
		}
		cl.stack = append(cl.stack, num1+num2)
	case parser.CalcParserSUB:
		num1, num2, err := cl.getLastTwoOperands()
		if err != nil {
			return
		}
		cl.stack = append(cl.stack, num1-num2)
	}
}

func (cl *calcListener) ExitNumber(c *parser.NumberContext) {
	num, err := strconv.Atoi(c.GetText())
	if err != nil {
		return
	}
	cl.stack = append(cl.stack, num)
}

func main() {
	inpStream := antlr.NewInputStream("1 + 2 * 3 + 20")
	lexer := parser.NewCalcLexer(inpStream)
	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewCalcParser(tokenStream)

	listener := &calcListener{}
	antlr.ParseTreeWalkerDefault.Walk(listener, p.Start())

	fmt.Println(listener.Pop())
}
