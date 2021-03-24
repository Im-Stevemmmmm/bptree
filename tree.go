package bptree

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	err          error
	defaultOrder = 4
	order        = defaultOrder
	queue        *Node
)

type Tree struct {
	Root *Node
}

type Record struct {
	Value []byte
}

type Node struct {
	Pointers []interface{}
	Keys     []int
	Parent   *Node
	IsLeaf   bool
	NumKeys  int
	Next     *Node
}

func NewTree() *Tree {
	return &Tree{}
}

func (t *Tree) Insert(key int, value []byte) error {
	var pointer *Record
	var leaf *Node

	if _, err := t.Find(key); err == nil {
		return errors.New("key already exists")
	}

	pointer, err := makeRecord(value)
	if err != nil {
		return err
	}

	if t.Root == nil {
		return t.startNewTree(key, pointer)
	}

	leaf = t.findLeaf(key)

	if leaf.NumKeys < order-1 {
		insertIntoLeaf(leaf, key, pointer)
		return nil
	}

	return t.insertIntoLeafAfterSplitting(leaf, key, pointer)
}

func (t *Tree) Find(key int) (*Record, error) {
	i := 0
	c := t.findLeaf(key)
	if c == nil {
		return nil, errors.New("key not found")
	}
	for i = 0; i < c.NumKeys; i++ {
		if c.Keys[i] == key {
			break
		}
	}
	if i == c.NumKeys {
		return nil, errors.New("key not found")
	}

	r, _ := c.Pointers[i].(*Record)

	return r, nil
}

func (t *Tree) FindAndPrint(key int) {
	r, err := t.Find(key)

	if err != nil || r == nil {
		fmt.Printf("Record not found under key %d.\n", key)
	} else {
		fmt.Printf("Record at %d -- key %d, value %s.\n", r, key, r.Value)
	}
}

func (t *Tree) FindAndPrintRange(keyStart, keyEnd int) {
	var i int
	arraySize := keyEnd - keyStart + 1
	returnedKeys := make([]int, arraySize)
	returnedPointers := make([]interface{}, arraySize)
	numFound := t.findRange(keyStart, keyEnd, returnedKeys, returnedPointers)
	if numFound == 0 {
		fmt.Println("None found,")
	} else {
		for i = 0; i < numFound; i++ {
			c, _ := returnedPointers[i].(*Record)
			fmt.Printf("Key: %d  Location: %d  Value: %s\n",
				returnedKeys[i],
				returnedPointers[i],
				c.Value)
		}
	}
}

func (t *Tree) PrintTree() {
	var n *Node
	i := 0
	rank := 0
	newRank := 0

	if t.Root == nil {
		fmt.Printf("Empty tree.\n")
		return
	}
	queue = nil
	enqueue(t.Root)
	for queue != nil {
		n = dequeue()
		if n != nil {
			if n.Parent != nil && n == n.Parent.Pointers[0] {
				newRank = t.pathToRoot(n)
				if newRank != rank {
					fmt.Printf("\n")
				}
			}

			for i = 0; i < n.NumKeys; i++ {

				fmt.Printf("%d ", n.Keys[i])
			}

			if !n.IsLeaf {
				for i = 0; i <= n.NumKeys; i++ {
					c, _ := n.Pointers[i].(*Node)
					enqueue(c)
				}
			}

			fmt.Printf(" | ")
		}
	}
	fmt.Printf("\n")
}

func (t *Tree) PrintLeaves() {
	if t.Root == nil {
		fmt.Printf("Empty tree.\n")
		return
	}

	var i int
	c := t.Root
	for !c.IsLeaf {
		c, _ = c.Pointers[0].(*Node)
	}

	for {
		for i = 0; i < c.NumKeys; i++ {

			fmt.Printf("%d ", c.Keys[i])
		}

		if c.Pointers[order-1] != nil {
			fmt.Printf(" | ")
			c, _ = c.Pointers[order-1].(*Node)
		} else {
			break
		}
	}
	fmt.Printf("\n")
}

