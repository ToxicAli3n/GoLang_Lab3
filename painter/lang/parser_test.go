package Lang

import (
	Painter "github.com/roman-mazur/architecture-lab-3/painter"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	p := &Parser{}

	ops, err := p.Parse(strings.NewReader("white"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(ops))
	}

	ops, err = p.Parse(strings.NewReader("green"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(ops))
	}

	ops, err = p.Parse(strings.NewReader("update"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(ops))
	}

	ops, err = p.Parse(strings.NewReader("bgrect 0.5 0.5 -0.5 -0.5"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(ops))
	}

	ops, err = p.Parse(strings.NewReader("figure 0.5 0.5"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(ops))
	}

	ops, err = p.Parse(strings.NewReader("move 0.5 0.5"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(ops))
	}

	ops, err = p.Parse(strings.NewReader("reset"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(ops))
	}

	ops, err = p.Parse(strings.NewReader("unknown"))
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if len(ops) != 0 {
		t.Errorf("Expected 0 operations, got %d", len(ops))
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
