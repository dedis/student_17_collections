package collection

// Interfaces

type Collection interface {
    // TODO: Add interface items here.
}

// Structs

type collection struct {
    root *node
    fields []Field
    Scope scope

    transaction bool
}

// Constructors

func EmptyCollection(fields... Field) (collection collection) {
    collection.fields = fields
    collection.Scope.All()

    collection.root = new(node)
    collection.root.known = true

    collection.root.branch()

    collection.placeholder(collection.root.children.left)
    collection.placeholder(collection.root.children.right)
    collection.update(collection.root)

    return
}

func EmptyVerifier(fields... Field) (verifier collection) {
    empty := EmptyCollection(fields...)

    verifier.fields = fields
    verifier.Scope.None()

    verifier.root = new(node)
    verifier.root.known = false
    verifier.root.label = empty.root.label

    return
}