func (t *Tree) Delete(key int) error {
	keyRecord, err := t.Find(key)
	if err != nil {
		return err
	}
	keyLeaf := t.findLeaf(key)
	if keyRecord != nil && keyLeaf != nil {
		t.deleteEntry(keyLeaf, key, keyRecord)
	}
	return nil
}

//
//
//
//
//
//
//
// Private Functions
//
//
//
//
//
//
//
func enqueue(newNode *Node) {
	var c *Node
	if queue == nil {
		queue = newNode
		queue.Next = nil
	} else {
		c = queue
		for c.Next != nil {
			c = c.Next
		}
		c.Next = newNode
		newNode.Next = nil
	}
}

func dequeue() *Node {
	n := queue
	queue = queue.Next
	n.Next = nil
	return n
}

func (t *Tree) height() int {
	h := 0
	c := t.Root
	for !c.IsLeaf {
		c, _ = c.Pointers[0].(*Node)
		h++
	}
	return h
}

func (t *Tree) pathToRoot(child *Node) int {
	length := 0
	c := child
	for c != t.Root {
		c = c.Parent
		length += 1
	}
	return length
}

func (t *Tree) findRange(keyStart, keyEnd int, returnedKeys []int, returnedPointers []interface{}) int {
	var i int
	numFound := 0

	n := t.findLeaf(keyStart)
	if n == nil {
		return 0
	}
	for i = 0; i < n.NumKeys && n.Keys[i] < keyStart; i++ {
		if i == n.NumKeys { // could be wrong
			return 0
		}
	}
	for ; i < n.NumKeys && n.Keys[i] <= keyEnd; i++ {
		returnedKeys[numFound] = n.Keys[i]
		returnedPointers[numFound] = n.Pointers[i]
		numFound += 1
	}
	n, _ = n.Pointers[order-1].(*Node)
	i = 0
	return numFound
}

func (t *Tree) findLeaf(key int) *Node {
	i := 0
	c := t.Root
	if c == nil {
		return c
	}
	for !c.IsLeaf {
		i = 0
		for i < c.NumKeys {
			if key >= c.Keys[i] {
				i += 1
			} else {
				break
			}
		}
		c, _ = c.Pointers[i].(*Node)
	}
	return c
}

func cut(length int) int {
	if length%2 == 0 {
		return length / 2
	}

	return length/2 + 1
}

//
//	INSERTION
//
func makeRecord(value []byte) (*Record, error) {
	newRecord := new(Record)
	newRecord.Value = value
	return newRecord, nil
}

func makeNode() (*Node, error) {
	newNode := new(Node)
	newNode.Keys = make([]int, order-1)
	if newNode.Keys == nil {
		return nil, errors.New("error: new node keys array")
	}
	newNode.Pointers = make([]interface{}, order)
	if newNode.Keys == nil {
		return nil, errors.New("error: new node pointers array")
	}
	newNode.IsLeaf = false
	newNode.NumKeys = 0
	newNode.Parent = nil
	newNode.Next = nil
	return newNode, nil
}

func makeLeaf() (*Node, error) {
	leaf, err := makeNode()
	if err != nil {
		return nil, err
	}
	leaf.IsLeaf = true
	return leaf, nil
}

func getLeftIndex(parent, left *Node) int {
	leftIndex := 0
	for leftIndex <= parent.NumKeys && parent.Pointers[leftIndex] != left {
		leftIndex += 1
	}
	return leftIndex
}

func insertIntoLeaf(leaf *Node, key int, pointer *Record) {
	var i, insertionPoint int

	for insertionPoint < leaf.NumKeys && leaf.Keys[insertionPoint] < key {
		insertionPoint += 1
	}

	for i = leaf.NumKeys; i > insertionPoint; i-- {
		leaf.Keys[i] = leaf.Keys[i-1]
		leaf.Pointers[i] = leaf.Pointers[i-1]
	}
	leaf.Keys[insertionPoint] = key
	leaf.Pointers[insertionPoint] = pointer
	leaf.NumKeys += 1
	return
}

