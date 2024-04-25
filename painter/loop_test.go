package Painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/shiny/screen"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) NewBuffer(image.Point) (screen.Buffer, error) {
	return nil, nil
}

func (m *Mock) NewWindow(*screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

func (m *Mock) Update(texture screen.Texture) {
	m.Called(texture)
}

func (m *Mock) NewTexture(size image.Point) (screen.Texture, error) {
	args := m.Called(size)
	return args.Get(0).(screen.Texture), args.Error(1)
}

func (m *Mock) Release() {
	m.Called()
}

func (m *Mock) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	m.Called(dp, src, sr)
}

func (m *Mock) Bounds() image.Rectangle {
	args := m.Called()
	return args.Get(0).(image.Rectangle)
}

func (m *Mock) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Called(dr, src, op)
}

func (m *Mock) Size() image.Point {
	args := m.Called()
	return args.Get(0).(image.Point)
}

func (m *Mock) Do(t screen.Texture) bool {
	args := m.Called(t)
	return args.Bool(0)
}

func TestLoop_Post_Successful(t *testing.T) {
	textureMock := new(Mock)
	receiverMock := new(Mock)
	screenMock := new(Mock)

	texture := image.Pt(400, 400)
	screenMock.On("NewTexture", texture).Return(textureMock, nil)
	receiverMock.On("Update", textureMock).Return()
	loop := Loop{
		Receiver: receiverMock,
	}

	loop.Start(screenMock)

	operationOne := new(Mock)
	textureMock.On("Bounds").Return(image.Rectangle{})
	operationOne.On("Do", textureMock).Return(true)

	assert.Empty(t, loop.mq.operations)
	loop.Post(operationOne)
	time.Sleep(1 * time.Second)
	assert.Empty(t, loop.mq.operations)

	operationOne.AssertCalled(t, "Do", textureMock)
	receiverMock.AssertCalled(t, "Update", textureMock)
	screenMock.AssertCalled(t, "NewTexture", image.Pt(400, 400))
}

func TestLoop_Post_Failed(t *testing.T) {
	textureMock := new(Mock)
	receiverMock := new(Mock)
	screenMock := new(Mock)

	texture := image.Pt(400, 400)
	screenMock.On("NewTexture", texture).Return(textureMock, nil)
	receiverMock.On("Update", textureMock).Return()
	loop := Loop{
		Receiver: receiverMock,
	}

	loop.Start(screenMock)

	operationOne := new(Mock)
	textureMock.On("Bounds").Return(image.Rectangle{})
	operationOne.On("Do", textureMock).Return(false)

	assert.Empty(t, loop.mq.operations)
	loop.Post(operationOne)
	time.Sleep(1 * time.Second)
	assert.Empty(t, loop.mq.operations)

	operationOne.AssertCalled(t, "Do", textureMock)
	receiverMock.AssertNotCalled(t, "Update", textureMock)
	screenMock.AssertCalled(t, "NewTexture", image.Pt(400, 400))
}

func TestLoop_Post_Multiple(t *testing.T) {
	textureMock := new(Mock)
	receiverMock := new(Mock)
	screenMock := new(Mock)

	texture := image.Pt(400, 400)
	screenMock.On("NewTexture", texture).Return(textureMock, nil)
	receiverMock.On("Update", textureMock).Return()
	loop := Loop{
		Receiver: receiverMock,
	}

	loop.Start(screenMock)

	operationOne := new(Mock)
	operationTwo := new(Mock)
	textureMock.On("Bounds").Return(image.Rectangle{})
	operationOne.On("Do", textureMock).Return(true)
	operationTwo.On("Do", textureMock).Return(true)

	assert.Empty(t, loop.mq.operations)
	loop.Post(operationOne)
	loop.Post(operationTwo)
	time.Sleep(1 * time.Second)
	assert.Empty(t, loop.mq.operations)

	operationOne.AssertCalled(t, "Do", textureMock)
	operationTwo.AssertCalled(t, "Do", textureMock)
	receiverMock.AssertCalled(t, "Update", textureMock)
	screenMock.AssertCalled(t, "NewTexture", image.Pt(400, 400))
}
