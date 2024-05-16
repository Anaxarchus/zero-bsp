package zerobsp

import "github.com/Anaxarchus/zero-gdscript/pkg/vector2"

func Defragment(tree *BspTree) []*BspTree {
	res := []*BspTree{tree}
	// iterate over all tree leaves
	for _, a := range tree.Leaves {
		for _, b := range tree.Leaves {
			if a == b {
				continue
			}
			// check if they are valid for merger
			merge := canMerge(a, b)

			// calculate a new rect from the joined areas
			xPos := min(a.Position.X, b.Position.X)
			yPos := max(a.Position.Y, b.Position.Y)
			width := 0.0
			height := 0.0
			if a.Position.X > b.Position.X {
				yPos = 
			}
			size := vector2.New(max(a.Position.X, b.Position.X))
			// create new trees from the new area and from the remainder area if there is a remainder
			// create a new tree and set root as new BspNode
		}
	}
	return res
}

func canMerge(a, b *BspNode) bool {
	// two nodes can merge if:
	return notIsBranch(a, b) && shareBorder(a, b) && shareCorner(a, b)
}

func shareBorder(a, b *BspNode) bool {
	relPos := a.Position.Sub(b.Position)
	return relPos.IsEqualApprox(a.Position)
}

func shareCorner(a, b *BspNode) bool {
	for _, c1 := range getCorners(a) {
		for _, c2 := range getCorners(b) {
			if c1.IsEqualApprox(c2) {
				return true
			}
		}
	}
	return false
}

func notIsBranch(a, b *BspNode) bool {
	return !(a.IsBranch && b.IsBranch)
}

func getCorners(node *BspNode) [4]vector2.Vector2 {
	return [4]vector2.Vector2{
		node.Position,
		node.Position.Add(vector2.New(node.Size.X, 0.0)),
		node.Position.Add(vector2.New(node.Size.X, node.Size.Y)),
		node.Position.Add(vector2.New(0.0, node.Size.Y)),
	}
}
