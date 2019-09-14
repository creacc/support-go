package codis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type driver interface {
	prepareSet(location string, value interface{}) []interface{}
	prepareGet(location string) []interface{}
	set(conn redis.Conn, args []interface{}) (interface{}, error)
	get(conn redis.Conn, args []interface{}) (interface{}, error)
}

type keyDriver struct {
	expire time.Duration
}

func (d keyDriver) prepareSet(location string, value interface{}) []interface{} {
	return []interface{}{location, value}
}

func (d keyDriver) prepareGet(location string) []interface{} {
	return []interface{}{location}
}

func (d keyDriver) set(conn redis.Conn, args []interface{}) (interface{}, error) {
	reply, err := conn.Do("SET", args...)
	if err != nil {
		return nil, err
	}
	if d.expire > 0 {
		_, err = conn.Do("PEXPIRE", d.expire)
		if err != nil {
			return nil, err
		}
	}
	return reply, nil
}

func (d keyDriver) get(conn redis.Conn, args []interface{}) (interface{}, error) {
	v, e := conn.Do("GET", args...)
	if e != nil {
		return nil, e
	}
	return v, nil
}

type hashDriver struct {
	key string
}

func (d hashDriver) prepareSet(location string, value interface{}) []interface{} {
	return []interface{}{d.key, location, value}
}

func (d hashDriver) prepareGet(location string) []interface{} {
	return []interface{}{d.key, location}
}

func (d hashDriver) set(conn redis.Conn, args []interface{}) (interface{}, error) {
	_, e := conn.Do("HSET", args...)
	if e != nil {
		return nil, e
	}
	return nil, nil
}

func (d hashDriver) get(conn redis.Conn, args []interface{}) (interface{}, error) {
	v, e := conn.Do("HGET", args...)
	if e != nil {
		return nil, e
	}
	return v, nil
}
