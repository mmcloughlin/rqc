package rqc

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	conn    redis.Conn
	builder Builder
}

func (suite *RedisTestSuite) SetupSuite() {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	suite.conn = conn
	suite.builder = Builder{
		Conn:      conn,
		Namespace: "rqc:test",
	}
}

func (suite *RedisTestSuite) TearDownSuite() {
	suite.conn.Close()
}

func (suite *RedisTestSuite) SetupTest() {
	_, err := suite.conn.Do("FLUSHDB")
	if err != nil {
		panic(err)
	}
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}

func (suite *RedisTestSuite) TestIntersect() {
	suite.conn.Do("SADD", "a", 1, 2, 3)
	suite.conn.Do("SADD", "b", 2, 3, 4)

	query := suite.builder.Select("a").Intersect("b")
	query.Run()

	result, err := redis.Strings(suite.conn.Do("ZRANGE", query.ResultKey, 0, -1))
	if err != nil {
		panic(err)
	}
	suite.Equal([]string{"2", "3"}, result, "Expected set intersection {2,3}")
}

func (suite *RedisTestSuite) TestComplement() {
	suite.conn.Do("SADD", "a", 1, 2, 3)
	suite.conn.Do("SADD", "b", 2, 3, 4)

	query := suite.builder.Select("a").Complement("b")
	query.Run()

	result, err := redis.Strings(suite.conn.Do("ZRANGE", query.ResultKey, 0, -1))
	if err != nil {
		panic(err)
	}
	suite.Equal([]string{"1"}, result, "Expected set complement {1}")
}

func (suite *RedisTestSuite) TestFilterGt() {
	suite.conn.Do("SADD", "a", 1, 2, 3)
	suite.conn.Do("ZADD", "r", 1, 1, 2, 2, 3, 3)

	query := suite.builder.Select("a").Filter("r", Gt(2.5))
	query.Run()

	result, err := redis.Strings(suite.conn.Do("ZRANGE", query.ResultKey, 0, -1))
	if err != nil {
		panic(err)
	}
	suite.Equal([]string{"3"}, result, "Expected filtered result {3}")
}

func (suite *RedisTestSuite) TestFilterLt() {
	suite.conn.Do("SADD", "a", 1, 2, 3)
	suite.conn.Do("ZADD", "r", 1, 1, 2, 2, 3, 3)

	query := suite.builder.Select("a").Filter("r", Lt(2.5))
	query.Run()

	result, err := redis.Strings(suite.conn.Do("ZRANGE", query.ResultKey, 0, -1))
	if err != nil {
		panic(err)
	}
	suite.Equal([]string{"1", "2"}, result, "Expected filtered result {1,2}")
}

func (suite *RedisTestSuite) TestFilterRange() {
	suite.conn.Do("SADD", "a", 1, 2, 3)
	suite.conn.Do("ZADD", "r", 1, 1, 2, 2, 3, 3)

	query := suite.builder.Select("a").Filter("r", Range{1.5, 2.5})
	query.Run()

	result, err := redis.Strings(suite.conn.Do("ZRANGE", query.ResultKey, 0, -1))
	if err != nil {
		panic(err)
	}
	suite.Equal([]string{"2"}, result, "Expected filtered result {2}")
}
