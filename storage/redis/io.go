package codis

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/support-go/session"
)

type context struct {
	d driver
	c *Client
}

type Writer struct {
	context
}

func (w *Writer) Set(location string, value interface{}) session.Status {
	stat := w.c.dial(func(conn redis.Conn) session.Status {
		w.d.set(conn, w.d.prepareSet(location, value))
		return session.NoError
	})
	if stat.Failed() {
		return session.RedisSetErr.NewStatus(stat, "Redis写入失败")
	}
	return session.NoError
}

type Reader struct {
	context
}

func (r *Reader) GetBool(location string) (reply bool, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.Bool(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetInt(location string) (reply int, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.Int(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetInts(location string) (reply []int, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.Ints(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetIntMap(location string) (reply map[string]int, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.IntMap(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetInt64(location string) (reply int64, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.Int64(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetInt64s(location string) (reply []int64, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.Int64s(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetInt64Map(location string) (reply map[string]int64, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.Int64Map(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetUint64(location string) (reply uint64, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.Uint64(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetFloat64(location string) (reply float64, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.Float64(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetFloat64s(location string) (reply []float64, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.Float64s(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetString(location string) (reply string, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.String(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetBytes(location string) (reply []byte, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.Bytes(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetValues(location string) (reply []interface{}, stat session.Status) {
	stat = r.c.dial(func(conn redis.Conn) session.Status {
		res, err := redis.Values(r.d.get(conn, r.d.prepareGet(location)))
		if stat = checkError(err); stat.Successful() {
			reply = res
		}
		return stat
	})
	return
}

func (r *Reader) GetJson(location string, value interface{}) session.Status {
	res, stat := r.GetBytes(location)
	if stat.Failed() {
		return stat
	}
	err := json.Unmarshal(res, value)
	if err != nil {
		return session.FormatErr.ByMessage("Redis读取Key不存在")
	}
	return session.NoError
}

func checkError(err error) session.Status {
	if err != nil {
		if err == redis.ErrNil {
			return session.NotFound.ByMessage("Redis读取Key不存在")
		}
		return session.RedisGetErr.NewStatus(err, "Redis读取失败")
	}
	return session.NoError
}
