package collection

import csha256 "crypto/sha256"

// Interfaces

type userupdate interface {
    Records() []Proof
    Check(ReadOnly) bool
    Apply(ReadWrite)
}

type ReadOnly interface {
    Get([]byte) Record
}

type ReadWrite interface {
    Get([]byte) Record
    Add([]byte, ... interface{}) error
    Set([]byte, ... interface{}) error
    SetField([]byte, int, interface{}) error
    Remove([]byte) error
}

// Structs

type proxy struct {
    collection *collection
    paths map[[csha256.Size]byte]bool
}

// Constructors

func (this *collection) proxy(keys [][]byte) (proxy proxy) {
    proxy.collection = this
    proxy.paths = make(map[[csha256.Size]byte]bool)

    for index := 0; index < len(keys); index++ {
        proxy.paths[sha256(keys[index])] = true
    }

    return
}

// Methods

func (this proxy) Get(key []byte) Record {
    if !(this.has(key)) {
        panic("Accessing undeclared key from update.")
    }

    record, _ := this.collection.Get(key).Record()
    return record
}

func (this proxy) Add(key []byte, values... interface{}) error {
    if !(this.has(key)) {
        panic("Accessing undeclared key from update.")
    }

    return this.collection.Add(key, values...)
}

func (this proxy) Set(key []byte, values... interface{}) error {
    if !(this.has(key)) {
        panic("Accessing undeclared key from update.")
    }

    return this.collection.Set(key, values...)
}

func (this proxy) SetField(key []byte, field int, value interface{}) error {
    if !(this.has(key)) {
        panic("Accessing undeclared key from update.")
    }

    return this.collection.SetField(key, field, value)
}

func (this proxy) Remove(key []byte) error {
    if !(this.has(key)) {
        panic("Accessing undeclared key from update.")
    }

    return this.collection.Remove(key)
}

// Private methods

func (this proxy) has(key []byte) bool {
    path := sha256(key)
    return this.paths[path]
}
