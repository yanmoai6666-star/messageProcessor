package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/example/message_processor/models"
	"github.com/example/message_processor/utils"
)

// 主程序入口
func main() {
	fmt.Println("Message Processor Starting...")

	// 创建测试用户
	user := models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	// 测试JSON序列化
	userJSON, err := user.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed to marshal user: %v", err)
	}

	fmt.Printf("User JSON: %s\n", userJSON)

	// 测试JSON反序列化
	var newUser models.User
	if err := json.Unmarshal(userJSON, &newUser); err != nil {
		log.Fatalf("Failed to unmarshal user: %v", err)
	}

	fmt.Printf("Unmarshaled User: %+v\n", newUser)
}

// ConvertToJSON 将任意类型转换为JSON字符串
func ConvertToJSON(data interface{}) (string, error) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// ParseFromJSON 从JSON字符串解析到目标类型
func ParseFromJSON(jsonStr string, target interface{}) error {
	return json.Unmarshal([]byte(jsonStr), target)
}