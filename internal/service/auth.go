package service

import (
	"auth-service/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"

	pb "auth-service/api/auth"
)

type AuthService struct {
	pb.UnimplementedAuthServer

	ua  *biz.AuthUseCase
	log *log.Helper
}

func NewAuthService(ua *biz.AuthUseCase, logger log.Logger) *AuthService {
	return &AuthService{ua: ua, log: log.NewHelper(logger)}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*empty.Empty, error) {
	_, err := s.ua.Register(ctx, strings.TrimSpace(req.Email), strings.TrimSpace(req.FirstName), strings.TrimSpace(req.LastName), strings.TrimSpace(req.Password))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthService) ProviderAuth(ctx context.Context, req *pb.ProviderAuthRequest) (*pb.ProviderAuthResponse, error) {
	redirect, err := s.ua.CreateState(req.Provider)
	if err != nil {
		return nil, err
	}
	return &pb.ProviderAuthResponse{
		Redirect: redirect,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.TokenResponse, error) {
	tok, err := s.ua.Login(ctx, strings.TrimSpace(req.Email), strings.TrimSpace(req.Password))
	if err != nil {
		return nil, err
	}
	return &pb.TokenResponse{AccessToken: tok.AccessToken, RefreshToken: tok.RefreshToken}, nil
}

func (s *AuthService) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateReply, error) {
	id, err := s.ua.Validate(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &pb.ValidateReply{
		UserId: id,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.TokenResponse, error) {
	tok, err := s.ua.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &pb.TokenResponse{AccessToken: tok.AccessToken, RefreshToken: tok.RefreshToken}, nil
}

func (s *AuthService) Callback(ctx context.Context, req *pb.CallbackRequest) (*pb.TokenResponse, error) {
	tok, err := s.ua.Callback(ctx, req.Provider, req.OauthCode, req.State)
	if err != nil {
		return nil, err
	}

	return &pb.TokenResponse{AccessToken: tok.AccessToken, RefreshToken: tok.RefreshToken}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*empty.Empty, error) {
	_, err := s.ua.ChangePassword(ctx, req.UserId, req.Password, strings.TrimSpace(req.NewPassword))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
