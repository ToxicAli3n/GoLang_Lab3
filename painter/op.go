package Painter

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/exp/shiny/screen"
)

// Operation визначає інтерфейс для операцій, які змінюють текстуру.
type Operation interface {
	Do(t screen.Texture) bool
}

// StateTweaker визначає інтерфейс для зміни стану малюнку.
type StateTweaker interface {
	SetState(sol *StatefulOperationList)
}

// StatefulOperationList утримує список операцій, які змінюють стан.
type StatefulOperationList struct {
	BgOperation      Operation
	BgRectOperation  Operation
	FigureOperations []Operation
}

// Do виконує всі операції в списку.
func (sol *StatefulOperationList) Do(t screen.Texture) bool {
	if sol.BgOperation != nil {
		sol.BgOperation.Do(t)
	} else {
		defaultFill(t, color.White)
	}
	if sol.BgRectOperation != nil {
		sol.BgRectOperation.Do(t)
	}
	for _, op := range sol.FigureOperations {
		op.Do(t)
	}
	return false
}

// Update змінює стан використовуючи StateTweaker.
func (sol *StatefulOperationList) Update(tweaker StateTweaker) {
	tweaker.SetState(sol)
}

// defaultFill зафарбовує текстуру у вказаний колір.
func defaultFill(t screen.Texture, c color.Color) {
	t.Fill(t.Bounds(), c, screen.Src)
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// OperationFill зафарбовує текстуру у вказаний колір.
type OperationFill struct {
	Color color.Color
}

func (op OperationFill) Do(t screen.Texture) bool {
	t.Fill(t.Bounds(), op.Color, screen.Src)
	return false
}

func (op OperationFill) SetState(sol *StatefulOperationList) {
	sol.BgOperation = op
}

// RelativePoint представляє точку відносно розмірів.
type RelativePoint struct {
	X, Y float64
}

func (p RelativePoint) ToAbs(size image.Point) image.Point {
	return image.Point{X: int(p.X * float64(size.X)), Y: int(p.Y * float64(size.Y))}
}

// OperationBGRect зафарбовує прямокутну область текстури.
type OperationBGRect struct {
	Min, Max RelativePoint
}

func (op OperationBGRect) Do(t screen.Texture) bool {
	rect := image.Rectangle{
		Min: op.Min.ToAbs(t.Size()),
		Max: op.Max.ToAbs(t.Size()),
	}
	t.Fill(rect, color.Black, draw.Src)
	return false
}

func (op OperationBGRect) SetState(sol *StatefulOperationList) {
	sol.BgRectOperation = op
}

// OperationFigure визначає операцію для фігури.
type OperationFigure struct {
	Center RelativePoint
}

func (op OperationFigure) Do(t screen.Texture) bool {
	centerAbs := op.Center.ToAbs(t.Size())
	drawT(t, centerAbs, 50, 40, color.RGBA{R: 0, G: 54, B: 206, A: 1})
	return false
}

func (op OperationFigure) SetState(sol *StatefulOperationList) {
	if sol.FigureOperations == nil {
		sol.FigureOperations = []Operation{}
	}
	sol.FigureOperations = append(sol.FigureOperations, op)
}

func drawT(t screen.Texture, center image.Point, hlen, hwidth int, c color.Color) {

	topHorizontal := image.Rect(center.X-hlen, center.Y-hwidth, center.X+hlen, center.Y)
	t.Fill(topHorizontal, c, draw.Src)

	// Нижній вертикальний прямокутник
	bottomVertical := image.Rect(center.X-hwidth/2, center.Y, center.X+hwidth/2, center.Y+hlen)
	t.Fill(bottomVertical, c, draw.Src)
}

// MoveTweaker зміщує фігури.
type MoveTweaker struct {
	Offset RelativePoint
}

func (tweaker MoveTweaker) SetState(sol *StatefulOperationList) {
	for i := range sol.FigureOperations {
		if fig, ok := sol.FigureOperations[i].(OperationFigure); ok {
			fig.Center.X = tweaker.Offset.X
			fig.Center.Y = tweaker.Offset.Y
			sol.FigureOperations[i] = fig
		}
	}
}

// ResetTweaker скидає стан до початкового.
type ResetTweaker struct{}

func (tweaker ResetTweaker) SetState(sol *StatefulOperationList) {
	blackFillOperation := OperationFill{Color: color.Black}
	sol.BgOperation = blackFillOperation
	sol.BgRectOperation = nil
	sol.FigureOperations = nil
}
