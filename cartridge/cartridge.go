package cartridge

import (
	"embed"
	"image"

	"github.com/TheMightyGit/marv/marvlib"
	"github.com/TheMightyGit/marv/marvtypes"
)

//go:embed "resources/*"
var Resources embed.FS

const (
	GfxBankFont = iota
	GfxBankTiles
	GfxBankOverlay
	GfxBankSmallFont
)
const (
	MapBankMap = iota
)
const (
	MapBankAreaMainMap = iota
	MapBankAreaPlayer
	MapBankAreaIgnore
	MapBankAreaOverlay
	MapBankAreaTitleScreenBG
	MapBankAreaTitleScreenLogo
	MapBankAreaTitleScreenGreenCross
)
const (
	SpriteMap = iota
	SpriteNPC0
	SpriteNPC1
	SpriteNPC2
	SpriteNPC3
	SpriteNPC4
	SpriteNPC5
	SpriteNPC6
	SpriteNPC7
	SpriteNPC8
	SpriteNPC9
	SpriteNPC10
	SpriteNPC11
	SpriteNPC12
	SpriteNPC13
	SpritePlayer
	SpriteSpeech
	SpriteTitleScreenBG
	SpriteTitleScreenGreenCross
	SpriteTitleScreenGreenCrossDetail
	SpriteTitleScreenLogo
	SpriteTitleScreenText
	SpriteOverlay
)

const (
	MAPPED_GAMEPAD_BIT_DPAD_UP    = uint16(1 << 0)
	MAPPED_GAMEPAD_BIT_DPAD_DOWN  = uint16(1 << 1)
	MAPPED_GAMEPAD_BIT_DPAD_LEFT  = uint16(1 << 2)
	MAPPED_GAMEPAD_BIT_DPAD_RIGHT = uint16(1 << 3)

	MAPPED_GAMEPAD_BIT_BUTTON_BOTTOM = uint16(1 << 4)
	MAPPED_GAMEPAD_BIT_BUTTON_LEFT   = uint16(1 << 5)
	MAPPED_GAMEPAD_BIT_BUTTON_TOP    = uint16(1 << 6)
	MAPPED_GAMEPAD_BIT_BUTTON_RIGHT  = uint16(1 << 7)

	MAPPED_GAMEPAD_BIT_SHOULDER_LEFT  = uint16(1 << 8)
	MAPPED_GAMEPAD_BIT_SHOULDER_RIGHT = uint16(1 << 9)

	MAPPED_GAMEPAD_BIT_SELECT = uint16(1 << 10)
	MAPPED_GAMEPAD_BIT_START  = uint16(1 << 11)
)

func isUp(e *InputType) bool {
	return marvlib.API.InputIsKeyDown(marvlib.KeyArrowUp) || e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_DPAD_UP != 0
}
func isDown(e *InputType) bool {
	return marvlib.API.InputIsKeyDown(marvlib.KeyArrowDown) || e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_DPAD_DOWN != 0
}
func isLeft(e *InputType) bool {
	return marvlib.API.InputIsKeyDown(marvlib.KeyArrowLeft) || e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_DPAD_LEFT != 0
}
func isRight(e *InputType) bool {
	return marvlib.API.InputIsKeyDown(marvlib.KeyArrowRight) || e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_DPAD_RIGHT != 0
}
func isButton(e *InputType) bool {
	return marvlib.API.InputIsKeyDown(marvlib.KeySpace) ||
		e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_BUTTON_BOTTOM != 0 ||
		e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_BUTTON_TOP != 0 ||
		e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_BUTTON_LEFT != 0 ||
		e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_BUTTON_RIGHT != 0
}

var suppressControllerRepeat bool

func isJustUp(e *InputType) bool {
	return marvlib.API.InputIsKeyJustPressed(marvlib.KeyArrowUp) || (!suppressControllerRepeat && e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_DPAD_UP != 0)
}
func isJustDown(e *InputType) bool {
	return marvlib.API.InputIsKeyJustPressed(marvlib.KeyArrowDown) || (!suppressControllerRepeat && e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_DPAD_DOWN != 0)
}
func isJustButton(e *InputType) bool {
	// space or any button
	return marvlib.API.InputIsKeyJustPressed(marvlib.KeySpace) ||
		(!suppressControllerRepeat && e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_BUTTON_BOTTOM != 0) ||
		(!suppressControllerRepeat && e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_BUTTON_TOP != 0) ||
		(!suppressControllerRepeat && e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_BUTTON_LEFT != 0) ||
		(!suppressControllerRepeat && e.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_BUTTON_RIGHT != 0)
}

