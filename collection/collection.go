package collection

import "errors"

// getter

type getter struct {
    collection *collection
    key []byte
}

// Methods

func (this getter) Record() (record, error) {
    return this.collection.getrecord(this.key)
}

func (this getter) Proof() (Proof, error) {
    return this.collection.getproof(this.key)
}

// collection

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

func (this *collection) Begin() {
    this.transaction = true
}

func (this *collection) End() {
    this.fix(this.root)
    this.collect()
    this.transaction = false
}

func (this *collection) Get(key []byte) getter {
    return getter{this, key}
}

func (this *collection) Add(key []byte, values [][]byte) error {
    path := sha256(key)
    store := this.Scope.match(path)

    depth := 0
    cursor := this.root

    if !(cursor.known) {
        return errors.New("Applying update to unknown subtree. Proof needed.")
    }

    for {
        step := bit(path[:], depth)

        if !(cursor.children.left.known) || !(cursor.children.right.known) {
            return errors.New("Applying update to unknown subtree. Proof needed.")
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
                cursor.children.right.label = collision.label
                cursor.children.right.key = collision.key
                cursor.children.right.values = collision.values

                cursor.children.left.values = this.placeholdervalues()
                this.update(cursor.children.left)
            } else {
                cursor.children.left.label = collision.label
                cursor.children.left.key = collision.key
                cursor.children.left.values = collision.values

                cursor.children.right.values = this.placeholdervalues()
                this.update(cursor.children.right)
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

func (this *collection) Verify(target interface{}) bool {
    switch value := target.(type) {
    case Proof:
        return this.verifyproof(value)
    case Update:
        return value.Valid() && value.Verify(this.verifyproof)
    default:
        panic("Verify accepts only targets of type Proof or Update.")
    }
}

func (this *collection) Apply(update Update) bool {
    get := func(key []byte) (record, error) {
        return this.Get(key).Record()
    }

    if update.Applicable(get) {
        set := func([]byte, [][]byte) error { // TODO: Substitute with this.Set when implemented
            return errors.New("Not implemented.")
        }

        remove := func([]byte) error { // TODO: Substitute with this.Remove when implemented
            return errors.New("Not implemented.")
        }

        update.Apply(get, this.Add, set, remove)
        return true
    } else {
        return false
    }
}

// Private methods

func (this *collection) getrecord(key []byte) (record, error) {
    path := sha256(key)

    depth := 0
    cursor := this.root

    for {
        if !(cursor.known) {
            return record{}, errors.New("Record lies in an unknown subtree.")
        }

        if cursor.leaf() {
            if (len(key) == len(cursor.key)) && match(key, cursor.key, 8 * len(key)) {
                return record{this, cursor.key, cursor.values}, nil
            } else {
                return record{this, []byte{}, [][]byte{}}, nil
            }
        } else {
            if bit(path[:], depth) {
                cursor = cursor.children.right
            } else {
                cursor = cursor.children.left
            }

            depth++
        }
    }
}

func (this *collection) getproof(key []byte) (Proof, error) {
    var proof Proof

    proof.root = this.root.label
    proof.key = key

    path := sha256(key)

    depth := 0
    cursor := this.root

    if !(cursor.known) {
        return proof, errors.New("Record lies in unknown subtree.")
    }

    for {
        if !(cursor.children.left.known) || !(cursor.children.right.known) {
            return proof, errors.New("Record lies in unknown subtree.")
        }

        var left dump
        var right dump

        left.leaf = cursor.children.left.leaf()
        left.label = cursor.children.left.label
        left.values = cursor.children.left.values

        right.leaf = cursor.children.right.leaf()
        right.label = cursor.children.right.label
        right.values = cursor.children.right.values

        if left.leaf {
            left.key = cursor.children.left.key
        } else {
            left.children.left = cursor.children.left.children.left.label
            left.children.right = cursor.children.left.children.right.label
        }

        if right.leaf {
            right.key = cursor.children.right.key
        } else {
            right.children.left = cursor.children.right.children.left.label
            right.children.right = cursor.children.right.children.right.label
        }

        proof.steps = append(proof.steps, step{left, right})

        if bit(path[:], depth) {
            cursor = cursor.children.right
        } else {
            cursor = cursor.children.left
        }

        depth++

        if cursor.leaf() {
            break
        }
    }

    return proof, nil
}

func (this *collection) verifyproof(proof Proof) bool {
    if this.root.inconsistent {
        panic("Verify called on inconsistent tree.")
    }

    path := sha256(proof.key)

    depth := 0
    cursor := this.root

    if !(cursor.known) {
        return false
    }

    for {
        if depth >= len(proof.steps) {
            return false
        }

        if !(this.match(cursor.children.left, &(proof.steps[depth].left))) || !(this.match(cursor.children.right, &(proof.steps[depth].right))) {
            return false
        }

        if bit(path[:], depth) {
            cursor = cursor.children.right
        } else {
            cursor = cursor.children.left
        }

        depth++

        if cursor.leaf() {
            break
        }
    }

    return true
}

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

func (this *collection) match(reference *node, dump *dump) bool {
    if (dump.label != reference.label) || !(dump.consistent()) {
        return false
    }

    if !(reference.known) {
        this.temporary = append(this.temporary, reference)
        reference.known = true

        reference.values = dump.values

        if dump.leaf {
            reference.key = dump.key
        } else {
            reference.children.left = new(node)
            reference.children.left.parent = reference
            reference.children.left.known = false
            reference.children.left.label = dump.children.left

            reference.children.right = new(node)
            reference.children.right.parent = reference
            reference.children.right.known = false
            reference.children.right.label = dump.children.right

            this.temporary = append(this.temporary, reference.children.left, reference.children.right)
        }
    }

    return true
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
