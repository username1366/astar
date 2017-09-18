package main

import (
	"fmt"
)

const (
	NumOfNeighbors = 8
	PassableCell   = 0
	StartCell      = 1
	Obstacles      = 2
	FinishCell     = 3
	Columns        = 8
	Rows           = 8
)

type Matrix struct {
	Data         [][]int // Input Data
	Cells        []Cell
	StartPoint   Coordinates
	FinishPoint  Coordinates
	CurrentPoint Coordinates
	Path         []Coordinates // traversed Path
	//OpenList     []Cell
	//ClosedList   []Cell
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
	Delimiter   bool
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
			if (i+1)%Columns == 0 {
				m.Cells[i].Delimiter = true
			}

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

	i := m.GetCellIndex(m.CurrentPoint)
	fmt.Printf("Current point: %v\n", m.GetCellIndex(m.CurrentPoint))
	StartPointIndex := m.GetCellIndex(m.StartPoint)
	// determine F, G, H for each neighbor
	for _, v := range m.Cells[i].Neighbors {

		NeighborCellIndex := m.GetCellIndex(v.c)
		//fmt.Printf("Neighbor: [%v, %v] ", v.c.x, v.c.y)
		fmt.Printf("Neighbor index: [%v] ", NeighborCellIndex)

		// determine G
		if v.Diagonal {
			m.Cells[NeighborCellIndex].G = m.Cells[StartPointIndex].G + 14
		} else {
			m.Cells[NeighborCellIndex].G = m.Cells[StartPointIndex].G + 10
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
		if m.Cells[NeighborCellIndex].Type == PassableCell && !m.Cells[NeighborCellIndex].ClosedList {
			m.Cells[NeighborCellIndex].OpenList = true
			m.Cells[NeighborCellIndex].ClosedList = false
			//m.OpenList = append(m.OpenList, m.Cells[NeighborCellIndex])
		}

		//fmt.Printf("Parent index : %v\n", v.ParentCellIndex)
		//fmt.Printf("Parent : [%v, %v]\n", m.Cells[v.ParentCellIndex].Coordinates.x, m.Cells[v.ParentCellIndex].Coordinates.y)
	}
	m.Cells[i].OpenList = false
	m.Cells[i].ClosedList = true
}

func (m *Matrix) Move() int {
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
	m.Cells[i].Type = 7
	m.Cells[i].OpenList = false
	m.Cells[i].ClosedList = true
	fmt.Printf("Chose [%v] cell index with F = %v for next step\n", i, n)
	//fmt.Printf("Index [%v, %v]\n", m.Cells[i].Coordinates.x, m.Cells[i].Coordinates.y)
	return i
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

func (m Matrix) Print() {
	for _, v := range m.Cells {
		fmt.Printf("%v (%v)[%v]\t\t", v.Type, v.G, v.CellIndex)
		if v.Delimiter {
			fmt.Printf("\n")
		}

	}
}

func (m Matrix) PrintCellNeighbors(i int) {
	fmt.Printf("%v\n", m.Cells[i].Neighbors)
}

func (m Matrix) PrintOpenList() {
	fmt.Printf("OpenList:\n")
	for _, v := range m.Cells {
		if v.OpenList {
			fmt.Printf("%v\n", v.CellIndex)
		}
	}
}

func (m Matrix) PrintClosedList() {
	fmt.Printf("ClosedList:\n")
	for _, v := range m.Cells {
		if v.ClosedList {
			fmt.Printf("%v\n", v.CellIndex)
		}
	}
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
		{0, 0, 0, 0, 0, 0, 0, 0}, // 7
		//  1, 2, 3, 4, 5, 6, 7
	}
	m.Construct()
	for i := 0; i < 10000; i++ {
		m.Print()
		m.EvaluateMovementCost()
		m.Move()
		m.PrintOpenList()
		m.PrintClosedList()
		CellIndex := m.GetCellIndex(m.CurrentPoint)
		if m.Cells[CellIndex].H == 10 {
			break
		}
	}

	m.Print()

	// m.Move() < m.FinishPoint

	//fmt.Println(m.StartPoint)
	//m.Print()

}
