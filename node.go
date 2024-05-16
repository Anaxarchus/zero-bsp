package zerobsp

import (
	"fmt"

	"github.com/Anaxarchus/zero-gdscript/pkg/vector2"
)

type BspNode struct {
	Position    vector2.Vector2
	Size        vector2.Vector2
	IsBranch    bool
	IsMajor     bool
	SplitWith   SplitMode
	Data        any
	DataRotated bool
	Major       *BspNode
	Minor       *BspNode
	Parent      *BspNode
}

type SplitMode int

const (
	SplitModeHorizontal SplitMode = iota
	SplitModeVertical
	SplitModeBalanced
	SplitModeAlternate
)

type RotationMode int

const (
	RotationModeNone RotationMode = iota
	RotationModeFit
	RotationModeAlign
	RotationModeUnalign
)

func NewBspNode(position, size vector2.Vector2) *BspNode {
	return &BspNode{
		Position: position,
		Size:     size,
		IsBranch: false,
	}
}

func (n *BspNode) GetArea() float64 {
	return n.Size.X * n.Size.Y
}

func (n *BspNode) GetCapacity() Capacity {
	cap := Capacity{Total: n.Size.X * n.Size.Y}
	if n.IsBranch {
		cap.Used = cap.Total
	}
	if n.Major != nil {
		cap = cap.Add(n.Major.GetCapacity())
	}
	if n.Minor != nil {
		cap = cap.Add(n.Minor.GetCapacity())
	}
	return cap
}

func (n *BspNode) Insert(data any, dataSize vector2.Vector2, rotationMode RotationMode, splitMode SplitMode) error {
	err := n.Split(dataSize, rotationMode, splitMode)
	if err != nil {
		return err
	}
	n.Data = data
	return nil
}

func (n *BspNode) TestFit(usage vector2.Vector2, rotationMode RotationMode) vector2.Vector2 {
	rFits := (usage.Y < n.Size.X || usage.X < n.Size.Y) && rotationMode != RotationModeNone
	if rFits && getShouldRotate(n.Size, usage, rotationMode) {
		usage = vector2.New(usage.Y, usage.X)
	}
	return n.Size.Sub(usage)
}

func (n *BspNode) Split(usage vector2.Vector2, rotationMode RotationMode, splitmode SplitMode) error {
	fits := usage.X < n.Size.X || usage.Y < n.Size.Y
	rFits := (usage.Y < n.Size.X || usage.X < n.Size.Y) && rotationMode != RotationModeNone
	if !fits && !rFits {
		return fmt.Errorf("error: usage does not fit in partition. Usage(x%f, y%f), Partition(x%f, y%f)", usage.X, usage.Y, n.Size.X, n.Size.Y)
	}
	if rFits && getShouldRotate(n.Size, usage, rotationMode) {
		usage = vector2.New(usage.Y, usage.X)
	}
	var splits [2]vector2.Vector2
	switch splitmode {
	case SplitModeHorizontal:
		splits = getHorizontalNodeSplit(n.Size, usage)
	case SplitModeVertical:
		splits = getVerticalNodeSplit(n.Size, usage)
	case SplitModeBalanced:
		a := getHorizontalNodeSplit(n.Size, usage)
		b := getVerticalNodeSplit(n.Size, usage)
		if (a[0].X*a[0].Y - a[1].X*a[1].Y) < (b[0].X*b[0].Y - b[1].X*b[1].Y) {
			splits = a
		} else {
			splits = b
		}
	case SplitModeAlternate:
		if n.SplitWith == SplitModeHorizontal {
			splits = getVerticalNodeSplit(n.Size, usage)
		} else if n.SplitWith == SplitModeVertical {
			splits = getHorizontalNodeSplit(n.Size, usage)
		}
	}

	majorSize := splits[0]
	majorPosition := n.Position.Add(vector2.New(0.0, usage.Y))
	minorSize := splits[1]
	minorPosition := n.Position.Add(vector2.New(usage.X, 0.0))
	if minorSize.X*minorSize.Y > majorSize.X*majorSize.Y {
		majorSize = splits[1]
		majorPosition = n.Position.Add(vector2.New(usage.X, 0.0))
		minorSize = splits[0]
		minorPosition = n.Position.Add(vector2.New(0.0, usage.Y))
	}

	n.Major = NewBspNode(majorPosition, majorSize)
	n.Major.Parent = n

	n.Minor = NewBspNode(minorPosition, minorSize)
	n.Minor.Parent = n

	n.IsBranch = true
	n.DataRotated = !fits
	n.Size = usage

	return nil
}

func getHorizontalNodeSplit(nodeSize, usageSize vector2.Vector2) [2]vector2.Vector2 {
	return [2]vector2.Vector2{
		vector2.New(nodeSize.X, nodeSize.Y-usageSize.Y),
		vector2.New(nodeSize.X-usageSize.X, usageSize.Y),
	}
}

func getVerticalNodeSplit(nodeSize, usageSize vector2.Vector2) [2]vector2.Vector2 {
	return [2]vector2.Vector2{
		vector2.New(usageSize.X, nodeSize.Y-usageSize.Y),
		vector2.New(nodeSize.X-usageSize.X, nodeSize.Y),
	}
}

func getShouldRotate(nodeSize, usageSize vector2.Vector2, rotationMode RotationMode) bool {
	switch rotationMode {
	case RotationModeNone:
		return false
	case RotationModeAlign:
		return (usageSize.X < usageSize.Y && nodeSize.X < nodeSize.Y) || (usageSize.Y < usageSize.X && nodeSize.Y < nodeSize.X)
	case RotationModeUnalign:
		return !(usageSize.X < usageSize.Y && nodeSize.X < nodeSize.Y) || (usageSize.Y < usageSize.X && nodeSize.Y < nodeSize.X)
	case RotationModeFit:
		return usageSize.X > nodeSize.X || usageSize.Y > nodeSize.Y
	default:
		return false
	}
}
