package handlers

import (
	"fmt"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	chatv1 "github.com/agynio/gateway/gen/agynio/api/chat/v1"
	"github.com/agynio/gateway/internal/chatgen"
)

func chatFromProto(chat *chatv1.Chat) (chatgen.Chat, error) {
	if chat == nil {
		return chatgen.Chat{}, fmt.Errorf("chat missing")
	}

	chatID, err := parseUUID(chat.GetId())
	if err != nil {
		return chatgen.Chat{}, fmt.Errorf("parse chat id: %w", err)
	}

	participants, err := chatParticipantsFromProto(chat.GetParticipants())
	if err != nil {
		return chatgen.Chat{}, err
	}

	createdAt := chat.GetCreatedAt()
	if createdAt == nil {
		return chatgen.Chat{}, fmt.Errorf("chat created_at missing")
	}

	var updatedAt *time.Time
	if updated := chat.GetUpdatedAt(); updated != nil {
		value := updated.AsTime().UTC()
		updatedAt = &value
	}

	return chatgen.Chat{
		Id:           chatID,
		Participants: participants,
		CreatedAt:    createdAt.AsTime().UTC(),
		UpdatedAt:    updatedAt,
	}, nil
}

func chatParticipantsFromProto(participants []*chatv1.ChatParticipant) ([]chatgen.ChatParticipant, error) {
	if len(participants) == 0 {
		return []chatgen.ChatParticipant{}, nil
	}

	items := make([]chatgen.ChatParticipant, 0, len(participants))
	for _, participant := range participants {
		converted, err := chatParticipantFromProto(participant)
		if err != nil {
			return nil, err
		}
		items = append(items, converted)
	}

	return items, nil
}

func chatParticipantFromProto(participant *chatv1.ChatParticipant) (chatgen.ChatParticipant, error) {
	if participant == nil {
		return chatgen.ChatParticipant{}, fmt.Errorf("chat participant missing")
	}

	participantID, err := parseUUID(participant.GetId())
	if err != nil {
		return chatgen.ChatParticipant{}, fmt.Errorf("parse participant id: %w", err)
	}

	joinedAt := participant.GetJoinedAt()
	if joinedAt == nil {
		return chatgen.ChatParticipant{}, fmt.Errorf("participant joined_at missing")
	}

	return chatgen.ChatParticipant{
		Id:       participantID,
		JoinedAt: joinedAt.AsTime().UTC(),
	}, nil
}

func chatMessageFromProto(message *chatv1.ChatMessage) (chatgen.ChatMessage, error) {
	if message == nil {
		return chatgen.ChatMessage{}, fmt.Errorf("chat message missing")
	}

	messageID, err := parseUUID(message.GetId())
	if err != nil {
		return chatgen.ChatMessage{}, fmt.Errorf("parse message id: %w", err)
	}

	chatID, err := parseUUID(message.GetChatId())
	if err != nil {
		return chatgen.ChatMessage{}, fmt.Errorf("parse chat id: %w", err)
	}

	senderID, err := parseUUID(message.GetSenderId())
	if err != nil {
		return chatgen.ChatMessage{}, fmt.Errorf("parse sender id: %w", err)
	}

	fileIDs, err := uuidSliceFromStrings(message.GetFileIds())
	if err != nil {
		return chatgen.ChatMessage{}, fmt.Errorf("parse file ids: %w", err)
	}

	createdAt := message.GetCreatedAt()
	if createdAt == nil {
		return chatgen.ChatMessage{}, fmt.Errorf("message created_at missing")
	}

	return chatgen.ChatMessage{
		Id:        messageID,
		ChatId:    chatID,
		SenderId:  senderID,
		Body:      message.GetBody(),
		FileIds:   fileIDs,
		CreatedAt: createdAt.AsTime().UTC(),
	}, nil
}

func chatMessagesFromProto(messages []*chatv1.ChatMessage) ([]chatgen.ChatMessage, error) {
	if len(messages) == 0 {
		return []chatgen.ChatMessage{}, nil
	}

	items := make([]chatgen.ChatMessage, 0, len(messages))
	for _, message := range messages {
		converted, err := chatMessageFromProto(message)
		if err != nil {
			return nil, err
		}
		items = append(items, converted)
	}

	return items, nil
}

func chatCreateToProto(request chatgen.ChatCreateRequest) *chatv1.CreateChatRequest {
	return &chatv1.CreateChatRequest{
		ParticipantIds: uuidSliceToStrings(request.ParticipantIds),
	}
}

func chatMessageCreateToProto(chatID openapi_types.UUID, request chatgen.ChatMessageCreateRequest) (*chatv1.SendMessageRequest, error) {
	body := stringValue(request.Body)
	fileIDs := []string{}
	if request.FileIds != nil {
		fileIDs = uuidSliceToStrings(*request.FileIds)
	}
	if body == "" && len(fileIDs) == 0 {
		return nil, fmt.Errorf("message body or fileIds required")
	}

	return &chatv1.SendMessageRequest{
		ChatId:  chatID.String(),
		Body:    body,
		FileIds: fileIDs,
	}, nil
}

func markAsReadToProto(chatID openapi_types.UUID, request chatgen.MarkAsReadRequest) *chatv1.MarkAsReadRequest {
	return &chatv1.MarkAsReadRequest{
		ChatId:     chatID.String(),
		MessageIds: uuidSliceToStrings(request.MessageIds),
	}
}

func uuidSliceFromStrings(values []string) ([]openapi_types.UUID, error) {
	items := make([]openapi_types.UUID, 0, len(values))
	for _, value := range values {
		parsed, err := parseUUID(value)
		if err != nil {
			return nil, err
		}
		items = append(items, parsed)
	}
	return items, nil
}

func uuidSliceToStrings(values []openapi_types.UUID) []string {
	items := make([]string, 0, len(values))
	for _, value := range values {
		items = append(items, value.String())
	}
	return items
}
