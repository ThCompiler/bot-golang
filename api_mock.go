package botgolang

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

type MockHandler struct {
	http.Handler
	logger *logrus.Logger
}

func (h *MockHandler) SendMessage(w http.ResponseWriter) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&Response{
		OK: true,
	})

	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("cannot encode json")
	}
}

func (h *MockHandler) sendErrorResponse(w http.ResponseWriter, description string) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&Response{
		OK:          false,
		Description: description,
	})

	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("cannot encode json")
	}
}

func (h *MockHandler) TokenError(w http.ResponseWriter) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&Response{
		OK:          false,
		Description: "Missing required parameter 'token'",
	})

	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("cannot encode json")
	}
}

func (h *MockHandler) GetEvents(w http.ResponseWriter) {
	events := `
		{
			"ok": true,
		  "events": [
			{
			  "eventId": 1,
			  "type": "newMessage",
			  "payload": {
				"msgId": "57883346846815030",
				"chat": {
				  "chatId": "681869378@chat.agent",
				  "type": "channel",
				  "title": "The best channel"
				},
				"from": {
				  "userId": "1234567890",
				  "firstName": "Name",
				  "lastName": "SurName"
				},
				"timestamp": 1546290000,
				"text": "Hello!",
				"parts": [
				  {
					"type": "sticker",
					"payload": {
					  "fileId": "2IWuJzaNWCJZxJWCvZhDYuJ5XDsr7hU"
					}
				  },
				  {
					"type": "mention",
					"payload": {
					  "userId": "1234567890",
					  "firstName": "Name",
					  "lastName": "SurName"
					}
				  },
				  {
					"type": "voice",
					"payload": {
					  "fileId": "IdjUEXuGdNhLKUfD5rvkE03IOax54cD"
					}
				  },
				  {
					"type": "file",
					"payload": {
					  "fileId": "ZhSnMuaOmF7FRez2jGWuQs5zGZwlLa0",
					  "type": "image",
					  "caption": "Last weekend trip"
					}
				  },
				  {
					"type": "forward",
					"payload": {
					  "message": {
						"msgId": "12354",
						"text": "test1"
						}
					}
				  },
				  {
					"type": "reply",
					"payload": {
					  "message": {
						"msgId": "12354",
						"text": "test"
						}
					}
				  }
				]
			  }
			},
			{
			  "eventId": 2,
			  "type": "editedMessage",
			  "payload": {
				"msgId": "57883346846815030",
				"chat": {
				  "chatId": "681869378@chat.agent",
				  "type": "channel",
				  "title": "The best channel"
				},
				"from": {
				  "userId": "1234567890",
				  "firstName": "Name",
				  "lastName": "SurName"
				},
				"timestamp": 1546290000,
				"text": "Hello!",
				"editedTimestamp": 1546290099
			  }
			},
			{
			  "eventId": 3,
			  "type": "deletedMessage",
			  "payload": {
				"msgId": "57883346846815030",
				"chat": {
				  "chatId": "681869378@chat.agent",
				  "type": "channel",
				  "title": "The best channel"
				},
				"timestamp": 1546290000
			  }
			},
			{
			  "eventId": 4,
			  "type": "pinnedMessage",
			  "payload": {
				"chat": {
				  "chatId": "681869378@chat.agent",
				  "type": "group",
				  "title": "The best group"
				},
				"from": {
				  "userId": "9876543210",
				  "firstName": "Name",
				  "lastName": "SurName"
				},
				"msgId": "6720509406122810000",
				"text": "Some important information!",
				"timestamp": 1564740530
			  }
			},
			{
			  "eventId": 5,
			  "type": "unpinnedMessage",
			  "payload": {
				"chat": {
				  "chatId": "681869378@chat.agent",
				  "type": "group",
				  "title": "The best group"
				},
				"msgId": "6720509406122810000",
				"timestamp": 1564740530
			  }
			},
			{
			  "eventId": 6,
			  "type": "newChatMembers",
			  "payload": {
				"chat": {
				  "chatId": "681869378@chat.agent",
				  "type": "group",
				  "title": "The best group"
				},
				"newMembers": [
				  {
					"userId": "1234567890",
					"firstName": "Name",
					"lastName": "SurName"
				  }
				],
				"addedBy": {
				  "userId": "9876543210",
				  "firstName": "Name",
				  "lastName": "SurName"
				}
			  }
			},
			{
			  "eventId": 7,
			  "type": "leftChatMembers",
			  "payload": {
				"chat": {
				  "chatId": "681869378@chat.agent",
				  "type": "group",
				  "title": "The best group"
				},
				"leftMembers": [
				  {
					"userId": "1234567890",
					"firstName": "Name",
					"lastName": "SurName"
				  }
				],
				"removedBy": {
				  "userId": "9876543210",
				  "firstName": "Name",
				  "lastName": "SurName"
				}
			  }
	        },
	        {
	            "eventId": 8,
	            "payload": {
	                "callbackData": "echo",
	                "from": {
					  "firstName": "Name",
					  "userId": "1234567890"
	                },
					"message": {
	          			"chat": {
	            			"chatId": "1234567890",
	            			"type": "private"
	          			},
	          			"from": {
	            			"firstName": "bot_name",
	            			"nick": "bot_nick",
	            			"userId": "bot_id"
	          			},
						"msgId": "6720509406122810000",
	          			"text": "Some important information!",
						"timestamp": 1564740530
	        		},
	                "queryId": "SVR:123456"
	            },
	            "type": "callbackQuery"
			}
		  ]
		}
	`

	_, err := w.Write([]byte(events))
	if err != nil {
		h.logger.Fatal("failed to write events")
	}
}

