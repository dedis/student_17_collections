package collection

// Interfaces

type Update interface {
    Valid() bool
    Verify(func(Proof) bool) bool
    Applicable(func([]byte) (record, error)) bool
    Apply(func([]byte) (record, error), func([]byte, [][]byte) error, func([]byte, [][]byte) error, func([]byte) error)
}

// Structs

// addupdate

type addupdate struct {
    key []byte
    values [][]byte

    proofs struct {
        key Proof
    }
}

// Constructors

func AddUpdate(key []byte, values [][]byte, proof Proof) (addupdate addupdate) {
    addupdate.key = key
    addupdate.values = values
    addupdate.proofs.key = proof

    return
}

// Interface

func (this addupdate) Valid() bool {
    return (len(this.key) == len(this.proofs.key.key)) && match(this.key, this.proofs.key.key, 8 * len(this.key))
}

func (this addupdate) Verify(verify func(Proof) bool) bool {
    return verify(this.proofs.key)
}

func (this addupdate) Applicable(get func([]byte) (record, error)) bool {
    record, error := get(this.key)

    if error != nil {
        return false
    }

    if record.Exists() {
        return false
    }

    return true
}

func (this addupdate) Apply(get func([]byte) (record, error), add func([]byte, [][]byte) error, set func([]byte, [][]byte) error, remove func([]byte) error) {
    add(this.key, this.values)
}
