package Lang

import (
	Painter "github.com/roman-mazur/architecture-lab-3/painter"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Input(t *testing.T) {
	type testCase struct {
		name string
		cmd  string
		err  string
	}

	var testTable = []testCase{
		{
			name: "duplicate command",
			cmd:  "white white",
			err:  "Invalid argument count",
		},
		{
			name: "invalid argument amount",
			cmd:  "reset white",
			err:  "Invalid argument count",
		},
		{
			name: "unknown command",
			cmd:  "hello",
			err:  "unknown command",
		},
		{
			name: "multiple command: unknown command",
			cmd:  "green\n bgrect 0.1 0.1 0.1 0.1\n hello",
			err:  "unknown command",
		},
		{
			name: "out of the range",
			cmd:  "move 3 3",
			err:  "value at pos 0 is not in [-1,1] range",
		},
		{
			name: "out of the range",
			cmd:  "bgrect 0.3 -8 0.5 0.3",
			err:  "value at pos 1 is not in [-1,1] range",
		},
		{
			name: "invalid argument",
			cmd:  "figure j -0.9",
			err:  "invalid argument at pos 0",
		},
		{
			name: "invalid amount of argument",
			cmd:  "bgrect 0.3 a",
			err:  "Invalid argument count",
		},
		{
			name: "move: invalid amount of argument",
			cmd:  "move 0.3 3 3 3",
			err:  "Invalid argument count",
		},
		{
			name: "miltiple cmd: duplicate cmd",
			cmd:  "move 0.3 0.1\n white white",
			err:  "Invalid argument count",
		},
		{
			name: "figure: out of range",
			cmd:  "move 0.3 0.1\n white\n figure 1.2 0.1",
			err:  "value at pos 0 is not in [-1,1] range",
		},
		{
			name: "invalid argument",
			cmd:  "reset\n figure j -0.9\n green\n update",
			err:  "invalid argument at pos 0",
		},
	}

	for _, test := range testTable {
		p := &Parser{}
		_, err := p.Parse(strings.NewReader(test.cmd))

		assert.Equal(t, test.err, err.Error(), test.name)
	}
}

func TestParser_Operation(t *testing.T) {
	type testCase struct {
		name     string
		cmd      string
		ops      []Painter.Operation
		figures  []*Painter.OperationFigure
		checkIdx int
	}

	testTable := []testCase{
		{
			name: "Simple consequent commands with state and without",
			cmd:  "white\ngreen\nupdate",
			ops: []Painter.Operation{
				&Painter.StatefulOperationList{},
				&Painter.StatefulOperationList{},
				Painter.UpdateOp,
			},
			checkIdx: -1,
		},
		{
			name: "Figures are parsed with correct position",
			cmd:  "white\nfigure 0.5 0.5\nfigure -0.1 0.933",
			ops: []Painter.Operation{
				&Painter.StatefulOperationList{},
				&Painter.StatefulOperationList{},
				&Painter.StatefulOperationList{},
			},
			figures: []*Painter.OperationFigure{

				{Center: Painter.RelativePoint{X: 0.5, Y: 0.5}},
				{Center: Painter.RelativePoint{X: -0.1, Y: 0.933}},
			},
			checkIdx: 2,
		},
		{
			name: "Figures have correct position after multiple move operations",
			cmd:  "figure 0.5 0.5\nfigure 0.4 0.35\nmove 0.2 0.2\nmove -0.1 0.1\nupdate",
			ops: []Painter.Operation{
				&Painter.StatefulOperationList{},
				&Painter.StatefulOperationList{},
				&Painter.StatefulOperationList{},
				&Painter.StatefulOperationList{},
				Painter.UpdateOp,
			},
			figures: []*Painter.OperationFigure{
				{Center: Painter.RelativePoint{X: 0.6, Y: 0.8}},
				{Center: Painter.RelativePoint{X: 0.5, Y: 0.65}},
			},
			checkIdx: 2,
		},
	}
	delta := 0.99

	for _, test := range testTable {
		p := &Parser{}
		res, err := p.Parse(strings.NewReader(test.cmd))
		assert.Nil(t, err)
		assert.Equal(t, len(test.ops), len(res))
		for idx, op := range res {
			assert.IsType(t, test.ops[idx], op)
		}
		if test.checkIdx != -1 {
			st, ok := res[test.checkIdx].(*Painter.StatefulOperationList)
			if !ok {
				panic("Test case is incorrect")
			}
			assert.Equal(t, len(test.figures), len(st.FigureOperations))
			for idx, figure := range test.figures {
				stFigure := st.FigureOperations[idx]
				assert.InDelta(t, figure.Center.X, stFigure.Center.X, delta)
				assert.InDelta(t, figure.Center.Y, stFigure.Center.Y, delta)
			}
		}
	}
}
