package collection

// Structs

type collection struct {
    root *node
    fields []Field
    Scope scope

    AutoCollect flag
    transaction struct {
        ongoing bool
        id uint64
    }
}

// Constructors

func EmptyCollection(fields... Field) (collection collection) {
    collection.fields = fields

    collection.Scope.All()
    collection.AutoCollect.Enable()

    collection.root = new(node)
    collection.root.known = true

    collection.root.branch()

    collection.placeholder(collection.root.children.left)
    collection.placeholder(collection.root.children.right)
    collection.update(collection.root)

    return
}

func EmptyVerifier(fields... Field) (verifier collection) {
    verifier.fields = fields

    verifier.Scope.None()
    verifier.AutoCollect.Enable()

    empty := EmptyCollection(fields...)

    verifier.root = new(node)
    verifier.root.known = false
    verifier.root.label = empty.root.label

    return
}

// Methods

func (this *collection) Clone() (collection collection) {
    if this.transaction.ongoing {
        panic("Cannot clone a collection while a transaction is ongoing.")
    }

    collection.root = new(node)

    collection.fields = make([]Field, len(this.fields))
    copy(collection.fields, this.fields)

    collection.Scope = this.Scope.clone()
    collection.AutoCollect = this.AutoCollect

    collection.transaction.ongoing = false
    collection.transaction.id = 0

    var explore func(*node, *node)
    explore = func(dstcursor *node, srccursor *node) {
        dstcursor.label = srccursor.label
        dstcursor.known = srccursor.known

        dstcursor.transaction.inconsistent = false
        dstcursor.transaction.backup = nil

        dstcursor.key = srccursor.key
        dstcursor.values = make([][]byte, len(srccursor.values))
        copy(dstcursor.values, srccursor.values)

        if !(srccursor.leaf()) {
            dstcursor.branch()
            explore(dstcursor.children.left, srccursor.children.left)
            explore(dstcursor.children.right, srccursor.children.right)
        }
    }

    explore(collection.root, this.root)

    return
}
