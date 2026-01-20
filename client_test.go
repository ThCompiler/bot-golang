package botgolang

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_Do_OK(t *testing.T) {
	client := NewApiMockClient(t)

	bytes, err := client.Do("/", url.Values{}, nil)

	require.NoError(t, err)
	require.JSONEq(t, `{"ok":true}`, string(bytes))
}

func TestClient_Do_Error(t *testing.T) {
	client := NewApiMockClient(t)
	client.token = "" // Clear token to trigger error

	expected := `{"ok":false, "description":"Missing required parameter 'token'"}`

	bytes, err := client.Do("/", url.Values{}, nil)

	require.EqualError(t, err, "error status from API: Missing required parameter 'token'")
	require.JSONEq(t, expected, string(bytes))
}

func TestClient_GetEvents_OK(t *testing.T) {
	expected := []*Event{
		{
			EventID: 1,
			Type:    NEW_MESSAGE,
			Payload: EventPayload{
				BaseEventPayload: BaseEventPayload{
					MsgID: "57883346846815030",
					Chat: Chat{
						ID:    "681869378@chat.agent",
						Type:  "channel",
						Title: "The best channel",
					},
					From: Contact{
						User:      User{"1234567890"},
						FirstName: "Name",
						LastName:  "SurName",
					},
					Text:      "Hello!",
					Timestamp: 1546290000,
				},
				Parts: []Part{
					{
						Type: STICKER,
						Payload: PartPayload{
							FileID: "2IWuJzaNWCJZxJWCvZhDYuJ5XDsr7hU",
						},
					},
					{
						Type: MENTION,
						Payload: PartPayload{
							FirstName: "Name",
							LastName:  "SurName",
							UserID:    "1234567890",
						},
					},
					{
						Type: VOICE,
						Payload: PartPayload{
							FileID: "IdjUEXuGdNhLKUfD5rvkE03IOax54cD",
						},
					},
					{
						Type: FILE,
						Payload: PartPayload{
							FileID:  "ZhSnMuaOmF7FRez2jGWuQs5zGZwlLa0",
							Caption: "Last weekend trip",
							Type:    "image",
						},
					},
					{
						Type: FORWARD,
						Payload: PartPayload{
							PartMessage: PartMessage{
								MsgID: "12354",
								Text:  "test1",
							},
						},
					},
					{
						Type: REPLY,
						Payload: PartPayload{
							PartMessage: PartMessage{
								MsgID: "12354",
								Text:  "test",
							},
						},
					},
				},
			},
		},
		{
			EventID: 2,
			Type:    EDITED_MESSAGE,
			Payload: EventPayload{
				BaseEventPayload: BaseEventPayload{
					MsgID: "57883346846815030",
					Chat: Chat{
						ID:    "681869378@chat.agent",
						Type:  "channel",
						Title: "The best channel",
					},
					From: Contact{
						User:      User{"1234567890"},
						FirstName: "Name",
						LastName:  "SurName",
					},
					Text:      "Hello!",
					Timestamp: 1546290000,
				},
			},
		},
		{
			EventID: 3,
			Type:    DELETED_MESSAGE,
			Payload: EventPayload{
				BaseEventPayload: BaseEventPayload{
					MsgID: "57883346846815030",
					Chat: Chat{
						ID:    "681869378@chat.agent",
						Type:  "channel",
						Title: "The best channel",
					},
					Timestamp: 1546290000,
				},
			},
		},
		{
			EventID: 4,
			Type:    PINNED_MESSAGE,
			Payload: EventPayload{
				BaseEventPayload: BaseEventPayload{
					MsgID: "6720509406122810000",
					Chat: Chat{
						ID:    "681869378@chat.agent",
						Type:  "group",
						Title: "The best group",
					},
					From: Contact{
						User:      User{"9876543210"},
						FirstName: "Name",
						LastName:  "SurName",
					},
					Text:      "Some important information!",
					Timestamp: 1564740530,
				},
			},
		},
		{
			EventID: 5,
			Type:    UNPINNED_MESSAGE,
			Payload: EventPayload{
				BaseEventPayload: BaseEventPayload{
					MsgID: "6720509406122810000",
					Chat: Chat{
						ID:    "681869378@chat.agent",
						Type:  "group",
						Title: "The best group",
					},
					Timestamp: 1564740530,
				},
			},
		},
		{
			EventID: 6,
			Type:    NEW_CHAT_MEMBERS,
			Payload: EventPayload{
				BaseEventPayload: BaseEventPayload{
					Chat: Chat{
						ID:    "681869378@chat.agent",
						Type:  "group",
						Title: "The best group",
					},
				},
				NewMembers: []Contact{
					{
						User:      User{"1234567890"},
						FirstName: "Name",
						LastName:  "SurName",
					},
				},
				AddedBy: Contact{
					User:      User{"9876543210"},
					FirstName: "Name",
					LastName:  "SurName",
				},
			},
		},
		{
			EventID: 7,
			Type:    LEFT_CHAT_MEMBERS,
			Payload: EventPayload{
				BaseEventPayload: BaseEventPayload{
					Chat: Chat{
						ID:    "681869378@chat.agent",
						Type:  "group",
						Title: "The best group",
					},
				},
				LeftMembers: []Contact{
					{
						User:      User{"1234567890"},
						FirstName: "Name",
						LastName:  "SurName",
					},
				},
				RemovedBy: Contact{
					User:      User{"9876543210"},
					FirstName: "Name",
					LastName:  "SurName",
				},
			},
		},
		{
			EventID: 8,
			Type:    CALLBACK_QUERY,
			Payload: EventPayload{
				CallbackData: "echo",
				CallbackMsg: BaseEventPayload{
					MsgID: "6720509406122810000",
					Chat: Chat{
						ID:   "1234567890",
						Type: "private",
					},
					From: Contact{
						User:      User{"bot_id"},
						FirstName: "bot_name",
					},
					Text:      "Some important information!",
					Timestamp: 1564740530,
				},
				BaseEventPayload: BaseEventPayload{
					From: Contact{
						User:      User{"1234567890"},
						FirstName: "Name",
					},
				},
				QueryID: "SVR:123456",
			},
		},
	}

	client := NewApiMockClient(t)

	events, err := client.GetEvents(0, 0)

	require.NoError(t, err)
	require.Equal(t, events, expected)
}

