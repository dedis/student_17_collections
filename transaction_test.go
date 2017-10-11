package collection

import "testing"
import "encoding/binary"

func TestTransactionBegin(test *testing.T) {
    ctx := testctx("[transaction.go]", test)

    collection := EmptyCollection()
    collection.Begin()

    if !(collection.transaction) {
        test.Error("[transaction.go]", "[begin]", "Begin() does not set the transaction flag.")
    }

    ctx.should_panic("[begin]", func() {
        collection.Begin()
    })
}

func TestTransactionRollback(test *testing.T) {
    ctx := testctx("[transaction.go]", test)

    collection := EmptyCollection()
    reference := EmptyCollection()

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Add(key)
        reference.Add(key)
    }

    collection.Begin()

    for index := 512; index < 1024; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Add(key)
    }

    collection.Rollback()
    ctx.verify.tree("[rollback]", &collection)

    if collection.root.label != reference.root.label {
        test.Error("[transaction.go]", "[rollback]", "Rollback() doesn't produce the same tree as before.")
    }

    collection.fix()

    if collection.root.label != reference.root.label {
        test.Error("[transaction.go]", "[rollback]", "Fixing after Rollback() has a non-null effect.")
    }

    if len(collection.temporary) != 0 {
        test.Error("[transaction.go]", "[rollback]", "Rollback() does not empty the temporary nodes slice.")
    }

    ctx.should_panic("[rollbackagain]", func() {
        collection.Rollback()
    })
}

func TestTransactionEnd(test *testing.T) {
    ctx := testctx("[transaction.go]", test)

    collection := EmptyCollection()

    collection.Begin()

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Add(key)
    }

    collection.End()

    ctx.verify.tree("[end]", &collection)

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        ctx.verify.key("[end]", &collection, key)
    }

    oldroot := collection.root.label
    collection.fix()

    if collection.root.label != oldroot {
        test.Error("[transaction.go]", "[end]", "Fixing after End() alters the tree root.")
    }

    if len(collection.temporary) != 0 {
        test.Error("[transaction.go]", "[end]", "End() does not empty the temporary nodes slice.")
    }

    ctx.should_panic("[endagain]", func() {
        collection.End()
    })
}

func TestTransactionConfirm(test *testing.T) {
    collection := EmptyCollection()
    reference := EmptyCollection()

    collection.transaction = true
    reference.transaction = true

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Add(key)
        reference.Add(key)
    }

    var explore func(*node) int
    explore = func(node *node) int {
        if node.leaf() {
            if node.transaction.backup != nil {
                return 1
            } else {
                return 0
            }
        } else {
            if node.transaction.backup != nil {
                return 1 + explore(node.children.left) + explore(node.children.right)
            } else {
                return explore(node.children.left) + explore(node.children.right)
            }
        }
    }

    count := explore(collection.root)
    if count < 512 {
        test.Error("[transaction.go]", "[backup]", "Not enough backups after transaction operations.")
    }

    collection.confirm()

    count = explore(collection.root)
    if count != 0 {
        test.Error("[transaction.go]", "[confirm]", "confirm() does not remove all the backups.")
    }

    collection.fix()
    reference.fix()

    if collection.root.label != reference.root.label {
        test.Error("[transaction.go]", "[confirm]", "confirm() does not only remove the backups, but it also alters the values of the nodes.")
    }
}

func TestTransactionFix(test *testing.T) {
    ctx := testctx("[transaction.go]", test)

    collection := EmptyCollection()

    collection.root.children.left.key = []byte("leftkey")
    collection.root.children.left.transaction.inconsistent = true
    collection.root.transaction.inconsistent = true

    collection.fix()
    ctx.verify.tree("[fix]", &collection)

    oldrootlabel := collection.root.label

    collection.root.children.right.key = []byte("rightkey")
    collection.root.children.right.transaction.inconsistent = true

    collection.fix()

    if collection.root.label != oldrootlabel {
        test.Error("[transaction.go]", "[fix]", "Fix should not visit nodes that are not marked as inconsistent.")
    }

    collection.root.transaction.inconsistent = true
    collection.fix()

    if collection.root.label == oldrootlabel {
        test.Error("[transaction.go]", "[fix]", "Fix should alter the label of the root of a collection tree.")
    }

    ctx.verify.tree("[fix]", &collection)
}

func TestTransactionCollect(test *testing.T) {
    collection := EmptyCollection()

    collection.root.children.left.branch()
    collection.root.children.right.branch()
    collection.root.children.right.children.left.branch()
    collection.root.children.right.children.right.branch()

    collection.temporary = append(collection.temporary, collection.root.children.left, collection.root.children.right)
    collection.collect()

    if collection.root.children.left.known || collection.root.children.right.known {
        test.Error("[transaction.go]", "[known]", "Collect doesn't make temporary nodes unknown.")
    }

    if (collection.root.children.left.children.left != nil) || (collection.root.children.left.children.right != nil) || (collection.root.children.right.children.left != nil) || (collection.root.children.right.children.right != nil) {
        test.Error("[transaction.go]", "[children]", "Collect does not prune children of temporary nodes.")
    }

    if len(collection.temporary) != 0 {
        test.Error("[transaction.go]", "[temporary]", "Collect does not empty temporary nodes list.")
    }
}
