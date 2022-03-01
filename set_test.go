package set

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	key      = "set"
	falseKey = "not_set"
)

func makeSet(n int) *Set {
	set := New()
	for i := 1; i <= n; i++ {
		member := fmt.Sprintf("member_%d", i)
		set.SAdd(key, member)
	}
	return set
}

func TestSet_SAdd(t *testing.T) {
	set := New()
	set.SAdd(key, "a")

	n := set.SAdd(key, "b")
	assert.Equal(t, 2, n)
}

func TestSet_SPop(t *testing.T) {
	set := makeSet(3)
	set.SPop(key, 3)

	res := set.SPop(falseKey, 1)
	assert.Equal(t, 0, len(res))
}

func TestSet_SIsMember(t *testing.T) {
	set := makeSet(3)

	isMember := set.SIsMember(key, "member_1")
	assert.Equal(t, true, isMember)
	isMember = set.SIsMember(key, "123")
	assert.Equal(t, false, isMember)
	isMember = set.SIsMember(falseKey, "123")
	assert.Equal(t, false, isMember)
}

func TestSet_SRem(t *testing.T) {
	set := makeSet(3)

	n := set.SRem(key, "member_1")
	assert.Equal(t, true, n)

	n = set.SRem(key, "member_1")
	assert.Equal(t, false, n)

	n = set.SRem(key, "member_2")
	assert.Equal(t, true, n)

	n = set.SRem(key, "member_3")
	assert.Equal(t, true, n)

	n = set.SRem(key, "member_4")
	assert.Equal(t, false, n)

	n = set.SRem(key, "ss")
	assert.Equal(t, false, n)

	n = set.SRem(falseKey, "abc")
	assert.Equal(t, false, n)
}

func TestSet_SUnion(t *testing.T) {
	set := makeSet(3)

	set.SAdd("set2", "member_1")
	set.SAdd("set2", "member_2")
	set.SAdd("set2", "member_3")

	members := set.SUnion(key, "set2")
	assert.Equal(t, 3, len(members))
}

func TestSet_SDiff(t *testing.T) {
	set := makeSet(3)
	set.SAdd("set2", "member_1")
	set.SAdd("set2", "member_4")
	set.SAdd("set2", "member_5")

	t.Run("normal situation", func(t *testing.T) {
		members := set.SDiff(key, "set2")
		assert.Equal(t, 2, len(members))
	})
	t.Run("one key", func(t *testing.T) {
		members := set.SDiff(key)
		assert.Equal(t, 3, len(members))
	})
}

func TestSet_SClear(t *testing.T) {
	set := makeSet(3)
	set.SClear(key)

	val := set.SMembers(key)
	assert.Equal(t, len(val), 0)
}

func TestSet_SKeyExists(t *testing.T) {
	set := makeSet(3)

	exists1 := set.SKeyExists(key)
	assert.Equal(t, exists1, true)

	set.SClear(key)

	exists2 := set.SKeyExists(key)
	assert.Equal(t, exists2, false)
}
