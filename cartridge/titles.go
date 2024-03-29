package cartridge

import (
	"image"
	"math"
	"math/rand"

	"github.com/TheMightyGit/marv/marvtypes"
)

type Titles struct {
	textArea marvtypes.MapBankArea
}

var (
	textFrames = [][5]string{
		[5]string{
			"                                        ",
			"              @TheMightyGit             ",
			"           ITCH.IO/JAM/GBJAM-9          ",
			"                                        ",
			"        \x1bF\x05<SPACE>\x1bF\x01/\x1bF\x05<BUTTON>\x1bF\x01 to start       ",
		},
		[5]string{
			"                                        ",
			"    Code, Art, Writing: TheMightyGit    ",
			"                                        ",
			"                                        ",
			"        \x1bF\x05<SPACE>\x1bF\x01/\x1bF\x05<BUTTON>\x1bF\x01 to start       ",
		},
		[5]string{
			"                                        ",
			"           Palette: Polyducks           ",
			"                                        ",
			"                                        ",
			"        \x1bF\x05<SPACE>\x1bF\x01/\x1bF\x05<BUTTON>\x1bF\x01 to start       ",
		},
		[5]string{
			"                                        ",
			"  In-Game font based on 'Kitchen Sink'  ",
			"       by Retroshark & Polyducks        ",
			"                                        ",
			"        \x1bF\x05<SPACE>\x1bF\x01/\x1bF\x05<BUTTON>\x1bF\x01 to start       ",
		},
		//		[5]string{
		//			"                                        ",
		//			"Music: 'Somewhere in the Elevator'      ",
		//			"       by Peachtea@You're Perfect Studio",
		//			"                                        ",
		//			"        \x1bF\x05<SPACE>\x1bF\x01/\x1bF\x05<BUTTON>\x1bF\x01 to start       ",
		//		},
		[5]string{
			"     Special thanks to playtesters:     ",
			"      Pipsissiwa, Dunk, Joe, Mikey      ",
			"    ( and anyone I've forgotten <3 )    ",
			"                                        ",
			"        \x1bF\x05<SPACE>\x1bF\x01/\x1bF\x05<BUTTON>\x1bF\x01 to start       ",
		},
		[5]string{
			"And for the Brighton Indie Coffee group ",
			"   for putting up with my constant dev  ",
			"         spam during this jam <3        ",
			"                                        ",
			"        \x1bF\x05<SPACE>\x1bF\x01/\x1bF\x05<BUTTON>\x1bF\x01 to start       ",
		},
	}
)

var (
	logoViewport             image.Rectangle
	greencrossViewport       image.Rectangle
	greencrossDetailViewport image.Rectangle
	f                        image.Point
)

func (t *Titles) Start() {

	API.SpritesGet(SpriteTitleScreenBG).ChangePos(viewport)
	API.SpritesGet(SpriteTitleScreenBG).Show(
		GfxBankTiles,
		API.MapBanksGet(MapBankMap).GetArea(MapBankAreaTitleScreenBG),
	)

	logoViewport = viewport
	logoViewport.Min.X += (16 * 1.5)
	logoViewport.Min.Y += (16 * 1)
	logoViewport.Max.X = 16 * 8
	logoViewport.Max.Y = 16 * 5
	API.SpritesGet(SpriteTitleScreenLogo).ChangePos(logoViewport)
	API.SpritesGet(SpriteTitleScreenLogo).Show(
		GfxBankTiles,
		API.MapBanksGet(MapBankMap).GetArea(MapBankAreaTitleScreenLogo),
	)

	greencrossViewport = viewport
	greencrossViewport.Max.X = 16 * 3
	greencrossViewport.Max.Y = 16 * 3
	API.SpritesGet(SpriteTitleScreenGreenCross).ChangePos(greencrossViewport)
	API.SpritesGet(SpriteTitleScreenGreenCross).Show(
		GfxBankTiles,
		API.MapBanksGet(MapBankMap).GetArea(MapBankAreaTitleScreenGreenCross),
	)
	greencrossViewport.Min.X = -2000

	greencrossDetailViewport = viewport
	greencrossDetailViewport.Max.X = 16
	greencrossDetailViewport.Max.Y = 16
	API.SpritesGet(SpriteTitleScreenGreenCrossDetail).ChangePos(greencrossDetailViewport)
	API.SpritesGet(SpriteTitleScreenGreenCrossDetail).Show(
		GfxBankTiles,
		API.MapBanksGet(MapBankMap).GetArea(MapBankAreaPlayer),
	)
	greencrossDetailViewport.Min.X = -2000

	t.textArea = API.MapBanksGet(MapBankMap).AllocArea(image.Point{40, 5})
	t.textArea.ClearWithColour(0, 0, 14, 16)
	v := viewport
	v.Min.Y += 114
	v.Max.Y = 5 * 6
	API.SpritesGet(SpriteTitleScreenText).ChangePos(v)
	API.SpritesGet(SpriteTitleScreenText).Show(
		GfxBankSmallFont,
		t.textArea,
	)
}

