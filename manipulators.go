package collection

import "errors"

// Private methods (collection) (manipulators)

func (this *collection) Add(key []byte, values... interface{}) error {
    if len(values) != len(this.fields) {
        panic("Wrong number of values provided.")
    }

    rawvalues := make([][]byte, len(this.fields))
    for index := 0; index < len(this.fields); index++ {
        rawvalues[index] = this.fields[index].Encode(values[index])
    }

    path := sha256(key)
    store := this.Scope.match(path)

    depth := 0
    cursor := this.root

    if !(cursor.known) {
        return errors.New("Applying update to unknown subtree. Proof needed.")
    }

    for {
        if !(cursor.children.left.known) || !(cursor.children.right.known) {
            return errors.New("Applying update to unknown subtree. Proof needed.")
        }

        step := bit(path[:], depth)
        depth++

        if step {
            cursor = cursor.children.right
        } else {
            cursor = cursor.children.left
        }

        if cursor.placeholder() {
            if this.transaction {
                cursor.backup()
            }

            cursor.key = key
            cursor.values = rawvalues
            this.update(cursor)

            break
        } else if cursor.leaf() {
            if equal(key, cursor.key) {
                return errors.New("Key collision.")
            }

            collision := *cursor
            collisionpath := sha256(collision.key)
            collisionstep := bit(collisionpath[:], depth)

            if this.transaction {
                cursor.backup()
            }

            cursor.key = []byte{}
            cursor.branch()

            if collisionstep {
                cursor.children.right.known = true
                cursor.children.right.label = collision.label
                cursor.children.right.key = collision.key
                cursor.children.right.values = collision.values

                this.placeholder(cursor.children.left)
            } else {
                cursor.children.left.known = true
                cursor.children.left.label = collision.label
                cursor.children.left.key = collision.key
                cursor.children.left.values = collision.values

                this.placeholder(cursor.children.right)
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
            cursor.transaction.inconsistent = true
        } else {
            this.update(cursor)
        }
    }

    if !(this.transaction) {
        this.collect()
    }

    return nil
}

func (this *collection) Set(key []byte, values... interface{}) error {
    if len(values) != len(this.fields) {
        panic("Wrong number of values provided.")
    }

    rawvalues := make([][]byte, len(this.fields))
    for index := 0; index < len(this.fields); index++ {
        rawvalues[index] = this.fields[index].Encode(values[index])
    }

    path := sha256(key)

    depth := 0
    cursor := this.root

    if !(cursor.known) {
        return errors.New("Applying update to unknown subtree. Proof needed.")
    }

    for {
        if !(cursor.children.left.known) || !(cursor.children.right.known) {
            return errors.New("Applying update to unknown subtree. Proof needed.")
        }

        step := bit(path[:], depth)
        depth++

        if step {
            cursor = cursor.children.right
        } else {
            cursor = cursor.children.left
        }

        if cursor.leaf() {
            if !(equal(cursor.key, key)) {
                return errors.New("Key not found.")
            } else {
                if this.transaction {
                    cursor.backup()
                }

                cursor.values = rawvalues
                this.update(cursor)

                break
            }
        }
    }

    for {
        if cursor.parent == nil {
            break
        }

        cursor = cursor.parent

        if this.transaction {
            cursor.transaction.inconsistent = true
        } else {
            this.update(cursor)
        }
    }

    if !(this.transaction) {
        this.collect()
    }

    return nil
}
