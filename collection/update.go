package collection

type updatekind uint

const(
    add updatekind = iota
    remove
    set
)

type update struct {
    kind updatekind
    key []byte
    values [][]byte

    // TODO: Add proofs: updates should be served also by verifiers.
}

// Constructors

func AddUpdate(key []byte, values [][]byte) update {
    return update{add, key, values}
}

func RemoveUpdate(key []byte) update {
    return update{remove, key, [][]byte{}}
}

func SetUpdate(key []byte, values [][]byte) update {
    return update{set, key, values}
}
