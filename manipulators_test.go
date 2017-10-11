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

    if len(transaction.temporary) < 512 {
        test.Error("[manipulators.go]", "[transactioncollection]", "Not enough temporary nodes listed in a collection without scope.")
    }

    ctx.should_panic("[wrongvalues]", func() {
        collection.Add([]byte("panickey"))
    })

    ctx.should_panic("[wrongvalues]", func() {
        keycollision.Add([]byte("panickey"), uint64(13))
    })
}
