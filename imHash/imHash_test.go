package imHash

import "testing"

type siPair struct {
	string
	int
}

func createTestingHash() *StringHash {
	sh := NewStringHash()
	sh = sh.Insert("one", 1)
	sh = sh.Insert("two", 2)
	sh = sh.Insert("three", 3)
	sh = sh.Insert("four", 4)
	sh = sh.Insert("five", 5)
	sh = sh.Insert("six", 6)
	sh = sh.Insert("seven", 7)
	sh = sh.Insert("eight", 8)
	sh = sh.Insert("nine", 9)
	sh = sh.Insert("ten", 10)
	return sh
}

func TestSimpleInsertFind(t *testing.T) {
	sh := createTestingHash()
	sh2 := sh.Insert("two", 3)
	k, v := sh.Find("one")
	if k != "one" || v.(int) != 1 {
		t.Error("Expected one/1 but got ", k, "/", v)
	}
	k, v = sh.Find("two")
	if k != "two" || v.(int) != 2 {
		t.Error("Expected two/2 but got ", k, "/", v, "Likely not immutable")
	}
	k, v = sh2.Find("two")
	if k != "two" || v.(int) != 3 {
		t.Error("Expected two/3 but got ", k, "/", v)
	}

}

// The following three triplets hash down to the same value:
// bcquipper == feantu == cfcoseismic
// dbsuperobligation == dfChloridella == bedistingue
// edYiddishist == bdupwards == ffwrybill

var collisionTriples = [][]*siPair{{&siPair{"bcquipper", 20},
	&siPair{"feantu", 21},
	&siPair{"cfcoseismic", 22}},
	{&siPair{"dbsuperobligation", 30},
		&siPair{"dfChloridella", 31},
		&siPair{"bedistingue", 32}},
	{&siPair{"edYiddishist", 40},
		&siPair{"bdupwards", 41},
		&siPair{"ffwrybill", 42}}}

func TestHashCollisions(t *testing.T) {
	sh := createTestingHash()
	for i := 0; i < len(collisionTriples); i++ {
		for j := 0; j < len(collisionTriples[i]); j++ {
			k, v := collisionTriples[i][j].string, collisionTriples[i][j].int
			sh = sh.Insert(k, v)
		}
	}
	for i := 0; i < len(collisionTriples); i++ {
		for j := 0; j < len(collisionTriples[i]); j++ {
			k := collisionTriples[i][len(collisionTriples[i])-j-1].string
			v := collisionTriples[i][len(collisionTriples[i])-j-1].int
			kr, vr := sh.Find(k)
			if k != kr || v != vr {
				t.Error("Collisions not resolved properly.  Expected ", k, "/", v, "but got", kr, "/", vr)
			}
		}
	}

	// check for a change to work...
	sh2 := sh.Insert("ffwrybill", 52)
	collisionTriples[2][2] = &siPair{"ffwrybill", 52}
	for i := 0; i < len(collisionTriples); i++ {
		for j := 0; j < len(collisionTriples[i]); j++ {
			k := collisionTriples[i][len(collisionTriples[i])-j-1].string
			v := collisionTriples[i][len(collisionTriples[i])-j-1].int
			kr, vr := sh2.Find(k)
			if k != kr || v != vr {
				t.Error("Collision lists do not appear to update/resolve properly.  Expected ", k, "/", v, "but got", kr, "/", vr)
			}
		}
	}

	// check for immutability
	collisionTriples[2][2] = &siPair{"ffwrybill", 42}
	for i := 0; i < len(collisionTriples); i++ {
		for j := 0; j < len(collisionTriples[i]); j++ {
			k := collisionTriples[i][len(collisionTriples[i])-j-1].string
			v := collisionTriples[i][len(collisionTriples[i])-j-1].int
			kr, vr := sh.Find(k)
			if k != kr || v != vr {
				t.Error("Collisions lists do not appear immutable.  Expected ", k, "/", v, "but got", kr, "/", vr, "(sh @", sh, ", sh2 @ ", sh2)
			}
		}
	}

}