var (
	textIdx    = 0
	frameCount = 0
)

func (t *Titles) updateText(c uint8) {
	pos := image.Point{}
	for _, txt := range textFrames[textIdx%len(textFrames)] {
		posDelta := t.textArea.StringToMap(
			pos,
			c, 15,
			txt+"\n",
		)
		pos = pos.Add(posDelta)
	}
}

func (t *Titles) Update() {
	e := &InputType{
		MousePos:        API.InputMousePos(),
		MousePosDelta:   API.InputMousePosDelta(),
		MouseWheelDelta: API.InputMouseWheelDelta(),
		MousePressed:    API.InputMousePressed(),
		MouseHeld:       API.InputMouseHeld(),
		MouseReleased:   API.InputMouseReleased(),
		InputChars:      API.InputChars(),
		GamepadButtonStates: [4]marvtypes.GamepadState{
			{Unmapped: 0, Mapped: 0},
			{Unmapped: 0, Mapped: 0},
			{Unmapped: 0, Mapped: 0},
			{Unmapped: 0, Mapped: 0},
		},
	}

	if isButton(e) {
		API.SpritesGet(SpriteTitleScreenBG).SetEnabled(false)
		API.SpritesGet(SpriteTitleScreenLogo).SetEnabled(false)
		API.SpritesGet(SpriteTitleScreenText).SetEnabled(false)
		API.SpritesGet(SpriteTitleScreenGreenCross).SetEnabled(false)
		API.SpritesGet(SpriteTitleScreenGreenCrossDetail).SetEnabled(false)
		onTitleScreen = false
		suppressControllerRepeat = true
	}

	v := logoViewport
	v.Min.Y += int(math.Sin(float64(frameCount)/20.0) * 16)
	API.SpritesGet(SpriteTitleScreenLogo).ChangePos(v)

	greencrossViewport.Min.X -= 1
	greencrossDetailViewport.Min.X -= 1
	if greencrossViewport.Min.X < (viewport.Min.X - (16 * 3)) {
		greencrossViewport.Min.X = viewport.Min.X + 160 + (16 * 3)
		greencrossViewport.Min.Y = viewport.Min.Y + rand.Intn(112-48)
		greencrossDetailViewport.Min.Y = greencrossViewport.Min.Y + 16
		greencrossDetailViewport.Min.X = greencrossViewport.Min.X + 16
		f = image.Point{
			16 * rand.Intn(10),
			16 * rand.Intn(2),
		}
	}
	API.SpritesGet(SpriteTitleScreenGreenCross).ChangePos(greencrossViewport)
	API.SpritesGet(SpriteTitleScreenGreenCrossDetail).ChangePos(greencrossDetailViewport)
	if rand.Intn(20) == 0 {
		f.Y = 16 * rand.Intn(2)
	}
	API.SpritesGet(SpriteTitleScreenGreenCrossDetail).ChangeViewport(f)

	if frameCount%(60*4) == 0 {
		t.updateText(1)
	}
	if frameCount%(60*4) == 4 {
		t.updateText(5)
	}
	if frameCount%(60*4) == 8 {
		t.updateText(14)
	}
	if frameCount%(60*4) == (60*4)-18 {
		t.updateText(5)
	}
	if frameCount%(60*4) == (60*4)-14 {
		t.updateText(1)
	}
	if frameCount%(60*4) == (60*4)-11 {
		t.updateText(15)
		textIdx++
	}
	frameCount++
}
