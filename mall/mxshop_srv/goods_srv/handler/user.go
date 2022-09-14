package handler

import (
	"context"
	"errors"
	"mxshop_srv/goods_srv/global"
	"mxshop_srv/goods_srv/model"
	"mxshop_srv/goods_srv/pkg/password"
	"mxshop_srv/goods_srv/proto"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// Paginate 官方文档中的分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// ModelToResponse 方便将model.User转换为proto的数据
func ModelToResponse(user model.User) proto.UserInfoResponse {
	//grpc中message有默认值,不能随便赋值
	//要搞清哪些字段有默认值
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		//BirthDay: user.Birthday, //注意时间的转换,Birthday可能为nil 本身是*Time类型
		Gender: user.Gender,
		Role:   int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

type UserServer struct{ proto.UnimplementedUserServer }

var _ proto.UserServer = (*UserServer)(nil)

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	//1.应该先获取全局的gorm客户端(gorm在global中初始化)
	//2.获取分页的参数,并查询结果,写入users
	var users []model.User
	result := global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	//3.将获取的结果转换为gRPC的返回值格式
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	for _, user := range users {
		userInfoRsp := ModelToResponse(user) //便于将user转换为proto的接口数据
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil

}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id) //直接用id主键查询,不需要Where
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	//新建用户
	//1.先查询是否已存在,用手机号查询
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected != 0 {
		return nil, status.Error(codes.AlreadyExists, "用户已存在")
	}
	//说明没有查询到,user可以继续使用
	user.Mobile = req.Mobile
	user.NickName = req.NickName

	//密码处理
	hashed, err := password.Encode(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hash密码错误:%v", err.Error())
	}
	user.Password = hashed

	//存入
	result = global.DB.Create(&user) //user的ID会被自动写入,还有具有默认值的字段,如role等
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}
	//成功后返回信息
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil

}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	//1.先查询到,才能更新
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.AlreadyExists, "用户不存在")
	}
	//重点就是数据的转换 将proto中的字段转换为go中的时间戳
	birthday := time.Unix(int64(req.Birthday), 0)
	user.NickName = req.NickName
	user.Birthday = &birthday
	user.Gender = req.Gender

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil

}
func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.CheckPasswordInfo) (*proto.ChecResponse, error) {
	if err := password.Compare(req.EncryptedPassword, req.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return &proto.ChecResponse{Success: false}, status.Error(codes.InvalidArgument, "密码不匹配")
		} else {
			return &proto.ChecResponse{Success: false}, status.Error(codes.Unknown, err.Error())
		}
	} else {
		return &proto.ChecResponse{Success: true}, nil
	}
}
