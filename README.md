# set
Set is an in-memory Redis like [set](https://redis.io/commands#set) datastructure

Getting Started
===============

## Installing

To start using hash, install Go and run `go get`:

```sh
$ go get -u github.com/arriqaaq/set
```

This will retrieve the library.

## Usage

```go
package main

import "github.com/arriqaaq/set"

type kv struct{k,v string}

func main() {
    key:="set1"

    // SAdd (accepts any value)
    	s := set.New()
	s.SAdd(key, "a")
	s.SAdd(key, 1)
	s.SAdd(key, &kv{1,2})

    // SPop
    	s.SPop("set#2", 1)
	isMember := set.SIsMember(key, "a")
	assert.Equal(t, true, isMember)

    // SRem
	n := set.SRem(key, "member_1")
	assert.Equal(t, true, n)

    // SUnion
	set2 := makeSet(3)

	set2.SAdd("set2", "a")
	set2.SAdd("set2", "b")
	set2.SAdd("set2", "c")

	members := set.SUnion(key, "set2")
	assert.Equal(t, 2, len(members))

    // SDiff
    members := set.SDiff(key, "set2")
    assert.Equal(t, 0, len(members))

}
```

## Supported Commands

```go
SADD
SCARD
SDIFF
SDIFFSTORE
SINTER
SINTERSTORE
SISMEMBER
SMEMBERS
SMISMEMBER
SMOVE
SPOP
SRANDMEMBER
SREM
SUNION
SUNIONSTORE
```
