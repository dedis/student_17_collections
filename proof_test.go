package collection

import "testing"
import csha256 "crypto/sha256"
import "encoding/binary"

func TestProofDumpNode(test *testing.T) {
    stake64 := Stake64{}
    data := Data{}

    collection := EmptyCollection(stake64, data)
    collection.Add([]byte("mykey"), uint64(66), []byte("myvalue"))

    rootdump := dumpnode(collection.root)

    if rootdump.label != collection.root.label {
        test.Error("[proof.go]", "[dumpnode]", "dumpnode() sets wrong label on dump of internal node.")
    }

    if len(rootdump.key) != 0 {
        test.Error("[proof.go]", "[dumpnode]", "dumpnode() sets key on internal node.")
    }

    if len(rootdump.values) != 2 {
        test.Error("[proof.go]", "[dumpnode]", "dumpnode() sets the wrong number of values on internal node.")
    }

    if !equal(rootdump.values[0], collection.root.values[0]) || !equal(rootdump.values[1], collection.root.values[1]) {
        test.Error("[proof.go]", "[dumpnode]", "dumpnode() sets the wrong values on internal node.")
    }

    if (rootdump.children.left != collection.root.children.left.label) || (rootdump.children.right != collection.root.children.right.label) {
        test.Error("[proof.go]", "[dumpnode]", "dumpnode() sets the wrong children labels on internal node.")
    }

    var leaf *node

    if collection.root.children.left.placeholder() {
        leaf = collection.root.children.right
    } else {
        leaf = collection.root.children.left
    }

    leafdump := dumpnode(leaf)

    if leafdump.label != leaf.label {
        test.Error("[proof.go]", "[dumpnode]", "dumpnode() sets wrong label on dump of leaf.")
    }

    if !equal(leafdump.key, leaf.key) {
        test.Error("[proof.go]", "[dumpnode]", "dumpnode() sets wrong key on leaf.")
    }

    if len(leafdump.values) != 2 {
        test.Error("[proof.go]", "[dumpnode]", "dumpnode() sets the wrong number of values on leaf.")
    }

    if !equal(leafdump.values[0], leaf.values[0]) || !equal(leafdump.values[1], leaf.values[1]) {
        test.Error("[proof.go]", "[dumpnode]", "dumpnode() sets the wrong values on leaf.")
    }

    var empty [csha256.Size]byte

    if (leafdump.children.left != empty) || (leafdump.children.right != empty) {
        test.Error("[proof.go]", "[dumpnode]", "dumpnode() sets non-null children labels on leaf.")
    }
}

func TestProofDumpGetters(test *testing.T) {
    stake64 := Stake64{}
    data := Data{}

    collection := EmptyCollection(stake64, data)
    collection.Add([]byte("mykey"), uint64(66), []byte("myvalue"))

    rootdump := dumpnode(collection.root)

    var leaf *node

    if collection.root.children.left.placeholder() {
        leaf = collection.root.children.right
    } else {
        leaf = collection.root.children.left
    }

    leafdump := dumpnode(leaf)

    if rootdump.leaf() {
        test.Error("[proof.go]", "[dumpgetters]", "leaf() returns true on internal node.")
    }

    if !(leafdump.leaf()) {
        test.Error("[proof.go]", "[dumpgetters]", "leaf() returns false on leaf node.")
    }
}

func TestProofDumpConsistent(test *testing.T) {
    stake64 := Stake64{}
    data := Data{}

    collection := EmptyCollection(stake64, data)
    collection.Add([]byte("mykey"), uint64(66), []byte("myvalue"))

    rootdump := dumpnode(collection.root)

    var leaf *node

    if collection.root.children.left.placeholder() {
        leaf = collection.root.children.right
    } else {
        leaf = collection.root.children.left
    }

    leafdump := dumpnode(leaf)

    if !(rootdump.consistent()) {
        test.Error("[proof.go]", "[consistent]", "consistent() returns false on valid internal node.")
    }

    rootdump.label[0]++

    if rootdump.consistent() {
        test.Error("[proof.go]", "[consistent]", "consistent() returns true on invalid internal node.")
    }

    if !(leafdump.consistent()) {
        test.Error("[proof.go]", "[consistent]", "consistent() returns false on valid leaf.")
    }

    leafdump.label[0]++

    if leafdump.consistent() {
        test.Error("[proof.go]", "[consistent]", "consistent() returns true on invalid leaf.")
    }
}

