package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type Demo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type DemoRepo interface {
	GetDemo(ctx context.Context, id int64) (*Demo, error)
}

type DemoUseCase struct {
	repo DemoRepo
	log  *log.Helper
}

func NewDemoUseCase(repo DemoRepo, logger log.Logger) *DemoUseCase {
	return &DemoUseCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "usecase/address"))}
}

func (d *DemoUseCase) GetDemo(ctx context.Context, id int64) (*Demo, error) {
	return d.repo.GetDemo(ctx, id)
}
