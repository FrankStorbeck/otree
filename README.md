# tree

Otree is a simple library for ordered tree data structures. It follows the
description from
[wikipedia.org/wiki/Tree_(data_structure)](https://en.wikipedia.org/wiki/Tree_(data_structure)).

## Terminology

A `node` is a structure which may contain data and connections to other nodes.
Each node in a tree has zero or more `child` nodes, which are below it in the
`tree` (by convention, trees are drawn with descendants going downwards). A node
that has a child is called the child's `parent` node. All nodes have exactly one
parent, except the topmost root node, which has none. A node might have many
`ancestor` nodes, such as the parent's parent. Child nodes with the same parent
are `sibling` nodes. Typically siblings have an order, with the first one
conventionally drawn on the left.

An `internal` node is any node of a tree that has child nodes. Similarly, an
`external` node, or `leaf` node, is any node that does not have child nodes.

The `height` of a node is the length of the longest downward `path` to a leaf
from that node. The height of the root is the height of the tree. The `depth` of
a node is the length of the path to its root. When using zero-based counting,
the root node has depth zero, leaf nodes have height zero, and a tree with only
a single node has depth and height zero.

An `ancestor` is a node reachable by repeated proceeding from child to parent.
A `descendant` is a node reachable by repeated proceeding from parent to child.

The `degree` of a node is its number of children. A leaf has necessarily degree
zero. The `degree` of tree is the maximum degree of all nodes in the tree.

The `distance` is the number of edges along the path between two nodes.
The `level` of a node is the zero-based counting of edges along the path to the
root node. The `width` is the number of nodes in a level. The `breadth` is
number of leaves.

The `size` of a tree is the number of nodes in it.

An `ordered` tree is one in which an ordering is specified for the children of a
node.

## Example

```
func main() {
  tr := New("root")

  level1 := []*Node{NewNode(10), NewNode(11), NewNode(12), NewNode(13)}
  tr.LinkChildren(tr.Root(), 0, level1...)

  level2 := []*Node{NewNode(20), NewNode(21)}
  tr.LinkChildren(level1[0], 0, level2...)

  fmt.Println(tr.String()) // output: root[10[20 21] 11 12 13]
}
```
