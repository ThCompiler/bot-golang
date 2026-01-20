package botgolang

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage_AttachNewFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		fileName    string
	}{
		{
			name:        "valid_file",
			fileContent: "test content",
			fileName:    "test.txt",
		},
		{
			name:        "empty_file",
			fileContent: "",
			fileName:    "empty.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile, err := os.CreateTemp(tmpDir, tt.fileName)
			t.Cleanup(func() {
				err := tmpFile.Close()
				assert.NoError(t, err)
			})

			require.NoError(t, err)

			_, err = tmpFile.WriteString(tt.fileContent)
			require.NoError(t, err)

			msg := &Message{}
			msg.AttachNewFile(tmpFile)

			require.NotNil(t, msg.File)
			assert.Equal(t, tmpFile.Name(), msg.File.Name())
			assert.Equal(t, OtherFile, msg.ContentType)
		})
	}
}

func TestMessage_AttachExistingFile(t *testing.T) {
	tests := []struct {
		name   string
		fileID string
	}{
		{
			name:   "valid_file_id",
			fileID: "file123",
		},
		{
			name:   "empty_file_id",
			fileID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{}
			msg.AttachExistingFile(tt.fileID)

			assert.Equal(t, tt.fileID, msg.FileID)
			assert.Equal(t, OtherFile, msg.ContentType)
		})
	}
}

func TestMessage_AttachNewVoice(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		fileName    string
	}{
		{
			name:        "valid_voice_file",
			fileContent: "voice content",
			fileName:    "voice.ogg",
		},
		{
			name:        "empty_voice_file",
			fileContent: "",
			fileName:    "empty.m4a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile, err := os.CreateTemp(tmpDir, tt.fileName)
			t.Cleanup(func() {
				err := tmpFile.Close()
				assert.NoError(t, err)
			})

			require.NoError(t, err)

			_, err = tmpFile.WriteString(tt.fileContent)
			require.NoError(t, err)

			msg := &Message{}
			msg.AttachNewVoice(tmpFile)

			require.NotNil(t, msg.File)
			assert.Equal(t, tmpFile.Name(), msg.File.Name())
			assert.Equal(t, Voice, msg.ContentType)

		})
	}
}

func TestMessage_AttachExistingVoice(t *testing.T) {
	tests := []struct {
		name   string
		fileID string
	}{
		{
			name:   "valid_voice_id",
			fileID: "voice123",
		},
		{
			name:   "empty_voice_id",
			fileID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{}
			msg.AttachExistingVoice(tt.fileID)

			assert.Equal(t, tt.fileID, msg.FileID)
			assert.Equal(t, Voice, msg.ContentType)
		})
	}
}

func TestMessage_Send_WithNewFile(t *testing.T) {
	client := NewApiMockClient(t)

	tmpDir := t.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "test.txt")
	t.Cleanup(func() {
		err := tmpFile.Close()
		assert.NoError(t, err)
	})

	require.NoError(t, err)

	_, err = tmpFile.WriteString("test content")
	require.NoError(t, err)

	// Seek back to beginning for reading
	_, err = tmpFile.Seek(0, 0)
	require.NoError(t, err)

	msg := &Message{
		client:      &client,
		Chat:        Chat{ID: "chat123"},
		File:        NewUploadFileFromReader(tmpFile.Name(), tmpFile),
		ContentType: OtherFile,
	}

	err = msg.Send()
	assert.NoError(t, err)

}

func TestMessage_Send_WithExistingFile(t *testing.T) {
	client := NewApiMockClient(t)

	msg := &Message{
		client:      &client,
		Chat:        Chat{ID: "chat123"},
		FileID:      "file123",
		ContentType: OtherFile,
	}

	err := msg.Send()
	assert.NoError(t, err)
}

func TestMessage_Send_WithNewVoice(t *testing.T) {
	client := NewApiMockClient(t)

	tmpDir := t.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "voice.ogg")
	t.Cleanup(func() {
		err := tmpFile.Close()
		assert.NoError(t, err)
	})

	require.NoError(t, err)

	_, err = tmpFile.WriteString("voice content")
	require.NoError(t, err)

	// Seek back to beginning for reading
	_, err = tmpFile.Seek(0, 0)
	require.NoError(t, err)

	msg := &Message{
		client:      &client,
		Chat:        Chat{ID: "chat123"},
		File:        NewUploadFileFromReader(tmpFile.Name(), tmpFile),
		ContentType: Voice,
	}

	err = msg.Send()
	assert.NoError(t, err)

}

func TestMessage_Send_WithExistingVoice(t *testing.T) {
	client := NewApiMockClient(t)

	msg := &Message{
		client:      &client,
		Chat:        Chat{ID: "chat123"},
		FileID:      "voice123",
		ContentType: Voice,
	}

	err := msg.Send()
	assert.NoError(t, err)
}

func TestMessage_Send_UnknownContentType_AutoDetect(t *testing.T) {
	tests := []struct {
		name         string
		fileID       string
		expectedPath string
	}{
		{
			name:         "voice_file_id_starts_with_I",
			fileID:       "Ivoice123",
			expectedPath: "/messages/sendVoice",
		},
		{
			name:         "regular_file_id_does_not_start_with_I",
			fileID:       "file123",
			expectedPath: "/messages/sendFile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewApiMockClient(t)

			msg := &Message{
				client:      &client,
				Chat:        Chat{ID: "chat123"},
				FileID:      tt.fileID,
				ContentType: Unknown,
			}

			err := msg.Send()
			assert.NoError(t, err)
		})
	}
}

func TestMessage_Send_UnknownContentType_WithFile_Extension(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{
			name:     "aac_extension_detects_voice",
			fileName: "voice.aac",
		},
		{
			name:     "ogg_extension_detects_voice",
			fileName: "voice.ogg",
		},
		{
			name:     "m4a_extension_detects_voice",
			fileName: "voice.m4a",
		},
		{
			name:     "txt_extension_detects_file",
			fileName: "file.txt",
		},
		{
			name:     "pdf_extension_detects_file",
			fileName: "file.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewApiMockClient(t)

			reader := strings.NewReader("content")
			msg := &Message{
				client:      &client,
				Chat:        Chat{ID: "chat123"},
				File:        NewUploadFileFromReader(tt.fileName, reader),
				ContentType: Unknown,
			}

			err := msg.Send()
			assert.NoError(t, err)
		})
	}
}

func TestMessage_Send_Errors(t *testing.T) {
	tests := []struct {
		name        string
		message     *Message
		expectedErr string
	}{
		{
			name: "nil_client",
			message: &Message{
				Chat: Chat{ID: "chat123"},
			},
			expectedErr: "client is not inited",
		},
		{
			name: "empty_chat_id",
			message: &Message{
				client:      &Client{},
				ContentType: Text,
				Text:        "test",
			},
			expectedErr: "message should have chat id",
		},
		{
			name: "no_data",
			message: &Message{
				client:      &Client{},
				Chat:        Chat{ID: "chat123"},
				ContentType: Unknown,
			},
			expectedErr: "cannot send message or file without data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Send()

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestMessage_Send_TextMessage(t *testing.T) {
	client := NewApiMockClient(t)

	msg := &Message{
		client:      &client,
		Chat:        Chat{ID: "chat123"},
		Text:        "test message",
		ContentType: Text,
	}

	err := msg.Send()
	assert.NoError(t, err)
}
