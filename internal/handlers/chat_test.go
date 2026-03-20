package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	chatv1 "github.com/agynio/gateway/gen/agynio/api/chat/v1"
	"github.com/agynio/gateway/internal/chatgen"
)

type chatCall struct {
	method string
	assert func(any)
	resp   any
	err    error
}

type stubChatClient struct {
	t     *testing.T
	calls []chatCall
	idx   int
}

func (s *stubChatClient) Expect(call chatCall) {
	s.calls = append(s.calls, call)
}

func (s *stubChatClient) AssertDone() {
	if s.idx != len(s.calls) {
		s.t.Fatalf("expected %d calls, got %d", len(s.calls), s.idx)
	}
}

func (s *stubChatClient) nextCall(method string, req any) chatCall {
	if s.idx >= len(s.calls) {
		s.t.Fatalf("unexpected call %s", method)
	}
	call := s.calls[s.idx]
	s.idx++
	if call.method != method {
		s.t.Fatalf("unexpected method: got %s want %s", method, call.method)
	}
	if call.assert != nil {
		call.assert(req)
	}
	return call
}

func (s *stubChatClient) CreateChat(ctx context.Context, req *chatv1.CreateChatRequest, _ ...grpc.CallOption) (*chatv1.CreateChatResponse, error) {
	call := s.nextCall("CreateChat", req)
	resp, ok := call.resp.(*chatv1.CreateChatResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubChatClient) GetChats(ctx context.Context, req *chatv1.GetChatsRequest, _ ...grpc.CallOption) (*chatv1.GetChatsResponse, error) {
	call := s.nextCall("GetChats", req)
	resp, ok := call.resp.(*chatv1.GetChatsResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubChatClient) GetMessages(ctx context.Context, req *chatv1.GetMessagesRequest, _ ...grpc.CallOption) (*chatv1.GetMessagesResponse, error) {
	call := s.nextCall("GetMessages", req)
	resp, ok := call.resp.(*chatv1.GetMessagesResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubChatClient) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest, _ ...grpc.CallOption) (*chatv1.SendMessageResponse, error) {
	call := s.nextCall("SendMessage", req)
	resp, ok := call.resp.(*chatv1.SendMessageResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func (s *stubChatClient) MarkAsRead(ctx context.Context, req *chatv1.MarkAsReadRequest, _ ...grpc.CallOption) (*chatv1.MarkAsReadResponse, error) {
	call := s.nextCall("MarkAsRead", req)
	resp, ok := call.resp.(*chatv1.MarkAsReadResponse)
	if call.resp != nil && !ok {
		s.t.Fatalf("unexpected response type: %T", call.resp)
	}
	return resp, call.err
}

func TestChatGetChats(t *testing.T) {
	stub := &stubChatClient{t: t}
	chatID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	participantID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	createdAt := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	joinedAt := time.Date(2024, 1, 2, 3, 4, 6, 0, time.UTC)
	updatedAt := time.Date(2024, 1, 3, 4, 5, 6, 0, time.UTC)

	stub.Expect(chatCall{
		method: "GetChats",
		assert: func(req any) {
			input := req.(*chatv1.GetChatsRequest)
			if input.PageSize != 5 {
				t.Fatalf("unexpected page size: %d", input.PageSize)
			}
			if input.PageToken != "cursor" {
				t.Fatalf("unexpected page token: %s", input.PageToken)
			}
		},
		resp: &chatv1.GetChatsResponse{
			Chats: []*chatv1.Chat{
				{
					Id: chatID.String(),
					Participants: []*chatv1.ChatParticipant{
						{Id: participantID.String(), JoinedAt: timestamppb.New(joinedAt)},
					},
					CreatedAt: timestamppb.New(createdAt),
					UpdatedAt: timestamppb.New(updatedAt),
				},
			},
			NextPageToken: "next",
		},
	})

	pageSize := 5
	pageToken := "cursor"
	h := NewChatHandler(stub)
	resp, err := h.GetChats(context.Background(), chatgen.GetChatsRequestObject{Params: chatgen.GetChatsParams{
		PageSize:  &pageSize,
		PageToken: &pageToken,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, ok := resp.(chatgen.GetChats200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if list.NextPageToken == nil || *list.NextPageToken != "next" {
		t.Fatalf("unexpected next page token: %v", list.NextPageToken)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 chat, got %d", len(list.Items))
	}
	item := list.Items[0]
	if item.Id != openapi_types.UUID(chatID) {
		t.Fatalf("unexpected chat id: %s", item.Id.String())
	}
	if !item.CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected createdAt: %v", item.CreatedAt)
	}
	if item.UpdatedAt == nil || !item.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("unexpected updatedAt: %v", item.UpdatedAt)
	}
	if len(item.Participants) != 1 {
		t.Fatalf("expected 1 participant, got %d", len(item.Participants))
	}
	participant := item.Participants[0]
	if participant.Id != openapi_types.UUID(participantID) {
		t.Fatalf("unexpected participant id: %s", participant.Id.String())
	}
	if !participant.JoinedAt.Equal(joinedAt) {
		t.Fatalf("unexpected joinedAt: %v", participant.JoinedAt)
	}

	stub.AssertDone()
}

func TestChatCreateChat(t *testing.T) {
	stub := &stubChatClient{t: t}
	chatID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	participantID := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	createdAt := time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)

	stub.Expect(chatCall{
		method: "CreateChat",
		assert: func(req any) {
			input := req.(*chatv1.CreateChatRequest)
			if len(input.ParticipantIds) != 1 || input.ParticipantIds[0] != participantID.String() {
				t.Fatalf("unexpected participant ids: %v", input.ParticipantIds)
			}
		},
		resp: &chatv1.CreateChatResponse{
			Chat: &chatv1.Chat{
				Id: chatID.String(),
				Participants: []*chatv1.ChatParticipant{
					{Id: participantID.String(), JoinedAt: timestamppb.New(createdAt)},
				},
				CreatedAt: timestamppb.New(createdAt),
			},
		},
	})

	request := chatgen.CreateChatRequestObject{Body: &chatgen.ChatCreateRequest{
		ParticipantIds: []openapi_types.UUID{openapi_types.UUID(participantID)},
	}}

	h := NewChatHandler(stub)
	resp, err := h.CreateChat(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(chatgen.CreateChat201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Id != openapi_types.UUID(chatID) {
		t.Fatalf("unexpected chat id: %s", created.Id.String())
	}
	if len(created.Participants) != 1 {
		t.Fatalf("expected 1 participant, got %d", len(created.Participants))
	}

	stub.AssertDone()
}

func TestChatGetMessages(t *testing.T) {
	stub := &stubChatClient{t: t}
	chatID := uuid.MustParse("55555555-5555-5555-5555-555555555555")
	messageID := uuid.MustParse("66666666-6666-6666-6666-666666666666")
	senderID := uuid.MustParse("77777777-7777-7777-7777-777777777777")
	fileID := uuid.MustParse("88888888-8888-8888-8888-888888888888")
	createdAt := time.Date(2024, 7, 8, 9, 10, 11, 0, time.UTC)

	stub.Expect(chatCall{
		method: "GetMessages",
		assert: func(req any) {
			input := req.(*chatv1.GetMessagesRequest)
			if input.ChatId != chatID.String() {
				t.Fatalf("unexpected chat id: %s", input.ChatId)
			}
			if input.PageSize != 3 {
				t.Fatalf("unexpected page size: %d", input.PageSize)
			}
			if input.PageToken != "token" {
				t.Fatalf("unexpected page token: %s", input.PageToken)
			}
		},
		resp: &chatv1.GetMessagesResponse{
			Messages: []*chatv1.ChatMessage{
				{
					Id:        messageID.String(),
					ChatId:    chatID.String(),
					SenderId:  senderID.String(),
					Body:      "hello",
					FileIds:   []string{fileID.String()},
					CreatedAt: timestamppb.New(createdAt),
				},
			},
			NextPageToken: "next",
			UnreadCount:   2,
		},
	})

	pageSize := 3
	pageToken := "token"
	h := NewChatHandler(stub)
	resp, err := h.GetChatMessages(context.Background(), chatgen.GetChatMessagesRequestObject{
		ChatId: openapi_types.UUID(chatID),
		Params: chatgen.GetChatMessagesParams{PageSize: &pageSize, PageToken: &pageToken},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, ok := resp.(chatgen.GetChatMessages200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if list.UnreadCount != 2 {
		t.Fatalf("unexpected unread count: %d", list.UnreadCount)
	}
	if list.NextPageToken == nil || *list.NextPageToken != "next" {
		t.Fatalf("unexpected next page token: %v", list.NextPageToken)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 message, got %d", len(list.Items))
	}
	message := list.Items[0]
	if message.Id != openapi_types.UUID(messageID) {
		t.Fatalf("unexpected message id: %s", message.Id.String())
	}
	if message.ChatId != openapi_types.UUID(chatID) {
		t.Fatalf("unexpected chat id: %s", message.ChatId.String())
	}
	if message.SenderId != openapi_types.UUID(senderID) {
		t.Fatalf("unexpected sender id: %s", message.SenderId.String())
	}
	if message.Body != "hello" {
		t.Fatalf("unexpected body: %s", message.Body)
	}
	if len(message.FileIds) != 1 || message.FileIds[0] != openapi_types.UUID(fileID) {
		t.Fatalf("unexpected file ids: %v", message.FileIds)
	}
	if !message.CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected createdAt: %v", message.CreatedAt)
	}

	stub.AssertDone()
}

func TestChatSendMessage(t *testing.T) {
	stub := &stubChatClient{t: t}
	chatID := uuid.MustParse("99999999-9999-9999-9999-999999999999")
	messageID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	senderID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	fileID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	createdAt := time.Date(2024, 8, 9, 10, 11, 12, 0, time.UTC)
	body := "ping"
	fileIDs := []openapi_types.UUID{openapi_types.UUID(fileID)}

	stub.Expect(chatCall{
		method: "SendMessage",
		assert: func(req any) {
			input := req.(*chatv1.SendMessageRequest)
			if input.ChatId != chatID.String() {
				t.Fatalf("unexpected chat id: %s", input.ChatId)
			}
			if input.Body != body {
				t.Fatalf("unexpected body: %s", input.Body)
			}
			if len(input.FileIds) != 1 || input.FileIds[0] != fileID.String() {
				t.Fatalf("unexpected file ids: %v", input.FileIds)
			}
		},
		resp: &chatv1.SendMessageResponse{
			Message: &chatv1.ChatMessage{
				Id:        messageID.String(),
				ChatId:    chatID.String(),
				SenderId:  senderID.String(),
				Body:      body,
				FileIds:   []string{fileID.String()},
				CreatedAt: timestamppb.New(createdAt),
			},
		},
	})

	request := chatgen.SendChatMessageRequestObject{
		ChatId: openapi_types.UUID(chatID),
		Body: &chatgen.ChatMessageCreateRequest{
			Body:    &body,
			FileIds: &fileIDs,
		},
	}

	h := NewChatHandler(stub)
	resp, err := h.SendChatMessage(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := resp.(chatgen.SendChatMessage201JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if created.Id != openapi_types.UUID(messageID) {
		t.Fatalf("unexpected message id: %s", created.Id.String())
	}
	if len(created.FileIds) != 1 || created.FileIds[0] != openapi_types.UUID(fileID) {
		t.Fatalf("unexpected file ids: %v", created.FileIds)
	}

	stub.AssertDone()
}

func TestChatMarkAsRead(t *testing.T) {
	stub := &stubChatClient{t: t}
	chatID := uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")
	messageID := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")

	stub.Expect(chatCall{
		method: "MarkAsRead",
		assert: func(req any) {
			input := req.(*chatv1.MarkAsReadRequest)
			if input.ChatId != chatID.String() {
				t.Fatalf("unexpected chat id: %s", input.ChatId)
			}
			if len(input.MessageIds) != 1 || input.MessageIds[0] != messageID.String() {
				t.Fatalf("unexpected message ids: %v", input.MessageIds)
			}
		},
		resp: &chatv1.MarkAsReadResponse{ReadCount: 1},
	})

	request := chatgen.MarkChatAsReadRequestObject{
		ChatId: openapi_types.UUID(chatID),
		Body: &chatgen.MarkAsReadRequest{
			MessageIds: []openapi_types.UUID{openapi_types.UUID(messageID)},
		},
	}

	h := NewChatHandler(stub)
	resp, err := h.MarkChatAsRead(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ack, ok := resp.(chatgen.MarkChatAsRead200JSONResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if ack.ReadCount != 1 {
		t.Fatalf("unexpected read count: %d", ack.ReadCount)
	}

	stub.AssertDone()
}
