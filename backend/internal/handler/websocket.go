package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking in production
		return true
	},
}

// WebSocketMessage represents a message sent over WebSocket.
type WebSocketMessage struct {
	Type      string    `json:"type"`
	Content   string    `json:"content,omitempty"`
	MessageID string    `json:"messageId,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// MessageCompleteResponse is sent when a message stream is complete.
type MessageCompleteResponse struct {
	Type        string            `json:"type"`
	MessageID   string            `json:"messageId"`
	FullContent string            `json:"fullContent"`
	CodeBlocks  []model.CodeBlock `json:"codeBlocks"`
	Timestamp   time.Time         `json:"timestamp"`
}

// ErrorResponse is sent when an error occurs.
type ErrorResponse struct {
	Type      string    `json:"type"`
	Error     string    `json:"error"`
	Code      string    `json:"code,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// WebSocketHandler handles WebSocket connections.
type WebSocketHandler struct {
	chatService *service.ChatService
	logger      zerolog.Logger
}

// NewWebSocketHandler creates a new WebSocketHandler.
func NewWebSocketHandler(chatService *service.ChatService, logger zerolog.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		chatService: chatService,
		logger:      logger,
	}
}

// HandleConnection handles a WebSocket connection.
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	projectIDParam := c.Query("projectId")
	if projectIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "projectId is required"})
		return
	}

	projectID, err := uuid.Parse(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid projectId"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to upgrade connection")
		return
	}
	defer conn.Close()

	h.logger.Info().Str("projectId", projectID.String()).Msg("WebSocket connection established")

	// Create a mutex for thread-safe writes to the WebSocket
	var writeMu sync.Mutex

	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error().Err(err).Msg("unexpected WebSocket close")
			}
			break
		}

		h.logger.Debug().
			Str("type", msg.Type).
			Str("projectId", projectID.String()).
			Msg("received message")

		// Handle different message types
		switch msg.Type {
		case "ping":
			h.sendPong(conn, &writeMu)
		case "chat_message":
			h.handleChatMessage(c.Request.Context(), conn, &writeMu, projectID, msg)
		default:
			h.sendError(conn, &writeMu, "unknown message type", "UNKNOWN_TYPE")
		}
	}

	h.logger.Info().Str("projectId", projectID.String()).Msg("WebSocket connection closed")
}

func (h *WebSocketHandler) sendPong(conn *websocket.Conn, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	response := WebSocketMessage{
		Type:      "pong",
		Timestamp: time.Now().UTC(),
	}
	conn.WriteJSON(response)
}

func (h *WebSocketHandler) handleChatMessage(ctx context.Context, conn *websocket.Conn, mu *sync.Mutex, projectID uuid.UUID, msg WebSocketMessage) {
	messageID := uuid.New().String()

	// Send message_start
	h.sendMessageStart(conn, mu, messageID)

	// Process message through chat service with streaming
	onChunk := func(chunk string) {
		mu.Lock()
		defer mu.Unlock()

		chunkMsg := WebSocketMessage{
			Type:      "message_chunk",
			MessageID: messageID,
			Content:   chunk,
			Timestamp: time.Now().UTC(),
		}
		if err := conn.WriteJSON(chunkMsg); err != nil {
			h.logger.Error().Err(err).Msg("failed to write chunk to WebSocket")
		}
	}

	// Create a context with timeout for the Claude API call
	chatCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	result, err := h.chatService.ProcessMessage(chatCtx, projectID, msg.Content, onChunk)
	if err != nil {
		h.logger.Error().Err(err).
			Str("projectId", projectID.String()).
			Msg("failed to process message")
		h.sendError(conn, mu, "Failed to generate response", "AI_ERROR")
		return
	}

	// Send message_complete with code blocks
	h.sendMessageComplete(conn, mu, messageID, result.Content, result.CodeBlocks)
}

func (h *WebSocketHandler) sendMessageStart(conn *websocket.Conn, mu *sync.Mutex, messageID string) {
	mu.Lock()
	defer mu.Unlock()

	startMsg := WebSocketMessage{
		Type:      "message_start",
		MessageID: messageID,
		Timestamp: time.Now().UTC(),
	}
	conn.WriteJSON(startMsg)
}

func (h *WebSocketHandler) sendMessageComplete(conn *websocket.Conn, mu *sync.Mutex, messageID string, content string, codeBlocks []model.CodeBlock) {
	mu.Lock()
	defer mu.Unlock()

	completeMsg := MessageCompleteResponse{
		Type:        "message_complete",
		MessageID:   messageID,
		FullContent: content,
		CodeBlocks:  codeBlocks,
		Timestamp:   time.Now().UTC(),
	}
	conn.WriteJSON(completeMsg)
}

func (h *WebSocketHandler) sendError(conn *websocket.Conn, mu *sync.Mutex, errorMsg string, code string) {
	mu.Lock()
	defer mu.Unlock()

	response := ErrorResponse{
		Type:      "error",
		Error:     errorMsg,
		Code:      code,
		Timestamp: time.Now().UTC(),
	}
	conn.WriteJSON(response)
}