func (t *Tree) insertIntoLeafAfterSplitting(leaf *Node, key int, pointer *Record) error {
	var newLeaf *Node
	var insertionIndex, split, newKey, i, j int
	var err error

	newLeaf, err = makeLeaf()
	if err != nil {
		return nil
	}

	tempKeys := make([]int, order)
	tempPointers := make([]interface{}, order)

	for insertionIndex < order-1 && leaf.Keys[insertionIndex] < key {
		insertionIndex += 1
	}

	for i = 0; i < leaf.NumKeys; i++ {
		if j == insertionIndex {
			j += 1
		}
		tempKeys[j] = leaf.Keys[i]
		tempPointers[j] = leaf.Pointers[i]
		j += 1
	}

	tempKeys[insertionIndex] = key
	tempPointers[insertionIndex] = pointer

	leaf.NumKeys = 0

	split = cut(order - 1)

	for i = 0; i < split; i++ {
		leaf.Pointers[i] = tempPointers[i]
		leaf.Keys[i] = tempKeys[i]
		leaf.NumKeys += 1
	}

	j = 0
	for i = split; i < order; i++ {
		newLeaf.Pointers[j] = tempPointers[i]
		newLeaf.Keys[j] = tempKeys[i]
		newLeaf.NumKeys += 1
		j += 1
	}

	newLeaf.Pointers[order-1] = leaf.Pointers[order-1]
	leaf.Pointers[order-1] = newLeaf

	for i = leaf.NumKeys; i < order-1; i++ {
		leaf.Pointers[i] = nil
	}
	for i = newLeaf.NumKeys; i < order-1; i++ {
		newLeaf.Pointers[i] = nil
	}

	newLeaf.Parent = leaf.Parent
	newKey = newLeaf.Keys[0]

	return t.insertIntoParent(leaf, newKey, newLeaf)
}

func insertIntoNode(n *Node, leftIndex, key int, right *Node) {
	var i int
	for i = n.NumKeys; i > leftIndex; i-- {
		n.Pointers[i+1] = n.Pointers[i]
		n.Keys[i] = n.Keys[i-1]
	}
	n.Pointers[leftIndex+1] = right
	n.Keys[leftIndex] = key
	n.NumKeys += 1
}

func (t *Tree) insertIntoNodeAfterSplitting(oldNode *Node, leftIndex, key int, right *Node) error {
	var i, j, split, kPrime int
	var newNode, child *Node
	var tempKeys []int
	var tempPointers []interface{}
	var err error

	tempPointers = make([]interface{}, order+1)
	tempKeys = make([]int, order)

	for i = 0; i < oldNode.NumKeys+1; i++ {
		if j == leftIndex+1 {
			j += 1
		}
		tempPointers[j] = oldNode.Pointers[i]
		j += 1
	}

	j = 0
	for i = 0; i < oldNode.NumKeys; i++ {
		if j == leftIndex {
			j += 1
		}
		tempKeys[j] = oldNode.Keys[i]
		j += 1
	}

	tempPointers[leftIndex+1] = right
	tempKeys[leftIndex] = key

	split = cut(order)
	newNode, err = makeNode()
	if err != nil {
		return err
	}
	oldNode.NumKeys = 0
	for i = 0; i < split-1; i++ {
		oldNode.Pointers[i] = tempPointers[i]
		oldNode.Keys[i] = tempKeys[i]
		oldNode.NumKeys += 1
	}
	oldNode.Pointers[i] = tempPointers[i]
	kPrime = tempKeys[split-1]
	j = 0
	for i += 1; i < order; i++ {
		newNode.Pointers[j] = tempPointers[i]
		newNode.Keys[j] = tempKeys[i]
		newNode.NumKeys += 1
		j += 1
	}
	newNode.Pointers[j] = tempPointers[i]
	newNode.Parent = oldNode.Parent
	for i = 0; i <= newNode.NumKeys; i++ {
		child, _ = newNode.Pointers[i].(*Node)
		child.Parent = newNode
	}

	return t.insertIntoParent(oldNode, kPrime, newNode)
}

