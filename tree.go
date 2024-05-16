package zerobsp

import (
	"fmt"

	"github.com/Anaxarchus/zero-gdscript/pkg/vector2"
)

type PlacementMode int

const (
	PlacementModeFirst PlacementMode = iota
	PlacementModeWorst
	PlacementModeBest
)

type BspTree struct {
	Root     *BspNode
	Leaves   LeafList
	Capacity Capacity
}

type LeafList []*BspNode

type Capacity struct {
	Total float64
	Used  float64
}

func (c Capacity) Add(capacity Capacity) Capacity {
	c.Total += capacity.Total
	c.Used += capacity.Used
	return c
}

func NewBspTree(width, height float64) *BspTree {
	root := NewBspNode(vector2.Zero(), vector2.New(width, height))
	return &BspTree{
		Root:   root,
		Leaves: LeafList{root},
	}
}

func (t *BspTree) Insert(data any, dataSize vector2.Vector2, placementMode PlacementMode, rotationMode RotationMode, splitMode SplitMode) error {
	var node *BspNode
	switch placementMode {
	case PlacementModeFirst:
		node = t.findFirst(dataSize, rotationMode)
	case PlacementModeBest:
		node = t.findBest(dataSize, rotationMode)
	case PlacementModeWorst:
		node = t.findWorst(dataSize, rotationMode)

	}
	if node == nil {
		return fmt.Errorf("no nodes we're found")
	} else {
		err := node.Insert(data, dataSize, rotationMode, splitMode)
		if err != nil {
			return err
		}
		t.Leaves.Remove(node)
		t.Leaves.Add(node.Major)
		t.Leaves.Add(node.Minor)
		return nil
	}
}

func (t *BspTree) findFirst(dataSize vector2.Vector2, rotationMode RotationMode) *BspNode {
	for _, node := range t.Leaves {
		fit := node.TestFit(dataSize, rotationMode)
		if fit.X > 0 && fit.Y > 0 {
			return node
		}
	}
	return nil
}

func (t *BspTree) findWorst(dataSize vector2.Vector2, rotationMode RotationMode) *BspNode {
	var worst *BspNode
	var score float64
	for _, node := range t.Leaves {
		fit := node.TestFit(dataSize, rotationMode)
		if fit.X*fit.Y > score {
			score = fit.X * fit.Y
			worst = node
		}
	}
	return worst
}

func (t *BspTree) findBest(dataSize vector2.Vector2, rotationMode RotationMode) *BspNode {
	var worst *BspNode
	var score float64
	for _, node := range t.Leaves {
		fit := node.TestFit(dataSize, rotationMode)
		if fit.X*fit.Y < score {
			score = fit.X * fit.Y
			worst = node
		}
	}
	return worst
}

func (t *BspTree) GetCapacity() Capacity {
	return t.Root.GetCapacity()
}

func (t *BspTree) GetEffeciency() float64 {
	cap := t.GetCapacity()
	return cap.Used / cap.Total
}

// Function to remove an element from a slice by pointer
func (ls *LeafList) Remove(leaf *BspNode) {
	index := -1
	for i, p := range *ls {
		if p == leaf {
			index = i
			break
		}
	}
	if index != -1 {
		// Remove the element at the found index
		*ls = append((*ls)[:index], (*ls)[index+1:]...)
	}
}

// Function to add an element to the slice, ensuring it doesn't already exist
func (ls *LeafList) Add(leaf *BspNode) {
	for _, p := range *ls {
		if p == leaf {
			// Element already exists, do nothing
			return
		}
	}
	// Add the new element
	*ls = append(*ls, leaf)
}
