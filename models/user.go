package models

import (
	"encoding/json"
	"time"
)

// User 用户模型
// 包含用户基本信息
// 实现了json.Marshaler和json.Unmarshaler接口

// User 用户结构体
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // 密码不序列化到JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MarshalJSON 自定义JSON序列化方法
// 实现json.Marshaler接口
func (u User) MarshalJSON() ([]byte, error) {
	// 创建一个匿名结构体用于JSON序列化
	type Alias User
	return json.Marshal(&struct {
		Alias
		// 添加自定义字段
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (Alias)(u),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	})
}

// UnmarshalJSON 自定义JSON反序列化方法
// 实现json.Unmarshaler接口
func (u *User) UnmarshalJSON(data []byte) error {
	// 创建一个匿名结构体用于JSON反序列化
	type Alias User
	aux := &struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 解析时间字符串
	createdAt, err := time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return err
	}
	u.CreatedAt = createdAt

	updatedAt, err := time.Parse(time.RFC3339, aux.UpdatedAt)
	if err != nil {
		return err
	}
	u.UpdatedAt = updatedAt

	return nil
}

// ToJSON 将User转换为JSON字符串
func (u User) ToJSON() (string, error) {
	bytes, err := u.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON 从JSON字符串创建User
func FromJSON(jsonStr string) (User, error) {
	var user User
	if err := json.Unmarshal([]byte(jsonStr), &user); err != nil {
		return User{}, err
	}
	return user, nil
}