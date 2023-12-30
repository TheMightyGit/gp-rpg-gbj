package cartridge

import (
	"image"
)

func isFloor(tx, ty uint8) bool {
	if ty == 15 { // all triggers are walkable
		return true
	} else if ty >= 6 {
		if tx <= 3 {
			return true
		}
	}
	return false
}

type Map struct {
}

func (m *Map) Start() {
	API.SpritesGet(SpriteMap).ChangePos(viewport)
	API.SpritesGet(SpriteMap).Show(GfxBankTiles, API.MapBanksGet(MapBankMap).GetArea(MapBankAreaMainMap))
}

func (m *Map) Update() {
	API.SpritesGet(SpriteMap).ChangeViewport(
		image.Point{}.Add(c.Pos),
	)
}

func (m *Map) Get(pos image.Point) (uint8, uint8, uint8, uint8) {
	return API.MapBanksGet(MapBankMap).GetArea(MapBankAreaMainMap).Get(pos)
}

func (m *Map) IsFloorAt(pos image.Point) bool {
	tx, ty, _, _ := m.Get(pos)
	return isFloor(tx, ty)
}

func (m *Map) IsNPCAt(pos image.Point) *NPC {
	for i := NPC_0_SUSAN; i < NPC_END; i++ {
		npcPos := image.Point{
			npcs[i].Pos.X / 16,
			npcs[i].Pos.Y / 16,
		}
		if pos == npcPos {
			return npcs[i]
		}
	}
	return nil
}

func (m *Map) GetTriggerAt(pos image.Point) int {
	tx, ty, _, _ := m.Get(pos)
	if ty == 15 {
		return int(tx)
	}
	return FLAG_NO_FLAG
}

func (m *Map) UpdateBrightnessFrom(pos image.Point) {
	intensity := -1

	// intensity is based on floor type
	tx, ty, _, _ := m.Get(pos)

	if tx == 0 && ty == 8 {
		intensity = 0 // max
	}
	if tx == 1 && ty == 9 {
		intensity = 1
	}
	if tx == 2 && ty == 9 {
		intensity = 2
	}
	if tx == 1 && ty == 8 {
		intensity = 3 // min
	}

	if intensity > -1 {
		for i := 0; i < SpriteOverlay; i++ {
			API.SpritesGet(i).ChangePalette(intensity)
		}
	}
}