func handleEvent(e *InputType) {

	if e.GamepadButtonStates[0].Mapped == 0 {
		suppressControllerRepeat = false
	}

	if p.State == PlayerStateIdle {
		tilePos := image.Point{p.Pos.X / 16, p.Pos.Y / 16}

		if isUp(e) {
			if npc := m.IsNPCAt(tilePos.Add(image.Point{0, -1})); npc != nil {
				p.Say(npc.Conversation)
			} else if m.IsFloorAt(tilePos.Add(image.Point{0, -1})) {
				p.State = PlayerStateMovingUp
			}
		}
		if isDown(e) {
			if npc := m.IsNPCAt(tilePos.Add(image.Point{0, 1})); npc != nil {
				p.Say(npc.Conversation)
			} else if m.IsFloorAt(tilePos.Add(image.Point{0, 1})) {
				p.State = PlayerStateMovingDown
			}
		}
		if isLeft(e) {
			if npc := m.IsNPCAt(tilePos.Add(image.Point{-1, 0})); npc != nil {
				p.Say(npc.Conversation)
			} else if m.IsFloorAt(tilePos.Add(image.Point{-1, 0})) {
				p.State = PlayerStateMovingLeft
			}
		}
		if isRight(e) {
			if npc := m.IsNPCAt(tilePos.Add(image.Point{1, 0})); npc != nil {
				p.Say(npc.Conversation)
			} else if m.IsFloorAt(tilePos.Add(image.Point{1, 0})) {
				p.State = PlayerStateMovingRight
			}
		}
	}

	if p.State == PlayerStateSpeaking {
		if isJustButton(e) {
			suppressControllerRepeat = true
			p.State = PlayerStateSpeakingNext
		}
	}

	if p.State == PlayerStateSpeakingSelectWord {
		if isJustUp(e) {
			suppressControllerRepeat = true
			if s.InvIdx > 0 {
				s.InvIdx--
			}
		}
		if isJustDown(e) {
			suppressControllerRepeat = true
			if s.InvIdx < len(Inventory)-1 {
				s.InvIdx++
			}
		}
		if isJustButton(e) {
			suppressControllerRepeat = true
			p.State = PlayerStateSpeakingGuessWord
		}
	}

	if p.State == PlayerStateSpeakingIncorrectWord {
		if isJustButton(e) {
			suppressControllerRepeat = true
			p.State = PlayerStateSpeakingSelectWord
		}
	}
}

var (
	viewport = image.Rect(80, 28-4, 160, 144)
	m        *Map
	p        *Player
	c        *Camera
	s        *Speech
	t        *Titles
	npcs     []*NPC
)

const (
	NPC_0_SUSAN = iota
	NPC_1_JUSTIN
	NPC_2_HELEN
	NPC_3_TERRY
	NPC_4_TERI
	NPC_5_PHIL_MCGLASS
	NPC_6_BAR_PATRON_1
	NPC_7_BAR_PATRON_2
	NPC_8_BAR_PATRON_3
	NPC_9_JOSS_TWIN
	NPC_10_JESS_TWIN
	NPC_11_DOG
	NPC_12_JILL
	NPC_13_BAR_PATRON_4
	NPC_END
)

