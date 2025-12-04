# Function Search Results

This file contains all functions defined in the message_processor repository, organized by file path.

## message_processor/main.go

```go
// main 程序入口函数，测试User模型的JSON序列化和反序列化
func main()

// ConvertToJSON 将任意类型转换为JSON字符串
func ConvertToJSON(data interface{}) (string, error)

// ParseFromJSON 从JSON字符串解析到目标类型
func ParseFromJSON(jsonStr string, target interface{}) error
```

## message_processor/models/user.go

```go
// (u User) MarshalJSON 自定义JSON序列化方法，实现json.Marshaler接口
func (u User) MarshalJSON() ([]byte, error)

// (u *User) UnmarshalJSON 自定义JSON反序列化方法，实现json.Unmarshaler接口
func (u *User) UnmarshalJSON(data []byte) error

// (u User) ToJSON 将User转换为JSON字符串
func (u User) ToJSON() (string, error)

// FromJSON 从JSON字符串创建User
func FromJSON(jsonStr string) (User, error)
```

## message_processor/models/config.go

```go
// (c Config) MarshalJSON 自定义JSON序列化方法
func (c Config) MarshalJSON() ([]byte, error)

// (c *Config) UnmarshalJSON 自定义JSON反序列化方法
func (c *Config) UnmarshalJSON(data []byte) error

// LoadConfigFromFile 从文件加载配置
func LoadConfigFromFile(filePath string) (Config, error)

// (c Config) SaveConfigToFile 将配置保存到文件
func (c Config) SaveConfigToFile(filePath string) error
```

## message_processor/utils/json.go

```go
// JSONMarshal 将数据序列化为JSON字节数组
func JSONMarshal(data interface{}) ([]byte, error)

// JSONMarshalIndent 将数据序列化为带缩进的JSON字节数组
func JSONMarshalIndent(data interface{}, prefix, indent string) ([]byte, error)

// JSONUnmarshal 将JSON字节数组反序列化为数据
func JSONUnmarshal(data []byte, v interface{}) error

// JSONMarshalToString 将数据序列化为JSON字符串
func JSONMarshalToString(data interface{}) (string, error)

// JSONUnmarshalFromString 将JSON字符串反序列化为数据
func JSONUnmarshalFromString(jsonStr string, v interface{}) error

// ToJSONString 将数据转换为格式化的JSON字符串
func ToJSONString(data interface{}) (string, error)

// PrettyPrintJSON 格式化打印JSON数据
func PrettyPrintJSON(data interface{}) error

// SaveJSONToFile 将数据保存为JSON文件
func SaveJSONToFile(data interface{}, filePath string) error

// LoadJSONFromFile 从JSON文件加载数据
func LoadJSONFromFile(filePath string, v interface{}) error

// ValidateJSON 验证JSON字符串是否有效
func ValidateJSON(jsonStr string) error

// JSONToMap 将JSON字符串转换为map
func JSONToMap(jsonStr string) (map[string]interface{}, error)

// MapToJSON 将map转换为JSON字符串
func MapToJSON(m map[string]interface{}) (string, error)
```

## message_processor/utils/helpers.go

```go
// TruncateString 截断字符串到指定长度
func TruncateString(s string, maxLen int) string

// SnakeToCamel 将蛇形命名转换为驼峰命名
func SnakeToCamel(s string) string

// CamelToSnake 将驼峰命名转换为蛇形命名
func CamelToSnake(s string) string

// FormatTime 将时间格式化为标准格式
func FormatTime(t time.Time) string

// ParseTime 解析时间字符串
func ParseTime(s string) (time.Time, error)

// GetTimeAgo 获取时间差的友好描述
func GetTimeAgo(t time.Time) string

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) (string, error)

// GenerateRandomID 生成随机ID
func GenerateRandomID() (string, error)

// Min 返回两个整数中的较小值
func Min(a, b int) int

// Max 返回两个整数中的较大值
func Max(a, b int) int

// Clamp 将值限制在指定范围内
func Clamp(value, min, max int) int
```

## message_processor/api/handlers.go

