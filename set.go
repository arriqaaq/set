package set

// helpful to not write everywhere struct{}{}
var keyExists = struct{}{}

type (
	Set struct {
		records map[string]set
	}

	set map[interface{}]struct{}
)

/*
Related commands

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
*/

// newSet creates a new non-threadsafe Set.
func newSet() set {
	s := make(map[interface{}]struct{})

	return s
}

// add includes the specified items (one or more) to the set. The underlying
// Set s is modified. If passed nothing it silently returns.
func (s set) add(items ...interface{}) {
	if len(items) == 0 {
		return
	}

	for _, item := range items {
		s[item] = keyExists
	}
}

// remove deletes the specified items from the set.  The underlying Set s is
// modified. If passed nothing it silently returns.
func (s set) remove(items ...interface{}) {
	if len(items) == 0 {
		return
	}

	for _, item := range items {
		delete(s, item)
	}
}

// has looks for the existence of items passed. It returns false if nothing is
// passed. For multiple items it returns true only if all of  the items exist.
func (s set) has(items ...interface{}) bool {
	// assume checked for empty item, which not exist
	if len(items) == 0 {
		return false
	}

	exist := true
	for _, item := range items {
		if _, exist = s[item]; !exist {
			break
		}
	}
	return exist
}

// each traverses the items in the Set, calling the provided function for each
// set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s set) each(f func(item interface{}) bool) {
	for item := range s {
		if !f(item) {
			break
		}
	}
}

// Copy returns a new Set with a copy of s.
func (s set) copy() set {
	u := newSet()
	for item := range s {
		u.add(item)
	}
	return u
}

// list returns a slice of all items. There is also StringSlice() and
// IntSlice() methods for returning slices of type string or int.
func (s set) list() []interface{} {
	list := make([]interface{}, 0, len(s))

	for item := range s {
		list = append(list, item)
	}

	return list
}

