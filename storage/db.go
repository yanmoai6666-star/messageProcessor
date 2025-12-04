package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL驱动

	"github.com/example/message_processor/models"
)

// DB 数据库接口
type DB interface {
	// 连接管理
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context) error

	// 用户相关操作
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error

	// 事务管理
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
}

// Tx 事务接口
type Tx interface {
	Commit() error
	Rollback() error

	// 用户操作（事务中）
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
}

// PostgresDB PostgreSQL数据库实现
type PostgresDB struct {
	db     *sql.DB
	config DBConfig
}

// DBConfig 数据库配置
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB 创建新的PostgreSQL数据库实例
func NewPostgresDB(config DBConfig) *PostgresDB {
	return &PostgresDB{
		config: config,
	}
}

// Connect 连接到数据库
func (p *PostgresDB) Connect(ctx context.Context) error {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.config.Host, p.config.Port, p.config.User, p.config.Password,
		p.config.DBName, p.config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	p.db = db
	return nil
}

// Disconnect 断开数据库连接
func (p *PostgresDB) Disconnect(ctx context.Context) error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// Ping 测试数据库连接
func (p *PostgresDB) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// CreateUser 创建用户
func (p *PostgresDB) CreateUser(ctx context.Context, user *models.User) error {
	// 数据库查询，与JSON序列化无关
	query := `
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := p.db.QueryRowContext(ctx, query, 
		user.Username, user.Email, user.Password, 
		user.CreatedAt, user.UpdatedAt).Scan(&user.ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("user already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID 根据ID获取用户
func (p *PostgresDB) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (p *PostgresDB) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := p.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUser 更新用户
func (p *PostgresDB) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, updated_at = $3
		WHERE id = $4
	`

	user.UpdatedAt = time.Now()

	result, err := p.db.ExecContext(ctx, query, 
		user.Username, user.Email, user.UpdatedAt, user.ID)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser 删除用户
func (p *PostgresDB) DeleteUser(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"

	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// BeginTx 开始事务
func (p *PostgresDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := p.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &PostgresTx{tx: tx}, nil
}

// PostgresTx PostgreSQL事务实现
type PostgresTx struct {
	tx *sql.Tx
}

// Commit 提交事务
func (t *PostgresTx) Commit() error {
	return t.tx.Commit()
}

// Rollback 回滚事务
func (t *PostgresTx) Rollback() error {
	return t.tx.Rollback()
}

// CreateUser 在事务中创建用户
func (t *PostgresTx) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := t.tx.QueryRowContext(ctx, query, 
		user.Username, user.Email, user.Password, 
		user.CreatedAt, user.UpdatedAt).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user in transaction: %w", err)
	}

	return nil
}

// GetUserByID 在事务中根据ID获取用户
func (t *PostgresTx) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := t.tx.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user in transaction: %w", err)
	}

	return &user, nil
}