package cartridge

import (
	"image"
	"math/rand"

	"github.com/TheMightyGit/marv/marvlib"
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
	marvlib.API.SpritesGet(n.SpriteIdx).ChangePos(image.Rectangle{n.Pos.Add(viewport.Min), image.Point{16, 16}})
	marvlib.API.SpritesGet(n.SpriteIdx).Show(GfxBankTiles, marvlib.API.MapBanksGet(MapBankMap).GetArea(MapBankAreaPlayer))
	marvlib.API.SpritesGet(n.SpriteIdx).ChangeViewport(image.Point{16 * n.AnimColumn, 0})
}

func (n *NPC) Update() {
	marvlib.API.SpritesGet(n.SpriteIdx).ChangePos(
		image.Rectangle{
			n.Pos.Add(viewport.Min).Sub(c.Pos),
			image.Point{16, 16},
		},
	)
	marvlib.API.SpritesGet(n.SpriteIdx).ChangeViewport(image.Point{16 * n.AnimColumn, 16 * (n.Frame % 2)})
	if rand.Intn(60) == 0 {
		n.Frame++
	}
}