func (s set) size() int {
	return len(s)
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given t set.
func (s set) Merge(t set) {
	t.each(func(item interface{}) bool {
		s[item] = keyExists
		return true
	})
}

// it's not the opposite of Merge.
// Separate removes the set items containing in t from set s. Please aware that
func (s set) Separate(t set) {
	s.remove(t.list()...)
}

// union is the merger of multiple sets. It returns a new set with all the
// elements present in all the sets that are passed.
//
// The dynamic type of the returned set is determined by the first passed set's
// implementation of the New() method.
func union(set1, set2 set, sets ...set) set {
	u := set1.copy()
	set2.each(func(item interface{}) bool {
		u.add(item)
		return true
	})
	for _, set := range sets {
		set.each(func(item interface{}) bool {
			u.add(item)
			return true
		})
	}

	return u
}

// difference returns a new set which contains items which are in in the first
// set but not in the others. Unlike the difference() method you can use this
// function separately with multiple sets.
func difference(set1, set2 set, sets ...set) set {
	s := set1.copy()
	s.Separate(set2)
	for _, set := range sets {
		s.Separate(set)
	}
	return s
}

// intersection returns a new set which contains items that only exist in all given sets.
func intersection(set1, set2 set, sets ...set) set {
	all := union(set1, set2, sets...)
	result := union(set1, set2, sets...)

	all.each(func(item interface{}) bool {
		if !set1.has(item) || !set2.has(item) {
			result.remove(item)
		}

		for _, set := range sets {
			if !set.has(item) {
				result.remove(item)
			}
		}
		return true
	})
	return result
}

func New() *Set {
	return &Set{
		records: make(map[string]set),
	}
}

func (s *Set) SAdd(key string, member interface{}) int {
	if !s.exists(key) {
		s.records[key] = make(set)
	}
	s.records[key].add(member)
	return s.records[key].size()
}

func (s *Set) SPop(key string, count int) []interface{} {
	var val []interface{}
	if !s.exists(key) || count <= 0 {
		return val
	}

	for k := range s.records[key] {
		s.records[key].remove(k)
		val = append(val, k)

		count--
		if count == 0 {
			break
		}
	}
	return val
}

func (s *Set) SIsMember(key string, member interface{}) bool {
	return s.fieldExists(key, member)
}

func (s *Set) SRandMember(key string, count int) []interface{} {
	var val []interface{}
	if !s.exists(key) || count == 0 {
		return val
	}

	if count > 0 {
		for k := range s.records[key] {
			val = append(val, k)
			if len(val) == count {
				break
			}
		}
	} else {
		count = -count
		randomVal := func() interface{} {
			// map indexes are not sorted and will return random
			// key when ranged through
			for k := range s.records[key] {
				return k
			}
			return nil
		}

		for count > 0 {
			val = append(val, randomVal())
			count--
		}
	}
	return val
}

func (s *Set) SRem(key string, member interface{}) bool {
	if !s.exists(key) {
		return false
	}

	if ok := s.records[key].has(member); ok {
		s.records[key].remove(member)
		return true
	}

	return false
}

func (s *Set) SMove(src, dst string, member interface{}) bool {
	if !s.fieldExists(src, member) {
		return false
	}

	if !s.exists(dst) {
		s.records[dst] = make(set)
	}

	s.records[src].remove(member)
	s.records[dst].add(member)

	return true
}

func (s *Set) SCard(key string) int {
	if !s.exists(key) {
		return 0
	}

	return s.records[key].size()
}

func (s *Set) SMembers(key string) (val []interface{}) {
	if !s.exists(key) {
		return
	}

	val = s.records[key].list()
	return
}

func (s *Set) SUnion(keys ...string) (val []interface{}) {
	if len(keys) == 0 {
		return
	}
	if len(keys) == 1 {
		if s.exists(keys[0]) {
			return s.records[keys[0]].list()
		}
		return
	}

	m := make([]set, 0, 1)
	for _, k := range keys {
		if s.exists(k) {
			m = append(m, s.records[k])
		}
	}

	un := union(m[0], m[1], m[2:]...)
	for k := range un {
		val = append(val, k)
	}

	return
}

func (s *Set) SUnionStore(key string, keys ...string) int {
	un := s.SUnion(keys...)
	for _, k := range un {
		s.SAdd(key, k)
	}
	return len(un)
}

func (s *Set) SDiff(keys ...string) (val []interface{}) {
	if len(keys) == 0 {
		return
	}
	if len(keys) == 1 {
		if s.exists(keys[0]) {
			return s.records[keys[0]].list()
		}
		return
	}

	m := make([]set, 0, 1)
	for _, k := range keys {
		if s.exists(k) {
			m = append(m, s.records[k])
		}
	}

	diff := difference(m[0], m[1], m[2:]...)
	for k := range diff {
		val = append(val, k)
	}

	return
}

func (s *Set) SDiffStore(key string, keys ...string) int {
	diff := s.SDiff(keys...)
	for _, k := range diff {
		s.SAdd(key, k)
	}
	return len(diff)
}

func (s *Set) SInter(keys ...string) (val []interface{}) {
	if len(keys) == 0 {
		return
	}
	if len(keys) == 1 {
		if s.exists(keys[0]) {
			return s.records[keys[0]].list()
		}
		return
	}

	inter := make([]set, 0, 1)
	for _, k := range keys {
		if s.exists(k) {
			inter = append(inter, s.records[k])
		}
	}

	m := intersection(inter[0], inter[1], inter[2:]...)
	for v := range m {
		val = append(val, v)
	}
	return
}

func (s *Set) SInterStore(key string, keys ...string) int {
	if len(keys) == 0 {
		return 0
	}
	if len(keys) == 1 {
		if s.exists(keys[0]) {
			return len(s.records[keys[0]])
		}
		return 0
	}

	m := s.SInter(keys...)
	for _, k := range m {
		s.SAdd(key, k)
	}
	return len(m)
}

func (s *Set) SKeyExists(key string) (ok bool) {
	return s.exists(key)
}

func (s *Set) SClear(key string) {
	if s.SKeyExists(key) {
		delete(s.records, key)
	}
}

func (s *Set) exists(key string) bool {
	_, exist := s.records[key]
	return exist
}

func (s *Set) fieldExists(key string, field interface{}) bool {
	fieldSet, exist := s.records[key]
	if !exist {
		return false
	}
	ok := fieldSet.has(field)
	return ok
}

func (s *Set) Keys() []string {
	keys := make([]string, 0, len(s.records))
	for k := range s.records {
		keys = append(keys, k)
	}
	return keys
}
