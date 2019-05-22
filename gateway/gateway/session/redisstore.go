package session

import (
	"time"

	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

// RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

// NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//initialize and return a new RedisStore struct
	if client == nil {
		panic("No client provided!")
	}
	return &RedisStore{client, sessionDuration}
}

// Store implementation

// Save saves the provided `sessionState` and associated SessionID to the store.
//
// The `sessionState` parameter is typically a pointer to a struct containing all the data you want to be
// associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	sesJson, err := json.Marshal(sessionState)
	if err != nil {
		return fmt.Errorf("error marshaling sessionState into json:\n%s", err.Error())
	}
	err = rs.Client.Set(getRedisKey(sid), sesJson, rs.SessionDuration).Err()
	if err != nil {
		return fmt.Errorf("error setting session state:\n%s", err.Error())
	}
	return nil
}

// Get populates `sessionState` with the data previously saved for the given SessionID.
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	pipe := rs.Client.Pipeline()
	res := pipe.Get(getRedisKey(sid))
	expErr := pipe.Expire(getRedisKey(sid), rs.SessionDuration).Err()
	_, pipeErr := pipe.Exec()
	if res.Err() != nil {
		return ErrStateNotFound
	}
	if expErr != nil {
		return fmt.Errorf("error changing expiration of session <%s>:\n%s", sid, expErr.Error())
	}
	if pipeErr != nil {
		return fmt.Errorf("error getting sid <%s>:\n%v", string(sid), pipeErr.Error())
	}
	err := json.Unmarshal([]byte(res.Val()), sessionState)
	if err != nil {
		return fmt.Errorf("error unmarshaling sessionState: %s", err.Error())
	}
	return nil
}

// Exists determines if the session id is in the session store.
func (rs *RedisStore) Exists(sid SessionID) (bool, error) {

	res := rs.Client.Exists(getRedisKey(sid))

	if res.Err() != nil {
		return false, res.Err()
	}
	exRes := res.Val()
	return exRes == 1, nil
}

// Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	err := rs.Client.Del(getRedisKey(sid)).Err()
	if err != nil {
		return fmt.Errorf("error deleting the session <%s>:\n%s", sid, err.Error())
	}
	return nil
}

// getRedisKey() returns the redis key to use for the SessionID.
func getRedisKey(sid SessionID) string {
	// convert the SessionID to a string and add the prefix "sid:" to keep
	// SessionID keys separate from other keys that might end up in this
	// redis instance.
	return "sid:" + sid.String()
}
