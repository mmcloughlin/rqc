package rqc

import (
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
)

type Builder struct {
	Conn redis.Conn
}

func (b *Builder) Select(key string) *Selection {
	return &Selection{
		Builder:          b,
		BaseKey:          key,
		ResultKey:        "result",
		IntersectionKeys: []string{key},
	}
}

type Selection struct {
	Builder          *Builder
	BaseKey          string
	ResultKey        string
	IntersectionKeys []string
}

// Intersect specifies another set key to intersect with.
func (s Selection) Intersect(key string) Selection {
	s.IntersectionKeys = append(s.IntersectionKeys, key)
	return s
}

func (s Selection) Generate() string {
	intersectionKeyArgs := strings.Join(s.IntersectionKeys, "', '")
	return fmt.Sprintf("redis.call('ZINTERSTORE', '%s', %d, '%s')\n",
		s.ResultKey, len(s.IntersectionKeys), intersectionKeyArgs)
}

func (s Selection) Script() *redis.Script {
	code := s.Generate()
	return redis.NewScript(0, code)
}

func (s Selection) Run() {
	script := s.Script()
	_, err := script.Do(s.Builder.Conn)
	if err != nil {
		panic(err)
	}
}
