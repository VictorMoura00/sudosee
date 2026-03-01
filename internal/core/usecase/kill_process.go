package usecase

import (
	"github.com/VictorMoura00/sudosee/internal/core/domain"
)

type KillProcessUseCase struct {
	repo domain.ProcessRepository
}

func NewKillProcessUseCase(repo domain.ProcessRepository) *KillProcessUseCase {
	return &KillProcessUseCase{repo: repo}
}

func (uc *KillProcessUseCase) Execute(pid int) error {
	return uc.repo.Terminate(pid)
}