func Start() {
	//marv.SfxBanks[0].PlayLooped()
	marvlib.API.MidBanksGet(0).Play()

	m = &Map{}
	c = &Camera{}
	p = &Player{
		Pos: image.Point{32, 16},
	}
	s = &Speech{}
	t = &Titles{}

	npcs = make([]*NPC, NPC_END-NPC_0_SUSAN)

	npcs[NPC_0_SUSAN] = &NPC{
		Pos:          image.Point{16 * 12, 16 * 18},
		SpriteIdx:    SpriteNPC0,
		Name:         "Susan",
		AnimColumn:   6,
		Conversation: ConversationSusan,
	}
	npcs[NPC_1_JUSTIN] = &NPC{
		Pos:          image.Point{16 * 19, 16 * 10},
		SpriteIdx:    SpriteNPC1,
		Name:         "Justin",
		AnimColumn:   2,
		Conversation: ConversationJustin,
	}
	npcs[NPC_2_HELEN] = &NPC{
		Pos:          image.Point{16 * 16, 16 * 19},
		SpriteIdx:    SpriteNPC2,
		Name:         "Helen",
		AnimColumn:   7,
		Conversation: ConversationHelen,
	}
	npcs[NPC_3_TERRY] = &NPC{
		Pos:          image.Point{16 * 28, 16 * 6},
		SpriteIdx:    SpriteNPC3,
		Name:         "Terry",
		AnimColumn:   9,
		Conversation: ConversationTerry,
	}
	npcs[NPC_4_TERI] = &NPC{
		Pos:          image.Point{16 * 37, 16 * 2},
		SpriteIdx:    SpriteNPC4,
		Name:         "Teri",
		AnimColumn:   6,
		Conversation: ConversationTeri,
	}
	npcs[NPC_5_PHIL_MCGLASS] = &NPC{
		Pos:          image.Point{16 * 33, 16 * 37},
		SpriteIdx:    SpriteNPC5,
		Name:         "Phil McGlass",
		AnimColumn:   2,
		Conversation: ConversationPhil,
	}
	npcs[NPC_6_BAR_PATRON_1] = &NPC{
		Pos:          image.Point{16 * 31, 16 * 30},
		SpriteIdx:    SpriteNPC6,
		Name:         "Laura",
		AnimColumn:   6,
		Conversation: ConversationBarPatronLaura,
	}
	npcs[NPC_7_BAR_PATRON_2] = &NPC{
		Pos:          image.Point{16 * 25, 16 * 36},
		SpriteIdx:    SpriteNPC7,
		Name:         "Jeff",
		AnimColumn:   2,
		Conversation: ConversationBarPatronJeffInBeerGarden,
	}
	npcs[NPC_8_BAR_PATRON_3] = &NPC{
		Pos:          image.Point{16 * 37, 16 * 33},
		SpriteIdx:    SpriteNPC8,
		Name:         "Todd",
		AnimColumn:   5,
		Conversation: ConversationBarPatronTodd,
	}
	npcs[NPC_9_JOSS_TWIN] = &NPC{
		Pos:          image.Point{16 * 58, 16 * 15},
		SpriteIdx:    SpriteNPC9,
		Name:         "Joss",
		AnimColumn:   3,
		Conversation: ConversationJossTwin,
	}
	npcs[NPC_10_JESS_TWIN] = &NPC{
		Pos:          image.Point{16 * 60, 16 * 16},
		SpriteIdx:    SpriteNPC10,
		Name:         "Jess",
		AnimColumn:   4,
		Conversation: ConversationJessTwin,
	}
	npcs[NPC_11_DOG] = &NPC{
		Pos:          image.Point{16 * 29, 16 * 9},
		SpriteIdx:    SpriteNPC11,
		Name:         "Dog",
		AnimColumn:   8,
		Conversation: ConversationDog,
	}
	npcs[NPC_12_JILL] = &NPC{
		Pos:          image.Point{16 * 13, 16 * 21},
		SpriteIdx:    SpriteNPC12,
		Name:         "Jill",
		AnimColumn:   6,
		Conversation: ConversationJill,
	}
	npcs[NPC_13_BAR_PATRON_4] = &NPC{
		Pos:          image.Point{16 * 24, 16 * 33},
		SpriteIdx:    SpriteNPC13,
		Name:         "Tristan",
		AnimColumn:   5,
		Conversation: ConversationBarPatronTristanInBeerGarden,
	}

	c.Start()
	m.Start()
	p.Start()
	s.Start()
	for i := NPC_0_SUSAN; i < NPC_END; i++ {
		npcs[i].Start()
	}

	marvlib.API.SpritesGet(SpriteOverlay).ChangePos(image.Rectangle{
		image.Point{0, 0},
		image.Point{320, 200},
	})
	marvlib.API.SpritesGet(SpriteOverlay).Show(
		GfxBankOverlay,
		marvlib.API.MapBanksGet(MapBankMap).GetArea(MapBankAreaOverlay),
	)
	//	marv.Sprites[SpriteOverlay].ChangeViewport(
	//		image.Point{256 - 32, 28},
	//	)

	p.Say(IntroText)

	t.Start()
	// marv.ModBanks[0].Play()
}

var onTitleScreen = true

func Update() {
	if onTitleScreen {
		t.Update()
	} else {
		e := &InputType{
			MousePos:        marvlib.API.InputMousePos(),
			MousePosDelta:   marvlib.API.InputMousePosDelta(),
			MouseWheelDelta: marvlib.API.InputMouseWheelDelta(),
			MousePressed:    marvlib.API.InputMousePressed(),
			MouseHeld:       marvlib.API.InputMouseHeld(),
			MouseReleased:   marvlib.API.InputMouseReleased(),
			InputChars:      marvlib.API.InputChars(),
			GamepadButtonStates: [4]marvtypes.GamepadState{
				{Unmapped: 0, Mapped: 0},
				{Unmapped: 0, Mapped: 0},
				{Unmapped: 0, Mapped: 0},
				{Unmapped: 0, Mapped: 0},
			},
		}
		handleEvent(e)
		c.Update()
		m.Update()
		p.Update()
		for i := NPC_0_SUSAN; i < NPC_END; i++ {
			npcs[i].Update()
		}
		s.Update()
	}
}
