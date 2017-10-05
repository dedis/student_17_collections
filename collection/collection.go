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

    collection.root.children.left.values = collection.placeholdervalues()
    collection.root.children.right.values = collection.placeholdervalues()

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

func (this *collection) placeholdervalues() [][]byte {
    values := make([][]byte, len(this.values))

    for index := 0; index < len(this.values); index++ {
        values[index] = this.values[index].Placeholder()
    }

    return values
}

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
    path := sha256(key)
    store := this.Scope.match(path)

    depth := 0
    cursor := this.root

    if !(cursor.known) {
        return errors.New("Applying update to unknown subtree. Proof needed.") // TODO: first check if a proof was provided. If so, and the proof is valid, use that to expand the tree with nodes from the proof, setting them temporary only if the key lies outside the scope of the collection. If the proof is absent or invalid, return this error.
    }

    for {
        step := bit(path[:], depth)

        if !(cursor.children.left.known) || !(cursor.children.right.known) {
            return errors.New("Applying update to unknown subtree. Proof needed.") // TODO: first check if a proof was provided. If so, and the proof is valid, use that to expand the tree with nodes from the proof, setting them temporary only if the key lies outside the scope of the collection. If the proof is absent or invalid, return this error.
        }

        if step {
            cursor = cursor.children.right
        } else {
            cursor = cursor.children.left
        }

        depth++

        if cursor.placeholder() {
            cursor.key = key
            cursor.values = values
            this.update(cursor)

            break
        } else if cursor.leaf() {
            if (len(key) == len(cursor.key)) && match(key, cursor.key, 8 * len(key)) {
                return errors.New("Key collision.")
            }

            collision := *cursor
            collisionpath := sha256(collision.key)
            collisionstep := bit(collisionpath[:], depth)

            cursor.key = []byte{}
            cursor.children.left = new(node)
            cursor.children.right = new(node)

            cursor.children.left.known = true
            cursor.children.right.known = true

            cursor.children.left.parent = cursor
            cursor.children.right.parent = cursor

            if collisionstep {
                cursor.children.right.key = collision.key
                cursor.children.right.values = collision.values
                cursor.children.left.values = this.placeholdervalues()
            } else {
                cursor.children.left.key = collision.key
                cursor.children.left.values = collision.values
                cursor.children.right.values = this.placeholdervalues()
            }

            if !store {
                this.temporary = append(this.temporary, cursor.children.left, cursor.children.right)
            }
        }
    }

    for {
        if cursor.parent == nil {
            break
        }

        cursor = cursor.parent

        if this.transaction {
            cursor.inconsistent = true
        } else {
            this.update(cursor)
        }
    }

    if !(this.transaction) {
        this.collect()
    }

    return nil
}

func (this *collection) applyremove(key []byte) error {
    return errors.New("Not implemented.")
}

func (this *collection) applyset(key []byte, values [][]byte) error {
    return errors.New("Not implemented.")
}
