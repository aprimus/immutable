package imHash

import "fmt"

import "immutable/imList"

const _DEBUGPRINTS = false
const _PACKAGENAME = "imHash"

type bitflag uint32
type HashType uint32

const _BFLength = 32 // length of Bitflag type

const _MASK HashType = 31
const _ONE HashType = 1
const _SHIFTBY = 5
const _LEVELS = 6

type KeyType interface {
	BTHash() HashType
}

type ValueType interface{}

type kvPair struct {
	key   interface{}
	value interface{}
}

type imNode struct {
	id HashType
	bitflag
	kids []*imNode
	kvs  *imList.IMList
}

type IMHash struct {
	root *imNode
}

var NumNodes int

func NewHash() *IMHash {
	nbt := new(IMHash)
	tallyNode()
	root := new(imNode)
	nbt.root = root
	return nbt
}

func (bt *IMHash) Insert(key KeyType, value ValueType) *IMHash {
	nbt := new(IMHash)
	nbt.root = bt.root.copy()
	nbt.root.insert(key, value, key.BTHash(), 0)
	return nbt
}

func (bt *IMHash) Find(key KeyType) (interface{}, ValueType) {
	kv := bt.root.find(key.BTHash(), key, 0)
	if kv != nil {
		k := kv.key
		v := kv.value
		return k, v
	}
	return nil, nil
}

func tallyNode() {
	NumNodes++
}

func newNode(h HashType) *imNode {
	tallyNode()
	nn := new(imNode)
	nn.id = h
	nn.bitflag = bitflag(uint64(h))
	return nn
}

func (btn *imNode) copy() *imNode {
	tallyNode()
	newKids := make([]*imNode, popcount(btn.bitflag))
	copy(newKids, btn.kids)
	nbtn := &imNode{btn.id, btn.bitflag, newKids, btn.kvs}
	return nbtn
}

func (imn *imNode) find(hval HashType, key interface{}, level uint32) *kvPair {
	if imn.id == hval {
		return findKVInList(imn.kvs, key)
	}
	modifiedHV := hval >> (level * _SHIFTBY)
	desiredKid := modifiedHV & _MASK
	doesKidExist := _ONE & HashType(imn.bitflag>>desiredKid)
	if doesKidExist == _ONE {
		smallerKids := imn.bitflag << (_BFLength - desiredKid)
		numSmallerKids := popcount(smallerKids)
		targetNode := imn.kids[numSmallerKids]
		return targetNode.find(hval, key, (level + 1))
	}
	return nil
}

// Because we are operating in immutable land, and we know that we're
// doing updating, the node passed in HAS ALREADY BEEN DUPLICATED.  That
// means we can modify it.  However, when we call the function recursively,
// we need to make sure that the new receiving struct is also virgin.

func (imn *imNode) insert(key interface{}, value ValueType, hval HashType, level uint32) {
	if imn.id == hval {
		imn.kvs = addOrUpdateKVList(imn.kvs, key, value)
	} else {
		modifiedHV := hval >> (level * _SHIFTBY)
		// desiredKid is in 0.._BFLENGTH
		desiredKid := modifiedHV & _MASK
		doesKidExist := _ONE & HashType(imn.bitflag>>desiredKid)
		if doesKidExist == _ONE {
			smallerKids := imn.bitflag << (_BFLength - desiredKid)
			numSmallerKids := popcount(smallerKids)
			// Make the duplicated new child to call next
			targetNode := imn.kids[numSmallerKids].copy()
			imn.kids[numSmallerKids] = targetNode
			targetNode.insert(key, value, hval, (level + 1))
		} else { // Need to make a new node for the kid
			numExistingKids := popcount(imn.bitflag)
			numNewKids := numExistingKids + 1
			newBitFlag := imn.bitflag | bitflag(_ONE<<desiredKid)
			imn.bitflag = newBitFlag
			newKids := make([]*imNode, numNewKids)
			smallerKids := imn.bitflag << (_BFLength - desiredKid)
			numSmallerKids := popcount(smallerKids)
			myIndex := numSmallerKids
			for i := 0; i < numSmallerKids; i++ {
				newKids[i] = imn.kids[i]
			}
			for i := numSmallerKids + 1; i < numNewKids; i++ {
				newKids[i] = imn.kids[i-1]
			}
			tallyNode()
			newContentNode := &imNode{hval, 0, nil, nil}
			newContentNode.kvs = imList.New().Push(&kvPair{key, value})
			newKids[myIndex] = newContentNode
			imn.kids = newKids
		}
	}
}

/* UTILITY FUNCTIONs */

func popcount(v bitflag) (numOnes int) {
	numOnes = 0
	for i := 0; i < _BFLength; i++ {
		currentBit := v & 1
		if currentBit == 1 {
			numOnes++
		}
		v = v >> 1
	}
	return
}

func dp(strings ...interface{}) {
	if _DEBUGPRINTS {
		fmt.Println("(debug)", _PACKAGENAME, strings)
	}
}

func makeKeyFinder(k interface{}) func(interface{}) bool {
	return func(toTest interface{}) bool {
		kvp := toTest.(*kvPair)
		if kvp.key == k {
			return true
		}
		return false
	}
}

func addOrUpdateKVList(l *imList.IMList, key interface{}, value ValueType) *imList.IMList {
	newPair := &kvPair{key, value}
	nl := l.UpdateOrInsert(newPair, makeKeyFinder(key))
	return nl
}

func findKVInList(l *imList.IMList, key interface{}) *kvPair {
	e := l.Fetch(makeKeyFinder(key))
	if e != nil {
		return e.(*kvPair)
	}
	println("Failed to find ", key.(string), " in list")
	return nil
}

func removeKVfromList(l *imList.IMList, key interface{}) *imList.IMList {
	return l.RemoveByFunc(makeKeyFinder(key))
}
