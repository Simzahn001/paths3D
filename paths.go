/*
Package paths3D is a simple library written in Go made to handle 3D pathfinding. All you need to do is generate a Grid,
specify which cells aren't walkable and what height the cells are, optionally change the cost on specific cells, and
finally get a path from one cell to another. For a simple guide, take a look at how-to.md in the github repo: https://github.com/Simzahn001/paths3D
*/
package paths

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

// A Cell represents a point on a Grid map. It has an X and Y value for the position, a Cost, which influences which Cells are
// ideal for paths, Walkable, which indicates if the tile can be walked on or should be avoided, a Rune, which indicates
// which rune character the Cell is represented by, and a HeightLevel (default: 0), which represents the height of this cell.
type Cell struct {
	X, Y, HeightLevel int
	Cost              float64
	Walkable          bool
	Rune              rune
}

func (cell Cell) String() string {
	return fmt.Sprintf("X:%d Y:%d Height:%d Cost:%f Walkable:%t Rune:%s(%d)", cell.X, cell.Y, cell.HeightLevel, cell.Cost, cell.Walkable, string(cell.Rune), int(cell.Rune))
}

// Grid represents a "map" composed of individual Cells at each point in the map.
// Data is a 2D array of Cells.
// CellWidth and CellHeight indicate the size of Cells for Cell Position <-> World Position translation.
type Grid struct {
	Data [][]*Cell
}

// NewGrid returns a new Grid of (gridWidth x gridHeight) size.
func NewGrid(gridWidth, gridHeight int) *Grid {

	m := &Grid{}

	for y := 0; y < gridHeight; y++ {
		m.Data = append(m.Data, []*Cell{})
		for x := 0; x < gridWidth; x++ {
			m.Data[y] = append(m.Data[y], &Cell{
				X:           x,
				Y:           y,
				HeightLevel: 0,
				Cost:        1,
				Walkable:    true,
				Rune:        ' ',
			})
		}
	}
	return m
}

// NewGridFromStringArrays creates a Grid map from a 1D array of strings. Each string becomes a row of Cells, each
// with one rune as its character.
func NewGridFromStringArrays(arrays []string) *Grid {

	m := &Grid{}

	for y := 0; y < len(arrays); y++ {
		m.Data = append(m.Data, []*Cell{})
		stringLine := []rune(arrays[y])
		for x := 0; x < len(arrays[y]); x++ {
			m.Data[y] = append(m.Data[y], &Cell{
				X:           x,
				Y:           y,
				HeightLevel: 0,
				Cost:        1,
				Walkable:    true,
				Rune:        stringLine[x],
			})
		}
	}

	return m

}

// NewGridFromRuneArrays creates a Grid map from a 2D array of runes. Each individual Rune becomes a Cell in the resulting Grid.
func NewGridFromRuneArrays(arrays [][]rune) *Grid {

	m := &Grid{}

	for y := 0; y < len(arrays); y++ {
		m.Data = append(m.Data, []*Cell{})
		for x := 0; x < len(arrays[y]); x++ {
			m.Data[y] = append(m.Data[y], &Cell{
				X:           x,
				Y:           y,
				HeightLevel: 0,
				Cost:        1,
				Walkable:    true,
				Rune:        arrays[y][x],
			})
		}
	}

	return m

}

// AddHeightMap adds a height to the grid via a key-value map. All runes, the map contains, do have an assigned height.
// This height is applied to ALL cells with this rune. After the execution of this method, letters aren't bound to the height;
// they are no pointers. If you change a letter, the height will stay the same.
// Keep in mind, that cell runes are case-sensitive.
func (m *Grid) AddHeightMap(profile map[rune]int) {

	//loop trough all cells
	for _, cell := range m.AllCells() {
		//check if the map contains the rune of the cell
		heightLevel, exists := profile[cell.Rune]
		if exists {
			cell.HeightLevel = heightLevel
		}
	}

}

// DataToString returns a string, used to easily identify the Grid map.
func (m *Grid) DataToString() string {
	s := ""
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			s += string(m.Data[y][x].Rune) + " "
		}
		s += "\n"
	}
	return s
}

