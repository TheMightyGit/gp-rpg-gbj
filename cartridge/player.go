package cartridge

import (
	"image"
	"strings"
)

const (
	PlayerStateIdle = iota
	PlayerStateMovingLeft
	PlayerStateMovingRight
	PlayerStateMovingUp
	PlayerStateMovingDown
	PlayerStateSpeaking
	PlayerStateSpeakingNext
	PlayerStateSpeakingSelectWord
	PlayerStateSpeakingGuessWord
	PlayerStateSpeakingIncorrectWord
	PlayerEndGame
)

type Player struct {
	Pos           image.Point
	State         int
	Frame         int
	conversation  *Conversation
	lineIdx       int
	correctAnswer string
	animOffset    int
}

func (p *Player) Say(conversation *Conversation) {
	if !conversation.Completed {
		p.State = PlayerStateSpeaking
		p.conversation = conversation
		p.lineIdx = -1
		p.SayNext()
	}
}

// len of string with escape code stripped.
func lenStripped(t string) int {
	l := len(t)
	l -= strings.Count(t, "\x1bF") * 3
	l -= strings.Count(t, "\x1bB") * 3
	l -= strings.Count(t, "\x1bS") * 2
	return l
}

func (p *Player) getDisplayableText(rawString string) (string, string) {
	correctAnswer := ""
	cookedStrings := []string{"", "", "", ""}
	if strings.Contains(rawString, ":") {
		rawString = strings.Split(rawString, ":")[1]
	}
	row := 0
	for _, word := range strings.Fields(rawString) {
		w := word
		if strings.HasPrefix(w, "_") && strings.HasSuffix(w, "_") {
			correctAnswer = strings.Trim(w, "_")
			if s.InvIdx < 1 || s.InvIdx >= len(Inventory) {
				w = "\x1bF\x0f\x1bB\x05  \x18\x19  \x1bF\x0e\x1bB\x0f"
			} else {
				w = "\x1bF\x0f\x1bB\x05" + Inventory[s.InvIdx] + "\x1bF\x0e\x1bB\x0f"
			}
		} else if strings.HasPrefix(w, "+") && strings.HasSuffix(w, "+") {
			newWord := strings.Trim(w, "+")
			Inventory = append(Inventory, newWord)
			newWordFlashCounter = 60 * 6
			if newWordText == "" {
				newWordText = newWord
			} else {
				newWordText += ", " + newWord
			}
			w = "\x1bF\x05" + newWord + "\x1bF\x0e"
		}
		if strings.HasPrefix(w, "*") {
			w = "\x1bF\x01" + w
		}
		if strings.HasSuffix(w, "*") {
			w = w + "\x1bF\x0e"
		}
		if (len(cookedStrings[row]) + lenStripped(w)) > 24 {
			row++
		}
		if cookedStrings[row] != "" {
			cookedStrings[row] += " "
		}
		// fmt.Println(w)
		cookedStrings[row] += w
	}
	return strings.Join(cookedStrings, "\n"), correctAnswer
}

func (p *Player) SayNext() bool {
	p.lineIdx++
	return p.SayAgain()
}

func (p *Player) SayAgain() bool {
	if p.lineIdx < 0 {
		p.lineIdx = 0
	}
	var displayText string
	if p.lineIdx < len(p.conversation.Words) {
		speakerName := "Doctor"
		if strings.Contains(p.conversation.Words[p.lineIdx], ":") {
			bits := strings.Split(p.conversation.Words[p.lineIdx], ":")
			speakerName = bits[0]
		}
		displayText, p.correctAnswer = p.getDisplayableText(p.conversation.Words[p.lineIdx])
		s.Say(displayText, speakerName)
		p.State = PlayerStateSpeaking
		if p.correctAnswer != "" {
			s.ShowKitBag()
			p.State = PlayerStateSpeakingSelectWord
		} else {
			s.HideKitBag()
		}
		return true
	} else {
		// default to idle, but Then may change that.
		p.State = PlayerStateIdle
		prevConversation := p.conversation
		p.conversation = nil
		p.lineIdx = 0
		if prevConversation.Then != nil {
			prevConversation.Then(prevConversation.Completed)
			if p.conversation != nil {
				// we've jumped into another conversation, so
				// pretend it extends this one otherwise we get
				// a rendering issue.
				return true
			}
		}
		s.HideSpeech()
		s.HideKitBag()
		return false
	}
}

func (p *Player) Start() {
	API.SpritesGet(SpritePlayer).ChangePos(image.Rectangle{p.Pos.Add(viewport.Min), image.Point{16, 16}})
	API.SpritesGet(SpritePlayer).Show(GfxBankTiles, API.MapBanksGet(MapBankMap).GetArea(MapBankAreaPlayer))
	API.SpritesGet(SpritePlayer).ChangeViewport(image.Point{16 * 1, 16 * 1})
}

func (p *Player) guessedCorrectWord() bool {
	return s.InvIdx >= 0 && s.InvIdx < len(Inventory) && Inventory[s.InvIdx] == p.correctAnswer
}

func (p *Player) removeCorrectWordFromInventory() {
	newInventory := []string{}
	for _, word := range Inventory {
		if word != p.correctAnswer {
			newInventory = append(newInventory, word)
		}
	}
	Inventory = newInventory
}

