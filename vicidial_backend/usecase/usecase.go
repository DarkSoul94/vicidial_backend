package usecase

import (
	"context"

	"github.com/DarkSoul94/vicidial_backend/vicidial_backend"
)

// Usecase ...
type Usecase struct {
	repo vicidial_backend.Repository
}

// NewUsecase ...
func NewUsecase(repo vicidial_backend.Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

// HelloWorld ...
func (u *Usecase) HelloWorld(c context.Context) {
	println("Hello")
}