// Visualise return a string visualisation of the grid's cell heights.
// Not walkable blocks are represented by a blank space.
// An error is returned if there are more than 26 height levels. The heightmap is filled though,
// but every layer after the 26th is still visualised with the letter 'z'
func (m *Grid) Visualise() (visualisation []string, error error) {

	//check if more than 26 different height levels are contained.
	//If so, the visualisation is kinda buggy, because heights bigger than 26 will be
	//displayed with the same character.
	heights := m.GetHeightLevels()
	if len(m.GetHeightLevels()) > 26 {
		error = errors.New("there are more than 26 height levels. All levels after the 26th will not be displayed correctly")
	}

	//sort array
	sort.Ints(heights)

	//creating a map with all height levels and letters
	letters := make(map[int]rune)
	ascii := 97
	for _, height := range heights {
		letters[height] = rune(ascii)
		ascii++
		if ascii > 122 {
			ascii = 121
		}
	}

	//create the strings
	visualisation = []string{}
	for y := 0; y < m.Height(); y++ {
		var currentString = strings.Builder{}
		for x := 0; x < m.Width(); x++ {
			cell := m.Get(x, y)

			//non-walkable cells should be represented my a blank space
			if !cell.Walkable {
				currentString.WriteString(" ")
			} else {
				currentString.WriteString(string(letters[cell.HeightLevel]))
			}

		}
		visualisation = append(visualisation, currentString.String())
	}

	return visualisation, error
}

func (m *Grid) VisualisePath(path *Path) []string {
	//check if path is nil
	if path == nil {
		return nil
	}

	//get the grid as string
	visualisation, _ := m.Visualise()

	//go through each cell of the path and set the cell to the letter '#'
	for _, cell := range path.Cells {
		rowAsRuneArray := []rune(visualisation[cell.Y])
		rowAsRuneArray[cell.X] = '#'
		visualisation[cell.Y] = string(rowAsRuneArray)
	}

	return visualisation
}

// Get returns a pointer to the Cell in the x and y position provided.
func (m *Grid) Get(x, y int) *Cell {
	if x < 0 || y < 0 || x >= m.Width() || y >= m.Height() {
		return nil
	}
	return m.Data[y][x]
}

// Height returns the height of the Grid map.
func (m *Grid) Height() int {
	return len(m.Data)
}

// Width returns the width of the Grid map.
func (m *Grid) Width() int {
	return len(m.Data[0])
}

// GetAverageHeight returns the average height over all the Grid. Use math.Round to get an int
func (m *Grid) GetAverageHeight() float64 {

	sum := 0
	i := 1

	for _, cell := range m.AllCells() {
		sum += cell.HeightLevel
		i++
	}

	return float64(sum) / float64(i)

}

// GetMaxHeight returns the maximum height of the whole Grid
func (m *Grid) GetMaxHeight() int {

	var maxHeight = math.MinInt

	for _, cell := range m.AllCells() {
		if cell.HeightLevel > maxHeight {
			maxHeight = cell.HeightLevel
		}
	}

	return maxHeight
}

// GetMinHeight returns the minimum height of the whole Grid
func (m *Grid) GetMinHeight() int {

	var maxHeight = math.MaxInt

	for _, cell := range m.AllCells() {
		if cell.HeightLevel < maxHeight {
			maxHeight = cell.HeightLevel
		}
	}

	return maxHeight
}

// GetHeightLevels returns a list of all different height levels.
// use len() on the returned slice to get the amount of different height levels
func (m *Grid) GetHeightLevels() []int {

	var heightLevels = []int{}

	for _, cell := range m.AllCells() {
		//if height of the current cell is not yet contained
		if !containesInt(heightLevels, cell.HeightLevel) {
			heightLevels = append(heightLevels, cell.HeightLevel)
		}
	}

	return heightLevels
}

// CellsByRune returns a slice of pointers to Cells that all have the character provided.
func (m *Grid) CellsByRune(char rune) []*Cell {

	cells := make([]*Cell, 0)

	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			c := m.Get(x, y)
			if c.Rune == char {
				cells = append(cells, c)
			}
		}
	}

	return cells

}

// AllCells returns a single slice of pointers to all Cells contained in the Grid's 2D Data array.
func (m *Grid) AllCells() []*Cell {

	cells := make([]*Cell, 0)

	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			cells = append(cells, m.Get(x, y))
		}
	}

	return cells

}

// CellsByCost returns a slice of pointers to Cells that all have the Cost value provided.
func (m *Grid) CellsByCost(cost float64) []*Cell {

	cells := make([]*Cell, 0)

	for y := 0; y < m.Height(); y++ {

		for x := 0; x < m.Width(); x++ {

			c := m.Get(x, y)
			if c.Cost == cost {
				cells = append(cells, c)
			}

		}

	}

	return cells

}

