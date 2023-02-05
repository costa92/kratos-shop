package service

import (
	"context"
	"shop/app/user/service/internal/biz"

	pb "shop/api/user/service/v1"
)

type DemoService struct {
	pb.UnimplementedDemoServer
	dc *biz.DemoUseCase
}

func NewDemoService(demoUseCase *biz.DemoUseCase) *DemoService {
	return &DemoService{
		dc: demoUseCase,
	}
}

func (s *DemoService) CreateDemo(ctx context.Context, req *pb.CreateDemoRequest) (*pb.CreateDemoReply, error) {
	return &pb.CreateDemoReply{}, nil
}
func (s *DemoService) UpdateDemo(ctx context.Context, req *pb.UpdateDemoRequest) (*pb.UpdateDemoReply, error) {
	return &pb.UpdateDemoReply{}, nil
}
func (s *DemoService) DeleteDemo(ctx context.Context, req *pb.DeleteDemoRequest) (*pb.DeleteDemoReply, error) {
	return &pb.DeleteDemoReply{}, nil
}
func (s *DemoService) GetDemo(ctx context.Context, req *pb.GetDemoRequest) (*pb.GetDemoReply, error) {
	res, err := s.dc.GetDemo(ctx, 112)
	if err != nil {
		return nil, err
	}
	return &pb.GetDemoReply{
		Id:       res.Id,
		Name:     res.Name,
		UserNick: "111",
	}, nil
}
func (s *DemoService) ListDemo(ctx context.Context, req *pb.ListDemoRequest) (*pb.ListDemoReply, error) {
	return &pb.ListDemoReply{}, nil
}
