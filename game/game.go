package game

import (
	"bufio"
	"fmt"
	"os"
)

type GameUI interface {
	Draw(*Level)
	GetInput() *Input
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	Quit
)

type Input struct {
	Typ InputType
}

type Tile rune

const (
	StoneWall Tile = '#'
	DirtFloor Tile = '.'
	Door      Tile = '|'
	Blank     Tile = 0
	Pending   Tile = -1
)

type Entity struct {
	X, Y int
}

type Player struct {
	Entity
}

type Level struct {
	Map    [][]Tile
	Player Player
}

func loadLevelFromFile(filename string) *Level {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	levelLines := make([]string, 0)
	longestRow := 0
	index := 0

	for scanner.Scan() {
		levelLines = append(levelLines, scanner.Text())
		if len(levelLines[index]) > longestRow {
			longestRow = len(levelLines[index])
		}
		index++
	}
	level := &Level{}
	level.Map = make([][]Tile, len(levelLines))

	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}

	for y := 0; y < len(level.Map); y++ {
		line := levelLines[y]
		for x, c := range line {
			var t Tile
			switch c {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case '#':
				t = StoneWall
			case '|':
				t = Door
			case '.':
				t = DirtFloor
			case 'P':
				t = Pending
				level.Player.X = x
				level.Player.Y = y
			default:
				panic("Invalid character in map")
			}
			level.Map[y][x] = t
		}

	}

	for y, row := range level.Map {
		for x, tile := range row {
			if tile == Pending {
			SearchLoop:
				for searchX := x - 1; searchX <= x+1; searchX++ {
					for searchY := y - 1; searchY <= y+1; searchY++ {
						searchTile := level.Map[searchX][searchY]
						switch searchTile {
						case DirtFloor:
							level.Map[y][x] = DirtFloor
							break SearchLoop
						}

					}
				}
			}
		}
	}

	return level

}

func handleInput(level *Level, input *Input) {
	switch input.Typ {
	case Up:
		fmt.Println("Up")
		level.Player.Y--
	case Down:
		fmt.Println("Down")
		level.Player.Y++
	case Left:
		fmt.Println("Left")
		level.Player.X--
	case Right:
		fmt.Println("Right")
		level.Player.X++
	}
}

func Run(ui GameUI) {
	fmt.Println("Starting...")
	level := loadLevelFromFile("game/maps/level1.map")

	for {
		ui.Draw(level)
		input := ui.GetInput()

		if input != nil && input.Typ == Quit {
			return
		}

		handleInput(level, input)
	}

}