func (p *Player) triggerFlagsIfSteppedOnTrigger() {
	tilePos := image.Point{p.Pos.X / 16, p.Pos.Y / 16}
	triggerId := m.GetTriggerAt(tilePos)
	if triggerId > FLAG_NO_FLAG {
		Flags[triggerId]++
	}
	switch m.GetTriggerAt(tilePos) {
	case FLAG_0_GARB_OFF:
		p.animOffset = 0
	case FLAG_1_GARB_ON:
		if Flags[triggerId] == 1 {
			p.Say(GarbUp)
		}
		p.animOffset = 1
	case FLAG_2_WORD_CRICKET_BALL:
		if Flags[triggerId] == 1 {
			p.Say(WordMagazineCricketBall)
		}
	case FLAG_3_WORD_MOTH:
		if Flags[triggerId] == 1 {
			p.Say(WordRecordsMothAndButterFingers)
		}
	case FLAG_4_WORD_SOMETHING:
		if Flags[triggerId] == 1 {
			p.Say(WordRecordsSomething)
		}
	case FLAG_5_WORD_ANOTHER:
		if Flags[triggerId] == 1 {
			p.Say(WordReceptionPhoneAnother)
		}
	case FLAG_6_TERRY_HOUSE:
		if Flags[triggerId] == 1 {
			p.Say(HouseTerry)
		}
	case FLAG_7_TERI_HOUSE:
		if Flags[triggerId] == 1 {
			p.Say(HouseTeri)
		}
	case FLAG_8_STANDS:
		if Flags[triggerId] == 1 {
			p.Say(Stands)
		}
	case FLAG_9_PUB:
		if Flags[triggerId] == 1 {
			p.Say(Pub)
		}
	case FLAG_10_DEEP_FOREST:
		if Flags[triggerId] == 1 {
			p.Say(DeepForest)
		} else if Flags[triggerId] > 1 {
			p.Say(DeepForestAgain)
		}
	case FLAG_11_PUB_GOING_BEHIND_BAR:
		if Flags[triggerId] == 1 {
			p.Say(PubBehindBar)
		}
	case FLAG_12_TERI_LENDING_BOOK:
		if Flags[triggerId] == 1 {
			p.Say(TeriLendingBook)
		}
	case FLAG_13_TERRY_PICTURE_ON_WALL:
		if Flags[triggerId] == 1 {
			p.Say(TerryPictureOnWall)
		}
	case FLAG_14_STATUE:
		if Flags[triggerId] == 1 {
			p.Say(StatueFound)
		}
	case FLAG_15_LIBRARY:
		if Flags[triggerId] == 1 {
			p.Say(LibraryFound)
		}
	}
}

func (p *Player) Update() {
	if p.State == PlayerStateIdle {
		p.Frame = 0
		if completedFlag && p.Pos.X == 32 && p.Pos.Y == 16 { // completed and next to phone?
			p.Say(EndConversation)
		}
	}

	if p.State >= PlayerStateMovingLeft && p.State <= PlayerStateMovingDown {

		switch p.State {
		case PlayerStateMovingDown:
			p.Pos.Y += 2
		case PlayerStateMovingUp:
			p.Pos.Y -= 2
		case PlayerStateMovingLeft:
			p.Pos.X -= 2
		case PlayerStateMovingRight:
			p.Pos.X += 2
		}

		p.Frame++

		// if we're on an exact multiple of 16 in x and y then go idle.
		if p.Pos.X%16 == 0 && p.Pos.Y%16 == 0 {
			p.State = PlayerStateIdle
			tilePos := image.Point{p.Pos.X / 16, p.Pos.Y / 16}
			m.UpdateBrightnessFrom(tilePos)
			p.triggerFlagsIfSteppedOnTrigger()
		}
	}

	if p.State == PlayerStateSpeakingNext {
		if !p.SayNext() {
			p.Frame = 0
			s.HideSpeech()
			s.HideKitBag()
		}
	}

	if p.State == PlayerStateSpeakingSelectWord {
		p.SayAgain()
	}

	if p.State == PlayerStateSpeakingGuessWord {
		if p.guessedCorrectWord() {
			p.removeCorrectWordFromInventory()
			p.conversation.Completed = true
			p.SayNext()
		} else {
			if Inventory[s.InvIdx] == "<exit>" {
				s.HideKitBag()
				displayText, _ := p.getDisplayableText("Let me think on this and I'll get back to you.")
				s.Say(displayText, "Doctor")
				p.State = PlayerStateSpeakingIncorrectWord
				// exit convo by...
				// forcing Then to run by talking past end of convo
				p.lineIdx = 999999
			} else {
				// have another guess...
				s.HideKitBag()
				displayText, _ := p.getDisplayableText("That doesn't seem right. Let's rethink the diagnosis.")
				s.Say(displayText, "Doctor")
				p.State = PlayerStateSpeakingIncorrectWord
				p.lineIdx-- // go back a convo step (if poss) to regain context
			}
		}
	}

	API.SpritesGet(SpritePlayer).ChangePos(
		image.Rectangle{
			p.Pos.Add(viewport.Min).Sub(c.Pos).Sub(image.Point{Y: 2}),
			image.Point{16, 16},
		},
	)
	API.SpritesGet(SpritePlayer).ChangeViewport(
		image.Point{16 * p.animOffset, 16 * ((p.Frame / 4) % 4)},
	)
}
