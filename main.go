package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	tick         int
	active       bool
	spacePressed bool
	width        int
	height       int
	curGen       byte
	genA         []byte
	genB         []byte
}

func (g *Game) getGen(gen byte) []byte {
	if gen == 0 {
		return g.genA
	} else {
		return g.genB
	}
}

func (g *Game) current() []byte {
	return g.getGen(g.curGen)
}

func (g *Game) get(gen byte, x, y int) bool {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return false
	}

	a := g.width * y + x
	b := a % 8
	c := (a - b) / 8
	return (g.getGen(gen)[c] & (1 << b)) != 0
}

func (g *Game) set(gen byte, x, y int, value bool) {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return
	}

	a := g.width * y + x
	b := a % 8
	c := (a - b) / 8

	if value {
		g.getGen(gen)[c] |= 1 << b
	} else {
		g.getGen(gen)[c] &= ^(1 << b)
	}
}

func (g *Game) neighbourCount(gen byte, x, y int) byte {
	var i byte
	for x2 := -1; x2 <= 1; x2++ {
		for y2 := -1; y2 <= 1; y2++ {
			if (x2 != 0 || y2 != 0) && g.get(gen, x + x2, y + y2) {
				i++
			}
		}
	}

	return i
}

func (g *Game) reset(gen byte) {
	data := g.getGen(gen)
	for i := range data {
		data[i] = 0
	}
}

func (g *Game) doTick() {
	var newGen byte = 0
	if g.curGen == 0 {
		newGen = 1
	}

	g.reset(newGen)

	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			c := g.neighbourCount(g.curGen, x, y)

			if (c == 2 && g.get(g.curGen, x, y)) || c == 3 {
				g.set(newGen, x, y, true)
			}
		}
	}

	g.curGen = newGen
}

func (g *Game) Update(screen *ebiten.Image) error {
	drawMode := 0
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		drawMode = 1
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		drawMode = -1
	}

	if drawMode != 0 {
		x, y := ebiten.CursorPosition()
		g.set(g.curGen, x, y, drawMode > 0)
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.spacePressed = true
	} else if g.spacePressed {
		g.spacePressed = false
		g.active = !g.active
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.reset(g.curGen)
	}

	if g.active && g.tick == 0 {
		g.doTick()
	}

	g.tick++
	if g.tick > 10 {
		g.tick = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var bgColor uint8 = 255
	if !g.active {
		bgColor = 240
	}

	screen.Fill(color.RGBA{bgColor, bgColor, bgColor, 255})

	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if g.get(g.curGen, x, y) {
				screen.Set(x, y, color.Black)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outwideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

func main() {
	g := &Game{
		width:  40,
		height: 30,
		genA:   make([]byte, 40 * 30 / 8),
		genB:   make([]byte, 40 * 30 / 8),
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("golife")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
