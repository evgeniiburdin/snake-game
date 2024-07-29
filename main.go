package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
	"gopkg.in/yaml.v2"
)

const (
	GAME_SPEED_COEFFICIENT = 868.75
	MAX_SNAKE_SPEED        = 150.0
)

func (game *SnakeGame) Render() {
	var out strings.Builder

	CLIWidth, CLIHeight := termbox.Size()

	fieldWidth := len(game.Field) * 3
	fieldHeight := len(game.Field)
	xOffset := (CLIWidth - fieldWidth) / 2
	yOffset := (CLIHeight - fieldHeight - 5) / 2

	out.WriteString(strings.Repeat("\n", yOffset))

	actualGameFieldSize := len(game.Field)*3 + 2
	out.WriteString(strings.Repeat(" ", CLIWidth+(actualGameFieldSize-7)/2) + "Score: " + strconv.Itoa(game.Snake.Length) + "\n\n")

	out.WriteString(game.RenderGlyphs["left_up_corner"])
	for i := 0; i < len(game.Field); i++ {
		out.WriteString(game.RenderGlyphs["up"])
	}
	out.WriteString(game.RenderGlyphs["right_up_corner"] + "\n")

	for i := 0; i < fieldHeight; i++ {
		out.WriteString(game.RenderGlyphs["left"])
		for j := 0; j < len(game.Field); j++ {
			if game.Field[i][j] == 1 {
				out.WriteString(fmt.Sprintf("\033["+game.RenderGlyphs["snake_color"]) + "m" + game.RenderGlyphs["snake_body"] + "\033[0m")
			} else if game.Field[i][j] == 2 {
				out.WriteString("\033[" + game.RenderGlyphs["apple_color"] + "m" + game.RenderGlyphs["apple"] + "\033[0m")
			} else {
				out.WriteString(game.RenderGlyphs["empty_cell"])
			}
		}
		out.WriteString(game.RenderGlyphs["right"] + "\n")
	}

	out.WriteString(game.RenderGlyphs["left_down_corner"])
	for i := 0; i < len(game.Field); i++ {
		out.WriteString(game.RenderGlyphs["down"])
	}
	out.WriteString(game.RenderGlyphs["right_down_corner"] + "\n")

	out.WriteString("\n")

	controlsHints := []string{
		"← - left",
		"→ - right",
		"↑ - up",
		"↓ - down",
		"Ctrl + P - pause",
		"ESC - exit",
	}

	for _, hint := range controlsHints {
		actualFieldSize := len(game.Field)*3 + 2
		out.WriteString(strings.Repeat(" ", CLIWidth+1+(actualFieldSize-len(hint))/2) + hint + "\n")
	}

	if game.IsPaused {
		pauseMessage := "GAME PAUSED (Press P to continue)"
		xPauseOffset := (CLIWidth - len(pauseMessage)) / 2
		out.WriteString(strings.Repeat(" ", xPauseOffset) + pauseMessage + "\n")
	}

	_, err := game.Output.Write([]byte("\033[H\033[2J"))
	if err != nil {
		panic(err)
	}

	lines := strings.Split(out.String(), "\n")
	for i := range lines {
		lines[i] = strings.Repeat(" ", xOffset) + lines[i]
	}

	_, err = game.Output.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		panic(err)
	}
}

type Config struct {
	LeftUpCorner    string `yaml:"left_up_corner"`
	Up              string `yaml:"up"`
	RightUpCorner   string `yaml:"right_up_corner"`
	Left            string `yaml:"left"`
	Right           string `yaml:"right"`
	LeftDownCorner  string `yaml:"left_down_corner"`
	Down            string `yaml:"down"`
	RightDownCorner string `yaml:"right_down_corner"`
	SnakeBody       string `yaml:"snake_body"`
	EmptyCell       string `yaml:"empty"`
	Apple           string `yaml:"apple"`
	SnakeColor      string `yaml:"snake_color"`
	FieldSize       int    `yaml:"field_size"`
	AppleColor      string `yaml:"apple_color"`
	//Output        string `yaml:"output"`
}

