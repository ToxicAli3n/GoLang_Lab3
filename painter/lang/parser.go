package Lang

import (
	"bufio"
	"fmt"
	"github.com/roman-mazur/architecture-lab-3/painter"
	"image/color"
	"io"
	"strconv"
	"strings"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
	// Зберігає стан малюнку у спеціальній операції.
	state Painter.StatefulOperationList
}

func (p *Parser) Parse(in io.Reader) ([]Painter.Operation, error) {
	var res []Painter.Operation

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		commandLine := scanner.Text()
		op, err := p.process(commandLine)

		if err != nil {
			return nil, err
		} else if op != nil {
			res = append(res, op)
		}
	}

	return res, nil
}

type countError struct{}

func (e countError) Error() string {
	return "Invalid argument count"
}

func (p *Parser) process(cmd string) (Painter.Operation, error) {
	var tweaker Painter.StateTweaker

	fields := strings.Fields(cmd)
	switch fields[0] {
	case "white":
		if len(fields) > 1 {
			return nil, countError{}
		}
		tweaker = Painter.OperationFill{Color: color.White}
	case "green":
		if len(fields) > 1 {
			return nil, countError{}
		}
		tweaker = Painter.OperationFill{Color: color.RGBA{G: 0xff, A: 0xff}}
	case "update":
		if len(fields) > 1 {
			return nil, countError{}
		}
		return Painter.UpdateOp, nil
	case "bgrect":
		args, err := processArguments(fields[1:], 4)
		if err != nil {
			return nil, err
		}
		tweaker = Painter.OperationBGRect{
			Min: Painter.RelativePoint{X: args[0], Y: args[1]},
			Max: Painter.RelativePoint{X: args[2], Y: args[3]},
		}
	case "figure":
		args, err := processArguments(fields[1:], 2)
		if err != nil {
			return nil, err
		}
		tweaker = Painter.OperationFigure{
			Center: Painter.RelativePoint{X: args[0], Y: args[1]},
		}
	case "move":
		args, err := processArguments(fields[1:], 2)
		if err != nil {
			return nil, err
		}
		tweaker = Painter.MoveTweaker{
			Offset: Painter.RelativePoint{X: args[0], Y: args[1]},
		}
	case "reset":
		if len(fields) > 1 {
			return nil, countError{}
		}
		tweaker = Painter.ResetTweaker{}
	default:
		return nil, fmt.Errorf("unknown command")
	}

	if nil != tweaker {
		p.state.Update(tweaker)
	}
	return &p.state, nil
}

func processArguments(args []string, requiredLen int) ([]float64, error) {
	if len(args) != requiredLen {
		return nil, countError{}
	}
	var processed []float64
	for idx, arg := range args {
		num, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid argument at pos %d", idx)
		}
		if num >= -1 && num <= 1 {
			processed = append(processed, num)
		} else {
			return nil, fmt.Errorf("value at pos %d is not in [-1,1] range", idx)
		}
	}

	return processed, nil
}
