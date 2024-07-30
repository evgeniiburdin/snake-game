# 🐍 Snake Game

This repository contains a terminal-based Snake game written in Go, using the termbox-go library for handling terminal input and output. The game features a configurable field size, customizable graphics, and real-time gameplay speed adjustments based on the snake's length.

## Features

- **Configurable Field Size**: Adjust the size of the playing field through a configuration file.
- **Customizable Graphics**: Define your own symbols and colors for the snake, apple, and field boundaries.
- **Real-time Speed Adjustment**: The game speed changes dynamically based on the length of the snake.
- **Pause and Resume**: Pause the game at any time and resume later.
- **Simple Controls**: Use the arrow keys to navigate the snake.

## Installation

To run this game, you need to have Go installed. You can download it from [the official website](https://golang.org/dl/).

1. Clone the repository:

    ```bash
    git clone https://github.com/evgeniiburdin/snake-game.git
   cd snake-game
    ```

2. Install dependencies:

    ```bash
    go get github.com/nsf/termbox-go
    go get gopkg.in/yaml.v2
    ```

3. Build the game:

    ```bash
    go build -o snake-game main.go
    ```

4. Run the game:

    ```bash
    ./snake
    ```

## 🍎 Configuration

The game can be customized using a `cfg.yaml` file. Below is an example configuration file:

```yaml
left_up_corner: "┌"
up: "───"
right_up_corner: "┐"
left: "│"
right: "│"
left_down_corner: "└"
down: "───"
right_down_corner: "┘"
snake_body: "███"
empty: "   "
apple: "███"

#31-36
snake_color: "36"
apple_color: "31"

#10-50 recommended
field_size: 13
```
snake_color and apple_color use ANSI color codes.
field_size defines the width and height of the playing field.

## Controls

Arrow Keys: Move the snake.

Ctrl + P: Pause the game.

ESC: Exit the game.

## Contribution

Feel free to fork this repository, create feature branches, and submit pull requests. Bug reports and feature requests are welcome.
