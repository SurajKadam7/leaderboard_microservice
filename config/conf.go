package conf

import (
	"encoding/json"
	"errors"
	"io"
)

type config struct {
	Key          string `json:"key"`
	Port         string `json:"port"`
	RedisAddress string `json:"redisAdd"`
	ConsulAdress string `json:"consulAdd"`
	PoolSize     int    `json:"poolSize"`
}

func New(r io.Reader) (c config, err error) {
	c = config{}

	err = json.NewDecoder(r).Decode(&c)
	if err != nil {
		return c, err
	}

	if ok, err := c.IsValid(); !ok {
		return c, err
	}

	return c, err
}

func (c *config) IsValid() (ok bool, err error) {
	if c.Key == "" {
		return ok, errors.New("key is not present")
	}
	if c.Port == "" {
		c.Port = ":8080"
	}

	if c.RedisAddress == "" {
		return ok, errors.New("address is not present")
	}

	if c.PoolSize == 0 {
		c.PoolSize = 50
	}

	return true, nil
}
