package collection

type record struct {
    collection *collection
    key []byte
    values [][]byte
}

// Getters

func (this *record) Exists() bool {
    return len(this.key) > 0
}

func (this *record) Key() []byte {
    return this.key
}

func (this *record) Values() [][]byte {
    return this.values
}
