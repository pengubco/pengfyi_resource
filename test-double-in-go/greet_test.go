package test_double_in_go

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// write stub manually
type GreetingStub struct {
	GreetingFunc func(string) (string, error)
}

func (g GreetingStub) Greeting(a string) (string, error) {
	if g.GreetingFunc != nil {
		return g.GreetingFunc(a)
	}
	return "", errors.New("not implemented")
}

func Test_Greeting_UsingManualStub(t *testing.T) {
	g := GreetingStub{GreetingFunc: func(s string) (string, error) {
		switch s {
		case "hello":
			return "world", nil
		case "HELLO":
			return "", errors.New("uppercase nota allowed")
		default:
			return "", nil
		}
	}}
	assert.True(t, Greet(g, "hello"))
	assert.False(t, Greet(g, "HELLO"))
}

// use a library, testify mock.
type GreetingMock struct {
	mock.Mock
}

func (m *GreetingMock) Greeting(a string) (string, error) {
	args := m.Called(a)
	return args.String(0), args.Error(1)
}

func Test_Greet_UsingTestifyMock(t *testing.T) {
	m := new(GreetingMock)
	m.On("Greeting", "hello").Return("world", nil)
	m.On("Greeting", "HELLO").Return("", errors.New("uppercase nota allowed"))
	assert.True(t, Greet(m, "hello"))
	assert.False(t, Greet(m, "HELLO"))
}

// use mockery to generate testify mock automatically
func Test_Greet_Using_Mockery(t *testing.T) {
	m := NewMockGreeter(t)
	m.On("Greeting", "hello").Return("world", nil)
	m.On("Greeting", "HELLO").Return("", errors.New("uppercase nota allowed"))
	assert.True(t, Greet(m, "hello"))
	assert.False(t, Greet(m, "HELLO"))
}
