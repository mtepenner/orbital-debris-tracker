package spatial

import "math"

type Point struct {
	ID string
	X  float64
	Y  float64
	Z  float64
}

type Box struct {
	CenterX float64
	CenterY float64
	CenterZ float64
	Half    float64
}

func (box Box) Contains(point Point) bool {
	return math.Abs(point.X-box.CenterX) <= box.Half &&
		math.Abs(point.Y-box.CenterY) <= box.Half &&
		math.Abs(point.Z-box.CenterZ) <= box.Half
}

func (box Box) IntersectsSphere(center Point, radius float64) bool {
	dx := math.Max(math.Abs(center.X-box.CenterX)-box.Half, 0)
	dy := math.Max(math.Abs(center.Y-box.CenterY)-box.Half, 0)
	dz := math.Max(math.Abs(center.Z-box.CenterZ)-box.Half, 0)
	return dx*dx+dy*dy+dz*dz <= radius*radius
}

type Octree struct {
	boundary Box
	capacity int
	points   []Point
	divided  bool
	children [8]*Octree
}

func NewOctree(boundary Box, capacity int) *Octree {
	return &Octree{boundary: boundary, capacity: capacity}
}

func (tree *Octree) Insert(point Point) bool {
	if !tree.boundary.Contains(point) {
		return false
	}
	if len(tree.points) < tree.capacity && !tree.divided {
		tree.points = append(tree.points, point)
		return true
	}
	if !tree.divided {
		tree.subdivide()
	}
	for _, child := range tree.children {
		if child.Insert(point) {
			return true
		}
	}
	return false
}

func (tree *Octree) QueryRadius(center Point, radius float64, out *[]Point) {
	if !tree.boundary.IntersectsSphere(center, radius) {
		return
	}
	for _, point := range tree.points {
		if distance(point, center) <= radius {
			*out = append(*out, point)
		}
	}
	if !tree.divided {
		return
	}
	for _, child := range tree.children {
		child.QueryRadius(center, radius, out)
	}
}

func (tree *Octree) subdivide() {
	half := tree.boundary.Half / 2
	centers := [8][3]float64{
		{tree.boundary.CenterX - half, tree.boundary.CenterY - half, tree.boundary.CenterZ - half},
		{tree.boundary.CenterX + half, tree.boundary.CenterY - half, tree.boundary.CenterZ - half},
		{tree.boundary.CenterX - half, tree.boundary.CenterY + half, tree.boundary.CenterZ - half},
		{tree.boundary.CenterX + half, tree.boundary.CenterY + half, tree.boundary.CenterZ - half},
		{tree.boundary.CenterX - half, tree.boundary.CenterY - half, tree.boundary.CenterZ + half},
		{tree.boundary.CenterX + half, tree.boundary.CenterY - half, tree.boundary.CenterZ + half},
		{tree.boundary.CenterX - half, tree.boundary.CenterY + half, tree.boundary.CenterZ + half},
		{tree.boundary.CenterX + half, tree.boundary.CenterY + half, tree.boundary.CenterZ + half},
	}
	for index, center := range centers {
		tree.children[index] = NewOctree(Box{CenterX: center[0], CenterY: center[1], CenterZ: center[2], Half: half}, tree.capacity)
	}
	tree.divided = true
	existing := tree.points
	tree.points = nil
	for _, point := range existing {
		tree.Insert(point)
	}
}

func distance(left Point, right Point) float64 {
	dx := left.X - right.X
	dy := left.Y - right.Y
	dz := left.Z - right.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
