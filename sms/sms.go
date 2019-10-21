package sms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
	"pkg.glorieux.io/mantra/internal/log"
)

// Handler handles messages
type Handler func(*Message) error

// ShortMessageService is a service for sending and receiving short messages
type ShortMessageService struct {
	accountSid     string
	authToken      string
	messageHandler Handler

	decoder    *schema.Decoder
	httpClient *http.Client
}

// New returns a new short service message
func New(sid, token string, handler Handler) *ShortMessageService {
	return &ShortMessageService{
		accountSid:     sid,
		authToken:      token,
		messageHandler: handler,

		decoder:    schema.NewDecoder(),
		httpClient: &http.Client{},
	}
}

// Send sends a message
func (sms *ShortMessageService) Send(msg *Message) {
	msgData := url.Values{}
	msgData.Set("From", msg.Sender)
	msgData.Set("To", msg.Recipient)
	msgData.Set("Body", msg.Content)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", sms.accountSid),
		strings.NewReader(msgData.Encode()),
	)
	if err != nil {
		log.Error(err)
		return
	}
	req.SetBasicAuth(sms.accountSid, sms.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := sms.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			log.Info("Sent message to: ", data["to"])
		}
	} else {
		log.Warn(resp.Status)
	}
}

func (sms *ShortMessageService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	message := new(Message)
	err = sms.decoder.Decode(message, r.PostForm)
	if err != nil {
		log.Error(err)
	}

	err = sms.messageHandler(message)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
	}
}

// Stop stops the ShortMessageService
func (*ShortMessageService) Stop() error {
	return nil
}
