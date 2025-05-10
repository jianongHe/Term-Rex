package game

import "github.com/nsf/termbox-go"

// Dino represents the player character
type Dino struct {
	X, Y       int
	jumping    bool
	jumpFrames int
}

func NewDino() *Dino {
	return &Dino{X: 2, Y: height - 2}
}

func (d *Dino) Update() {
	if d.jumping {
		d.Y = height - 3
		d.jumpFrames--
		if d.jumpFrames <= 0 {
			d.jumping = false
			d.Y = height - 2
		}
	}
}

func (d *Dino) Jump() {
	if !d.jumping {
		d.jumping = true
		d.jumpFrames = fps / 4
	}
}

var dinoSprite = Sprite{
	"  _  ",
	" /_\\ ",
	"/___\\",
}

func (d *Dino) Draw() {
	h := len(dinoSprite)
	startY := d.Y - (h - 1)
	dinoSprite.Draw(d.X, startY, termbox.ColorGreen, termbox.ColorDefault)
}
