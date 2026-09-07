package event

import(
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

type EventType string

const(
	PageView EventType="page_view"
	Click EventType="click"
	Purchase EventType="purchase"
	SignUp EventType="sign_up"
	CustomEvent EventType="custom_event"
)

type Event struct(
	ID string 'json:"id"'
	Type EventType 'json:"type"'
	UserID string 'json:"user_id"'
	SessionID string 'json:"session_id"'
	Timestamp time.Time 'json:"timestamp"'
	Properties json.RawMessage 'json:"properties"'
	Metadata map[string]interface{} 'json:"metadata,omitempty"'
)

func NewEvent(eventType EventType, UserID, SessionID string, properties interface{}) *Event{
	props, _ := json.Marshal(properties)

	return &Event{
		ID: uuid.New().String(),
		Type: eventType,
		UserID: userID,
		SessionID: sessionID,
		Timestamp: time.now(),
		Properties: props,
		Metadata: map[string]interface{}{
            "source": "api",
            "version": "1.0",
		},
	}
}

func (e *Event) ToJSON() ([]byte,error){
	return json.Marshal(e)
}

func FromJSON(data []byte) (e *Event, error){
	var event Event
	if err != json.UnMarshal(data, &event); err != nil{
		return nil, err
	}
	return &event, nil
}