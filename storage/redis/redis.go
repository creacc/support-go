package codis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/support-go/session"
	"github.com/support-go/utils/log"
	"time"
)

type Client struct {
	conf Config
}

func NewClient(conf Config) *Client {
	return &Client{
		conf: conf,
	}
}

func (c *Client) KeyWriter() *Writer {
	return c.KeyWriterWithExpire(-1)
}

func (c *Client) KeyWriterWithExpire(expire time.Duration) *Writer {
	return c.newWriter(keyDriver{expire: expire})
}

func (c *Client) KeyReader() *Reader {
	return c.newReader(keyDriver{})
}

func (c *Client) HashWriter(key string) *Writer {
	return c.newWriter(hashDriver{})
}

func (c *Client) HashReader(key string) *Reader {
	return c.newReader(hashDriver{})
}

func (c *Client) newReader(d driver) *Reader {
	return &Reader{context: c.newContext(d)}
}

func (c *Client) newWriter(d driver) *Writer {
	return &Writer{context: c.newContext(d)}
}

func (c *Client) newContext(d driver) context {
	return context{
		d: d,
		c: c,
	}
}

func (c *Client) dial(handler func(conn redis.Conn) session.Status) session.Status {
	conn, err := redis.Dial("tcp", c.conf.getAddress())
	if err != nil {
		log.Logger.Error("redis dial err %v", err)
		return session.RedisDialErr.NewStatus(err, "Redis连接失败")
	}
	defer conn.Close()
	return handler(conn)
}
