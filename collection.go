package collection

// Structs

type collection struct {
    root *node
    fields []Field
    Scope scope

    AutoCollect flag
    transaction struct {
        ongoing bool
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
