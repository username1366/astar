package main

import (
	"fmt"
)

type Matrix struct {
	Data         [][]int // Input Data
	Cells        []Cell
	StartPoint   Coordinates
	FinishPoint  Coordinates
	CurrentPoint Coordinates
	Path         []Coordinates
	OpenList     []Cell
	ClosedList   []Cell
}

type Cell struct {
	CellIndex   int
	Type        int // Cell type
	G           int // movement cost to StartPoint
	H           int // estimated movement cost to FinishPoint
	F           int
	StartPoint  bool // is StartPoint?
	FinishPoint bool // is FinishPoint?
	OpenList    bool
	ClosedList  bool
	Coordinates Coordinates // x, y Coordinates
	Neighbors   []Neighbor  // Cell's Neighbors
}

type Coordinates struct {
	x int
	y int
}

type Neighbor struct {
	ParentCellIndex int
	Diagonal        bool
	c               Coordinates // Neighbor Cell coordinates
}

const (
	NumOfNeighbors = 8
	PassableCell   = 0
	StartCell      = 1
	Obstacles      = 2
	FinishCell     = 3
)

func (m *Matrix) Construct() {

	var i int = 0                                      // Cells iterator
	m.Cells = make([]Cell, len(m.Data)*len(m.Data[0])) // allocate memory for Cells slice

	// fill up m.Cells
	for y := 0; y < len(m.Data); y++ {
		for x := 0; x < len(m.Data[0]); x++ {
			m.Cells[i].Type = m.Data[y][x]
			m.Cells[i].Coordinates = Coordinates{x, y}
			m.Cells[i].Neighbors = make([]Neighbor, NumOfNeighbors)
			m.Cells[i].CellIndex = i

			// determine 8 neighbors for each Cell
			for n := 0; n < NumOfNeighbors; n++ {
				switch n {
				case 0:
					m.Cells[i].Neighbors[n].c.x = x - 1
					m.Cells[i].Neighbors[n].c.y = y - 1
					m.Cells[i].Neighbors[n].Diagonal = true
				case 1:
					m.Cells[i].Neighbors[n].c.x = x
					m.Cells[i].Neighbors[n].c.y = y - 1
				case 2:
					m.Cells[i].Neighbors[n].c.x = x + 1
					m.Cells[i].Neighbors[n].c.y = y - 1
					m.Cells[i].Neighbors[n].Diagonal = true
				case 3:
					m.Cells[i].Neighbors[n].c.x = x + 1
					m.Cells[i].Neighbors[n].c.y = y
				case 4:
					m.Cells[i].Neighbors[n].c.x = x + 1
					m.Cells[i].Neighbors[n].c.y = y + 1
					m.Cells[i].Neighbors[n].Diagonal = true
				case 5:
					m.Cells[i].Neighbors[n].c.x = x
					m.Cells[i].Neighbors[n].c.y = y + 1
				case 6:
					m.Cells[i].Neighbors[n].c.x = x - 1
					m.Cells[i].Neighbors[n].c.y = y + 1
					m.Cells[i].Neighbors[n].Diagonal = true
				case 7:
					m.Cells[i].Neighbors[n].c.x = x - 1
					m.Cells[i].Neighbors[n].c.y = y
				}
			}

			// setup StartPoint, FinishPoint and add StartPoint to OpenList and Path
			if m.Data[y][x] == StartCell {
				m.Cells[i].StartPoint = true
				m.Cells[i].OpenList = true
				m.StartPoint = Coordinates{x, y}
				m.CurrentPoint = m.StartPoint
				m.OpenList = append(m.OpenList, m.Cells[i])
				m.Path = append(m.Path, m.StartPoint)
			} else if m.Data[y][x] == FinishCell {
				m.Cells[i].FinishPoint = true
				m.FinishPoint = Coordinates{x, y}
			}

			i++
		}
	}
}

