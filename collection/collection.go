package collection

import "errors"

type collection struct {
    root *node
    values []Value
    Scope scope

    transaction bool
    temporary []*node
}

// Constructors

func EmptyCollection(values... Value) (collection collection) {
    collection.Scope.All()

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

func (this *collection) Apply(update update) error {
    switch update.kind {
    case add:
        return this.applyadd(update.key, update.values) // TODO: use proofs if available
    case remove:
        return this.applyremove(update.key) // TODO: use proofs if available
    case set:
        return this.applyset(update.key, update.values) // TODO: use proofs if available
    }

    panic("Wrong update kind value.")
}

func (this *collection) Begin() {
    this.transaction = true
}

func (this *collection) End() {
    this.fix(this.root)
    this.collect()
    this.transaction = false
}

// Private methods

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

func (this *collection) fix(node *node) {
    if node.inconsistent {
        if node.leaf() {
            this.update(node)
        } else {
            this.fix(node.children.left)
            this.fix(node.children.right)
            this.update(node)
        }

        node.inconsistent = false
    }
}

func (this *collection) collect() {
    for index := 0; index < len(this.temporary); index++ {
        this.temporary[index].known = false

        this.temporary[index].key = []byte{}
        this.temporary[index].values = [][]byte{}

        this.temporary[index].children.left = nil
        this.temporary[index].children.right = nil

        this.temporary[index] = nil
    }

    this.temporary = this.temporary[:0]
}

func (this *collection) applyadd(key []byte, values [][]byte) error {
    return errors.New("Not implemented.")
}

func (this *collection) applyremove(key []byte) error {
    return errors.New("Not implemented.")
}

func (this *collection) applyset(key []byte, values [][]byte) error {
    return errors.New("Not implemented.")
}
