package collection

import "testing"

func TestNodeGetters(test *testing.T) {
    root := node{}
    left := node{}
    right := node{}

    root.children.left = &left
    root.children.right = &right

    left.parent = &root
    right.parent = &root

    if !(root.root()) {
        test.Error("[root]", "Getter root() of root node returns false.")
    }

    if left.root() || right.root() {
        test.Error("[root]", "Getter root() of non-root node returns true.")
    }

    if root.leaf() {
        test.Error("[leaf]", "Getter leaf() of non-leaf node returns true.")
    }

    if !(left.leaf()) || !(right.leaf()) {
        test.Error("[leaf]", "Getter leaf() of leaf node returns false.")
    }

    if root.placeholder() {
        test.Error("[placeholder]", "Getter placeholder() of non-leaf node returns true.")
    }

    if !(left.placeholder()) || !(right.placeholder()) {
        test.Error("[placeholder]", "Getter placeholder() of placeholder leaf node returns false.")
    }

    left.key = []byte("leftkey")
    right.key = []byte("rightkey")

    if left.placeholder() || right.placeholder() {
        test.Error("[placeholder]", "Getter placeholder() of non-placeholder leaf node returns true.")
    }
}
