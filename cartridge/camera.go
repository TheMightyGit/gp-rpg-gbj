package cartridge

import (
	"image"
)

type Camera struct {
	Pos image.Point
}

func (c *Camera) Start() {
}

// func distance(p1 image.Point, p2 image.Point) float64 {
// 	return math.Sqrt(float64((p2.X-p1.X)*(p2.X-p1.X) + (p2.Y-p1.Y)*(p2.Y-p1.Y)))
// }

func (c *Camera) Update() {

	xSpeed, ySpeed := 1, 1

	camCentre := c.Pos.Add(image.Point{80, 72})

	target := p.Pos
	target = target.Add(image.Point{8, 8})

	// change focus if in conversation so we're
	// no behind text.
	if p.State >= PlayerStateSpeaking && p.State <= PlayerStateSpeakingIncorrectWord {
		target = target.Add(image.Point{0, 32})
	}
	if p.State >= PlayerStateSpeakingSelectWord && p.State <= PlayerStateSpeakingIncorrectWord {
		target = target.Add(image.Point{48, 0})
	}

	/*
		if distance(target, camCentre) > 48 {
			// if one of the directions is only one away then do NOT set speed to 2
			// as it will cause a per frame jitter by overshooting.
			if oneAway := target.Sub(camCentre); !(oneAway.X == 1 || oneAway.X == -1 || oneAway.Y == 1 || oneAway.Y == -1) {
				speed = 2
			}
		}
	*/

	if camCentre.X-target.X > 32 || camCentre.X-target.X < -32 {
		xSpeed = 2
	}
	if camCentre.Y-target.Y > 32 || camCentre.Y-target.Y < -32 {
		ySpeed = 2
	}

	// head towards player
	if camCentre.X < target.X {
		c.Pos.X += xSpeed
	}
	if camCentre.X > target.X {
		c.Pos.X -= xSpeed
	}
	if camCentre.Y < target.Y {
		c.Pos.Y += ySpeed
	}
	if camCentre.Y > target.Y {
		c.Pos.Y -= ySpeed
	}

	// clamps
	if c.Pos.X < 0 {
		c.Pos.X = 0
	}
	if c.Pos.Y < 0 {
		c.Pos.Y = 0
	}
}
