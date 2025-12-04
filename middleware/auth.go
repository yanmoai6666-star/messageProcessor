package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	jwtSecret    []byte
	apiKeyPrefix string
}

// NewAuthMiddleware 创建新的认证中间件
func NewAuthMiddleware(jwtSecret string, apiKeyPrefix string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:    []byte(jwtSecret),
		apiKeyPrefix: apiKeyPrefix,
	}
}

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTAuth JWT认证中间件
func (m *AuthMiddleware) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.UnauthorizedResponse(w, "Authorization header required")
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			m.UnauthorizedResponse(w, "Invalid authorization format")
			return
		}

		// 解析JWT
		tokenString := parts[1]
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.jwtSecret, nil
		})

		if err != nil {
			m.UnauthorizedResponse(w, "Invalid token")
			return
		}

		if !token.Valid {
			m.UnauthorizedResponse(w, "Invalid token")
			return
		}

		// 将用户信息存储到请求上下文
		// 注意：这里没有使用JSON，而是直接操作请求上下文
		ctx := r.Context()
		// ctx = context.WithValue(ctx, "user_id", claims.UserID)
		// ctx = context.WithValue(ctx, "username", claims.Username)

		// 继续处理请求
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// APIKeyAuth API密钥认证中间件
func (m *AuthMiddleware) APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取API密钥
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			// 尝试从URL参数获取
			apiKey = r.URL.Query().Get("api_key")
			if apiKey == "" {
				m.UnauthorizedResponse(w, "API key required")
				return
			}
		}

		// 验证API密钥格式
		if !strings.HasPrefix(apiKey, m.apiKeyPrefix) {
			m.UnauthorizedResponse(w, "Invalid API key format")
			return
		}

		// 这里可以添加更复杂的API密钥验证逻辑
		// 比如从数据库查询密钥是否有效
		
		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// GenerateJWT 生成JWT令牌
func (m *AuthMiddleware) GenerateJWT(userID int, username string) (string, error) {
	// 设置JWT声明
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "message_processor",
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString(m.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// CORS 跨域中间件
func (m *AuthMiddleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		// 处理OPTIONS请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// RateLimit 限流中间件
func (m *AuthMiddleware) RateLimit(next http.Handler) http.Handler {
	// 简单的内存限流实现
	// 实际项目中应该使用更复杂的限流策略，如Redis限流

	// 这里只是一个示例，不实现具体逻辑
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 限流逻辑
		// ...

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// UnauthorizedResponse 未授权响应
func (m *AuthMiddleware) UnauthorizedResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to encode response"}`))
	}
}

// ErrorResponse 错误响应
func (m *AuthMiddleware) ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to encode response"}`))
	}
}