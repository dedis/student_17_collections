package collection

import csha256 "crypto/sha256"

type dump struct {
    label [csha256.Size]byte
    leaf bool

    key []byte
    values [][]byte
}

type step struct {
    left dump
    right dump
}

type proof struct {
    root [csha256.Size]byte
    key []byte
    steps []step
}

// Getters

func (this *proof) Key() []byte {
    return this.key
}

func (this *proof) Values() [][]byte{
    if len(this.steps) == 0 {
        return [][]byte{}
    }

    path := sha256(this.key)
    depth := len(this.steps) - 1

    if bit(path[:], depth) {
        if (len(this.steps[depth].right.key) == len(this.key)) && match(this.steps[depth].right.key, this.key, 8 * len(this.key)) {
            return this.steps[depth].right.values
        } else {
            return [][]byte{}
        }
    } else {
        if (len(this.steps[depth].left.key) == len(this.key)) && match(this.steps[depth].left.key, this.key, 8 * len(this.key)) {
            return this.steps[depth].left.values
        } else {
            return [][]byte{}
        }
    }
}

// Methods

func (this *proof) Match() bool {
    if len(this.steps) == 0 {
        return false
    }

    path := sha256(this.key)
    depth := len(this.steps) - 1

    if bit(path[:], depth) {
        return (len(this.steps[depth].right.key) == len(this.key)) && match(this.steps[depth].right.key, this.key, 8 * len(this.key))
    } else {
        return (len(this.steps[depth].left.key) == len(this.key)) && match(this.steps[depth].left.key, this.key, 8 * len(this.key))
    }
}
