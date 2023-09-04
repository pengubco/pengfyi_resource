package test_double_in_go

//go:generate mockery --name Greeter --filename greet_mock.go --inpackage-suffix --inpackage
type Greeter interface {
	Greeting(a string) (string, error)
}

// Greet is the function we are going to test.
func Greet(g Greeter, a string) bool {
	_, err := g.Greeting(a)
	return err == nil
}
