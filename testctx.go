package collection

import "testing"

// testctxstruct

type testctxstruct struct {
    file string
    test *testing.T

    verify testctxverifier
}

// Constructors

func testctx(file string, test *testing.T) testctxstruct {
    return testctxstruct{file, test, testctxverifier{file, test}}
}

// Methods

func (this testctxstruct) should_panic(prefix string, function func()) {
    defer func() {
        if recover() == nil {
            this.test.Error(this.file, prefix, "Function provided did not panic.")
        }
    }()

    function()
}

// testctxverifier

type testctxverifier struct {
    file string
    test *testing.T
}

// Methods

func (this testctxverifier) node(prefix string, collection *collection, node *node) {
    if !(node.known) {
        return
    }

    if node.leaf() {
        if (node.children.left != nil) || (node.children.right != nil) {
            this.test.Error(this.file, prefix, "Leaf node has one or more children.")
            return
        }

        if node.label != sha256(true, node.key, node.values) {
            this.test.Error(this.file, prefix, "Wrong leaf node label.")
            return
        }
    } else {
        if (node.children.left == nil) || (node.children.right == nil) {
            this.test.Error(this.file, prefix, "Internal node is missing one or more children.")
            return
        }

        if (node.children.left.parent != node) || (node.children.right.parent != node) {
            this.test.Error(this.file, prefix, "Children of internal node don't have its parent correctly set.")
            return
        }

        if node.label != sha256(false, node.values, node.children.left.label[:], node.children.right.label[:]) {
            this.test.Error(this.file, prefix, "Wrong internal node label.")
            return
        }

        if node.children.left.known && node.children.right.known {
            for index := 0; index < len(collection.fields); index++ {
                value := collection.fields[index].Parent(node.children.left.values[index], node.children.right.values[index])
                if !equal(value, node.values[index]) {
                    this.test.Error(this.file, prefix, "One or more internal node values conflict with the corresponding children values.")
                    return
                }
            }
        }
    }
}

func (this testctxverifier) treerecursion(prefix string, collection *collection, node *node, path []bool) {
    this.node(prefix, collection, node)

    if node.leaf() {
        if !(node.placeholder()) {
            for index := 0; index < len(path); index++ {
                keyhash := sha256(node.key)
                if path[index] != bit(keyhash[:], index) {
                    this.test.Error(this.file, prefix, "Leaf node on wrong path.")
                }
            }
        }
    } else {
        leftpath := make([]bool, len(path))
        rightpath := make([]bool, len(path))

        copy(leftpath, path)
        copy(rightpath, path)

        leftpath = append(leftpath, false)
        rightpath = append(rightpath, true)

        this.treerecursion(prefix, collection, node.children.left, leftpath)
        this.treerecursion(prefix, collection, node.children.right, rightpath)
    }
}

func (this testctxverifier) tree(prefix string, collection *collection) {
    this.treerecursion(prefix, collection, collection.root, []bool{})
}

func (this testctxverifier) keyrecursion(key []byte, node *node) bool {
    if node.leaf() {
        return equal(node.key, key)
    } else {
        return (this.keyrecursion(key, node.children.left)) || (this.keyrecursion(key, node.children.right))
    }
}

func (this testctxverifier) key(prefix string, collection *collection, key []byte) {
    if !(this.keyrecursion(key, collection.root)) {
        this.test.Error(this.file, prefix, "Node not found.")
    }
}