// CellsByWalkable returns a slice of pointers to Cells that all have the Cost value provided.
func (m *Grid) CellsByWalkable(walkable bool) []*Cell {

	cells := make([]*Cell, 0)

	for y := 0; y < m.Height(); y++ {

		for x := 0; x < m.Width(); x++ {

			c := m.Get(x, y)
			if c.Walkable == walkable {
				cells = append(cells, c)
			}

		}

	}

	return cells

}

// CellsByHeightLevel returns a slice of pointers to Cells that all have the height level provided.
func (m *Grid) CellsByHeightLevel(heightLevel int) []*Cell {
	cells := make([]*Cell, 1)

	for _, cell := range m.AllCells() {
		cells = append(cells, cell)
	}

	return cells
}

// SetWalkable sets walkability across all cells in the Grid with the specified rune.
func (m *Grid) SetWalkable(char rune, walkable bool) {

	for y := 0; y < m.Height(); y++ {

		for x := 0; x < m.Width(); x++ {
			cell := m.Get(x, y)
			if cell.Rune == char {
				cell.Walkable = walkable
			}
		}

	}

}

// SetHeightLevel sets the height level for all cells in the Grid with the specified rune.
func (m *Grid) SetHeightLevel(char rune, heightLevel int) {

	for _, cell := range m.AllCells() {
		if cell.Rune == char {
			cell.HeightLevel = heightLevel
		}
	}

}

// SetCost sets the movement cost across all cells in the Grid with the specified rune.
func (m *Grid) SetCost(char rune, cost float64) {

	for y := 0; y < m.Height(); y++ {

		for x := 0; x < m.Width(); x++ {
			cell := m.Get(x, y)
			if cell.Rune == char {
				cell.Cost = cost
			}
		}

	}

}