func TestProofGetters(test *testing.T) {
    proof := Proof{}
    proof.key = []byte("mykey")

    if !equal(proof.Key(), []byte("mykey")) {
        test.Error("[proof.go]", "[proofgetters]", "Key() returns wrong key.")
    }
}

func TestProofMatchValues(test *testing.T) {
    collision := func(key []byte, bits int) []byte {
        target := sha256(key)
        sample := make([]byte, 8)

        for index := 0;; index++ {
            binary.BigEndian.PutUint64(sample, uint64(index))
            hash := sha256(sample)
            if match(hash[:], target[:], bits) && !match(hash[:], target[:], bits + 1) {
                return sample
            }
        }
    }

    stake64 := Stake64{}
    data := Data{}

    firstkey := []byte("mykey")
    secondkey := collision(firstkey, 5)

    collection := EmptyCollection(stake64, data)
    collection.Add(firstkey, uint64(66), []byte("firstvalue"))
    collection.Add(secondkey, uint64(99), []byte("secondvalue"))

    proof := Proof{}
    proof.collection = &collection
    proof.key = firstkey
    proof.root = dumpnode(collection.root)

    path := sha256(firstkey)
    cursor := collection.root

    for depth := 0; depth < 6; depth++ {
        proof.steps = append(proof.steps, step{dumpnode(cursor.children.left), dumpnode(cursor.children.right)})

        if bit(path[:], depth) {
            cursor = cursor.children.right
        } else {
            cursor = cursor.children.left
        }
    }

    if !(proof.Match()) {
        test.Error("[proof.go]", "[match]", "Proof Match() returns false on matching key.")
    }

    firstvalues, err := proof.Values()

    if err != nil {
        test.Error("[proof.go]", "[values]", "Proof Values() returns error on matching key.")
    }

    if len(firstvalues) != 2 {
        test.Error("[proof.go]", "[values]", "Proof Values() returns wrong number of values.")
    }

    if (firstvalues[0].(uint64) != 66) || !equal(firstvalues[1].([]byte), []byte("firstvalue")) {
        test.Error("[proof.go]", "[values]", "Proof Values() returns wrong values.")
    }

    proof.key = secondkey

    if !(proof.Match()) {
        test.Error("[proof.go]", "[match]", "Proof Match() returns false on matching key.")
    }

    secondvalues, err := proof.Values()

    if err != nil {
        test.Error("[proof.go]", "[values]", "Proof Values() returns error on matching key.")
    }

    if len(secondvalues) != 2 {
        test.Error("[proof.go]", "[values]", "Proof Values() returns wrong number of values.")
    }

    if (secondvalues[0].(uint64) != 99) || !equal(secondvalues[1].([]byte), []byte("secondvalue")) {
        test.Error("[proof.go]", "[values]", "Proof Values() returns wrong values.")
    }

    proof.key = []byte("wrongkey")

    if proof.Match() {
        test.Error("[proof.go]", "[match]", "Proof Match() returns true on non-matching key.")
    }

    _, err = proof.Values()

    if err == nil {
        test.Error("[proof.go]", "[values]", "Proof Values() does not yield an error on non-matching key.")
    }

    proof.key = firstkey

    proof.steps[5].left.values[0] = make([]byte, 7)
    proof.steps[5].right.values[0] = make([]byte, 7)

    _, err = proof.Values()

    if err == nil {
        test.Error("[proof.go]", "[values]", "Proof Values() does not yield an error on a record with ill-formed values.")
    }

    proof.steps[5].left.values = [][]byte{make([]byte, 8)}
    proof.steps[5].left.values = [][]byte{make([]byte, 8)}

    _, err = proof.Values()

    if err == nil {
        test.Error("[proof.go]", "[values]", "Proof Values() does not yield an error on a record with wrong number of values.")
    }

    proof.steps = []step{}

    if proof.Match() {
        test.Error("[proof.go]", "[match]", "Proof Match() returns true on a proof with no steps.")
    }

    _, err = proof.Values()

    if err == nil {
        test.Error("[proof.go]", "[values]", "Proof Values() does not yield an error on a proof with no steps.")
    }
}
