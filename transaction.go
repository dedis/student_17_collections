package collection

// Methods (collection) (transaction methods)

func (this *collection) Begin() {
    if this.transaction {
        panic("Transaction already in progress.")
    }

    this.transaction = true
}

func (this *collection) Rollback() {
    if !(this.transaction) {
        panic("Transaction not in progress")
    }

    var explore func(*node)
    explore = func(node *node) {
        if node.transaction.inconsistent || (node.transaction.backup != nil) {
            node.restore()

            if !(node.leaf()) {
                explore(node.children.left)
                explore(node.children.right)
            }
        }
    }

    explore(this.root)

    for index := 0; index < len(this.temporary); index++ {
        this.temporary[index] = nil
    }

    this.temporary = this.temporary[:0]
    this.transaction = false
}

func (this *collection) End() {
    if !(this.transaction) {
        panic("Transaction not in progress.")
    }

    this.confirm()
    this.fix()
    this.collect()

    this.transaction = false
}

// Private methods (collection) (transaction methods)

func (this *collection) confirm() {
    var explore func(*node)
    explore = func(node *node) {
        if node.transaction.inconsistent || (node.transaction.backup != nil) {
            node.transaction.backup = nil

            if !(node.leaf()) {
                explore(node.children.left)
                explore(node.children.right)
            }
        }
    }

    explore(this.root)
}

func (this *collection) fix() {
    var explore func(*node)
    explore = func(node *node) {
        if node.transaction.inconsistent {
            if !(node.leaf()) {
                explore(node.children.left)
                explore(node.children.right)
            }

            this.update(node)
            node.transaction.inconsistent = false
        }
    }

    explore(this.root)
}

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
