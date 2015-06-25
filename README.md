# rqc

[![Build status](https://img.shields.io/travis/mmcloughlin/rqc.svg?style=flat-square)](https://travis-ci.org/mmcloughlin/rqc)

Redis query compiler (rqc) generates lua scripts to execute selection queries
on redis sets.

> Still under development.

## Installation

```
go get github.com/mmcloughlin/rqc
```

## Usage

Create a query builder with

```go
builder := Builder{
	Conn:      conn,
	Namespace: "queries",
}
```

Here `Conn` is expected to be a [redigo](github.com/garyburd/redigo/redis)
redis connection. `Namespace` is a prefix for all intermediate keys produced
in query execution.

## Acknowledgements

There are a few similar projects out there and I learned a lot from digging
around in their source code:

* [Zoom](https://github.com/albrow/zoom) is an awesome library that offers
  similar funtionality at a higher level

* [django redis engine](https://github.com/MirkoRossini/django-redis-engine)
  also contains very similar ideas
