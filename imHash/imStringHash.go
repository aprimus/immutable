package imHash

type StringHash struct {
	*IMHash
}

func NewStringHash() *StringHash {
	bt := NewHash()
	return &StringHash{bt}
}

func (ih *StringHash) Insert(key string, value ValueType) *StringHash {
	nih := NewStringHash()
	nih.root = ih.root.copy()
	nih.root.insert(key, value, hashstr(key), 0)
	return nih
}

func (sh *StringHash) Find(key string) (string, ValueType) {
	hval := hashstr(key)
	kv := sh.root.find(HashType(hval), key, 0)
	if kv != nil {
		k := kv.key.(string)
		if kv.key != k {
			println("In imStringHash, returned key is not a match:", key, kv.key)
		}
		v := kv.value
		return k, v
	}
	return "", nil
}

/* There is no theoretical / formal reason why this is
an optimal hash mechanism.  However, after tweaking values
this appears to be fairly uniform and have a low number of
collisions.  Inserting the 250k dictionary from my macbook
resulted in 6 collisions.  Statistically, the expected number
of colisions on inserting 250k items randomly, one at a time,
into 2^32 slots is ~=7.  Thus, I'm ok with this hash.

Performance is also quite reasonable, as each letter
requires only 5 bit operations
*/

func hashstr(s string) HashType {
	var result uint32
	result = 0
	bytes := []byte(s)

	for n, b := range bytes {
		result ^= uint32(b) ^ (uint32(n*int(b)) << 8)
		tmp := result >> 27
		result = result << 5
		result = result | tmp
	}

	return HashType(result)

}
