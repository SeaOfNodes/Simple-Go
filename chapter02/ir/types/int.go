package types

type IntType struct {
	Value int
}

func NewIntType(value int) Type {
	return &IntType{Value: value}
}

func (in *IntType) Simple() bool   { return false }
func (in *IntType) Constant() bool { return true }
