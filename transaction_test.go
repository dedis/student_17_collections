package collection

import "testing"

func TestTransactionCollect(test *testing.T) {
    collection := EmptyCollection()

    collection.root.children.left.branch()
    collection.root.children.right.branch()
    collection.root.children.right.children.left.branch()
    collection.root.children.right.children.right.branch()

    collection.temporary = append(collection.temporary, collection.root.children.left, collection.root.children.right)
    collection.collect()

    if collection.root.children.left.known || collection.root.children.right.known {
        test.Error("[known]", "Collect doesn't make temporary nodes unknown.")
    }

    if (collection.root.children.left.children.left != nil) || (collection.root.children.left.children.right != nil) || (collection.root.children.right.children.left != nil) || (collection.root.children.right.children.right != nil) {
        test.Error("[children]", "Collect does not prune children of temporary nodes.")
    }

    if len(collection.temporary) != 0 {
        test.Error("[temporary]", "Collect does not empty temporary nodes list.")
    }
}