func (h *MockHandler) SelfGet(w http.ResponseWriter, r *http.Request) {

	encoder := json.NewEncoder(w)

	if r.FormValue("chatId") == "" {
		err := encoder.Encode(&Response{
			OK:          false,
			Description: "Missing required parameter 'chatId'",
		})

		if err != nil {
			h.logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("cannot encode json")
		}
	}

	chats_getInfo := `{
		"about": "about user",
		"firstName": "User",
		"language": "en",
		"lastName": "Userov",
		"photo": [
			{
				"url": "https://rapi.myteaminternal/avatar/get?targetSn=test@test&size=1024"
			}
		],
		"type": "private",
		"ok": true
	}`

	_, err := w.Write([]byte(chats_getInfo))
	if err != nil {
		h.logger.Fatal("failed to write events")
	}
}

func (h *MockHandler) SendFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20) // 10MB max
		if err != nil && err != io.EOF {
			h.logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("cannot parse multipart form")
			h.sendErrorResponse(w, "Cannot parse multipart form")
			return
		}

		// Check if file field is present
		if r.MultipartForm == nil || r.MultipartForm.File["file"] == nil || len(r.MultipartForm.File["file"]) == 0 {
			h.sendErrorResponse(w, "Missing required parameter 'file'")
			return
		}

		// Check if file is empty
		fileHeader := r.MultipartForm.File["file"][0]
		if fileHeader.Size == 0 {
			h.sendErrorResponse(w, "File cannot be empty")
			return
		}
	} else if r.Method == http.MethodGet {
		// GET request must have fileId parameter
		if r.FormValue("fileId") == "" {
			h.sendErrorResponse(w, "Missing required parameter 'fileId'")
			return
		}
	} else {
		h.sendErrorResponse(w, "Invalid HTTP method")
		return
	}

	response := `{
		"ok": true,
		"msgId": "test123",
		"timestamp": 123456
	}`
	_, err := w.Write([]byte(response))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("cannot write response")
	}
}

func (h *MockHandler) SendVoice(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20) // 10MB max
		if err != nil && err != io.EOF {
			h.logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("cannot parse multipart form")
			h.sendErrorResponse(w, "Cannot parse multipart form")
			return
		}

		// Check if file field is present
		if r.MultipartForm == nil || r.MultipartForm.File["file"] == nil || len(r.MultipartForm.File["file"]) == 0 {
			h.sendErrorResponse(w, "Missing required parameter 'file'")
			return
		}

		// Check if file is empty
		fileHeader := r.MultipartForm.File["file"][0]
		if fileHeader.Size == 0 {
			h.sendErrorResponse(w, "File cannot be empty")
			return
		}
	} else if r.Method == http.MethodGet {
		// GET request must have fileId parameter
		if r.FormValue("fileId") == "" {
			h.sendErrorResponse(w, "Missing required parameter 'fileId'")
			return
		}
	} else {
		h.sendErrorResponse(w, "Invalid HTTP method")
		return
	}

	response := `{
		"ok": true,
		"msgId": "voice123",
		"timestamp": 123456
	}`
	_, err := w.Write([]byte(response))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("cannot write response")
	}
}

func (h *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.FormValue("token") == "":
		h.TokenError(w)
		return
	case r.URL.Path == "/messages/sendText":
		h.SendMessage(w)
		return
	case r.URL.Path == "/messages/sendFile":
		h.SendFile(w, r)
		return
	case r.URL.Path == "/messages/sendVoice":
		h.SendVoice(w, r)
		return
	case r.URL.Path == "/events/get":
		h.GetEvents(w)
		return
	case r.URL.Path == "/self/get":
		h.SelfGet(w, r)
		return
	default:
		encoder := json.NewEncoder(w)
		err := encoder.Encode(&Response{
			OK: true,
		})

		if err != nil {
			h.logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("cannot encode response")
		}
	}
}

func NewApiMockClient(t *testing.T) Client {
	t.Helper()

	testServer := httptest.NewServer(&MockHandler{logger: logrus.New()})
	t.Cleanup(testServer.Close)

	return Client{
		baseURL: testServer.URL,
		token:   "test_token",
		client:  http.DefaultClient,
		logger:  &logrus.Logger{},
	}
}