func TestClient_GetInfo_OK(t *testing.T) {
	client := NewApiMockClient(t)

	info, err := client.GetChatInfo("id_1234")
	require.NoError(t, err)
	require.NotEmpty(t, info.ID)
}

func TestClient_GetInfo_Error(t *testing.T) {
	require := require.New(t)

	require.NoError(nil)
}

func TestClient_Do_WithFileUpload(t *testing.T) {
	client := NewApiMockClient(t)

	content := strings.NewReader("test file content")
	uploadFile := NewUploadFileFromReader("test.txt", content)

	bytes, err := client.Do("/messages/sendFile", url.Values{}, uploadFile)

	require.NoError(t, err)
	require.JSONEq(t, `{"ok":true,"msgId":"test123","timestamp":123456}`, string(bytes))
}

func TestClient_DoWithContext_WithFileUpload(t *testing.T) {
	client := NewApiMockClient(t)

	content := strings.NewReader("test voice content")
	uploadFile := NewUploadFileFromReader("voice.ogg", content)

	bytes, err := client.DoWithContext(context.Background(), "/messages/sendVoice", url.Values{}, uploadFile)

	require.NoError(t, err)
	require.JSONEq(t, `{"ok":true,"msgId":"voice123","timestamp":123456}`, string(bytes))
}

func TestClient_Do_WithEmptyFile_Error(t *testing.T) {
	client := NewApiMockClient(t)

	content := strings.NewReader("")
	uploadFile := NewUploadFileFromReader("empty.txt", content)

	bytes, err := client.Do("/messages/sendFile", url.Values{}, uploadFile)

	require.JSONEq(t, `{"ok":false,"description":"File cannot be empty"}`, string(bytes))
	require.EqualError(t, err, "error status from API: File cannot be empty")
}

func TestClient_Do_GET_WithoutFileID_Error(t *testing.T) {
	client := NewApiMockClient(t)

	bytes, err := client.Do("/messages/sendFile", url.Values{}, nil)

	require.JSONEq(t, `{"ok":false,"description":"Missing required parameter 'fileId'"}`, string(bytes))
	require.EqualError(t, err, "error status from API: Missing required parameter 'fileId'")
}

func TestClient_Do_SendVoice_Validation(t *testing.T) {
	tests := []struct {
		name     string
		file     UploadFile
		expected string
	}{
		{
			name:     "without_file",
			file:     nil,
			expected: `{"ok":false,"description":"Missing required parameter 'fileId'"}`,
		},
		{
			name:     "with_empty_file",
			file:     NewUploadFileFromReader("empty.ogg", strings.NewReader("")),
			expected: `{"ok":false,"description":"File cannot be empty"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := NewApiMockClient(t)

			bytes, err := client.Do("/messages/sendVoice", url.Values{}, tc.file)

			require.JSONEq(t, tc.expected, string(bytes))
			require.Error(t, err)
		})
	}
}
