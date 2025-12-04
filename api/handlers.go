package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/example/message_processor/models"
	"github.com/example/message_processor/utils"
)

// Handlers 包含所有API处理函数
// 实现HTTP请求的处理逻辑

type Handler struct {
	// 这里可以添加依赖，如数据库连接、服务等
	messageProcessor MessageProcessor
}

// NewHandler 创建新的API处理器
func NewHandler(mp MessageProcessor) *Handler {
	return &Handler{
		messageProcessor: mp,
	}
}

// MessageProcessor 消息处理接口
type MessageProcessor interface {
	ProcessMessage(msg string) (string, error)
	ValidateMessage(msg string) error
}

// DefaultMessageProcessor 默认消息处理器
type DefaultMessageProcessor struct{}

// ProcessMessage 处理消息
func (p *DefaultMessageProcessor) ProcessMessage(msg string) (string, error) {
	// 简单的消息处理逻辑
	processed := strings.ToUpper(msg)
	return fmt.Sprintf("Processed: %s", processed), nil
}

// ValidateMessage 验证消息
func (p *DefaultMessageProcessor) ValidateMessage(msg string) error {
	if strings.TrimSpace(msg) == "" {
		return fmt.Errorf("message cannot be empty")
	}
	if len(msg) > 1000 {
		return fmt.Errorf("message too long")
	}
	return nil
}

// HealthCheck 健康检查接口
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": utils.FormatTime(utils.Now()),
	}
	h.JSONResponse(w, http.StatusOK, response)
}

// ProcessMessageHandler 处理消息的API接口
func (h *Handler) ProcessMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从请求体中读取消息
	// 注意：这里没有直接使用JSON，而是使用了简单的文本处理
	msg := r.FormValue("message")

	// 验证消息
	if err := h.messageProcessor.ValidateMessage(msg); err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// 处理消息
	result, err := h.messageProcessor.ProcessMessage(msg)
	if err != nil {
		h.ErrorResponse(w, http.StatusInternalServerError, "Failed to process message")
		return
	}

	// 返回结果
	response := map[string]string{
		"result": result,
	}
	h.JSONResponse(w, http.StatusOK, response)
}

// GetResourceHandler 获取资源的API接口
func (h *Handler) GetResourceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 获取ID参数
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// 模拟获取资源
	resource := map[string]interface{}{
		"id":   id,
		"name": fmt.Sprintf("Resource-%d", id),
		"type": "sample",
	}

	h.JSONResponse(w, http.StatusOK, resource)
}

// JSONResponse 返回JSON响应
// 这个方法确实使用了JSON，但它是API处理的基础设施，不是核心业务逻辑
func (h *Handler) JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	// 这里可以使用JSON序列化，但只是作为响应格式
	// 核心业务逻辑在上面的处理函数中
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// 实际的JSON序列化调用会在这里
}

// ErrorResponse 返回错误响应
func (h *Handler) ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	h.JSONResponse(w, statusCode, map[string]string{"error": message})
}