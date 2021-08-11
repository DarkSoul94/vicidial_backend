package usecase

import (
	"context"
)

// Usecase ...
type Usecase struct {
}

// NewUsecase ...
func NewUsecase() *Usecase {
	return &Usecase{}
}

// HelloWorld ...
func (u *Usecase) HelloWorld(c context.Context) {
	println("Hello")
}
