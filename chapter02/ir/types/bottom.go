package types

type bottomType struct{}

func (b *bottomType) Simple() bool   { return true }
func (b *bottomType) Constant() bool { return false }

var BottomType = &bottomType{}