type SnakeNode struct {
	X    int
	Y    int
	Next *SnakeNode
	Prev *SnakeNode
}

type Snake struct {
	Head   *SnakeNode
	Tail   *SnakeNode
	Length int
}

func NewSnake(length, startX, startY int) *Snake {
	head := &SnakeNode{
		X: startX,
		Y: startY,
	}
	current := head

	for i := 1; i < length; i++ {
		newNode := &SnakeNode{
			X: current.X - 1,
			Y: current.Y,
		}
		current.Next = newNode
		newNode.Prev = current
		current = newNode
	}

	tail := current

	return &Snake{
		Head:   head,
		Tail:   tail,
		Length: length,
	}
}

type SnakeGame struct {
	Snake            *Snake
	Field            [][]int
	SnakeDirection   string
	Result           string
	UpdateEvery      time.Duration
	SpeedCoefficient float64
	RenderGlyphs     map[string]string
	Output           io.Writer
	IsPaused         bool
}

func NewSnakeGame(output io.Writer, cfg *Config) *SnakeGame {
	field := make([][]int, cfg.FieldSize)

	for i := 0; i < cfg.FieldSize; i++ {
		field[i] = make([]int, cfg.FieldSize)
	}

	snake := NewSnake(3, cfg.FieldSize/2, cfg.FieldSize/2)
	game := &SnakeGame{
		Snake:          snake,
		Field:          field,
		SnakeDirection: "right",
		Result:         "in-game",
		RenderGlyphs: map[string]string{
			"left_up_corner":    cfg.LeftUpCorner,
			"up":                cfg.Up,
			"right_up_corner":   cfg.RightUpCorner,
			"left":              cfg.Left,
			"right":             cfg.Right,
			"left_down_corner":  cfg.LeftDownCorner,
			"down":              cfg.Down,
			"right_down_corner": cfg.RightDownCorner,
			"snake_body":        cfg.SnakeBody,
			"empty_cell":        cfg.EmptyCell,
			"apple":             cfg.Apple,
			"snake_color":       cfg.SnakeColor,
			"apple_color":       cfg.AppleColor,
		},
		Output: output,
	}
	game.Initialize()

	game.SpeedCoefficient = GAME_SPEED_COEFFICIENT
	game.AdjustGameSpeed()

	return game
}

func (game *SnakeGame) Initialize() {
	currentCell := game.Snake.Head

	for currentCell != nil {
		game.Field[currentCell.Y][currentCell.X] = 1
		currentCell = currentCell.Next
	}

	game.SpawnApple()
}

func (game *SnakeGame) SpawnApple() {
	rand.Seed(time.Now().UnixNano())
	randX := rand.Intn(len(game.Field))
	randY := rand.Intn(len(game.Field))

	for game.Field[randX][randY] != 0 {
		randX = rand.Intn(len(game.Field))
		randY = rand.Intn(len(game.Field))
	}

	game.Field[randX][randY] = 2
}

func (game *SnakeGame) AdjustGameSpeed() {
	maxSpeed := MAX_SNAKE_SPEED

	fieldArea := len(game.Field) * len(game.Field)
	minSpeed := -(3.0/16.0)*float64(fieldArea) + game.SpeedCoefficient

	normalizedLength := float64(game.Snake.Length) / float64(fieldArea)

	speed := minSpeed - (minSpeed-maxSpeed)*normalizedLength
	game.UpdateEvery = time.Millisecond * time.Duration(speed)
}

