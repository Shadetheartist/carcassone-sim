package road

import (
	"beeb/carcassonne/tile"
	"errors"
)

type Road struct {
	head *Node
	tail *Node
}

type Node struct {
	prev  *Node
	next  *Node
	Value *tile.Tile
}

func (l *Road) First() *Node {
	return l.head
}
func (n *Node) Next() *Node {
	return n.next
}
func (n *Node) Prev() *Node {
	return n.prev
}

// Create new node with value
func (l *Road) Push(v *tile.Tile) *Road {
	n := &Node{Value: v}
	if l.head == nil {
		l.head = n // First node
	} else {
		l.tail.next = n // Add after prev last node
		n.prev = l.tail // Link back to prev last node
	}
	l.tail = n // reset tail to newly added node
	return l
}
func (l *Road) Find(name string) *Node {
	found := false
	var ret *Node = nil
	for n := l.First(); n != nil && !found; n = n.Next() {
		if n.Value.Name == name {
			found = true
			ret = n
		}
	}
	return ret
}
func (l *Road) Delete(name string) bool {
	success := false
	node2del := l.Find(name)
	if node2del != nil {
		prev_node := node2del.prev
		next_node := node2del.next
		// Remove this node
		prev_node.next = node2del.next
		next_node.prev = node2del.prev
		success = true
	}
	return success
}

var errEmpty = errors.New("ERROR - List is empty")

// Pop last item from list
func (l *Road) Pop() (v *tile.Tile, err error) {
	if l.tail == nil {
		err = errEmpty
	} else {
		v = l.tail.Value
		l.tail = l.tail.prev
		if l.tail == nil {
			l.head = nil
		}
	}
	return v, err
}
