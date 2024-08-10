package types

type Type interface {
	Simple() bool
	Constant() bool
}
