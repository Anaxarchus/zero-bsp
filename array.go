package zerobsp

import "github.com/Anaxarchus/zero-gdscript/pkg/vector2"

type TreeArray struct {
	Trees []*BspTree
	Size  vector2.Vector2
}

type SelectMode int

const (
	SelectModeFirst      SelectMode = iota // fills a bin completely before moving on to the next
	SelectModeSequential                   // rotates to the next bin with every part
	SelectModeBalance                      // Tries to keep all the bins evenly filled
)

func NewTreeArray(size vector2.Vector2) *TreeArray {
	return &TreeArray{
		Trees: []*BspTree{},
		Size:  size,
	}
}
