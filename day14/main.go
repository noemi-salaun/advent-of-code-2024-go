package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strconv"
)

// var height = 7
// var width = 11
var height = 103
var width = 101

func main() {
	var bots, err = loadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	for s := range 10000 {
		for _, b := range *bots {
			b.move(1)
		}

		drawBots(s+1, bots)
		fmt.Println(s + 1)
	}
}

func drawPng(sec int, data *[][]bool) {

	// Créer une nouvelle image avec les dimensions du tableau
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Remplir l'image avec les valeurs du tableau
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Transformer la valeur du tableau en couleur
			var col color.Color
			if (*data)[y][x] {
				col = color.RGBA{0, 0, 0, 255} // Rouge
			} else {
				col = color.RGBA{255, 255, 255, 255} // Noir
			}

			// Définir la couleur du pixel
			img.Set(x, y, col)
		}
	}

	// Créer un fichier PNG
	outputFile, err := os.Create(fmt.Sprintf("day14/output/%d.png", sec))
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Encoder l'image en PNG et l'écrire dans le fichier
	err = png.Encode(outputFile, img)
	if err != nil {
		panic(err)
	}
}

var colors = []string{
	"\033[41m \033[0m",  // Rouge
	"\033[42m \033[0m",  // Vert
	"\033[43m \033[0m",  // Jaune
	"\033[44m \033[0m",  // Bleu
	"\033[45m \033[0m",  // Magenta
	"\033[46m \033[0m",  // Cyan
	"\033[47m \033[0m",  // Blanc
	"\033[100m \033[0m", // Noir brillant
	"\033[101m \033[0m", // Rouge brillant
	"\033[102m \033[0m", // Vert brillant
	"\033[103m \033[0m", // Jaune brillant
	"\033[104m \033[0m", // Bleu brillant
	"\033[105m \033[0m", // Magenta brillant
	"\033[106m \033[0m", // Cyan brillant
	"\033[107m \033[0m", // Blanc brillant
}

func printColor(n int) {
	val := n % len(colors)
	fmt.Print(colors[val])
}

func printBots(bots *[]*bot) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var found = false
			var foundI int
			for i, b := range *bots {
				if b.position.x == x && b.position.y == y {
					found = true
					foundI = i
					break
				}
			}
			if found {
				printColor(foundI)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func drawBots(sec int, bots *[]*bot) {
	var grid [][]bool
	for y := 0; y < height; y++ {
		var line []bool
		for x := 0; x < width; x++ {
			var found = false
			for _, b := range *bots {
				if b.position.x == x && b.position.y == y {
					found = true
					break
				}
			}
			line = append(line, found)
		}
		grid = append(grid, line)
	}

	drawPng(sec, &grid)
}

type quadrant struct {
	min      vector2
	max      vector2
	botCount int
}

type vector2 struct {
	x, y int
}

type bot struct {
	position vector2
	velocity vector2
}

func (b *bot) isInQuadrant(q *quadrant) bool {
	return b.position.x >= q.min.x && b.position.x <= q.max.x && b.position.y >= q.min.y && b.position.y <= q.max.y
}

func (b *bot) move(seconds int) {
	b.position.x = (b.position.x + b.velocity.x*seconds) % width
	if b.position.x < 0 {
		b.position.x += width
	}

	b.position.y = (b.position.y + b.velocity.y*seconds) % height
	if b.position.y < 0 {
		b.position.y += height
	}
}

func newBot(line string) *bot {
	// p=0,4 v=3,-3
	var re = regexp.MustCompile(`^p=(\d+),(\d+) v=(-?\d+),(-?\d+)$`)

	var newBot bot

	matches := re.FindAllStringSubmatch(line, -1)
	newBot.position.x, _ = strconv.Atoi(matches[0][1])
	newBot.position.y, _ = strconv.Atoi(matches[0][2])
	newBot.velocity.x, _ = strconv.Atoi(matches[0][3])
	newBot.velocity.y, _ = strconv.Atoi(matches[0][4])

	return &newBot
}

func loadInput() (*[]*bot, error) {
	readFile, err := os.Open("day14/input.txt")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	var bots []*bot

	for fileScanner.Scan() {
		line := fileScanner.Text()

		bots = append(bots, newBot(line))
	}

	return &bots, nil
}

func displayFrame(data [][]int) {
	// Effacer l'écran
	fmt.Print("\033[H\033[2J")

	for _, row := range data {
		for _, cell := range row {
			// Afficher un caractère ou une couleur en fonction de la valeur
			switch cell {
			case 0:
				fmt.Print("\033[41m \033[0m") // Rouge
			case 1:
				fmt.Print("\033[42m \033[0m") // Vert
			case 2:
				fmt.Print("\033[44m \033[0m") // Bleu
			default:
				fmt.Print(" ") // Espace
			}
		}
		fmt.Println()
	}
}
