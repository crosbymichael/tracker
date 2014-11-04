package redis

import (
	"encoding/json"
	"fmt"

	"github.com/crosbymichael/tracker/peer"
	"github.com/crosbymichael/tracker/registry"
	"github.com/garyburd/redigo/redis"
)

type RedisRegistry struct {
	pool *redis.Pool
}

func New(addr, pass string) registry.Registry {
	return &RedisRegistry{
		pool: newPool(addr, pass),
	}
}

func (r *RedisRegistry) FetchPeers() ([]*peer.Peer, error) {
	out := []*peer.Peer{}

	keys, err := redis.Strings(r.do("KEYS", "tracker:peer:*"))
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		data, err := redis.String(r.do("GET", key))
		if err != nil {
			return nil, err
		}

		var p *peer.Peer
		if err := json.Unmarshal([]byte(data), &p); err != nil {
			return nil, err
		}

		out = append(out, p)
	}

	return out, nil
}

func (r *RedisRegistry) SavePeer(p *peer.Peer, ttl int) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	key := r.getKey(p)
	if _, err := r.do("SETEX", key, ttl, string(data)); err != nil {
		return err
	}

	return nil
}

func (r *RedisRegistry) DeletePeer(p *peer.Peer) error {
	key := r.getKey(p)

	if _, err := r.do("DEL", key); err != nil {
		return err
	}

	return nil
}

func (r *RedisRegistry) Close() error {
	return r.pool.Close()
}

func (r *RedisRegistry) getKey(p *peer.Peer) string {
	return fmt.Sprintf("tracker:peer:%s", p.Hash())
}

func (r *RedisRegistry) do(cmd string, args ...interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return conn.Do(cmd, args...)
}

func newPool(addr, pass string) *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}

		if pass != "" {
			if _, err := c.Do("AUTH", pass); err != nil {
				return nil, err
			}
		}

		return c, nil
	}, 10)
}
