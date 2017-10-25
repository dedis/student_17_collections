package collection

import "testing"
import "encoding/binary"
import "math/rand"

func TestManipulatorsAdd(test *testing.T) {
    ctx := testctx("[manipulators.go]", test)

    stake64 := Stake64{}
    collection := EmptyCollection(stake64)

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Add(key, uint64(rand.Uint32()))
        ctx.verify.tree("[stakecollection]", &collection)
    }

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        ctx.verify.key("[stakecollection]", &collection, key)
    }

    unknownroot := EmptyCollection()
    unknownroot.root.known = false

    error := unknownroot.Add([]byte("key"))
    if error == nil {
        test.Error("[manipulators.go]", "[unknownroot]", "Add should yield an error on a collection with unknown root.")
    }

    unknownrootchildren := EmptyCollection()
    unknownrootchildren.root.children.left.known = false
    unknownrootchildren.root.children.right.known = false

    error = unknownrootchildren.Add([]byte("key"))
    if error == nil {
        test.Error("[manipulators.go]", "[unknownrootchildren]", "Add should yield an error on a collection with unknown root children.")
    }

    keycollision := EmptyCollection()
    keycollision.Add([]byte("key"))

    error = keycollision.Add([]byte("key"))
    if error == nil {
        test.Error("[manipulators.go]", "[keycollision]", "Add should yield an error on key collision.")
    }

    transaction := EmptyCollection(stake64)
    transaction.Scope.None()
    transaction.transaction = true

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))
        transaction.Add(key, uint64(rand.Uint32()))
    }

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))
        ctx.verify.key("[transactioncollection]", &transaction, key)
    }

    ctx.should_panic("[wrongvalues]", func() {
        collection.Add([]byte("panickey"))
    })

    ctx.should_panic("[wrongvalues]", func() {
        keycollision.Add([]byte("panickey"), uint64(13))
    })
}

func TestManipulatorsSet(test *testing.T) {
    ctx := testctx("[manipulators.go]", test)

    stake64 := Stake64{}
    collection := EmptyCollection(stake64)

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Add(key, uint64(index))
    }

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Set(key, uint64(2 * index))
        ctx.verify.tree("[stakecollection]", &collection)
    }

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        ctx.verify.values("[set]", &collection, key, uint64(index * 2))
    }

    unknownroot := EmptyCollection(stake64)
    unknownroot.root.known = false

    error := unknownroot.Set([]byte("key"), uint64(0))
    if error == nil {
        test.Error("[manipulators.go]", "[unknownroot]", "Set should yield an error on a collection with unknown root.")
    }

    unknownrootchildren := EmptyCollection(stake64)
    unknownrootchildren.root.children.left.known = false
    unknownrootchildren.root.children.right.known = false

    error = unknownrootchildren.Set([]byte("key"), uint64(0))
    if error == nil {
        test.Error("[manipulators.go]", "[unknownrootchildren]", "Set should yield an error on a collection with unknown root children.")
    }

    error = collection.Set([]byte("key"), uint64(13))
    if error == nil {
        test.Error("[manipulators.go]", "[notfound]", "Set should yield error when prompted to alter a value that does not exist.")
    }

    transaction := EmptyCollection(stake64)
    transaction.Scope.None()
    transaction.transaction = true

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))
        transaction.Add(key, uint64(index))
    }

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))
        transaction.Set(key, uint64(2 * index))
        transaction.Set(key, Same{})
    }

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))
        ctx.verify.values("[transactioncollection]", &transaction, key, uint64(2 * index))
    }

    ctx.should_panic("[wrongvalues]", func() {
        collection.Set([]byte("panickey"))
    })

    ctx.should_panic("[wrongvalues]", func() {
        collection.Set([]byte("panickey"), uint64(13), uint64(44))
    })
}

func TestManipulatorsSetField(test *testing.T) {
    ctx := testctx("[manipulators.go]", test)

    data := Data{}
    collection := EmptyCollection(data, data)

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Add(key, []byte{}, []byte{})
    }

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.SetField(key, index % 2, []byte("x"))
    }

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        if index % 2 == 0 {
            ctx.verify.values("[setfield]", &collection, key, []byte("x"), []byte{})
        } else {
            ctx.verify.values("[setfield]", &collection, key, []byte{}, []byte("x"))
        }
    }

    ctx.should_panic("[fieldoutofrange]", func() {
        collection.SetField([]byte("key"), 5, []byte("data"))
    })
}

func TestManipulatorsRemove(test *testing.T) {
    ctx := testctx("[manipulators.go]", test)

    collection := EmptyCollection()
    reference := EmptyCollection()

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Add(key)
    }

    collection.Add([]byte("test"))
    reference.Add([]byte("test"))

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Remove(key)
        ctx.verify.tree("[remove]", &collection)
    }

    if collection.root.label != reference.root.label {
        test.Error("[manipulators.go]", "[remove]", "Label is not path-independent.")
    }

    unknownroot := EmptyCollection()
    unknownroot.root.known = false

    error := unknownroot.Remove([]byte("key"))
    if error == nil {
        test.Error("[manipulators.go]", "[unknownroot]", "Remove should yield an error on a collection with unknown root.")
    }

    unknownrootchildren := EmptyCollection()
    unknownrootchildren.root.children.left.known = false
    unknownrootchildren.root.children.right.known = false

    error = unknownrootchildren.Remove([]byte("key"))
    if error == nil {
        test.Error("[manipulators.go]", "[unknownrootchildren]", "Remove should yield an error on a collection with unknown root children.")
    }

    error = collection.Remove([]byte("wrongkey"));

    if error == nil {
        test.Error("[manipulators.go]", "[notfound]", "Remove should yield an error when provided with a key that does not lie in the collection.")
    }

    transaction := EmptyCollection()

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))
        transaction.Add(key)
    }

    transaction.Scope.None()
    transaction.transaction = true

    for index := 0; index < 512; index += 2 {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))
        transaction.Remove(key)
    }

    for index := 0; index < 512; index += 2 {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))
        ctx.verify.nokey("[transactioncollection]", &transaction, key)
    }

    for index := 1; index < 512; index += 2 {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))
        ctx.verify.key("[transactioncollection]", &transaction, key)
    }

    for index := 1; index < 512; index += 2 {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))
        transaction.Remove(key)
    }

    collision := func(key []byte, bits int) []byte {
        target := sha256(key)
        sample := make([]byte, 8)

        for index := 0;; index++ {
            binary.BigEndian.PutUint64(sample, uint64(index))
            hash := sha256(sample)
            if match(hash[:], target[:], bits) {
                return sample
            }
        }
    }

    collisionkey := []byte("mykey")
    collidingkey := collision(collisionkey, 8)

    transaction.Add(collisionkey)
    transaction.Add(collidingkey)

    transaction.Remove(collisionkey)
    transaction.Remove(collidingkey)

    transaction.fix()
    transaction.Collect()
    transaction.transaction = false

    empty := EmptyCollection()
    if transaction.root.label != empty.root.label {
        test.Error("[manipulators.go]", "[transaction]", "Transaction on collection doesn't produce empty root after removing all records.")
    }
}