// GetPathFromCells returns a Path, from the starting Cell to the destination Cell. diagonals controls whether moving diagonally
// is acceptable when creating the Path. wallsBlockDiagonals indicates whether to allow diagonal movement "through" walls that are
// positioned diagonally.
func (m *Grid) GetPathFromCells(start, dest *Cell, stepHeight int, diagonals, wallsBlockDiagonals bool) *Path {

	openNodes := minHeap{}
	heap.Push(&openNodes, &Node{Cell: dest, Cost: dest.Cost})

	checkedNodes := make([]*Cell, 0)

	hasBeenAdded := func(cell *Cell) bool {

		for _, c := range checkedNodes {
			if cell == c {
				return true
			}
		}
		return false

	}

	path := &Path{StepHeight: stepHeight}

	if !start.Walkable || !dest.Walkable {
		return nil
	}

	for {

		// If the list of openNodes (nodes to check) is at 0, then we've checked all Nodes, and so the function can quit.
		if len(openNodes) == 0 {
			break
		}

		node := heap.Pop(&openNodes).(*Node)

		// If we've reached the start, then we've constructed our Path going from the destination to the start; we just have
		// to loop through each Node and go up, adding it and its parents recursively to the path.
		if node.Cell == start {

			var t = node
			for true {
				path.Cells = append(path.Cells, t.Cell)
				t = t.Parent
				if t == nil {
					break
				}
			}

			break
		}

		// Otherwise, we add the current node's neighbors to the list of cells to check, and list of cells that have already been
		// checked (so we don't get nodes being checked multiple times).
		if node.Cell.X > 0 {
			c := m.Get(node.Cell.X-1, node.Cell.Y)
			n := &Node{c, node, c.Cost + node.Cost}
			if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (node.Cell.HeightLevel-n.Cell.HeightLevel) <= stepHeight {
				heap.Push(&openNodes, n)
				checkedNodes = append(checkedNodes, n.Cell)
			}
		}
		if node.Cell.X < m.Width()-1 {
			c := m.Get(node.Cell.X+1, node.Cell.Y)
			n := &Node{c, node, c.Cost + node.Cost}
			if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (node.Cell.HeightLevel-n.Cell.HeightLevel) <= stepHeight {
				heap.Push(&openNodes, n)
				checkedNodes = append(checkedNodes, n.Cell)
			}
		}

		if node.Cell.Y > 0 {
			c := m.Get(node.Cell.X, node.Cell.Y-1)
			n := &Node{c, node, c.Cost + node.Cost}
			if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (node.Cell.HeightLevel-n.Cell.HeightLevel) <= stepHeight {
				heap.Push(&openNodes, n)
				checkedNodes = append(checkedNodes, n.Cell)
			}
		}
		if node.Cell.Y < m.Height()-1 {
			c := m.Get(node.Cell.X, node.Cell.Y+1)
			n := &Node{c, node, c.Cost + node.Cost}
			if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (node.Cell.HeightLevel-n.Cell.HeightLevel) <= stepHeight {
				heap.Push(&openNodes, n)
				checkedNodes = append(checkedNodes, n.Cell)
			}
		}

		// Do the same thing for diagonals.
		if diagonals {

			diagonalCost := .414 // Diagonal movement is slightly slower, so we should prioritize straightaways if possible

			up := m.Get(node.Cell.X, node.Cell.Y-1).Walkable
			down := m.Get(node.Cell.X, node.Cell.Y+1).Walkable
			left := m.Get(node.Cell.X-1, node.Cell.Y).Walkable
			right := m.Get(node.Cell.X+1, node.Cell.Y).Walkable

			if node.Cell.X > 0 && node.Cell.Y > 0 {
				c := m.Get(node.Cell.X-1, node.Cell.Y-1)
				n := &Node{c, node, c.Cost + node.Cost + diagonalCost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (!wallsBlockDiagonals || (left && up)) && (node.Cell.HeightLevel-n.Cell.HeightLevel) <= stepHeight {
					heap.Push(&openNodes, n)
					checkedNodes = append(checkedNodes, n.Cell)
				}
			}

			if node.Cell.X < m.Width()-1 && node.Cell.Y > 0 {
				c := m.Get(node.Cell.X+1, node.Cell.Y-1)
				n := &Node{c, node, c.Cost + node.Cost + diagonalCost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (!wallsBlockDiagonals || (right && up)) && (node.Cell.HeightLevel-n.Cell.HeightLevel) <= stepHeight {
					heap.Push(&openNodes, n)
					checkedNodes = append(checkedNodes, n.Cell)
				}
			}

			if node.Cell.X > 0 && node.Cell.Y < m.Height()-1 {
				c := m.Get(node.Cell.X-1, node.Cell.Y+1)
				n := &Node{c, node, c.Cost + node.Cost + diagonalCost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (!wallsBlockDiagonals || (left && down)) && (node.Cell.HeightLevel-n.Cell.HeightLevel) <= stepHeight {
					heap.Push(&openNodes, n)
					checkedNodes = append(checkedNodes, n.Cell)
				}
			}

			if node.Cell.X < m.Width()-1 && node.Cell.Y < m.Height()-1 {
				c := m.Get(node.Cell.X+1, node.Cell.Y+1)
				n := &Node{c, node, c.Cost + node.Cost + diagonalCost}
				if n.Cell.Walkable && !hasBeenAdded(n.Cell) && (!wallsBlockDiagonals || (right && down)) && (node.Cell.HeightLevel-n.Cell.HeightLevel) <= stepHeight {
					heap.Push(&openNodes, n)
					checkedNodes = append(checkedNodes, n.Cell)
				}
			}

		}

	}

	return path

}

// GetPath returns a Path, from the starting cell's X and Y to the ending cell's X and Y. diagonals controls whether
// moving diagonally is acceptable when creating the Path. wallsBlockDiagonals indicates whether to allow diagonal movement "through" walls
// that are positioned diagonally. This is essentially just a smoother way to get a Path from GetPathFromCells().
func (m *Grid) GetPath(startX, startY, endX, endY float64, stepHeight int, diagonals bool, wallsBlockDiagonals bool) *Path {

	sc := m.Get(int(startX), int(startY))
	ec := m.Get(int(endX), int(endY))

	if sc != nil && ec != nil {
		return m.GetPathFromCells(sc, ec, stepHeight, diagonals, wallsBlockDiagonals)
	}
	return nil
}

// DataAsStringArray returns a 2D array of runes for each Cell in the Grid. The first axis is the Y axis.
func (m *Grid) DataAsStringArray() []string {

	data := []string{}

	for y := 0; y < m.Height(); y++ {
		data = append(data, "")
		for x := 0; x < m.Width(); x++ {
			data[y] += string(m.Data[y][x].Rune)
		}
	}

	return data

}

// DataAsRuneArrays returns a 2D array of runes for each Cell in the Grid. The first axis is the Y axis.
func (m *Grid) DataAsRuneArrays() [][]rune {

	runes := [][]rune{}

	for y := 0; y < m.Height(); y++ {
		runes = append(runes, []rune{})
		for x := 0; x < m.Width(); x++ {
			runes[y] = append(runes[y], m.Data[y][x].Rune)
		}
	}

	return runes

}

// A Path is a struct that represents a path, or sequence of Cells from point A to point B. The Cells list is the list of Cells contained in the Path,
// and the CurrentIndex value represents the current step on the Path. Using Path.Next() and Path.Prev() advances and walks back the Path by one step.
type Path struct {
	Cells                    []*Cell
	CurrentIndex, StepHeight int
}

// TotalCost returns the total cost of the Path (i.e. is the sum of all the Cells in the Path).
func (p *Path) TotalCost() float64 {

	cost := 0.0
	for _, cell := range p.Cells {
		cost += cell.Cost
	}
	return cost

}

// Reverse reverses the Cells in the Path.
func (p *Path) Reverse() {

	np := []*Cell{}

	for c := len(p.Cells) - 1; c >= 0; c-- {
		np = append(np, p.Cells[c])
	}

	p.Cells = np

}

// Restart restarts the Path, so that calling path.Current() will now return the first Cell in the Path.
func (p *Path) Restart() {
	p.CurrentIndex = 0
}

// Current returns the current Cell in the Path.
func (p *Path) Current() *Cell {
	return p.Cells[p.CurrentIndex]

}

// Next returns the next cell in the path. If the Path is at the end, Next() returns nil.
func (p *Path) Next() *Cell {

	if p.CurrentIndex < len(p.Cells)-1 {
		return p.Cells[p.CurrentIndex+1]
	}
	return nil

}

// Advance advances the path by one cell.
func (p *Path) Advance() {

	p.CurrentIndex++
	if p.CurrentIndex >= len(p.Cells) {
		p.CurrentIndex = len(p.Cells) - 1
	}

}

// Prev returns the previous cell in the path. If the Path is at the start, Prev() returns nil.
func (p *Path) Prev() *Cell {

	if p.CurrentIndex > 0 {
		return p.Cells[p.CurrentIndex-1]
	}
	return nil

}

// Same returns if the Path shares the exact same cells as the other specified Path.
func (p *Path) Same(otherPath *Path) bool {

	if p == nil || otherPath == nil || len(p.Cells) != len(otherPath.Cells) {
		return false
	}

	for i := range p.Cells {
		if len(otherPath.Cells) <= i || p.Cells[i] != otherPath.Cells[i] {
			return false
		}
	}

	return true

}

// Length returns the length of the Path (how many Cells are in the Path).
func (p *Path) Length() int {
	return len(p.Cells)
}

// Get returns the Cell of the specified index in the Path. If the index is outside of the
// length of the Path, it returns -1.
func (p *Path) Get(index int) *Cell {
	if index < len(p.Cells) {
		return p.Cells[index]
	}
	return nil
}

// Index returns the index of the specified Cell in the Path. If the Cell isn't contained
// in the Path, it returns -1.
func (p *Path) Index(cell *Cell) int {
	for i, c := range p.Cells {
		if c == cell {
			return i
		}
	}
	return -1
}

// SetIndex sets the index of the Path, allowing you to safely manually manipulate the Path
// as necessary. If the index exceeds the bounds of the Path, it will be clamped.
func (p *Path) SetIndex(index int) {

	if index >= len(p.Cells) {
		p.CurrentIndex = len(p.Cells) - 1
	} else if index < 0 {
		p.CurrentIndex = 0
	} else {
		p.CurrentIndex = index
	}

}

// IsAtStart returns if the Path's current index is 0, the first Cell in the Path.
func (p *Path) IsAtStart() bool {
	return p.CurrentIndex == 0
}

// IsAtEnd returns if the Path's current index is the last Cell in the Path.
func (p *Path) IsAtEnd() bool {
	return p.CurrentIndex >= len(p.Cells)-1
}

// Node represents the node a path, it contains the cell it represents.
// Also contains other information such as the parent and the cost.
type Node struct {
	Cell   *Cell
	Parent *Node
	Cost   float64
}

type minHeap []*Node

func (mH minHeap) Len() int           { return len(mH) }
func (mH minHeap) Less(i, j int) bool { return mH[i].Cost < mH[j].Cost }
func (mH minHeap) Swap(i, j int)      { mH[i], mH[j] = mH[j], mH[i] }
func (mH *minHeap) Pop() interface{} {
	old := *mH
	n := len(old)
	x := old[n-1]
	*mH = old[0 : n-1]
	return x
}

func (mH *minHeap) Push(x interface{}) {
	*mH = append(*mH, x.(*Node))
}

// check if a int is contained in a array
// bc go has no build in function for this
func containesInt(array []int, i int) bool {
	for _, current := range array {
		if current == i {
			return true
		}
	}
	return false
}
