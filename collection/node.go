package collection

import csha256 "crypto/sha256"

type node struct {
    label [csha256.Size]byte

    known bool
    inconsistent bool

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

func (this *node) placeholder() bool {
    return this.leaf() && len(this.key) == 0
}
