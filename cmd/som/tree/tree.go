package tree

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Named is an interface for stuff that has a name
type Named interface {
	GetName() string
}

// MapNode is a node in the tree data structure
type MapNode[T Named] struct {
	Parent   *MapNode[T]
	Children orderedMap[*MapNode[T]]
	Value    T
	Depth    int
}

// NewNode creates a new tree node
func NewNode[T Named](value T, depth int) *MapNode[T] {
	return &MapNode[T]{
		Children: newOrderedMap[*MapNode[T]](),
		Value:    value,
		Depth:    depth,
	}
}

// MapTree is a tree data structure
type MapTree[T Named] struct {
	Root  *MapNode[T]
	Nodes orderedMap[*MapNode[T]]
}

// NewTree creates a new tree node
func NewTree[T Named](value T) *MapTree[T] {
	root := NewNode(value, 0)
	nodes := newOrderedMap[*MapNode[T]]()
	nodes.Set(value.GetName(), root)

	return &MapTree[T]{
		Root:  root,
		Nodes: nodes,
	}
}

// Ancestors returns a slice of all ancestors (i.e. recursive parents),
// and an ok bool whether the requested node was found.
//
// The first ancestor is the direct parent, while the last ancestor is the root node.
func (t *MapTree[T]) Ancestors(name string) ([]*MapNode[T], bool) {
	res := []*MapNode[T]{}
	curr, ok := t.Nodes.Get(name)
	if !ok {
		return res, false
	}
	for curr.Parent != nil {
		res = append(res, curr.Parent)
		curr = curr.Parent
	}
	return res, true
}

// Descendants returns a slice of all descendants (i.e. recursive children),
// and an ok bool whether the requested node was found.
//
// Descendants in the returned slice have undefined order.
func (t *MapTree[T]) Descendants(name string) ([]*MapNode[T], bool) {
	res := []*MapNode[T]{}
	curr, ok := t.Nodes.Get(name)
	if !ok {
		return res, false
	}
	desc := t.descendants(curr, res)
	return desc, true
}

func (t *MapTree[T]) descendants(n *MapNode[T], res []*MapNode[T]) []*MapNode[T] {
	for _, name := range n.Children.Keys() {
		child, _ := n.Children.Get(name)
		res = append(res, child)
		res = t.descendants(child, res)
	}
	return res
}

// AddTree adds a sub-tree without children
func (t *MapTree[T]) Add(parent *MapNode[T], child T) (*MapNode[T], error) {
	name := child.GetName()
	if _, ok := t.Nodes.Get(name); ok {
		return nil, fmt.Errorf("duplicate key '%s'", name)
	}

	node := NewNode(child, parent.Depth+1)
	node.Parent = parent
	parent.Children.Set(name, node)
	t.Nodes.Set(name, node)

	return node, nil
}

// AddNode adds a sub-tree
func (t *MapTree[T]) AddNode(parent *MapNode[T], child *MapNode[T]) error {
	name := child.Value.GetName()
	if _, ok := t.Nodes.Get(name); ok {
		return fmt.Errorf("duplicate key '%s'", name)
	}

	child.Parent = parent
	parent.Children.Set(name, child)
	t.Nodes.Set(name, child)

	return nil
}

// Aggregate aggregates values over the tree
func Aggregate[T Named, V any](t *MapTree[T], values map[string]V, zero V, fn func(a, b V) V) {
	aggregate(t.Root, values, zero, fn)
}

func aggregate[T Named, V any](nd *MapNode[T], values map[string]V, zero V, fn func(a, b V) V) V {
	agg, ok := values[nd.Value.GetName()]
	if !ok {
		agg = zero
	}
	for _, name := range nd.Children.Keys() {
		child, _ := nd.Children.Get(name)
		v := aggregate(child, values, zero, fn)
		agg = fn(agg, v)
	}
	values[nd.Value.GetName()] = agg
	return agg
}

// TreeFormatter formats trees
type TreeFormatter[T Named] struct {
	NameFunc     func(t *MapNode[T], indent int) string
	Indent       int
	prefixNone   string
	prefixEmpty  string
	prefixNormal string
	prefixLast   string
}

// NewTreeFormatter creates a new TreeFormatter
func NewTreeFormatter[T Named](
	nameFunc func(t *MapNode[T], indent int) string,
	indent int,
) TreeFormatter[T] {
	return TreeFormatter[T]{
		NameFunc:     nameFunc,
		Indent:       indent,
		prefixNone:   strings.Repeat(" ", indent),
		prefixEmpty:  "│" + strings.Repeat(" ", indent-1),
		prefixNormal: "├" + strings.Repeat("─", indent-1),
		prefixLast:   "└" + strings.Repeat("─", indent-1),
	}
}

// FormatTree formats a tree
func (f *TreeFormatter[T]) FormatTree(t *MapTree[T]) string {
	sb := strings.Builder{}
	f.formatTree(&sb, t.Root, 0, false, "")
	return sb.String()
}

func (f *TreeFormatter[T]) formatTree(sb *strings.Builder, t *MapNode[T], depth int, last bool, prefix string) {
	pref := prefix
	if depth > 0 {
		pref = prefix + f.createPrefix(last)
	}
	fmt.Fprint(sb, pref)
	fmt.Fprintf(sb, "%s", f.NameFunc(t, utf8.RuneCountInString(pref)))
	fmt.Fprint(sb, "\n")

	if depth > 0 {
		pref = prefix + f.createPrefixEmpty(last)
	}

	names := append([]string{}, t.Children.Keys()...)
	//sort.Strings(names)
	for i, name := range names {
		last := i == len(names)-1
		child, _ := t.Children.Get(name)
		f.formatTree(sb, child, depth+1, last, pref)
	}
}

func (f *TreeFormatter[T]) createPrefix(last bool) string {
	if last {
		return f.prefixLast
	}
	return f.prefixNormal
}

func (f *TreeFormatter[T]) createPrefixEmpty(last bool) string {
	if last {
		return f.prefixNone
	}
	return f.prefixEmpty
}
