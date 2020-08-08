package service

import (
	pb "Entry_Task/src/public/protos"
	"Entry_Task/src/server/tcp/repositories"
	"context"
)

type UserService interface {
	pb.UserServiceServer
}

type UserServiceManager struct {
	repo  repositories.UserRepository
	redis repositories.RedisStorage
}

func NewUserServiceManager(repo repositories.UserRepository, redis repositories.RedisStorage) UserService {
	return &UserServiceManager{repo: repo, redis: redis}
}

//grpc apis for log in request
func (u *UserServiceManager) LogIn(ctx context.Context, in *pb.LogInRequest) (*pb.LogInReply, error) {
	user, err := u.repo.SelectByUser(in.GetUsername(), in.GetPassword())

	return &pb.LogInReply{
		UserId:         user.UserId,
		Username:       user.Username,
		Password:       user.Password,
		Nickname:       user.Nickname,
		ProfilePicture: user.ProfilePicture}, err
}

//grpc apis for modifying Nickname
func (u *UserServiceManager) ModifyNickName(ctx context.Context, in *pb.ModifyNickNameRequest) (*pb.UsrOpReply, error) {
	err := u.repo.UpdateNickName(in.GetUserId(), in.GetNickname())
	if err != nil {
		return &pb.UsrOpReply{
			Message: err.Error()}, err
	}
	return &pb.UsrOpReply{Message: "Nickname is modified successfully"}, nil
}

//grpc apis for uploading profile picture
func (u *UserServiceManager) UploadProfile(ctx context.Context, in *pb.UploadProfileRequest) (*pb.UsrOpReply, error) {
	err := u.repo.UpdateProfile(in.GetUserId(), in.GetProfilePicture())
	if err != nil {
		return &pb.UsrOpReply{
			Message: err.Error()}, err
	}
	return &pb.UsrOpReply{Message: "Your profile picture has been uploaded successfully"}, nil
}

//read user's info from redis
func (u *UserServiceManager) LogInCache(ctx context.Context, in *pb.LogInCacheRequest) (*pb.LogInReply, error) {
	user, err := u.redis.GetByToken(in.GetToken())

	return &pb.LogInReply{
		UserId:         user.UserId,
		Username:       user.Username,
		Password:       user.Password,
		Nickname:       user.Nickname,
		ProfilePicture: user.ProfilePicture}, err
}

//update nickname from redis
func (u *UserServiceManager) ModifyNickNameCache(ctx context.Context, in *pb.ModifyNickNameCacheRequest) (*pb.UsrOpReply, error) {
	err := u.redis.UpdateNickNameByToken(in.GetToken(), in.GetNickname())
	if err != nil {
		return &pb.UsrOpReply{Message: err.Error()}, err
	}
	return &pb.UsrOpReply{Message: "User's nickname has been updated in cache"}, nil
}

//upload profile cache from redis
func (u *UserServiceManager) UploadProfileCache(ctx context.Context, in *pb.UploadProfileCacheRequest) (*pb.UsrOpReply, error) {
	err := u.redis.UploadProfileByToken(in.GetToken(), in.GetProfilePicture())
	if err != nil {
		return &pb.UsrOpReply{Message: err.Error()}, err
	}
	return &pb.UsrOpReply{Message: "User's profile picture has been uploaded in cache"}, nil
}

//write data from mysql to redis
func (u *UserServiceManager) MysqlToCache(ctx context.Context, in *pb.MysqlToCacheRequest) (*pb.UsrOpReply, error) {
	err := u.redis.SetByUser(in.GetToken(), in.GetUsername(), in.GetPassword(),
		in.GetNickname(), in.GetProfilePicture(), in.GetUserId())
	if err != nil {
		return &pb.UsrOpReply{Message: err.Error()}, err
	}
	return &pb.UsrOpReply{Message: "User's info has been written into cache successfully"}, nil
}
