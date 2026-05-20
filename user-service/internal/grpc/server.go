package grpc

import (
	"context"

	userpb "user-service/gen/userpb"
	"user-service/internal/model"
	"user-service/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct {
	userpb.UnimplementedUserServiceServer
	service *service.UserService
}

func NewUserGRPCServer(service *service.UserService) *UserGRPCServer {
	return &UserGRPCServer{service: service}
}

func (s *UserGRPCServer) RegisterUser(ctx context.Context, req *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	user, err := s.service.Register(req.GetName(), req.GetEmail(), req.GetPassword(), req.GetRole())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return toUserResponse(user), nil
}

func (s *UserGRPCServer) LoginUser(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	user, err := s.service.Login(req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &userpb.LoginResponse{
		Message: "login successful",
		User:    toUserResponse(user),
	}, nil
}

func (s *UserGRPCServer) GetUserByID(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.UserResponse, error) {
	user, err := s.service.GetUserByID(req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return toUserResponse(user), nil
}

func (s *UserGRPCServer) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.UserResponse, error) {
	user, err := s.service.GetUserByEmail(req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	user.Password = ""
	return toUserResponse(user), nil
}

func (s *UserGRPCServer) GetAllUsers(ctx context.Context, req *userpb.GetAllUsersRequest) (*userpb.GetAllUsersResponse, error) {
	users, err := s.service.GetAllUsers()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toUsersResponse(users), nil
}

func (s *UserGRPCServer) UpdateUserName(ctx context.Context, req *userpb.UpdateUserNameRequest) (*userpb.UserResponse, error) {
	user, err := s.service.UpdateUserName(req.GetId(), req.GetName())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return toUserResponse(user), nil
}

func (s *UserGRPCServer) UpdateUserRole(ctx context.Context, req *userpb.UpdateUserRoleRequest) (*userpb.UserResponse, error) {
	user, err := s.service.UpdateUserRole(req.GetId(), req.GetRole())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return toUserResponse(user), nil
}

func (s *UserGRPCServer) ChangePassword(ctx context.Context, req *userpb.ChangePasswordRequest) (*userpb.MessageResponse, error) {
	if err := s.service.ChangePassword(req.GetId(), req.GetPassword()); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &userpb.MessageResponse{Message: "password changed successfully"}, nil
}

func (s *UserGRPCServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.MessageResponse, error) {
	if err := s.service.DeleteUser(req.GetId()); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &userpb.MessageResponse{Message: "user deleted successfully"}, nil
}

func (s *UserGRPCServer) CheckUserExists(ctx context.Context, req *userpb.CheckUserExistsRequest) (*userpb.CheckUserExistsResponse, error) {
	exists, err := s.service.CheckUserExists(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.CheckUserExistsResponse{Exists: exists}, nil
}

func (s *UserGRPCServer) CountUsers(ctx context.Context, req *userpb.CountUsersRequest) (*userpb.CountUsersResponse, error) {
	count, err := s.service.CountUsers()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.CountUsersResponse{Count: int32(count)}, nil
}

func (s *UserGRPCServer) GetUsersByRole(ctx context.Context, req *userpb.GetUsersByRoleRequest) (*userpb.GetAllUsersResponse, error) {
	users, err := s.service.GetUsersByRole(req.GetRole())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toUsersResponse(users), nil
}

func toUserResponse(user *model.User) *userpb.UserResponse {
	return &userpb.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.String(),
	}
}

func toUsersResponse(users []model.User) *userpb.GetAllUsersResponse {
	response := &userpb.GetAllUsersResponse{
		Users: make([]*userpb.UserResponse, 0, len(users)),
	}

	for _, user := range users {
		response.Users = append(response.Users, toUserResponse(&user))
	}

	return response
}
