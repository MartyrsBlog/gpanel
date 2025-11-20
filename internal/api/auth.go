package api

import (
	"net/http"
	"time"

	"gpanel/internal/config"
	"gpanel/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Token    string      `json:"token"`
	User     UserInfo    `json:"user"`
	ExpireIn int64       `json:"expire_in"`
}

// UserInfo 用户信息结构体
type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	// 验证用户名和密码
	var userInfo UserInfo
	var password string
	
	err := config.DB.QueryRow(
		"SELECT id, username, email, role, password FROM users WHERE username = ?",
		req.Username,
	).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Email, &userInfo.Role, &password)
	
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalid username or password",
			"data":    nil,
		})
		return
	}

	// 验证密码 (生产环境应使用 bcrypt)
	if password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalid username or password",
			"data":    nil,
		})
		return
	}

	// 生成 JWT token
	cfg, _ := config.Load()
	expireTime := time.Now().Add(time.Duration(cfg.Auth.TokenExpire) * time.Second)
	
	claims := &middleware.Claims{
		Username: userInfo.Username,
		Role:     userInfo.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gpanel",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Auth.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to generate token",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Login successful",
		"data": LoginResponse{
			Token:    tokenString,
			User:     userInfo,
			ExpireIn: int64(cfg.Auth.TokenExpire),
		},
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	// 在实际应用中，可以将 token 加入黑名单
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Logout successful",
		"data":    nil,
	})
}

// GetUserInfo 获取当前用户信息
func GetUserInfo(c *gin.Context) {
	username, _ := c.Get("username")

	var userInfo UserInfo
	err := config.DB.QueryRow(
		"SELECT id, username, email, role FROM users WHERE username = ?",
		username,
	).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Email, &userInfo.Role)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "User not found",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    userInfo,
	})
}

// CreateUser 创建用户 (仅管理员)
func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	// 检查用户名是否已存在
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Database error",
			"data":    nil,
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Username already exists",
			"data":    nil,
		})
		return
	}

	// 插入新用户
	_, err = config.DB.Exec(
		"INSERT INTO users (username, password, email, role) VALUES (?, ?, ?, ?)",
		req.Username,
		req.Password, // 生产环境应使用 bcrypt 哈希
		req.Email,
		req.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create user",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "User created successfully",
		"data":    nil,
	})
}

// UpdatePassword 修改密码
func UpdatePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
			"data":    nil,
		})
		return
	}

	username, _ := c.Get("username")

	// 获取当前用户密码
	var currentPassword string
	err := config.DB.QueryRow(
		"SELECT password FROM users WHERE username = ?",
		username,
	).Scan(&currentPassword)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "User not found",
			"data":    nil,
		})
		return
	}

	// 验证旧密码
	if currentPassword != req.OldPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Old password is incorrect",
			"data":    nil,
		})
		return
	}

	// 更新密码
	_, err = config.DB.Exec(
		"UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE username = ?",
		req.NewPassword, // 生产环境应使用 bcrypt 哈希
		username,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update password",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Password updated successfully",
		"data":    nil,
	})
}