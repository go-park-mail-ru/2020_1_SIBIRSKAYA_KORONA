package models

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
)

type Session struct {
	SID     string
	ID      uint
	Expires int32
}

func (ses *Session) ToProto() *proto.SessionMess {
	return &proto.SessionMess{
		Sid:     ses.SID,
		Uid:     uint64(ses.ID),
		Expires: ses.Expires,
	}
}

func CreateSessionFromProto(ses proto.SessionMess) *Session {
	return &Session{
		SID:     ses.Sid,
		ID:      uint(ses.Uid),
		Expires: ses.Expires,
	}
}
