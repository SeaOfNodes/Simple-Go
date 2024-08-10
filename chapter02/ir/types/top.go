package types

type topType struct{}

func (t *topType) Simple() bool   { return true }
func (t *topType) Constant() bool { return false }

var TopType = &topType{}
