package rqc

import (
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
)

type Builder struct {
	Conn redis.Conn

	// Namespace is a key prefix for intermediate keys constructed by the
	// builder.
	Namespace string
}

func (b *Builder) Key(name string) string {
	return b.Namespace + ":" + name
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
	Code             []string
}

// Intersect specifies another set key to intersect with.
func (s Selection) Intersect(key string) Selection {
	s.IntersectionKeys = append(s.IntersectionKeys, key)
	return s
}

func (s Selection) Complement(key string) Selection {
	id := fmt.Sprintf("diff(%s,%s)", s.BaseKey, key)
	diffKey := s.Builder.Key(id)

	sdiffCode := fmt.Sprintf("redis.call('SDIFFSTORE', '%s', '%s', '%s')",
		diffKey, s.BaseKey, key)
	s.Code = append(s.Code, sdiffCode)

	s.IntersectionKeys = append(s.IntersectionKeys, diffKey)
	return s
}

func (s Selection) Generate() string {
	code := strings.Join(s.Code, "\n") + "\n"

	intersectionKeyArgs := strings.Join(s.IntersectionKeys, "', '")
	code += fmt.Sprintf("redis.call('ZINTERSTORE', '%s', %d, '%s')\n",
		s.ResultKey, len(s.IntersectionKeys), intersectionKeyArgs)

	fmt.Println(code)

	return code
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
