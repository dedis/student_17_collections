package collection

import csha256 "crypto/sha256"

// mask

type mask struct {
    value []byte
    bits int
}

// Private methods

func (this *mask) match(buffer [csha256.Size]byte) bool {
    return match(buffer[:], this.value[:], this.bits)
}

// scope

type scope struct {
    masks []mask
    all bool
}

// Methods

func (this *scope) All() {
    this.all = true
    this.masks = []mask{}
}

func (this *scope) None() {
    this.all = false
    this.masks = []mask{}
}

func (this *scope) Add(value []byte, bits int) {
    this.masks = append(this.masks, mask{value, bits})
}

// Private methods

func (this *scope) match(buffer [csha256.Size]byte) bool {
    if len(this.masks) == 0 {
        return this.all
    }

    for index := 0; index < len(this.masks); index++ {
        if this.masks[index].match(buffer) {
            return true
        }
    }

    return false
}