func (m *Matrix) EvaluateMovementCost() {

	var i int
	if m.StartPoint == m.CurrentPoint {
		// for first step
		fmt.Println("StartPoint is equal to CurrentPoint\n")
		i = m.GetCellIndex(m.StartPoint)
	} else {
		// for other steps
	}

	// determine F, G, H for each neighbor
	for _, v := range m.Cells[i].Neighbors {

		NeighborCellIndex := m.GetCellIndex(v.c)
		fmt.Printf("Neighbor: [%v, %v] ", v.c.x, v.c.y)

		// determine G
		if v.Diagonal {
			m.Cells[NeighborCellIndex].G = m.Cells[(v.ParentCellIndex)].G + 14
		} else {
			m.Cells[NeighborCellIndex].G = m.Cells[(v.ParentCellIndex)].G + 10
		}

		// determine H
		var hx, hy int

		if hx = (m.FinishPoint.x - v.c.x); hx < 0 {
			hx *= -1
		}
		if hy = (m.FinishPoint.y - v.c.y); hy < 0 {
			hy *= -1
		}

		m.Cells[NeighborCellIndex].H = (hx + hy) * 10
		m.Cells[NeighborCellIndex].F = m.Cells[NeighborCellIndex].G + m.Cells[NeighborCellIndex].H // F = G + H

		fmt.Printf("G: %v + ", m.Cells[NeighborCellIndex].G)
		fmt.Printf("H: %v = ", m.Cells[NeighborCellIndex].H)
		fmt.Printf("F: %v\n", m.Cells[NeighborCellIndex].F)

		v.ParentCellIndex = i // save parrent index for this neighbor

		// add neighbors to OpenList
		if m.Cells[NeighborCellIndex].Type == PassableCell {
			m.Cells[NeighborCellIndex].OpenList = true
			m.OpenList = append(m.OpenList, m.Cells[NeighborCellIndex])
		}

		fmt.Printf("Parent : [%v, %v]\n", m.Cells[i].Coordinates.x, m.Cells[i].Coordinates.y)
	}
	m.Cells[i].OpenList = false
	//m.PrintOpenList()

}

func (m Matrix) PrintMatrix() {
	k := 0
	l := 0
	fmt.Printf("*---0-------1-------2-------3-------4-------5-------6-------7-> x\n|\n")
	for i := 0; i < len(m.Cells); i++ {
		if k == 0 {
			fmt.Printf("%d \t%v\t", l, m.Cells[i].Type)
			l++
		} else {
			fmt.Printf("\t%v\t", m.Cells[i].Type)
		}
		if k == 7 && i < len(m.Cells)-1 {
			fmt.Printf("\n|\n")
			k = 0
		} else {
			k++
		}
		if i == len(m.Cells)-1 {
			fmt.Printf("\n|\nv\n")
		}
	}
}

func (m Matrix) PrintCellNeighbors(i int) {
	fmt.Printf("%v\n", m.Cells[i].Neighbors)
}

func (m Matrix) PrintOpenList() {
	fmt.Printf("OpenList struct: \n")
	for _, v := range m.OpenList {
		fmt.Printf("[%v, %v]\n", v.Coordinates.x, v.Coordinates.y)
		fmt.Printf("G: %v, H: %v\n", v.G, v.H)
	}
}

func (m Matrix) Move() int {
	var i, n int
	n = len(m.Data) * len(m.Data[0]) * 4
	for k, v := range m.Cells {
		if v.OpenList {
			if v.F > n {
				//fmt.Println(v.F, ">", n)
			} else {
				//fmt.Println(v.F, "<", n)
				n = v.F
				i = k // save to i the cell index with smallest F
			}
		}
	}

	m.CurrentPoint = m.Cells[i].Coordinates // setup new CurrentPoint
	fmt.Printf("Chosed [%v] cell index with F = %v for next step\n", i, n)
	fmt.Printf("Index [%v, %v]\n", m.Cells[i].Coordinates.x, m.Cells[i].Coordinates.y)
	return i
}

func (m Matrix) GetCellIndex(c Coordinates) int {
	var CellIndex int
	for k, v := range m.Cells {
		if v.Coordinates == c {
			CellIndex = k
		}
	}
	return CellIndex
}

func main() {
	var m Matrix
	m.Data = [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0}, // 0
		{0, 0, 0, 0, 0, 0, 0, 0}, // 1
		{0, 0, 0, 0, 2, 0, 0, 0}, // 2
		{0, 0, 0, 0, 2, 0, 0, 0}, // 3
		{0, 0, 0, 0, 2, 0, 0, 0}, // 4
		{0, 1, 0, 0, 2, 0, 0, 0}, // 5
		{0, 0, 0, 0, 2, 0, 0, 3}, // 6
		{0, 0, 0, 0, 2, 0, 0, 0}, // 7
		//  1, 2, 3, 4, 5, 6, 7
	}
	m.Construct()
	m.PrintMatrix()
	m.EvaluateMovementCost()
	m.PrintOpenList()
	m.Move()

	fmt.Println(m.StartPoint)
}
