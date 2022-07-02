package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/adapter/presenter"
	"mashu.example/internal/adapter/utils"
	"mashu.example/internal/usecase/chat/create_direct_message"
	"mashu.example/internal/usecase/chat/load_message_history"
	"mashu.example/internal/usecase/repository"
	"mashu.example/pkg"
)

type wsRequestType string
type wsResponseType string
type wsMsgPayload map[string]interface{}
type wsMessageHandler func(uuid.UUID, wsMsgPayload)

const (
	WS_REQ_CREATE_DM    wsRequestType = "CREATE_DM"
	WS_REQ_SEND_MSG     wsRequestType = "SEND_MSG"
	WS_REQ_LOAD_HISTORY wsRequestType = "LOAD_HISTORY"
)

const (
	WS_RES_SEND_MSG     wsResponseType = "SEND_MSG"
	WS_RES_LOAD_HISTORY wsResponseType = "LOAD_HISTORY"

	WS_RES_SUCCESS wsResponseType = "SUCCESS"
	WS_RES_ERR     wsResponseType = "ERROR"
)

// struct for websocket request message
type wsRequestMessage struct {
	Type    wsRequestType `json:"type"`
	Payload wsMsgPayload  `json:"payload"`
}

// struct for websocket response message
type wsResponseMessage struct {
	Code    int            `json:"code"`
	Type    wsResponseType `json:"type"`
	Payload wsMsgPayload   `json:"payload"`
}

// create a websocket message response
func newWsSuccessResponse(code int, message string) *wsResponseMessage {
	return &wsResponseMessage{code, WS_RES_SUCCESS, wsMsgPayload{"message": message}}
}

// create a websocket err response
func newWsErrResponse(errCode int, message string) *wsResponseMessage {
	return &wsResponseMessage{errCode, WS_RES_ERR, wsMsgPayload{"err": message}}
}

type websocketHandler struct {
	clients         map[uuid.UUID]*utils.WebSocketClient
	wsMsgHandlerMap map[wsRequestType]wsMessageHandler

	userRepo repository.UserRepo
	chatRepo repository.ChatRepo
}

func RegisterWebsocketApi(
	e *gin.Engine,
	userRepo repository.UserRepo,
	chatRepo repository.ChatRepo,
) {
	h := newWebSocketHandler(userRepo, chatRepo)

	e.GET("/websocket", h.handleConnection)
}

func (h *websocketHandler) handleConnection(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "empty userId",
		})
		return
	}

	upgrader := pkg.NewWebSocketUpgrader()

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		logrus.Panic("failed to construct websocket connection")
		return
	}

	defer ws.Close()

	// save ws client into memory
	clientId := uuid.MustParse(userId)
	newClient := utils.NewWebSocketClient(clientId, ws)
	fmt.Println("new client:", newClient)

	h.clients[clientId] = newClient

	// send message to connector
	message := fmt.Sprintf("user %s connected", newClient.UserId)
	fmt.Println(message)
	ws.WriteJSON(message)

	// infinite loop to handle incoming connection
	for {
		var req wsRequestMessage
		err := ws.ReadJSON(&req)
		if err != nil {
			msg := fmt.Sprintf("client %s disconnected", clientId)
			logrus.Info(msg)
			delete(h.clients, clientId)
			break
		}
		logrus.Info("websocket req message", req)

		if handler, ok := h.wsMsgHandlerMap[req.Type]; ok {
			handler(clientId, req.Payload)
		} else {
			response := "invalid command!"
			ws.WriteJSON(response)
		}
	}
}

