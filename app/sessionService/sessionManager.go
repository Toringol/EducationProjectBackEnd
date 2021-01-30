package sessionService

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Toringol/EducationProjectBackEnd/app/sessionService/session"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type sessionManager struct {
	redisConn redis.Conn
}

// NewSessionManager - create new SessionManager with connect to RedisStorage
func NewSessionManager(conn redis.Conn) *sessionManager {
	return &sessionManager{
		redisConn: conn,
	}
}

func (sm *sessionManager) Create(ctx context.Context, in *session.Session) (*session.SessionID, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	sessID := session.SessionID{
		ID: uid.String(),
	}

	dataSerialized, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	mkey := "sessions:" + sessID.ID

	result, err := redis.String(sm.redisConn.Do("SET", mkey, dataSerialized, "EX", 86400))
	if err != nil {
		return nil, err
	}
	if result != "OK" {
		return nil, fmt.Errorf("result not OK")
	}

	return &sessID, nil
}

func (sm *sessionManager) Check(ctx context.Context, in *session.SessionID) (*session.Session, error) {
	mkey := "sessions:" + in.ID
	data, err := redis.Bytes(sm.redisConn.Do("GET", mkey))
	if err != nil {
		return nil, err
	}

	sess := new(session.Session)

	err = json.Unmarshal(data, sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (sm *sessionManager) Delete(ctx context.Context, in *session.SessionID) (*session.Nothing, error) {
	mkey := "sessions:" + in.ID
	_, err := redis.Int(sm.redisConn.Do("DEL", mkey))
	if err != nil {
		return nil, err
	}

	return &session.Nothing{}, nil
}
