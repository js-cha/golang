package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func main() {
	ttt := TicTacToe{}
	ttt.Play()
}

type Player interface {
	MakePlay(g Grid)
}

type Mark uint8

const (
	Blank Mark = iota
	Nought
	Cross
)

func (m Mark) String() string {
	return []string{"-", "O", "X"}[m]
}

type Grid [][]Mark

func (g Grid) Print() {
	fmt.Println(g[0])
	fmt.Println(g[1])
	fmt.Println(g[2])
}

func (g Grid) Update(row, col int, mark Mark) (err error) {
	if g[row][col] == Blank {
		g[row][col] = mark
		return nil
	}
	return fmt.Errorf("position %d,%d is already occupied with %s", row, col, g[row][col])
}

func New() Grid {
	grid := Grid{
		{Blank, Blank, Blank},
		{Blank, Blank, Blank},
		{Blank, Blank, Blank},
	}

	return grid
}

func IsTicTacToeHorizontal(grid Grid) bool {
	return (grid[0][0] == grid[0][1] && grid[0][0] == grid[0][2] && grid[0][0] != Blank) ||
		(grid[1][0] == grid[1][1] && grid[1][0] == grid[1][2] && grid[1][0] != Blank) ||
		(grid[2][0] == grid[2][1] && grid[2][0] == grid[2][2] && grid[2][0] != Blank)
}

func IsTicTacToeVertical(grid Grid) bool {
	return (grid[0][0] == grid[1][0] && grid[0][0] == grid[2][0] && grid[0][0] != Blank) ||
		(grid[0][1] == grid[1][1] && grid[0][1] == grid[2][1] && grid[0][1] != Blank) ||
		(grid[0][2] == grid[1][2] && grid[0][2] == grid[2][2] && grid[0][2] != Blank)
}

func IsTicTacToeDiagonal(grid Grid) bool {
	return (grid[0][0] == grid[1][1] && grid[0][0] == grid[2][2] && grid[0][0] != Blank) ||
		(grid[0][2] == grid[1][1] && grid[0][2] == grid[2][0] && grid[0][2] != Blank)
}

func IsWinner(grid Grid) bool {
	return IsTicTacToeDiagonal(grid) || IsTicTacToeHorizontal(grid) || IsTicTacToeVertical(grid)
}

func IsBoardFull(grid Grid) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if grid[i][j] == Blank {
				return false
			}
		}
	}
	return true
}

type Person struct {
	Mark
}

func (p Person) MakePlay(g Grid) {
	fmt.Printf("info: %s's turn, please make a play on the grid, example: 0,1\n", p.Mark)
	g.Print()

	row, col, error := p.ChooseCoordinates()

	if error != nil {
		fmt.Println(error)
		g.Print()
		row, col, error = p.ChooseCoordinates()
	}

	err := g.Update(row, col, p.Mark)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Please choose different coordinates")
		p.MakePlay(g)
	}
}

func (p Person) ChooseCoordinates() (int, int, error) {
	var input string
	fmt.Scanf("%s\n", &input)
	coords := strings.Split(input, ",")
	row, err := strconv.Atoi(coords[0])

	if err != nil {
		fmt.Println("An error occured")
	}

	col, err := strconv.Atoi(coords[1])

	if err != nil {
		fmt.Println("An error occured")
	}

	if row > 2 || row < 0 || col > 2 || col < 0 {
		return 0, 0, fmt.Errorf("error: invalid coordinates, please try again")
	}

	return row, col, nil
}

type TicTacToe struct {
	grid Grid
	turn Player
	p1   Player
	p2   Player
}

func (t *TicTacToe) Setup() {
	t.grid = New()

	fmt.Println("Please select an option for Player 1:")
	fmt.Println("- Press [1] for Human")
	fmt.Println("- Press [2] for Bot")
	var p1 string
	fmt.Scanf("%s\n", &p1)

	if p1 == "1" {
		t.p1 = Person{Cross}
	} else {
		t.p1 = Bot{Cross}
	}

	fmt.Println("Please select an option for Player 2:")
	fmt.Println("- Press [1] for Human")
	fmt.Println("- Press [2] for Bot")
	var p2 string
	fmt.Scanf("%s\n", &p2)

	if p2 == "1" {
		t.p2 = Person{Nought}
	} else {
		t.p2 = Bot{Nought}
	}

	fmt.Println()
	fmt.Printf("Player 1 is %s\n", t.p1)
	fmt.Printf("Player 2 is %s\n", t.p2)
	fmt.Println("Starting game...")
	fmt.Println()
}

func (t TicTacToe) Play() {
	t.Setup()
	t.turn = t.p1
	for {
		if IsWinner(t.grid) {
			winner := t.p2
			if t.turn == t.p2 {
				winner = t.p1
			}
			t.grid.Print()
			fmt.Printf("Winner is %s\n", winner)
			return
		}

		if IsBoardFull(t.grid) {
			t.grid.Print()
			fmt.Println("Draw!")
			return
		}

		if t.turn == t.p1 {
			t.p1.MakePlay(t.grid)
			t.turn = t.p2
		} else {
			t.p2.MakePlay(t.grid)
			t.turn = t.p1
		}
	}
}

type Bot struct {
	Mark
}

func (b Bot) GetAvailableCoordinates(g Grid) []string {
	availableCoords := []string{}

	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g); j++ {
			if g[i][j] == Blank {
				coords := fmt.Sprintf("%d,%d", i, j)
				availableCoords = append(availableCoords, coords)
			}
		}
	}

	return availableCoords
}

func (b Bot) MakePlay(g Grid) {
	availableCoords := b.GetAvailableCoordinates(g)

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	randomIndex := r.Intn(len(availableCoords))
	randomCoords := availableCoords[randomIndex]
	coords := strings.Split(randomCoords, ",")

	row, _ := strconv.Atoi(coords[0])
	col, _ := strconv.Atoi(coords[1])

	fmt.Printf("info: %s's turn, please make a play on the grid, example: 0,1\n", b.Mark)
	fmt.Printf("info: %s makes a play on %d,%d\n", b.Mark, row, col)
	g.Update(row, col, b.Mark)
	g.Print()
}
