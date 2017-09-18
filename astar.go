package main

import (
	"fmt"
)

type Matrix struct {
	Data        [][]int // Input Data
	Cells       []Cell
	StartPoint  Coordinates //???
	FinishPoint Coordinates //???
	NextPoint   []Coordinates
}

type Cell struct {
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
	//CellIndex       int
	ParentCellIndex int
	Diagonal        bool
	//ParentCellCoordinates Coordinates
	c Coordinates // Neighbor Cell coordinates
}

func (m *Matrix) Construct() {
	var i int = 0                                      // Cells iterator
	m.Cells = make([]Cell, len(m.Data)*len(m.Data[0])) // allocate memory for Cells slice
	for y := 0; y < len(m.Data); y++ {
		for x := 0; x < len(m.Data[0]); x++ {
			if m.Data[y][x] == 1 {
				m.Cells[i].StartPoint = true
				m.Cells[i].OpenList = true
			}
			if m.Data[y][x] == 3 {
				m.Cells[i].FinishPoint = true
			}
			m.Cells[i].Type = m.Data[y][x]
			m.Cells[i].Coordinates.x = x
			m.Cells[i].Coordinates.y = y
			m.Cells[i].Neighbors = make([]Neighbor, 8)
			for n := 0; n < 8; n++ {
				//m.Cells[i].Neighbors[n].ParentCellIndex = i
				switch n {
				case 0:
					m.Cells[i].Neighbors[n].c.x = x - 1
					m.Cells[i].Neighbors[n].c.y = y - 1
					m.Cells[i].Neighbors[n].Diagonal = true
					//m.Cells[i].Neighbors[n].ParentCellCoordinates =
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
				default:
					fmt.Println("????")
				}

				//fmt.Printf("%v\n", m.GetCellIndex(m.Cells[i].Neighbors[n].c.x, m.Cells[i].Neighbors[n].c.y))
			}
			i++
		}
	}
}

func (m *Matrix) EvaluateMovementCost() {
	x, y := m.GetStartPoint()
	i := m.GetCellIndex(x, y)
	m.NextPoint = append(m.NextPoint, Coordinates{x, y})
	var NeighborCellIndex int
	//fmt.Printf("%v\n", m.Cells[i].OpenList)
	for _, v := range m.Cells[i].Neighbors {
		NeighborCellIndex = m.GetCellIndex(v.c.x, v.c.y)
		fmt.Printf("Neighbor: [%v, %v] ", v.c.x, v.c.y)
		if m.Cells[NeighborCellIndex].Type == 0 {
			m.Cells[NeighborCellIndex].OpenList = true
		}
		if v.Diagonal {
			m.Cells[NeighborCellIndex].G = m.Cells[(v.ParentCellIndex)].G + 14
		} else {
			m.Cells[NeighborCellIndex].G = m.Cells[(v.ParentCellIndex)].G + 10
		}

		fx, fy := m.GetFinishPoint()
		var hx, hy int
		if hx = (fx - v.c.x); hx < 0 {
			hx *= -1
		}
		if hy = (fy - v.c.y); hy < 0 {
			hy *= -1
		}
		m.Cells[NeighborCellIndex].H = (hx + hy) * 10
		m.Cells[NeighborCellIndex].F = m.Cells[NeighborCellIndex].G + m.Cells[NeighborCellIndex].H

		fmt.Printf("G: %v + ", m.Cells[NeighborCellIndex].G)
		fmt.Printf("H: %v = ", m.Cells[NeighborCellIndex].H)
		fmt.Printf("F: %v\n", m.Cells[NeighborCellIndex].F)

		v.ParentCellIndex = i
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
	fmt.Printf("OpenList: \n")
	for k, v := range m.Cells {
		if v.OpenList {
			fmt.Printf("[%v, %v]\n", m.Cells[k].Coordinates.x, m.Cells[k].Coordinates.y)
		}
	}
}

/*
func (m Matrix) GetOpenList(OpenList [][]int) {
	//OpenList
	for _, v := range m.Cells {
		if v.OpenList {
			//m.Cells[k].Coordinates.x, m.Cells[k].Coordinates.y
		}
	}
}*/

func (m Matrix) Move() int {
	var i, n int
	n = len(m.Data) * len(m.Data[0]) * 4
	for k, v := range m.Cells {
		if v.OpenList {
			if v.F > n {
				fmt.Println(v.F, ">", n)
			} else {
				fmt.Println(v.F, "<", n)
				n = v.F
				i = k
			}
		}
	}
	fmt.Printf("Chosed [%v] cell index with F = %v for next step\n", i, n)
	fmt.Printf("Index [%v, %v]\n", m.Cells[i].Coordinates.x, m.Cells[i].Coordinates.y)
	return i
}

func (m Matrix) GetCellIndex(x, y int) int {
	for k, v := range m.Cells {
		if v.Coordinates.x == x && v.Coordinates.y == y {
			return k
		}
	}
	return -1
}

func (m Matrix) GetStartPoint() (int, int) {
	for _, v := range m.Cells {
		if v.StartPoint {
			return v.Coordinates.x, v.Coordinates.y
		}
	}
	return -1, -1
}

func (m Matrix) GetFinishPoint() (int, int) {
	for _, v := range m.Cells {
		if v.FinishPoint {
			return v.Coordinates.x, v.Coordinates.y
		}
	}
	return -1, -1
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
	//x, y := m.GetStartPoint()
	//fmt.Printf("[%v, %v]\n\n", x, y)
	m.EvaluateMovementCost()
	m.PrintOpenList()
	m.Move()
	//fmt.Printf("%v", m.GetCellIndex(1, 5))
	//fmt.Println(m.NextPoint)
}
