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
    temporary []*node
}
