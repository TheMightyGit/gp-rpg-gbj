package cartridge

import (
	"image"
	"math/rand"
)

type NPC struct {
	Pos          image.Point
	SpriteIdx    int
	State        int
	Name         string
	AnimColumn   int
	Frame        int
	Conversation *Conversation
}

// TODO: how do we know when the conversation is successful?

func (n *NPC) Start() {
	API.SpritesGet(n.SpriteIdx).ChangePos(image.Rectangle{n.Pos.Add(viewport.Min), image.Point{16, 16}})
	API.SpritesGet(n.SpriteIdx).Show(GfxBankTiles, API.MapBanksGet(MapBankMap).GetArea(MapBankAreaPlayer))
	API.SpritesGet(n.SpriteIdx).ChangeViewport(image.Point{16 * n.AnimColumn, 0})
}

func (n *NPC) Update() {
	API.SpritesGet(n.SpriteIdx).ChangePos(
		image.Rectangle{
			n.Pos.Add(viewport.Min).Sub(c.Pos),
			image.Point{16, 16},
		},
	)
	API.SpritesGet(n.SpriteIdx).ChangeViewport(image.Point{16 * n.AnimColumn, 16 * (n.Frame % 2)})
	if rand.Intn(60) == 0 {
		n.Frame++
	}
}
