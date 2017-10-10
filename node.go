package collection

import csha256 "crypto/sha256"

type node struct {
    label [csha256.Size]byte

    known bool

    transaction struct {
        inconsistent bool
        backup *node
    }

    key []byte
    values [][]byte

    parent *node
    children struct {
        left *node
        right *node
    }
}

// Getters

func (this *node) root() bool {
    if this.parent == nil {
        return true
    } else {
        return false
    }
}

func (this *node) leaf() bool {
    if this.children.left == nil {
        return true
    } else {
        return false
    }
}

func (this * node) placeholder() bool {
    return this.leaf() && (len(this.key) == 0)
}

// Methods

func (this *node) backup() {
    if this.transaction.backup == nil {
        backup := new(node)
        (*backup) = (*this)
        this.transaction.backup = backup
    }
}

func (this *node) restore() {
    if this.transaction.backup != nil {
        backup := this.transaction.backup
        (*this) = (*backup)
        this.transaction.backup = nil
    }
}

func (this *node) branch() {
    this.children.left = new(node)
    this.children.right = new(node)

    this.children.left.parent = this
    this.children.right.parent = this
}

func (this *node) prune() {
    this.children.left = nil
    this.children.right = nil
}
