package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	stdLog "log"
)

type config struct {
	Key          string `json:"key"`
	Port         string `json:"port"`
	RedisAddress string `json:"redisAdd"`
	ConsulAdress string `json:"consulAdd"`
	PoolSize     int    `json:"poolSize"`
}

func New(r io.Reader) (config, error) {
	c := config{}
	err := json.NewDecoder(r).Decode(&c)

	if err != nil {
		panic(fmt.Errorf("not able unmarshal the kv err : %w", err))
	}

	if ok, err := c.IsValid(); !ok {
		panic(fmt.Sprint("Config Validation Faild : ", err))
	}

	stdLog.Printf("\nConfigurations : \n%+v", c)

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
