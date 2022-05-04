package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mashu.example/internal/adapter/utils"
	"mashu.example/internal/usecase/chat/create_direct_message"
	"mashu.example/internal/usecase/repository"
	"mashu.example/pkg"
)

type requestType string

const (
	WS_CREATE_DM requestType = "CREATE_DM"
	WS_SEND_MSG  requestType = "SEND_MSG"
)

type websocketRequestMessage struct {
	Type requestType `json:"type"`
	// UserName string            `json:"userName"`
	Payload map[string]string `json:"payload"`
}

type websocketErrorResponse struct {
	ErrCode int    `json:"errCode"`
	Message string `json:"message"`
}

type websocketHandler struct {
	clients         map[uuid.UUID]*utils.WebSocketClient
	wsMsgHandlerMap map[requestType]func(uuid.UUID, map[string]string)

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
		var req websocketRequestMessage
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

func (h *websocketHandler) createDM(userId uuid.UUID, payload map[string]string) {
	logrus.Info("start create DM")
	client := h.clients[userId]
	// jsonStr, err := json.Marshal(payload)
	// if err != nil {
	// 	errMsg := "failed to parse payload"
	// 	logrus.Error(errMsg)
	// 	client.Conn.WriteJSON(websocketErrorResponse{ErrCode: 400, Message: errMsg})
	// 	return
	// }

	// req := &create_direct_message.CreateDirectMessageUseCaseReq{}
	// if err := json.Unmarshal(jsonStr, req); err != nil {
	// 	errMsg := "failed to parse payload into struct"
	// 	logrus.Error(errMsg)
	// 	client.Conn.WriteJSON(websocketErrorResponse{ErrCode: 400, Message: errMsg})
	// 	return
	// }

	senderId := userId
	receiverId := uuid.MustParse(payload["targetUserId"])
	req := create_direct_message.NewCreateDirectMessageUseCaseReq(
		senderId,
		receiverId,
	)

	if !req.Validate() {
		errMsg := "bad request"
		logrus.Error(errMsg)
		client.Conn.WriteJSON(websocketErrorResponse{ErrCode: 400, Message: errMsg})
		return
	}

	res := create_direct_message.NewCreateDirectMessageUseCaseRes()
	uc := create_direct_message.NewCreateDirectMessageUseCase(
		h.chatRepo,
		h.userRepo,
		req,
		res,
	)
	uc.Execute()
	if res.Err != nil {
		fmt.Println(res.Err.Error())
		return
	}

	dm, err := h.chatRepo.GetDMByUserId(senderId, receiverId)
	if err != nil {
		fmt.Println(err.Error())
		client.Conn.WriteJSON("failed to get created DM")
		return
	}
	fmt.Printf("chatroom got! %+v\n", dm)

	client.Conn.WriteJSON("succeeded")
	logrus.Info("end create DM")
}

func (h *websocketHandler) sendMessage(userId uuid.UUID, payload map[string]string) {
	fmt.Println("send message")
}

func newWebSocketHandler(
	userRepo repository.UserRepo,
	chatRepo repository.ChatRepo,
) *websocketHandler {
	h := &websocketHandler{
		clients:         map[uuid.UUID]*utils.WebSocketClient{},
		wsMsgHandlerMap: map[requestType]func(uuid.UUID, map[string]string){},
		userRepo:        userRepo,
		chatRepo:        chatRepo,
	}

	h.wsMsgHandlerMap[WS_CREATE_DM] = h.createDM
	h.wsMsgHandlerMap[WS_SEND_MSG] = h.sendMessage

	return h
}
