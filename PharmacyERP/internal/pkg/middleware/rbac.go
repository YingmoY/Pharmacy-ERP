package middleware

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// casbinRule 对应 casbin_rule 表的结构
type casbinRule struct {
	ID    int64  `gorm:"column:id"`
	Ptype string `gorm:"column:ptype"`
	V0    string `gorm:"column:v0"`
	V1    string `gorm:"column:v1"`
	V2    string `gorm:"column:v2"`
}

func (casbinRule) TableName() string { return "casbin_rule" }

// rbacCache 用于缓存从数据库加载的 casbin 规则，减少频繁查询
type rbacCache struct {
	mu       sync.RWMutex
	rules    []casbinRule
	loadedAt time.Time
	ttl      time.Duration
}

var globalRBACCache = &rbacCache{ttl: 30 * time.Second}

// loadRules 从 DB 或缓存中获取所有 casbin_rule 记录
func (rc *rbacCache) loadRules(db *gorm.DB) ([]casbinRule, error) {
	rc.mu.RLock()
	// 缓存未过期则直接返回
	if time.Since(rc.loadedAt) < rc.ttl && rc.rules != nil {
		rules := rc.rules
		rc.mu.RUnlock()
		return rules, nil
	}
	rc.mu.RUnlock()

	// 重新从数据库加载，加写锁
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// 双重检查，防止并发情况下多次加载
	if time.Since(rc.loadedAt) < rc.ttl && rc.rules != nil {
		return rc.rules, nil
	}

	var rules []casbinRule
	if err := db.Find(&rules).Error; err != nil {
		return nil, err
	}
	rc.rules = rules
	rc.loadedAt = time.Now()
	return rules, nil
}

// CasbinRBACAuth RBAC权限校验中间件，基于casbin_rule表动态校验
func CasbinRBACAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户的角色（来自 JWT 中间件写入）
		roleStr, exists := GetCurrentUserRole(c)
		if !exists || roleStr == "" {
			core.FailWithStatus(c, 403, ecode.ErrPermission.Code, ecode.ErrPermission.Msg)
			c.Abort()
			return
		}

		method := c.Request.Method
		path := c.Request.URL.Path

		// 从缓存或数据库加载规则
		rules, err := globalRBACCache.loadRules(db)
		if err != nil {
			core.FailWithStatus(c, 500, ecode.ErrSystem.Code, ecode.ErrSystem.Msg)
			c.Abort()
			return
		}

		// 构建当前用户拥有的所有角色集合（主角色 + g 型继承角色）
		userRoles := map[string]bool{roleStr: true}

		// 查找 g 类型规则：v0=user:ID 关联的角色
		userID, hasUID := GetCurrentUserID(c)
		if hasUID && userID > 0 {
			userPrincipal := fmt.Sprintf("user:%d", userID)
			for _, r := range rules {
				if r.Ptype == "g" && r.V0 == userPrincipal && r.V1 != "" {
					userRoles[r.V1] = true
				}
			}
		}

		// 检查 p 类型规则，判断是否有访问权限
		allowed := false
		for _, r := range rules {
			if r.Ptype != "p" {
				continue
			}
			// v0 为角色代码，v1 为路径模式，v2 为 HTTP 方法
			if !userRoles[r.V0] {
				continue
			}
			// 方法匹配（不区分大小写，* 通配所有方法）
			if r.V2 != "*" && !strings.EqualFold(r.V2, method) {
				continue
			}
			// 使用 keyMatch2 进行路径匹配
			if keyMatch2(r.V1, path) {
				allowed = true
				break
			}
		}

		if !allowed {
			core.FailWithStatus(c, 403, ecode.ErrPermission.Code, ecode.ErrPermission.Msg)
			c.Abort()
			return
		}

		c.Next()
	}
}

// keyMatch2 实现casbin keyMatch2算法，将:param替换为正则[^/]+进行匹配
func keyMatch2(pattern, actual string) bool {
	// 将路径参数 :param 替换为正则 [^/]+，普通部分进行转义
	parts := strings.Split(pattern, "/")
	regexParts := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			// 路径参数，匹配非斜杠的任意非空字符
			regexParts = append(regexParts, `[^/]+`)
		} else {
			// 普通路径段，进行正则转义
			regexParts = append(regexParts, regexp.QuoteMeta(part))
		}
	}
	regexStr := "^" + strings.Join(regexParts, "/") + "$"
	matched, err := regexp.MatchString(regexStr, actual)
	if err != nil {
		return false
	}
	return matched
}
