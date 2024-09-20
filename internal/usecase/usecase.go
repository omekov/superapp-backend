package usecase

import (
	"github.com/omekov/dubaicarkzv2/internal/usecase/external"
	"github.com/omekov/dubaicarkzv2/internal/usecase/repository"
)

type UseCase struct {
	repo     repository.Repo
	external external.Client
	mrp      int
}

func NewUseCase(repo repository.Repo) UseCase {
	return UseCase{
		repo: repo,
		mrp:  3692,
	}
}

func (u UseCase) GetKGDFile() ([]byte, error) {
	return u.external.DownloadFile()
}
