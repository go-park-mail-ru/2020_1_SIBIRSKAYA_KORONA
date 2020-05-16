package models

//go:generate easyjson -all
type Event struct {
	ID        uint          `json:"-" gorm:"primary_key" `
	EventType string        `json:"eventType" gorm:"not null"`
	CreateAt  int64         `json:"createAt,omitempty" gorm:"not null"`
	IsRead    bool          `json:"isRead,omitempty"`
	Uid       uint          `json:"uid,omitempty" gorm:"not null"` // кому придет уведомление
	MakeUid   uint          `json:"-" gorm:"not null"`             // кто сделал действие
	MakeUsr   *User         `json:"makeUser,omitempty" gorm:"-" faker:"-"`
	MetaData  EventMetaData `json:"metaData" gorm:"foreignkey:eid" faker:"-"`
}

//easyjson:json
type Events []Event

func (event *Event) TableName() string {
	return "events"
}

type EventMetaData struct {
	ID         uint   `json:"-" gorm:"primary_key"`
	Eid        uint   `json:"-"`
	Uid        uint   `json:"-"` // над кем/чем совершено действие
	Usr        *User  `json:"user,omitempty" gorm:"-" faker:"-"`
	Bid        uint   `json:"bid,omitempty"`
	Cid        uint   `json:"cid,omitempty"`
	Tid        uint   `json:"tid,omitempty"`
	EntityData string `json:"entityData,omitempty"`
	Text       string `json:"text,omitempty"`
}

func (metaData *EventMetaData) TableName() string {
	return "event_meta_data"
}
