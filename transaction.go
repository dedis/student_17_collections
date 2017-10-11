package collection

// Private methods (collection) (transaction methods)

func (this *collection) collect() {
    for index := 0; index < len(this.temporary); index++ {
        this.temporary[index].known = false

        this.temporary[index].key = []byte{}
        this.temporary[index].values = [][]byte{}

        this.temporary[index].prune()
        this.temporary[index] = nil
    }

    this.temporary = this.temporary[:0]
}