func (game *SnakeGame) Play() {
	for game.Result == "in-game" {
		if game.IsPaused {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		game.MoveSnake()
		game.Render()

		time.Sleep(game.UpdateEvery)
	}
}

func (game *SnakeGame) MoveSnake() {
	fieldArea := len(game.Field) * len(game.Field)
	if game.Snake.Length == fieldArea {
		game.Result = "won"
		return
	}

	newHead := &SnakeNode{
		X: game.Snake.Head.X,
		Y: game.Snake.Head.Y,
	}

	switch game.SnakeDirection {
	case "up":
		newHead.Y--
	case "down":
		newHead.Y++
	case "left":
		newHead.X--
	case "right":
		newHead.X++
	}

	if newHead.X < 0 ||
		newHead.X >= len(game.Field) ||
		newHead.Y < 0 ||
		newHead.Y >= len(game.Field) ||
		game.Field[newHead.Y][newHead.X] == 1 {
		game.Result = "lost"
		return
	}

	newHead.Next = game.Snake.Head
	game.Snake.Head.Prev = newHead
	game.Snake.Head = newHead

	if game.Field[newHead.Y][newHead.X] == 2 {
		game.Snake.Length++
		game.AdjustGameSpeed()
		game.SpawnApple()
	} else {
		game.Field[game.Snake.Tail.Y][game.Snake.Tail.X] = 0
		game.Snake.Tail = game.Snake.Tail.Prev
	}

	game.Field[game.Snake.Head.Y][game.Snake.Head.X] = 1
}

func main() {
	config, err := loadGameConfig("cfg.yaml")
	if err != nil {
		fmt.Println("error loading config: ", err, "; trying to get to the right dir")
		execPath, err := os.Executable()
		if err != nil {
			fmt.Println("error pulling executable path: ", err.Error())
			return
		}
		err = os.Chdir(filepath.Dir(execPath))
		if err != nil {
			fmt.Println("error changing directory: ", err.Error())
			return
		}
		config, err = loadGameConfig("cfg.yaml")
		fmt.Println("OK!")
	}
	for {
		output := os.Stdout

		game := NewSnakeGame(output, &config)
		err = initializeGameControls(game)
		if err != nil {
			fmt.Println("error initializing game controls: ", err.Error())
			return
		}

		game.Play()

		deinitializeGameControls()

		clearTerminal()

		err = initializeGameControls(game)
		if err != nil {
			fmt.Println("error initializing game controls: ", err.Error())
			return
		}

		resultMessage := "YOU LOST"
		if game.Result == "won" {
			resultMessage = "YOU WON"
		}

		err = printTitle(fmt.Sprintf(resultMessage + "!! Press Enter to play again"))
		if err != nil {
			fmt.Println("error printing title: ", err.Error())
		}

		for {
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEnter {
				break
			}
		}
		termbox.Close()
	}
}

func initializeGameControls(game *SnakeGame) error {
	if err := termbox.Init(); err != nil {
		return err
	}

	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyArrowUp:
					game.ChangeSnakeDirection("up")
				case termbox.KeyArrowDown:
					game.ChangeSnakeDirection("down")
				case termbox.KeyArrowLeft:
					game.ChangeSnakeDirection("left")
				case termbox.KeyArrowRight:
					game.ChangeSnakeDirection("right")
				case termbox.KeyEsc:
					game.Result = "lost"
				case termbox.KeyCtrlP:
					game.IsPaused = !game.IsPaused
				}
			case termbox.EventError:
				panic(ev.Err)
			}
		}
	}()

	return nil
}

func (game *SnakeGame) ChangeSnakeDirection(newDirection string) {
	if (game.SnakeDirection == "up" && newDirection != "down") ||
		(game.SnakeDirection == "down" && newDirection != "up") ||
		(game.SnakeDirection == "left" && newDirection != "right") ||
		(game.SnakeDirection == "right" && newDirection != "left") {
		game.SnakeDirection = newDirection
	}
}

func deinitializeGameControls() {
	termbox.Close()
}

func printTitle(text string) error {
	CLIWidth, CLIHeight := termbox.Size()

	xOffset := (CLIWidth - len(text)) / 2
	yOffset := CLIHeight / 2

	for i, c := range text {
		termbox.SetCell(xOffset+i, yOffset, c, termbox.ColorDefault, termbox.ColorDefault)
	}

	err := termbox.Flush()
	if err != nil {
		return err
	}

	return nil
}

func loadGameConfig(filename string) (Config, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	configPath := filepath.Join(currentDir, filename)

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	config := &Config{}
	if err = yaml.Unmarshal(configFile, &config); err != nil {
		return Config{}, err
	}

	return *config, nil
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}
