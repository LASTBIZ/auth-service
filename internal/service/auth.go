package service

import (
	"context"

	pb "auth-service/api/auth"
)

type AuthService struct {
	pb.UnimplementedAuthServer
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) CreateAuth(ctx context.Context, req *pb.CreateAuthRequest) (*pb.CreateAuthReply, error) {
	return &pb.CreateAuthReply{}, nil
}
func (s *AuthService) UpdateAuth(ctx context.Context, req *pb.UpdateAuthRequest) (*pb.UpdateAuthReply, error) {
	return &pb.UpdateAuthReply{}, nil
}
func (s *AuthService) DeleteAuth(ctx context.Context, req *pb.DeleteAuthRequest) (*pb.DeleteAuthReply, error) {
	return &pb.DeleteAuthReply{}, nil
}
func (s *AuthService) GetAuth(ctx context.Context, req *pb.GetAuthRequest) (*pb.GetAuthReply, error) {
	return &pb.GetAuthReply{}, nil
}
func (s *AuthService) ListAuth(ctx context.Context, req *pb.ListAuthRequest) (*pb.ListAuthReply, error) {
	return &pb.ListAuthReply{}, nil
}
