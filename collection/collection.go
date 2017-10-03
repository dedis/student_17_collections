package collection

import "errors"

type collection struct {
    root *node
    values []Value
}

// Constructors

func EmptyCollection(values... Value) (collection collection) {
    collection.root = new(node)
    collection.root.children.left = new(node)
    collection.root.children.right = new(node)

    collection.root.known = true
    collection.root.children.left.known = true
    collection.root.children.right.known = true

    collection.root.children.left.parent = collection.root
    collection.root.children.right.parent = collection.root

    collection.values = values

    collection.root.children.left.values = make([][]byte, len(collection.values))
    collection.root.children.right.values = make([][]byte, len(collection.values))

    for index := 0; index < len(collection.values); index++ {
        collection.root.children.left.values[index] = collection.values[index].Placeholder()
        collection.root.children.right.values[index] = collection.values[index].Placeholder()
    }

    collection.update(collection.root.children.left)
    collection.update(collection.root.children.right)
    collection.update(collection.root)

    return
}

// Methods

func (this *collection) update(node *node) error {
    if !(node.known) {
        return errors.New("Update: updating an unknown node.")
    }

    if node.leaf() {
        node.label = sha256(true, node.key, node.values)
    } else {
        if !(node.children.left.known) || !(node.children.right.known) {
            return errors.New("Update: updating internal node with unknown children.")
        }

        node.values = make([][]byte, len(this.values))

        for index := 0; index < len(this.values); index++ {
            node.values[index] = this.values[index].Parent(node.children.left.values[index], node.children.right.values[index])
        }

        node.label = sha256(false, node.values, node.children.left.label[:], node.children.right.label[:])
    }

    return nil
}