func (t *Tree) insertIntoParent(left *Node, key int, right *Node) error {
	var leftIndex int
	parent := left.Parent

	if parent == nil {
		return t.insertIntoNewRoot(left, key, right)
	}
	leftIndex = getLeftIndex(parent, left)

	if parent.NumKeys < order-1 {
		insertIntoNode(parent, leftIndex, key, right)
		return nil
	}

	return t.insertIntoNodeAfterSplitting(parent, leftIndex, key, right)
}

func (t *Tree) insertIntoNewRoot(left *Node, key int, right *Node) error {
	t.Root, err = makeNode()
	if err != nil {
		return err
	}
	t.Root.Keys[0] = key
	t.Root.Pointers[0] = left
	t.Root.Pointers[1] = right
	t.Root.NumKeys += 1
	t.Root.Parent = nil
	left.Parent = t.Root
	right.Parent = t.Root
	return nil
}

func (t *Tree) startNewTree(key int, pointer *Record) error {
	t.Root, err = makeLeaf()
	if err != nil {
		return err
	}
	t.Root.Keys[0] = key
	t.Root.Pointers[0] = pointer
	t.Root.Pointers[order-1] = nil
	t.Root.Parent = nil
	t.Root.NumKeys += 1
	return nil
}

func getNeighbourIndex(n *Node) int {
	var i int

	for i = 0; i <= n.Parent.NumKeys; i++ {
		if reflect.DeepEqual(n.Parent.Pointers[i], n) {
			return i - 1
		}
	}

	return i
}

func removeEntryFromNode(n *Node, key int, pointer interface{}) *Node {
	var i, numPointers int

	for n.Keys[i] != key {
		i += 1
	}

	for i += 1; i < n.NumKeys; i++ {
		n.Keys[i-1] = n.Keys[i]
	}

	if n.IsLeaf {
		numPointers = n.NumKeys
	} else {
		numPointers = n.NumKeys + 1
	}

	i = 0
	for n.Pointers[i] != pointer {
		i += 1
	}
	for i += 1; i < numPointers; i++ {
		n.Pointers[i-1] = n.Pointers[i]
	}
	n.NumKeys -= 1

	if n.IsLeaf {
		for i = n.NumKeys; i < order-1; i++ {
			n.Pointers[i] = nil
		}
	} else {
		for i = n.NumKeys + 1; i < order; i++ {
			n.Pointers[i] = nil
		}
	}

	return n
}

func (t *Tree) adjustRoot() {
	var newRoot *Node

	if t.Root.NumKeys > 0 {
		return
	}

	if !t.Root.IsLeaf {
		newRoot, _ = t.Root.Pointers[0].(*Node)
		newRoot.Parent = nil
	} else {
		newRoot = nil
	}
	t.Root = newRoot

	return
}

func (t *Tree) coalesceNodes(n, neighbour *Node, neighbourIndex, kPrime int) {
	var i, j, neighbourInsertionIndex, nEnd int
	var tmp *Node

	if neighbourIndex == -1 {
		tmp = n
		n = neighbour
		neighbour = tmp
	}

	neighbourInsertionIndex = neighbour.NumKeys

	if !n.IsLeaf {
		neighbour.Keys[neighbourInsertionIndex] = kPrime
		neighbour.NumKeys += 1

		nEnd = n.NumKeys
		i = neighbourInsertionIndex + 1
		for j = 0; j < nEnd; j++ {
			neighbour.Keys[i] = n.Keys[j]
			neighbour.Pointers[i] = n.Pointers[j]
			neighbour.NumKeys += 1
			n.NumKeys -= 1
			i += 1
		}
		neighbour.Pointers[i] = n.Pointers[j]

		for i = 0; i < neighbour.NumKeys+1; i++ {
			tmp, _ = neighbour.Pointers[i].(*Node)
			tmp.Parent = neighbour
		}
	} else {
		i = neighbourInsertionIndex
		for j = 0; j < n.NumKeys; j++ {
			neighbour.Keys[i] = n.Keys[j]
			n.Pointers[i] = n.Pointers[j]
			neighbour.NumKeys += 1
		}
		neighbour.Pointers[order-1] = n.Pointers[order-1]
	}

	t.deleteEntry(n.Parent, kPrime, n)
}