func (h *websocketHandler) createDM(userId uuid.UUID, payload wsMsgPayload) {
	logrus.Info("start create DM")
	client := h.clients[userId]

	// payload validation
	type createDMPayload struct {
		TargetUserId string `json:"targetUserId" validate:"required"`
	}
	payloadByte, _ := json.Marshal(payload)
	p := &createDMPayload{}
	json.Unmarshal(payloadByte, p)
	err := validator.New().Struct(p)
	if err != nil {
		client.Conn.WriteJSON(newWsErrResponse(
			http.StatusBadRequest,
			fmt.Sprintf("failed to parse payload: %s", err),
		))
		return
	}

	senderId := userId
	receiverId := uuid.MustParse(payload["targetUserId"].(string))
	req := create_direct_message.NewCreateDirectMessageUseCaseReq(
		senderId,
		receiverId,
	)

	res := create_direct_message.NewCreateDirectMessageUseCaseRes()
	uc := create_direct_message.NewCreateDirectMessageUseCase(
		h.chatRepo,
		h.userRepo,
		req,
		res,
	)
	uc.Execute()
	if res.Err != nil {
		client.Conn.WriteJSON(newWsErrResponse(
			http.StatusConflict,
			res.Err.Error(),
		))
		return
	}

	client.Conn.WriteJSON(newWsSuccessResponse(
		http.StatusCreated,
		fmt.Sprintf("dm created, id: %s", res.DirectMessageId.String()),
	))

	logrus.Info("end of creating DM")
}

func (h *websocketHandler) sendMessage(userId uuid.UUID, payload wsMsgPayload) {
	senderClient := h.clients[userId]

	// payload validation
	type sendMessagePayload struct {
		TargetUserId string `json:"targetUserId" validate:"required"`
		Content      string `json:"content" validate:"required"`
	}
	payloadByte, _ := json.Marshal(payload)
	p := &sendMessagePayload{}
	json.Unmarshal(payloadByte, p)
	err := validator.New().Struct(p)
	if err != nil {
		senderClient.Conn.WriteJSON(newWsErrResponse(
			http.StatusBadRequest,
			fmt.Sprintf("failed to parse payload: %s", err)),
		)
		return
	}

	// get DM and add message
	senderId := userId
	receiverId, err := uuid.Parse(p.TargetUserId)
	if err != nil {
		senderClient.Conn.WriteJSON(newWsErrResponse(http.StatusBadRequest, "invalid user id"))
		return
	}
	dm, err := h.chatRepo.GetDMByUserId(senderId, receiverId)
	if err != nil {
		senderClient.Conn.WriteJSON(newWsErrResponse(http.StatusNotFound, "no DM created"))
		return
	}
	dm.AddMessage(senderId, p.Content)

	// send message
	senderClient.Conn.WriteJSON(&wsResponseMessage{
		Code: http.StatusOK,
		Type: WS_RES_SEND_MSG,
		Payload: wsMsgPayload{
			"message": p.Content,
		},
	})
	if receiverClient, ok := h.clients[receiverId]; ok {
		receiverClient.Conn.WriteJSON(&wsResponseMessage{
			Code: http.StatusOK,
			Type: WS_RES_SEND_MSG,
			Payload: wsMsgPayload{
				"message": p.Content,
			},
		})
	}
}

func (h *websocketHandler) loadMessageHistory(userId uuid.UUID, payload wsMsgPayload) {
	fmt.Println("start load history")
	userClient := h.clients[userId]

	req := load_message_history.NewLoadMessageHistoryUseCaseReq(userId)
	res := load_message_history.NewLoadMessageHistoryUseCaseRes()
	uc := load_message_history.NewLoadMessageHistoryUseCase(h.userRepo, h.chatRepo, req, res)

	uc.Execute()
	if res.Err != nil {
		userClient.Conn.WriteJSON(newWsErrResponse(
			http.StatusInternalServerError,
			res.Err.Error(),
		))
		return
	}

	p := presenter.NewMessageHistoryPresenter(res)
	vm := p.BuildViewModel()

	userClient.Conn.WriteJSON(&wsResponseMessage{
		Code: http.StatusOK,
		Type: WS_RES_LOAD_HISTORY,
		Payload: wsMsgPayload{
			"history": vm,
		},
	})
}

func newWebSocketHandler(
	userRepo repository.UserRepo,
	chatRepo repository.ChatRepo,
) *websocketHandler {
	h := &websocketHandler{
		clients:         map[uuid.UUID]*utils.WebSocketClient{},
		wsMsgHandlerMap: map[wsRequestType]wsMessageHandler{},
		userRepo:        userRepo,
		chatRepo:        chatRepo,
	}

	h.wsMsgHandlerMap[WS_REQ_CREATE_DM] = h.createDM
	h.wsMsgHandlerMap[WS_REQ_SEND_MSG] = h.sendMessage
	h.wsMsgHandlerMap[WS_REQ_LOAD_HISTORY] = h.loadMessageHistory

	return h
}
