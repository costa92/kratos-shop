package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"shop/app/user/service/internal/biz"
)

var _ biz.DemoRepo = (*DemoRepo)(nil)

type DemoRepo struct {
	data *Data
	log  *log.Helper
}

func NewDemoRepo(data *Data, logger log.Logger) biz.DemoRepo {
	return &DemoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/address")),
	}
}

func (d DemoRepo) GetDemo(ctx context.Context, id int64) (*biz.Demo, error) {
	return &biz.Demo{
		Id:   id,
		Name: "11",
	}, nil
}