func (t *Tree) redistributeNodes(n, neighbour *Node, neighbourIndex, kPrimeIndex, kPrime int) {
	var i int
	var tmp *Node

	if neighbourIndex != -1 {
		if !n.IsLeaf {
			n.Pointers[n.NumKeys+1] = n.Pointers[n.NumKeys]
		}
		for i = n.NumKeys; i > 0; i-- {
			n.Keys[i] = n.Keys[i-1]
			n.Pointers[i] = n.Pointers[i-1]
		}
		if !n.IsLeaf { // why the second if !n.IsLeaf
			n.Pointers[0] = neighbour.Pointers[neighbour.NumKeys]
			tmp, _ = n.Pointers[0].(*Node)
			tmp.Parent = n
			neighbour.Pointers[neighbour.NumKeys] = nil
			n.Keys[0] = kPrime
			n.Parent.Keys[kPrimeIndex] = neighbour.Keys[neighbour.NumKeys-1]
		} else {
			n.Pointers[0] = neighbour.Pointers[neighbour.NumKeys-1]
			neighbour.Pointers[neighbour.NumKeys-1] = nil
			n.Keys[0] = neighbour.Keys[neighbour.NumKeys-1]
			n.Parent.Keys[kPrimeIndex] = n.Keys[0]
		}
	} else {
		if n.IsLeaf {
			n.Keys[n.NumKeys] = neighbour.Keys[0]
			n.Pointers[n.NumKeys] = neighbour.Pointers[0]
			n.Parent.Keys[kPrimeIndex] = neighbour.Keys[1]
		} else {
			n.Keys[n.NumKeys] = kPrime
			n.Pointers[n.NumKeys+1] = neighbour.Pointers[0]
			tmp, _ = n.Pointers[n.NumKeys+1].(*Node)
			tmp.Parent = n
			n.Parent.Keys[kPrimeIndex] = neighbour.Keys[0]
		}
		for i = 0; i < neighbour.NumKeys-1; i++ {
			neighbour.Keys[i] = neighbour.Keys[i+1]
			neighbour.Pointers[i] = neighbour.Pointers[i+1]
		}
		if !n.IsLeaf {
			neighbour.Pointers[i] = neighbour.Pointers[i+1]
		}
	}
	n.NumKeys += 1
	neighbour.NumKeys -= 1

	return
}

func (t *Tree) deleteEntry(n *Node, key int, pointer interface{}) {
	var minKeys, neighbourIndex, kPrimeIndex, kPrime, capacity int
	var neighbour *Node

	n = removeEntryFromNode(n, key, pointer)

	if n == t.Root {
		t.adjustRoot()
		return
	}

	if n.IsLeaf {
		minKeys = cut(order - 1)
	} else {
		minKeys = cut(order) - 1
	}

	if n.NumKeys >= minKeys {
		return
	}

	neighbourIndex = getNeighbourIndex(n)

	if neighbourIndex == -1 {
		kPrimeIndex = 0
	} else {
		kPrimeIndex = neighbourIndex
	}

	kPrime = n.Parent.Keys[kPrimeIndex]

	if neighbourIndex == -1 {
		neighbour, _ = n.Parent.Pointers[1].(*Node)
	} else {
		neighbour, _ = n.Parent.Pointers[neighbourIndex].(*Node)
	}

	if n.IsLeaf {
		capacity = order
	} else {
		capacity = order - 1
	}

	if neighbour.NumKeys+n.NumKeys < capacity {
		t.coalesceNodes(n, neighbour, neighbourIndex, kPrime)
		return
	} else {
		t.redistributeNodes(n, neighbour, neighbourIndex, kPrimeIndex, kPrime)
		return
	}

}
