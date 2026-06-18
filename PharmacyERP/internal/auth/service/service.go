package service

import (
	"context"
	"errors"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/config"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	rolemodel "github.com/YingmoY/PharmacyERP/internal/role/model"
	rolerepo "github.com/YingmoY/PharmacyERP/internal/role/repository"
	"github.com/YingmoY/PharmacyERP/internal/user/model"
	userrepo "github.com/YingmoY/PharmacyERP/internal/user/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginResult 登录成功的返回结构
type LoginResult struct {
	Token       string
	ExpiresIn   int64
	TokenType   string
	User        model.UserDTO
	Roles       []rolemodel.Role
	Permissions []rolemodel.Permission
}

// AuthService 认证服务接口
type AuthService interface {
	Login(ctx context.Context, username, password, ip, userAgent string) (*LoginResult, error)
	Logout(ctx context.Context, userID int64) error
	GetCurrentUser(ctx context.Context, userID int64) (*model.UserDTO, []rolemodel.Role, []rolemodel.Permission, error)
	ChangePassword(ctx context.Context, userID int64, oldPwd, newPwd string) error
}

type authService struct {
	db       *gorm.DB
	userRepo userrepo.UserRepository
	roleRepo rolerepo.RoleRepository
	mqClient *mq.Client
	jwtCfg   *config.JWTConfig
	log      *zap.Logger
}

// NewAuthService 创建认证服务实例
func NewAuthService(
	db *gorm.DB,
	userRepo userrepo.UserRepository,
	roleRepo rolerepo.RoleRepository,
	mqClient *mq.Client,
	jwtCfg *config.JWTConfig,
	log *zap.Logger,
) AuthService {
	return &authService{
		db:       db,
		userRepo: userRepo,
		roleRepo: roleRepo,
		mqClient: mqClient,
		jwtCfg:   jwtCfg,
		log:      log,
	}
}

// Login 用户登录，校验账号密码，生成 JWT token，记录登录日志
func (s *authService) Login(ctx context.Context, username, password, ip, userAgent string) (*LoginResult, error) {
	// 查询用户
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 异步记录失败登录日志
			s.asyncSaveLoginLog(ctx, 0, username, ip, userAgent, false, "用户不存在")
			return nil, ecode.ErrUserNotFound
		}
		return nil, ecode.ErrSystem
	}

	// 校验账号状态
	if user.Status != 1 {
		s.asyncSaveLoginLog(ctx, user.ID, username, ip, userAgent, false, "账号已禁用")
		return nil, ecode.ErrUserDisabled
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.asyncSaveLoginLog(ctx, user.ID, username, ip, userAgent, false, "密码错误")
		return nil, ecode.ErrPasswordWrong
	}

	// 获取用户角色
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		s.log.Error("获取用户角色失败", zap.Int64("user_id", user.ID), zap.Error(err))
		roles = []rolemodel.Role{}
	}

	// 取第一个角色作为 JWT 中的主角色
	primaryRole := ""
	if len(roles) > 0 {
		primaryRole = roles[0].Code
	}

	permissions, err := s.userRepo.GetUserPermissions(ctx, user.ID)
	if err != nil {
		s.log.Error("鑾峰彇鐢ㄦ埛鏉冮檺澶辫触", zap.Int64("user_id", user.ID), zap.Error(err))
		permissions = []rolemodel.Permission{}
	}

	// 生成 JWT token
	token, err := middleware.GenerateToken(
		user.ID,
		primaryRole,
		s.jwtCfg.Secret,
		s.jwtCfg.ExpireTime,
		s.jwtCfg.Issuer,
	)
	if err != nil {
		s.log.Error("生成 token 失败", zap.Int64("user_id", user.ID), zap.Error(err))
		return nil, ecode.ErrSystem
	}

	// 更新最后登录时间和 IP
	now := time.Now()
	s.db.WithContext(ctx).Model(user).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ip,
	})

	// 异步记录成功登录日志
	s.asyncSaveLoginLog(ctx, user.ID, username, ip, userAgent, true, "")

	// 异步发布登录事件到消息队列
	s.asyncPublishLoginEvent(user.ID, username, ip)

	return &LoginResult{
		Token:       token,
		ExpiresIn:   int64(s.jwtCfg.ExpireTime.Seconds()),
		TokenType:   "Bearer",
		User:        user.ToDTO(),
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

// Logout 用户登出（记录日志，实际 token 失效依赖客户端）
func (s *authService) Logout(ctx context.Context, userID int64) error {
	s.log.Info("用户登出", zap.Int64("user_id", userID))
	// 异步发布登出事件
	if s.mqClient != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					s.log.Error("发布登出事件 panic", zap.Any("recover", r))
				}
			}()
			_ = s.mqClient.PublishLogEvent(context.Background(), mq.LogEvent{
				BusinessType: "auth",
				BusinessID:   "",
				Action:       "logout",
				OperatorID:   userID,
				Detail:       map[string]interface{}{"user_id": userID},
			})
		}()
	}
	return nil
}

// GetCurrentUser 获取当前用户信息及其角色
func (s *authService) GetCurrentUser(ctx context.Context, userID int64) (*model.UserDTO, []rolemodel.Role, []rolemodel.Permission, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, ecode.ErrUserNotFound
		}
		return nil, nil, nil, ecode.ErrSystem
	}

	roles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		s.log.Error("获取用户角色失败", zap.Int64("user_id", userID), zap.Error(err))
		roles = []rolemodel.Role{}
	}

	permissions, err := s.userRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		s.log.Error("鑾峰彇鐢ㄦ埛鏉冮檺澶辫触", zap.Int64("user_id", userID), zap.Error(err))
		permissions = []rolemodel.Permission{}
	}

	dto := user.ToDTO()
	return &dto, roles, permissions, nil
}

// ChangePassword 修改当前用户密码，需要验证旧密码
func (s *authService) ChangePassword(ctx context.Context, userID int64, oldPwd, newPwd string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ecode.ErrUserNotFound
		}
		return ecode.ErrSystem
	}

	// 校验旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPwd)); err != nil {
		return ecode.ErrPasswordWrong
	}

	// 生成新密码哈希，cost=10
	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), 10)
	if err != nil {
		s.log.Error("生成密码哈希失败", zap.Error(err))
		return ecode.ErrSystem
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}

// asyncSaveLoginLog 异步保存登录日志到数据库
func (s *authService) asyncSaveLoginLog(ctx context.Context, userID int64, username, ip, userAgent string, success bool, failReason string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error("保存登录日志 panic", zap.Any("recover", r))
			}
		}()
		log := &model.LoginLog{
			UserID:    userID,
			Username:  username,
			IP:        ip,
			UserAgent: userAgent,
			Success:   success,
			Message:   failReason,
		}
		if err := s.db.WithContext(context.Background()).Create(log).Error; err != nil {
			s.log.Error("保存登录日志失败", zap.Error(err))
		}
	}()
}

// asyncPublishLoginEvent 异步发布登录事件到消息队列
func (s *authService) asyncPublishLoginEvent(userID int64, username, ip string) {
	if s.mqClient == nil {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error("发布登录事件 panic", zap.Any("recover", r))
			}
		}()
		_ = s.mqClient.PublishLogEvent(context.Background(), mq.LogEvent{
			BusinessType: "auth",
			BusinessID:   "",
			Action:       "login",
			OperatorID:   userID,
			Detail:       map[string]interface{}{"username": username, "ip": ip},
		})
	}()
}