```go
// NewHandler 创建新的API处理器
func NewHandler(mp MessageProcessor) *Handler

// (p *DefaultMessageProcessor) ProcessMessage 处理消息
func (p *DefaultMessageProcessor) ProcessMessage(msg string) (string, error)

// (p *DefaultMessageProcessor) ValidateMessage 验证消息
func (p *DefaultMessageProcessor) ValidateMessage(msg string) error

// (h *Handler) HealthCheck 健康检查接口
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request)

// (h *Handler) ProcessMessageHandler 处理消息的API接口
func (h *Handler) ProcessMessageHandler(w http.ResponseWriter, r *http.Request)

// (h *Handler) GetResourceHandler 获取资源的API接口
func (h *Handler) GetResourceHandler(w http.ResponseWriter, r *http.Request)

// (h *Handler) JSONResponse 返回JSON响应
func (h *Handler) JSONResponse(w http.ResponseWriter, statusCode int, data interface{})

// (h *Handler) ErrorResponse 返回错误响应
func (h *Handler) ErrorResponse(w http.ResponseWriter, statusCode int, message string)
```

## message_processor/middleware/auth.go

```go
// NewAuthMiddleware 创建新的认证中间件
func NewAuthMiddleware(jwtSecret string, apiKeyPrefix string) *AuthMiddleware

// (m *AuthMiddleware) JWTAuth JWT认证中间件
func (m *AuthMiddleware) JWTAuth(next http.Handler) http.Handler

// (m *AuthMiddleware) APIKeyAuth API密钥认证中间件
func (m *AuthMiddleware) APIKeyAuth(next http.Handler) http.Handler

// (m *AuthMiddleware) GenerateJWT 生成JWT令牌
func (m *AuthMiddleware) GenerateJWT(userID int, username string) (string, error)

// (m *AuthMiddleware) CORS 跨域资源共享中间件
func (m *AuthMiddleware) CORS(next http.Handler) http.Handler

// (m *AuthMiddleware) RateLimit 速率限制中间件
func (m *AuthMiddleware) RateLimit(next http.Handler) http.Handler

// (m *AuthMiddleware) UnauthorizedResponse 未授权响应
func (m *AuthMiddleware) UnauthorizedResponse(w http.ResponseWriter, message string)

// (m *AuthMiddleware) ErrorResponse 错误响应
func (m *AuthMiddleware) ErrorResponse(w http.ResponseWriter, status int, message string)
```

## message_processor/storage/db.go

```go
// NewPostgresDB 创建新的PostgreSQL数据库实例
func NewPostgresDB(config DBConfig) *PostgresDB

// (p *PostgresDB) Connect 连接到数据库
func (p *PostgresDB) Connect(ctx context.Context) error

// (p *PostgresDB) Disconnect 断开数据库连接
func (p *PostgresDB) Disconnect(ctx context.Context) error

// (p *PostgresDB) Ping 测试数据库连接
func (p *PostgresDB) Ping(ctx context.Context) error

// (p *PostgresDB) CreateUser 创建用户
func (p *PostgresDB) CreateUser(ctx context.Context, user *models.User) error

// (p *PostgresDB) GetUserByID 通过ID获取用户
func (p *PostgresDB) GetUserByID(ctx context.Context, id int) (*models.User, error)

// (p *PostgresDB) GetUserByUsername 通过用户名获取用户
func (p *PostgresDB) GetUserByUsername(ctx context.Context, username string) (*models.User, error)

// (p *PostgresDB) UpdateUser 更新用户信息
func (p *PostgresDB) UpdateUser(ctx context.Context, user *models.User) error

// (p *PostgresDB) DeleteUser 删除用户
func (p *PostgresDB) DeleteUser(ctx context.Context, id int) error

// (p *PostgresDB) BeginTx 开始事务
func (p *PostgresDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)

// (t *PostgresTx) Commit 提交事务
func (t *PostgresTx) Commit() error

// (t *PostgresTx) Rollback 回滚事务
func (t *PostgresTx) Rollback() error

// (t *PostgresTx) CreateUser 在事务中创建用户
func (t *PostgresTx) CreateUser(ctx context.Context, user *models.User) error

// (t *PostgresTx) GetUserByID 在事务中通过ID获取用户
func (t *PostgresTx) GetUserByID(ctx context.Context, id int) (*models.User, error)
```

## message_processor/cmd/server/main.go

```go
// main 服务器入口函数
func main()

// setupRouter 设置HTTP路由
func setupRouter(handler *api.Handler, authMiddleware *middleware.AuthMiddleware) *http.ServeMux

// loadConfig 加载配置
func loadConfig(configFile string) (*AppConfig, error)
```