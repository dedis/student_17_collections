package collection

import "errors"
import csha256 "crypto/sha256"

// dump

type dump struct {
    label [csha256.Size]byte

    key []byte
    values [][]byte

    children struct {
        left [csha256.Size]byte
        right [csha256.Size]byte
    }
}

// Constructors

func dumpnode(node *node) (dump dump) {
    if node.leaf() {
        dump.label = node.label
        dump.key = node.key
        dump.values = node.values
    } else {
        dump.label = node.label
        dump.values = node.values
        dump.children.left = node.children.left.label
        dump.children.right = node.children.right.label
    }

    return
}

// Getters

func (this *dump) leaf() bool {
    var empty [csha256.Size]byte
    return (this.children.left == empty) && (this.children.right == empty)
}

// Methods

func (this *dump) consistent() bool {
    if this.leaf() {
        return this.label == sha256(true, this.key[:], this.values)
    } else {
        return this.label == sha256(false, this.values, this.children.left[:], this.children.right[:])
    }
}

// step

type step struct {
    left dump
    right dump
}

// Proof

type Proof struct {
    collection *collection
    key []byte
    
    root dump
    steps []step
}

// Getters

func (this *Proof) Key() []byte {
    return this.key
}

// Methods

func (this *Proof) Match() bool {
    if len(this.steps) == 0 {
        return false
    }

    path := sha256(this.key)
    depth := len(this.steps) - 1

    if bit(path[:], depth) {
        return equal(this.key, this.steps[depth].right.key)
    } else {
        return equal(this.key, this.steps[depth].left.key)
    }
}

func (this *Proof) Values() ([]interface{}, error) {
    if len(this.steps) == 0 {
        return []interface{}{}, errors.New("Proof has no steps.")
    }

    path := sha256(this.key)
    depth := len(this.steps) - 1

    match := false
    var rawvalues [][]byte

    if bit(path[:], depth) {
        if equal(this.key, this.steps[depth].right.key) {
            match = true
            rawvalues = this.steps[depth].right.values
        }
    } else {
        if equal(this.key, this.steps[depth].left.key) {
            match = true
            rawvalues = this.steps[depth].left.values
        }
    }

    if !match {
        return []interface{}{}, errors.New("No match found.")
    }

    if len(rawvalues) != len(this.collection.fields) {
        return []interface{}{}, errors.New("Wrong number of values.")
    }

    var values []interface{}

    for index := 0; index < len(rawvalues); index++ {
        value, err := this.collection.fields[index].Decode(rawvalues[index])

        if err != nil {
            return []interface{}{}, err
        }

        values = append(values, value)
    }

    return values, nil
}
