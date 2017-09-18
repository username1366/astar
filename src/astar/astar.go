package main

import (
	"fmt"
	"os"
)

const (
	NumOfNeighbors = 8
	PassableCell   = "-"
	StartCell      = "@"
	Obstacles      = "#"
	FinishCell     = "☗"
	Step           = "4"
)

type Matrix struct {
	Data         [][]string // Input Data
	Cells        []Cell
	StartPoint   Coordinates
	FinishPoint  Coordinates
	CurrentPoint Coordinates
}

type Cell struct {
	CellIndex   int
	Type        string // Cell type
	G           int    // movement cost to StartPoint
	H           int    // estimated movement cost to FinishPoint
	F           int
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
	Diagonal bool
	c        Coordinates // Neighbor Cell coordinates
}

func (m *Matrix) Construct() {

	var i int = 0                                      // Cells iterator
	m.Cells = make([]Cell, len(m.Data)*len(m.Data[0])) // allocate memory for Cells slice

	// fill up m.Cells slice
	for y := 0; y < len(m.Data); y++ {
		for x := 0; x < len(m.Data[0]); x++ {
			m.Cells[i].Type = m.Data[y][x]
			m.Cells[i].Coordinates = Coordinates{x, y}
			m.Cells[i].Neighbors = make([]Neighbor, NumOfNeighbors)
			m.Cells[i].CellIndex = i
			if (i+1)%len(m.Data[0]) == 0 {
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
				m.Cells[i].OpenList = true
				m.StartPoint = Coordinates{x, y}
				m.CurrentPoint = m.StartPoint
			} else if m.Data[y][x] == FinishCell {
				m.FinishPoint = Coordinates{x, y}
			}

			i++
		}
	}
}

func (m *Matrix) EvaluateMovementCost() {

	i := m.GetCellIndex(m.CurrentPoint)
	if os.Getenv("DEBUG") == "1" {
		fmt.Printf("Current point: %v\n", m.GetCellIndex(m.CurrentPoint))
	}
	StartPointIndex := m.GetCellIndex(m.StartPoint)
	// determine F, G, H for each neighbor
	for _, v := range m.Cells[i].Neighbors {

		NeighborCellIndex := m.GetCellIndex(v.c)
		if os.Getenv("DEBUG") == "1" {
			fmt.Printf("Neighbor index: [%v] ", NeighborCellIndex)
		}

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

		if os.Getenv("DEBUG") == "1" {
			fmt.Printf("G: %v + ", m.Cells[NeighborCellIndex].G)
			fmt.Printf("H: %v = ", m.Cells[NeighborCellIndex].H)
			fmt.Printf("F: %v\n", m.Cells[NeighborCellIndex].F)
		}

		// add neighbors to OpenList
		if m.Cells[NeighborCellIndex].Type == PassableCell && !m.Cells[NeighborCellIndex].ClosedList {
			m.Cells[NeighborCellIndex].OpenList = true
			m.Cells[NeighborCellIndex].ClosedList = false
		}
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
				if os.Getenv("DEBUG") == "1" {
					fmt.Println(v.F, ">", n)
				}
			} else {
				if os.Getenv("DEBUG") == "1" {
					fmt.Println(v.F, "<", n)
				}
				n = v.F
				i = k // save to i the cell index with smallest F
			}
		}
	}

	m.CurrentPoint = m.Cells[i].Coordinates // setup new CurrentPoint
	m.Cells[i].Type = Step
	m.Cells[i].OpenList = false
	m.Cells[i].ClosedList = true
	if os.Getenv("DEBUG") == "1" {
		fmt.Printf("Chose [%v] cell index with F = %v for next step\n", i, n)
	}
	return i
}

func (m Matrix) Print() {
	for _, v := range m.Cells {
		fmt.Printf("%v (%v)[%v]\t\t", v.Type, v.G, v.CellIndex)
		if v.Delimiter {
			fmt.Printf("\n")
		}

	}
}

func (m Matrix) PrintResult() {
	for _, v := range m.Cells {
		if v.Type == Step {
			fmt.Printf("+\t")
		} else {
			fmt.Printf("%v\t", v.Type)
		}
		if v.Delimiter {
			fmt.Printf("\n")
		}

	}
	fmt.Println("\n")
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
	m.Data = [][]string{
		{"-", "-", "-", "-", "#", "-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "-", "-", "-", "-"}, // 0
		{"-", "-", "-", "-", "#", "-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "-", "-", "-", "-"}, // 1
		{"-", "-", "-", "-", "#", "-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "-", "-", "-", "-"}, // 2
		{"-", "-", "@", "-", "#", "-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "-", "-", "-", "-"}, // 3
		{"-", "-", "-", "-", "#", "-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "-", "-", "-", "-"}, // 4
		{"-", "-", "-", "-", "#", "-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "-", "-", "-", "-"}, // 5
		{"-", "-", "-", "-", "#", "-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "#", "-", "-", "-"}, // 6
		{"-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "#", "-", "-", "#", "-", "-", "-", "-", "-", "#", "-", "-", "-"}, // 7
		{"-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "#", "-", "-", "-", "-", "-", "#", "-", "-", "-"}, // 8
		{"-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "-", "-", "-", "-", "-", "#", "-", "-", "☗"}, // 9
		{"-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "-", "-", "-", "-", "-", "#", "-", "-", "-"}, // 10
		{"-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "-", "-", "-", "-", "-", "#", "-", "-", "-"}, // 11
		{"-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "-", "-", "-", "-", "-", "#", "-", "-", "-"}, // 12
		{"-", "-", "-", "-", "-", "-", "-", "-", "#", "#", "-", "-", "-", "-", "-", "-", "-", "-", "-", "-", "-", "-", "-"}, // 13

	}
	m.Construct()
	for i := 0; i < len(m.Data)*len(m.Data[0])*4; i++ {
		m.EvaluateMovementCost()
		m.Move()

		if os.Getenv("DEBUG") == "1" {
			m.Print()
			m.PrintOpenList()
			m.PrintClosedList()
		}
		fmt.Println("Step:", i)
		m.PrintResult()

		CellIndex := m.GetCellIndex(m.CurrentPoint)
		if m.Cells[CellIndex].H == 10 {
			fmt.Println("I've found FinishPoint!\n")
			break
		}
	}
}
