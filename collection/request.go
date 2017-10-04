package collection

type requestkind uint

const(
    add requestkind = iota
    remove
    update
)

type request struct {
    kind requestkind
    key []byte
    values [][]byte

    // TODO: Add proofs: requests should be served also by verifiers.
}

// Constructors

func AddRequest(key []byte, values [][]byte) request {
    return request{add, key, values}
}

func RemoveRequest(key []byte) request {
    return request{remove, key, [][]byte{}}
}

func UpdateRequest(key []byte, values [][]byte) request {
    return request{update, key, values}
}
