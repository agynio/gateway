package handlers

import (
	"context"
	"fmt"

	chatv1 "github.com/agynio/gateway/gen/agynio/api/chat/v1"
	"github.com/agynio/gateway/internal/chatgen"
)

const chatBasePath = "/chat/v1"

func ChatBasePath() string {
	return chatBasePath
}

// ChatHandler implements chatgen.StrictServerInterface for the Chat API.
type ChatHandler struct {
	client chatv1.ChatServiceClient
}

func NewChatHandler(client chatv1.ChatServiceClient) *ChatHandler {
	if client == nil {
		panic("chat client is required")
	}
	return &ChatHandler{client: client}
}

var _ chatgen.StrictServerInterface = (*ChatHandler)(nil)

func (h *ChatHandler) GetChats(ctx context.Context, request chatgen.GetChatsRequestObject) (chatgen.GetChatsResponseObject, error) {
	resp, err := h.client.GetChats(ctx, &chatv1.GetChatsRequest{
		PageSize:  pageSizeFromParam(request.Params.PageSize),
		PageToken: stringValue(request.Params.PageToken),
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items := make([]chatgen.Chat, 0, len(resp.GetChats()))
	for _, chat := range resp.GetChats() {
		converted, err := chatFromProto(chat)
		if err != nil {
			return nil, responseProblem(err)
		}
		items = append(items, converted)
	}

	payload := chatgen.PaginatedChats{
		Items:         items,
		NextPageToken: stringPtr(resp.GetNextPageToken()),
	}

	return chatgen.GetChats200JSONResponse(payload), nil
}

func (h *ChatHandler) CreateChat(ctx context.Context, request chatgen.CreateChatRequestObject) (chatgen.CreateChatResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := h.client.CreateChat(ctx, chatCreateToProto(*request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	chat := resp.GetChat()
	if chat == nil {
		return nil, responseProblem(fmt.Errorf("create chat response missing chat"))
	}

	converted, err := chatFromProto(chat)
	if err != nil {
		return nil, responseProblem(err)
	}

	return chatgen.CreateChat201JSONResponse(converted), nil
}

func (h *ChatHandler) GetChatMessages(ctx context.Context, request chatgen.GetChatMessagesRequestObject) (chatgen.GetChatMessagesResponseObject, error) {
	resp, err := h.client.GetMessages(ctx, &chatv1.GetMessagesRequest{
		ChatId:    request.ChatId.String(),
		PageSize:  pageSizeFromParam(request.Params.PageSize),
		PageToken: stringValue(request.Params.PageToken),
	})
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	items, err := chatMessagesFromProto(resp.GetMessages())
	if err != nil {
		return nil, responseProblem(err)
	}

	payload := chatgen.PaginatedChatMessages{
		Items:         items,
		NextPageToken: stringPtr(resp.GetNextPageToken()),
		UnreadCount:   int(resp.GetUnreadCount()),
	}

	return chatgen.GetChatMessages200JSONResponse(payload), nil
}

func (h *ChatHandler) SendChatMessage(ctx context.Context, request chatgen.SendChatMessageRequestObject) (chatgen.SendChatMessageResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	messageRequest, err := chatMessageCreateToProto(request.ChatId, *request.Body)
	if err != nil {
		return nil, requestProblem(err)
	}

	resp, err := h.client.SendMessage(ctx, messageRequest)
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	message := resp.GetMessage()
	if message == nil {
		return nil, responseProblem(fmt.Errorf("send message response missing message"))
	}

	converted, err := chatMessageFromProto(message)
	if err != nil {
		return nil, responseProblem(err)
	}

	return chatgen.SendChatMessage201JSONResponse(converted), nil
}

func (h *ChatHandler) MarkChatAsRead(ctx context.Context, request chatgen.MarkChatAsReadRequestObject) (chatgen.MarkChatAsReadResponseObject, error) {
	if request.Body == nil {
		panic("validated request body is unexpectedly nil")
	}

	resp, err := h.client.MarkAsRead(ctx, markAsReadToProto(request.ChatId, *request.Body))
	if err != nil {
		return nil, grpcErrorToProblem(err)
	}

	payload := chatgen.MarkAsReadResponse{
		ReadCount: int(resp.GetReadCount()),
	}

	return chatgen.MarkChatAsRead200JSONResponse(payload), nil
}
