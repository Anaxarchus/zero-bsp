package zerobsp

import "github.com/Anaxarchus/zero-gdscript/pkg/vector2"

//func Defragment(tree *BspTree) []*BspTree {
//}

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
