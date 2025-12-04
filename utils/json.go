package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONUtils 提供JSON处理的工具函数集合
// 包含序列化、反序列化、格式化等功能

// JSONMarshal 将数据序列化为JSON字节数组
func JSONMarshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// JSONMarshalIndent 将数据序列化为带缩进的JSON字节数组
func JSONMarshalIndent(data interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(data, prefix, indent)
}

// JSONUnmarshal 将JSON字节数组反序列化为数据
func JSONUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// JSONMarshalToString 将数据序列化为JSON字符串
func JSONMarshalToString(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// JSONUnmarshalFromString 将JSON字符串反序列化为数据
func JSONUnmarshalFromString(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// ToJSONString 将数据转换为格式化的JSON字符串
func ToJSONString(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// PrettyPrintJSON 格式化打印JSON数据
func PrettyPrintJSON(data interface{}) error {
	jsonStr, err := ToJSONString(data)
	if err != nil {
		return err
	}
	fmt.Println(jsonStr)
	return nil
}

// SaveJSONToFile 将数据保存为JSON文件
func SaveJSONToFile(data interface{}, filePath string) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonBytes, 0644)
}

// LoadJSONFromFile 从JSON文件加载数据
func LoadJSONFromFile(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// ValidateJSON 验证JSON字符串是否有效
func ValidateJSON(jsonStr string) error {
	var js interface{}
	return json.Unmarshal([]byte(jsonStr), &js)
}

// JSONToMap 将JSON字符串转换为map
func JSONToMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}

// MapToJSON 将map转换为JSON字符串
func MapToJSON(m map[string]interface{}) (string, error) {
	return JSONMarshalToString(m)
}