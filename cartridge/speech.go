package cartridge

import (
	"fmt"
	"image"
	"math/rand"
	"strings"

	"github.com/TheMightyGit/marv/marvtypes"
)

var speechBox = &[9]image.Point{
	image.Point{28, 0},
	image.Point{29, 0},
	image.Point{30, 0},
	image.Point{28, 1},
	image.Point{29, 1},
	image.Point{30, 1},
	image.Point{28, 2},
	image.Point{29, 2},
	image.Point{30, 2},
}

type Speech struct {
	area            marvtypes.MapBankArea
	topArea         marvtypes.MapBankArea
	topTextArea     marvtypes.MapBankArea
	bottomArea      marvtypes.MapBankArea
	bottomTextArea  marvtypes.MapBankArea
	scoreTextArea   marvtypes.MapBankArea
	newWordTextArea marvtypes.MapBankArea
	InvIdx          int
}

func (s *Speech) Start() {
	s.area = API.MapBanksGet(MapBankMap).AllocArea(image.Point{27, 18})
	// the 6x8 grid aligns better offset by -1
	v := viewport
	v.Min.X--
	v.Max.X++
	API.SpritesGet(SpriteSpeech).ChangePos(v)
	API.SpritesGet(SpriteSpeech).Show(GfxBankFont, s.area)

	s.topArea = s.area.GetSubArea(image.Rect(10, 0, 27, 12))
	s.topTextArea = s.topArea.GetSubArea(image.Rect(1, 1, 16, 11))

	s.bottomArea = s.area.GetSubArea(image.Rect(0, 12, 27, 18))
	s.bottomTextArea = s.bottomArea.GetSubArea(image.Rect(1, 1, 26, 5))

	s.scoreTextArea = s.area.GetSubArea(image.Rect(0, 17, 27, 18))
	s.newWordTextArea = s.area.GetSubArea(image.Rect(0, 13, 27, 15))

	s.area.ClearWithColour(0, 0, 16, 16)
}

func (s *Speech) ShowSpeech(speakerName string) {
	s.bottomArea.ClearWithColour(1, 1, 0, 1)
	s.bottomArea.ClearWithColour(1, 1, 0, 1)
	s.bottomArea.DrawBox(
		image.Rectangle{image.Point{0, 0}, image.Point{26, 5}},
		speechBox, 15, 16,
	)
	s.bottomTextArea.ClearWithColour(0, 0, 15, 14)
	s.SetSpeakerName(speakerName)
}

func (s *Speech) SetSpeakerName(speakerName string) {
	xPos := 2
	if speakerName != "Doctor" && speakerName != "Dad" {
		xPos = 21 - len(speakerName)
	}
	s.bottomArea.StringToMap(
		image.Point{xPos, 0},
		15, 16,
		"\x1bC\x20\x01 "+speakerName+" \x1bC\x21\x01",
	)
	s.bottomArea.StringToMap(
		image.Point{xPos + 1, 0},
		5, 15,
		" "+speakerName+" ",
	)
}

func (s *Speech) ShowKitBag() {
	s.topArea.ClearWithColour(1, 1, 0, 1)
	s.topArea.ClearWithColour(1, 1, 0, 1)
	s.topArea.DrawBox(
		image.Rectangle{image.Point{0, 0}, image.Point{16, 11}},
		speechBox, 11, 16,
	)
	s.topTextArea.ClearWithColour(0, 0, 14, 15)
	s.topTextArea.StringToMap(
		image.Point{},
		14, 11,
		"\x1bF\x02Lexical Kit Bag",
	)
	pos := image.Point{0, 1}
	var (
		cFg uint8
		cBg uint8
	)

	for idx, invItem := range Inventory {

		if idx < s.InvIdx-8 {
			continue
		}

		cFg, cBg = 14, 11
		if idx == s.InvIdx {
			// cFg, cBg = cBg, cFg
			cFg, cBg = 15, 5
		}
		posDelta := s.topTextArea.StringToMap(
			pos,
			cFg, cBg,
			invItem+"\n",
		)
		pos = pos.Add(posDelta)
	}

	if s.InvIdx > 8 {
		s.topTextArea.Set(image.Point{14, 1}, 8, 1, 14, 15)
	}
	if len(Inventory) > 9 && s.InvIdx < len(Inventory)-1 {
		s.topTextArea.Set(image.Point{14, 9}, 9, 1, 14, 15)
	}
}

func (s *Speech) Say(txt string, speakerName string) {
	s.ShowSpeech(speakerName)
	var c uint8
	if strings.HasPrefix(txt, "PHONE:") {
		c = 2
	} else {
		c = 15
	}
	s.bottomTextArea.ClearWithColour(0, 0, 14, 15)
	s.bottomTextArea.StringToMap(
		image.Point{},
		14, c,
		txt,
	)
	// 'A button' symbol.
	s.bottomTextArea.Set(image.Point{24, 3}, 4, 1, 14, 15)
}

func (s *Speech) Update() {
	if curedFlashCounter > 0 && (p.State >= PlayerStateIdle && p.State <= PlayerStateMovingDown) {
		curedFlashCounter--
		s.UpdateScore()
	}
	if newWordFlashCounter > 0 && (p.State >= PlayerStateIdle && p.State <= PlayerStateMovingDown) {
		newWordFlashCounter--
		s.UpdateNewWordSplash()
	}
}

func (s *Speech) HideSpeech() {
	s.bottomArea.ClearWithColour(0, 0, 16, 16)
	s.UpdateScore()
}

var (
	NumCures = 0
	MaxCures = 11
)

var curedFlashCounter = 0
var (
	newWordFlashCounter = 0
	newWordText         = ""
)

func (s *Speech) CurePatient() {
	NumCures++
	curedFlashCounter = 60 * 3
	s.UpdateScore()
}

func center(t string) string {
	l := (27 - len(t)) / 2
	if l < 0 {
		l = 0
	}
	pad := strings.Repeat(" ", l)
	return pad + t + pad + " "
}

func (s *Speech) UpdateNewWordSplash() {
	if newWordFlashCounter > 0 {
		newWordFlashCounter--
		if newWordFlashCounter > 0 {
			var fg, bg uint8 = 14, 15
			txt := center("Words added to Kit Bag!") + "\n\x1bF\x05" + center(newWordText)
			s.newWordTextArea.StringToMap(image.Point{}, fg, bg, txt)
		} else {
			s.newWordTextArea.ClearWithColour(0, 0, 16, 16)
			newWordText = ""
		}
	}
}

var completedFlag bool

func (s *Speech) UpdateScore() {
	if NumCures == MaxCures && curedFlashCounter == 0 {
		if !completedFlag {
			p.Say(EveryoneCured)
			completedFlag = true
		}
	} else {
		txt := ""
		var fg, bg uint8 = 14, 16
		if curedFlashCounter > 0 {
			txt = fmt.Sprintf("        CURED:%d/%d +1          ", NumCures-1, MaxCures)
			if curedFlashCounter%20 < 10 {
				fg, bg = 15, 14
			}
		} else {
			txt = fmt.Sprintf("        CURED:%d/%d          ", NumCures, MaxCures)
		}
		s.scoreTextArea.StringToMap(image.Point{}, fg, bg, txt)
	}
}

func (s *Speech) HideKitBag() {
	s.topArea.ClearWithColour(0, 0, 16, 16)
	s.InvIdx = 0
}

type Conversation struct {
	Completed bool
	Words     []string
	Then      func(bool)
}

var Inventory = []string{
	"<exit>",
	"pencil",
	"cream",
	"ICU",
}

var GarbUp = &Conversation{
	Words: []string{
		"I better garb up as I'm heading outside.\n*dons hat* *picks up bag*",
		"It looks like cricket practice is underway over in the park.",
	},
}

var WordMagazineCricketBall = &Conversation{
	Words: []string{
		"This waiting room has a great magazine on the humble +cricket_ball+ and...",
		"advanced techniques for catching it.",
		"You never know when a cricket ball fact will get you out of a tight spot.",
	},
}

var WordRecordsMothAndButterFingers = &Conversation{
	Words: []string{
		"Oh, this filing cabinet is open. Well, it can't hurt to +bone+ up on some files.",
		"*shuffles files randomly*",
		"Ahh, yes, this one. Justin. Says here he's a big fan of the +moth+",
		"I'm not sure how that's medically relevant, but a doctor is always prepared.",
		"*shuffles files randomly*",
		"Ooops! I dropped them!", // MOVED: " I'm a real +butter_fingers+",
		"I best stop messing with these files.",
	},
}

var WordRecordsSomething = &Conversation{
	Words: []string{
		"This cabinet is filled with pamphlets for patients...",
		"\"Sneezing fits - do's and don'ts\" - one of our most popular pamphlets.",
		"\"Glue Teen Intolerance\" - Oh, I remember this ill-fated anti-drug campaign.",
		"\"10,000 ways to reach peak fitness!\" - this is a +huge_pamphlet+ *puff*",
		"I can't imagine any +patient+ was really helped by these.",
		"Anyway, I better put this all back and get on.",
	},
}

var WordReceptionPhoneAnother = &Conversation{
	Words: []string{
		"Reception Phone: You have 2 messages.",
		"Reception Phone: *beep*",
		"Reception Phone: Message 1...",
		"Reception Phone: \"Don't forget about the upcoming charity egg & +sausage+ race!\"",
		"That doesn't sound...",
		"Reception Phone: *beep*",
		"Reception Phone: Message 2...",
		"Reception Phone: \"Sorry, I meant charity egg & +spoon+ race!\"",
		"Reception Phone: \"Egg and SPOON. For charity.\"",
		"Reception Phone: \"Sorry, I was thinking about breakfast.\"",
		"Reception Phone: *beep*",
		"Reception Phone: You have no more messages.",
	},
}

var HouseTerry = &Conversation{
	Words: []string{
		"Ah, Terry's house. Perfect, I'll pop in.",
		"Interestingly, his twin lives next door.",
	},
}

var HouseTeri = &Conversation{
	Words: []string{
		"Ah, Teri's home. Great, I'll nip in.",
		"Teri both has twins and is a twin. Her twin lives next door.",
		"The kids will be around here somewhere.",
	},
}

var Stands = &Conversation{
	Words: []string{
		"These stands for the cricket ground were built when the floodlighting was fitted.",
		"Money well spent. Our team is the envy of the general area.",
		"General area supremacy is something of a local obsession.",
		"Looks like +graffiti+ is a problem, though...",
		"There are a lot of fresh looking 'Joss' and 'Jess' tags.",
		// "...there's some somewhat adult graffiti boasting about someone's bedroom prowess.",
		// "But also some additional graffiti leaving a 1 star review of said prowess.",
		// "These are not words that will likely be diagnostically important. I hope.",
	},
}

var Pub = &Conversation{
	Words: []string{
		"The Pixellated Arms pub. It has a great beer garden.",
		"The best pub in town...",
		"...the only pub in town.",
		// "It serves both real and faux ale, and something akin to food.",
	},
}

var DeepForest = &Conversation{
	Words: []string{
		"They call this 'The Last Forest'.",
		"Apparently no-one who enters the forest ever leaves impressed.",
	},
}
var DeepForestAgain = &Conversation{
	Words: []string{
		"I remain unimpressed.",
	},
}
var PubBehindBar = &Conversation{
	Words: []string{
		"I don't like going behind here without permission, but I need to check on Phil.",
	},
}
var TeriLendingBook = &Conversation{
	Words: []string{
		"Ah-ha, here's the book I wanted to read.",
		//"\"Plucking & Plectrums - The +guitar+ through the ages.\"",
		"Big Boat Bane - An +iceberg+ story",
		"I look forward to reading this adventure book when this pundemic is over.",
		// "A book of +shorter+ stories, for people who don't have time for short stories.",
	},
}
var TerryPictureOnWall = &Conversation{
	Words: []string{
		"It's a picture of Terry and Teri at the seaside as kids.",
		"Or Teri and Terry.",
		"I'm not entirely sure.",
		"That +shiny_jetski+ they're on looks fast and +furious+ , tho.",
	},
}
var LibraryFound = &Conversation{
	Words: []string{
		"Ah, the closed down library. A real shame.",
		"It had the slipperiest floor for +sock+ sliding and was much admired and attended.",
		"So much so that the kids didn't read any books, so it was forced to close.",
		"Even so, kids sometimes still hang out around the back.",
	},
}
var StatueFound = &Conversation{
	Words: []string{
		"The statue of the town founder.",
		"Or at least the plinth.",
		"The fund raising is ongoing so it should eventually be finished.",
		"It's tough to get enthusiasm to raise +money+ when...",
		"...no-one knows who founded the town.",
	},
}

// TODO:
// - X catching (CANNOT USE AS IT'S A 2ND WORD IN A CONVO!)
// - X butter_fingers

var ConversationSusan = &Conversation{
	Words: []string{
		"Susan: Doctor, doctor! I really hope you can help...",
		"How's that?",
		"Susan: Oh, don't you start as well!",
		"...",
		"Susan: ...I've got a _cricket_ball_ stuck in my bottom!",
		"That advanced catching techniques magazine has a lot to answer for.",
		"This would have stumped a lesser doctor...",
		"But I think we've caught this just in time...",
		"I declare I can get it all out...",
		"But you, uh...",
		"Susan: Yes?",
		"...may have the runs for a while.",
	},
	Then: func(success bool) {
		if success {
			s.CurePatient()
		}
	},
}

var ConversationJustin = &Conversation{
	Words: []string{
		"*You watch as Justin keeps bumping into the window.*",
		"*He finally notices you're there*",
		"Justin: Doctor, doctor! I can't help it, I just keep thinking I'm a _moth_",
		"I think you need a psychiatrist not a doctor.",
		"Justin: I know, but I was walking past and I saw your light was on...",
		"This pundemic is far worse than I thought.",
	},
	Then: func(success bool) {
		if success {
			s.CurePatient()
		}
	},
}

var ConversationHelen = &Conversation{
	Words: []string{
		"Helen: Doctor, doctor! They've dropped me from the cricket team - they call me _butter_fingers_",
		"Don't worry, what you have is not catching.", // FIXME: multiple guesses per convo is not supported!
		"Helen: Phew. Now all I have to do is find my mobile phone.",
		"Helen: I put it over on the stands but now I can't find it.",
		"Helen: The twins were playing around by the stands earlier...",
		"Helen: ...maybe they saw where I left my phone.",
	},
	Then: func(success bool) {
		if success {
			s.CurePatient()
		}
	},
}

var ConversationTeri = &Conversation{
	Words: []string{
		"Teri: Hi, I'm just making the pastry for some baking. The twins never seem to stop eating.",
		"...",
		"Teri: What was that pause for?",
		"Believe it or not I was waiting for a medically significant pun.",
		"Teri: Well... I'm baking a +disney+ princess and prince pie...",
		"Teri: ...but it's vegetarian, so no actual royalty in it.",
		"You can't trick my medical training, that was a regular pun.",
		"Teri: Haha, you're too good.",
		"Oh, the book you wanted to read is on the bedroom bookcase.",
		"Teri: I can't get it as I've got +butter_fingers+ right now.",
		"Could you get the book from the bedroom yourself?",
	},
}

func init() {
	ConversationTeri.Then = func(success bool) {
		// NOTE: we need this because we fill in no word during this converstion.
		ConversationTeri.Completed = true
	}
}

var ConversationTerry = &Conversation{
	Words: []string{
		"Terry: Doc, every time I drink a cup of hot chocolate I get a stabbing pain in the eye.",
		"Try taking the _spoon_ out first.",
		"Terry: What?",
		"Oh. That wasn't a medically significant pun?",
		"Terry: No! I actually do get pain.",
		"In that case it's probably a cracked tooth affecting your sinus. Best see a dentist.",
		"Terry: Thanks, Doc. I must admit though, my dog has been acting funny recently.",
	},
	Then: func(success bool) {
		if success {
			s.CurePatient()
		}
	},
}

var ConversationJossTwin = &Conversation{
	Words: []string{
		"Joss: It was Jess!",
		"You're not in trouble.",
		"Joss: Oh.",
		"What are you doing here?",
		"Joss: We thought we saw a medium to large long-winged bird of prey which typically has...",
		"Joss: ...a forked tail and frequently soars on updraughts of air.",
		"My birdwatching days are over, but that sure sounds like some kind of +kite+",
		"Joss: Yep.",
	},
}

func init() {
	ConversationJossTwin.Then = func(success bool) {
		// NOTE: we need this because we fill in no word during this converstion.
		ConversationJossTwin.Completed = true
	}
}

var ConversationJessTwin = &Conversation{
	Words: []string{
		"Jess: It was Joss!",
		"You're not in trouble.",
		"Jess: Ah.",
		"What are you doing here?",
		"Jess: Kicking an old cabinet. Want a go?",
		"*looks around*",
		"Yes.",
		"*you give it a good kick*",
		"Cabinet: *KER-CHANG-PLONK*",
		"Jess: YAY!",
		"Joss: I concur!",
		"That felt good, maybe just one more for the road...",
		"*you give it another good kick*",
		"Cabinet: *KER-CHANG-PLONK*",
		"Jess: YAY!",
		"Joss: I concur!",
	},
	Then: func(success bool) {
		EndConversation = EndingTwinsUnmasked
	},
}

var ConversationPhil = &Conversation{
	Words: []string{
		"Phil McGlass: Doctor! doctor! Every time I stand up quickly, I see Donald Duck, Goofy and Pluto.",
		"How long have you been getting these _disney_ spells?",
		"Phil McGlass: Are you taking the Mickey?",
		"Strictly diagnostically.",
	},
	Then: func(success bool) {
		if success {
			s.CurePatient()
		}
	},
}

var ConversationBarPatronLaura = &Conversation{
	Words: []string{
		"Laura: Doctor, doctor, help me! I'm getting shorter and shorter!",
		"Wait here and be a little _patient_",
	},
	Then: func(success bool) {
		if success {
			s.CurePatient()
		}
	},
}
var ConversationBarPatronJeffInBeerGarden = &Conversation{
	Words: []string{
		"Jeff: Doctor doctor, what can you give me for wind?",
		"Here, try this _kite_",
		"Jeff: What can you do for the +smell+ ?",
		"I can leave.",
	},
	Then: func(success bool) {
		if success {
			s.CurePatient()
		}
	},
}
var ConversationBarPatronTodd = &Conversation{
	Words: []string{
		//"Todd: Doctor! doctor! I just swallowed a harmonica.",
		//"Consider yourself lucky that you don't play the _guitar_",
		"Todd: Doctor! doctor! I have a lettuce stuck in my bum.",
		"*looks* Unfortunately it appears that this is just the tip of the _iceberg_",
		"...but I think it just needs a dressing.",
	},
	Then: func(success bool) {
		if success {
			s.CurePatient()
		}
	},
}
var ConversationDog = &Conversation{
	Words: []string{
		"Who's a good boy!? *pets the dog*",
		"Dog: Bork! Bork!",
		"Don't you mean Bark, Bark?",
		"Dog: Pedont.",
		"Interesting. This pundemic seems to be cross species.",
		"Dog: Cross? I'm bloody _furious_ !",
	},
	Then: func(success bool) {
		if success {
			s.CurePatient()
		}
	},
}
var ConversationJill = &Conversation{
	Words: []string{
		"Jill: Someone put _graffiti_ on my house last night!",
		"So why are you telling me?",
		"Jill: I can't understand the writing, was it you?",
	},
	Then: func(success bool) {
		if success {
			s.CurePatient()
		}
	},
}
var ConversationBarPatronTristanInBeerGarden = &Conversation{
	Words: []string{
		"Why is your nose swelling?",
		"Tristan: I bent over to _smell_ a brose.",
		"There is no b in rose.",
		"Tristan: There was a bee in this one.",
	},
	Then: func(success bool) {
		if success {
			s.CurePatient()
		}
	},
}

// Doctor! doctor! I think I am a pack of cards.
// Get up on the table so I can _deal_ with you.

// NPC: Doctor! doctor! I think I am a small bucket.
// Well, you are looking a little _pale_   <---- pale ale from the bar????

// TODO: we need a town statue!

// TODO: salad dressing joke!
// TODO: Cole's Law (colslaw)

var IntroText = &Conversation{
	Words: []string{
		"INFO: CONTROLS ARE <ARROW_KEYS> AND <SPACE>=\x14 OR <GAMEPAD_0>",
		// "DEBUG: TRIGGER DEBUG TILES ARE ON. SO IGNORE THOSE NUMBERS ON THE FLOOR ;)",
		"Before I leave the Health Centre to go on my doctors rounds...",
		"...I need to take some phone consultations.",
		// "No case is too tough for us Lexical Doctors because our rigorous application of...",
		// "...Differential Lexical Diagnosis means all maladies are shallow.",
	},
	Then: func(success bool) {
		p.Say(PhoneCall1)
	},
}

// var numCompletedPhonecalls = 0

func getIncompletePhonecall(previousCall *Conversation) *Conversation {
	// get all incomplete calls
	incompleteCalls := []*Conversation{}

	if !PhoneCall1.Completed && PhoneCall1 != previousCall {
		incompleteCalls = append(incompleteCalls, PhoneCall1)
	}
	if !PhoneCall2.Completed && PhoneCall2 != previousCall {
		incompleteCalls = append(incompleteCalls, PhoneCall2)
	}
	if !PhoneCall3.Completed && PhoneCall3 != previousCall {
		incompleteCalls = append(incompleteCalls, PhoneCall3)
	}

	// if no call is selected but the previousCall is incomplete
	// then we have to return it as the only option.
	if len(incompleteCalls) == 0 {
		if !previousCall.Completed {
			return previousCall
		} else {
			return nil
		}
	}

	// pick a random one
	return incompleteCalls[rand.Intn(len(incompleteCalls))]
}

var PhoneCall1 = &Conversation{
	Words: []string{
		"PHONE: *ring* *ring*",
		"Sid: Doctor, I've got a strawberry stuck in my ear!",
		"Don't worry, I have _cream_ for that.",
		"Sid: Will it make a mess?",
		"Only if you have meringue in there too. *click*",
	},
}

func init() {
	PhoneCall1.Then = func(success bool) {
		if nextCall := getIncompletePhonecall(PhoneCall1); nextCall != nil {
			p.Say(nextCall)
		} else {
			p.Say(IntroOutro)
		}
	}
}

var PhoneCall2 = &Conversation{
	Words: []string{
		"PHONE: *ring* *ring*",
		"Mr Jones: Doctor, my son has turned himself invisible!",
		"Quickly, take him to the _ICU_",
		"Mr Jones: Will do! Jeremy! Get in the car! *muffled* Oh, you are in the car. *click*",
		"Didn't see that one coming.",
	},
}

func init() {
	PhoneCall2.Then = func(success bool) {
		if nextCall := getIncompletePhonecall(PhoneCall2); nextCall != nil {
			p.Say(nextCall)
		} else {
			p.Say(IntroOutro)
		}
	}
}

var PhoneCall3 = &Conversation{
	Words: []string{
		"PHONE: *ring* *ring*",
		"Mrs Willis: Doctor, my son has swallowed a pen. What can I do?",
		"Use a _pencil_ until I can see him.",
		"Mrs Willis: I wish I could, but my other son is a constipated mathematician. *click*",
		"Rude. Well, I can't help everyone!",
	},
}

func init() {
	PhoneCall3.Then = func(success bool) {
		if nextCall := getIncompletePhonecall(PhoneCall3); nextCall != nil {
			p.Say(nextCall)
		} else {
			p.Say(IntroOutro)
		}
	}
}

//	&Conversation{
//		Words: []string{
//
//		"PHONE: *ring* *ring*",
//		"Tony: Doctor, I keep hearing a _ringing_ sound.",
//		"When you do hear it does it go away after the same number of rings?",
//		"Tony: Yes!",
//		"Have you heard it at all while speaking to me?",
//		"Tony: No!",
//		"Then answer your damn phone. *click*",
//		},
//	},

var EveryoneCured = &Conversation{
	Words: []string{
		"Ok, everyone is cured!",
		"I better get back to my office and phone the Health Dept back.",
	},
}

var IntroOutro = &Conversation{
	Words: []string{
		"Phone consults all done.\nTime for my rounds!",
		"PHONE: *ring* *ring*",
		"I'm not expecting another call... Hello?",
		"Unexpected Caller: It's the Health Department.",
		"Health Dept: We're seeing some unsettling health signs in the general population.",
		"Health Dept: People are exhibiting weirdly connected sets of lexical symptoms.",
		"Health Dept: We need all GPs to collect lexical differentials and attempt to treat.",
		"Health Dept: We fear this may be the start of a...",
		"Health Dept: ...pundemic.",
		"Gosh. I'll do my best.",
		"I was off on my rounds around town now anyway.",
		"Health Dept: Good luck out there. *click*",
	},
}

var EndingTwinsUnmasked = &Conversation{
	Words: []string{
		"I'll just call back the last number...",
		"Health Dept: Hello. It's the Health Department. How may I direct your call?",
		"I've cured as many people as I could...",
		"Cabinet: *KER-CHANG-PLONK*",
		"...and collected...",
		"Cabinet: *KER-CHANG-PLONK*",
		"Health Dept: *muffled* Jess, shhhhh! *muffled*",
		"...wait. Was that?",
		"It was!",
		"Joss!?",
		"Joss: Err...",
		"Joss: ...maybe?",
		"You prank called me as the Health Department?",
		"Joss: ...kinda.",
		"You pretended there was a pundemic?",
		"Joss: ...sorta.",
		"Do you have anything to say for yourself?",
		"Joss: Erm. That I'm a disenfranchised youth...",
		"Joss: ...with no real attachment to your societal constructs?",
		"Nice try, but right now you're going to do the following...",
		"Give Helen her phone back...",
		"...then go straight home. When you get home...",
		"...tell your mum you only deserve a small piece of pie.",
		"Joss: We're having pie tonight? But I love pie!",
		"But did you earn pie?",
		"Joss: No, sir.",
		"And tomorrow tell your uncle you'll feed and walk his dog for a week.",
		"Joss: Aww, but he's always borking at me!",
		"I know, but he's still your uncle.",
		"Joss: Oof.",
		"Not so funny now, is it?",
		"Joss: No, dad. Sorry, dad.",
		"Dad: Off you go, then. See you both at home.",
		"*click*",
		"Well, I still did good around town, and it feels good to do good.",
		//		"I guess some of the blame is on me. After all, I did buy Joss that...",
		//		"...government official themed voice changer he wanted.",
	},
	Then: func(success bool) {
		p.Say(Credits)
	},
}

var EndingTwinsHidden = &Conversation{
	Words: []string{
		"I'll just call back the last number...",
		"Health Dept: Hello. It's the Health Department. How may I direct your call?",
		"I've cured as many people as I could...",
		"*there's a background noise on the line that you don't recognise*",
		"*you ignore it, it's probably nothing*",
		"...and collected many useful diagnostic words.",
		"Health Dept: Excellent news. Well done you.",
		"Health Dept: Please email us those words, they'll help save lives.",
		"Health Dept: Thank you and goodbye. *click*",
		"Ah, it sure feels good to know I'm doing good!",
	},
	Then: func(success bool) {
		p.Say(Credits)
	},
}

var Credits = &Conversation{
	Words: []string{
		"CREDITS: Code, art, writing and game design by @TheMightyGit",
		"CREDITS: Font based on 'Kitchen Sink' by Retroshark & Polyducks",
		"CREDITS: Palette by Polyducks",
		"CREDITS: For GBJam-9 itch.io/jam/gbjam-9",
		// "CREDITS: But have you found both endings...?",
		"CREDITS: Press \x14 to restart.",
	},
	Then: func(success bool) {
		API.ConsoleReset()
		p.State = PlayerEndGame
	},
}

var EndConversation = EndingTwinsHidden

/*

Differential Lexical Diagnosis



-------------------------------------------------------------

Phone tutorial.

- pencil
- cream
- ICU
- hearing

// Patient: “Doctor, my son has swallowed a pen. What can I do?”
// Doctor: “Use a _pencil_ until I come see him.”

Phone:
Doctor, my son has swallowed a pen. What can I do?

A:
Use a _pencil_ until I come see him.


// How did the doctor cure the invisible man?
// He took him to the ICU.

Phone:
Doctor, my son has turned himself invisible!

A:
Quickly, take him to the _ICU_.


// Patient: “Doctor, doctor, I've got a strawberry stuck in my ear!”
// Doctor: “Don’t worry, I have some cream for that.”

Phone:
Doctor, doctor, I've got a strawberry stuck in my ear!

A:
Don’t worry, I have _cream_ for that.

// Doctor, I’m hearing a ringing sound?
// Then answer the phone.

Q:
Doctor, I keep _hearing_ a ringing sound?

A:
Then answer the phone.

-------------------------------------------------------------

Phone calls all done, now off on rounds.

-------------------------------------------------------------

- hearing (2)
- opening some windows (1)
- good handwriting (1)
- change (1)
- jogging (1)
- symptoms (1)
- silent (1)
- treating (1)

Q:
My neighbour told me about those exercise videos
by John Wick. I gave them a go, but I can't do
more than the first few exercises without
hurting myself and damaging my house.
What would you suggest?

A:
Get a _hearing_ test.

// The old man was sitting on the examining table in the doctor’s
// office, having his hearing checked. The doctor poked his light
// scope in the old man’s ear and said, “Hey, you have a suppository in your ear!”
// “Rats,” said the old man. “Now I know where my hearing aid went.”

Q:
Doctor, I need you to check my _hearing_.

A:
You have a suppository in your ear!

Q:
Rats. Now I know where my hearing aid went.


// A man goes to the doctor with a flatulence problem.
// The doctor asks, “How often do you pass gas?” and
// the man replies 10 to 15 times an hour. The doctor
// goes back to his office and returns with a pole with
// an iron hook. The man screams, “What are you going
// to do with that, Doc?”
// The doctor replies, “I’m going to open some windows.”

Q:
Doctor, I fart 10 to 15 times an hour. Can you help?

A:
Yes. I need a pole with an iron hook.

Q:
Oh my god! What do you need those for!?

A:
I'm _opening some windows_.


// How did you find that doctor was fake? She had good handwriting.

Q:
My previous doctor turned out to be a fake.

A:
Oh my word. How did they find out?

Q:
She had _good handwriting_.


// Doctor: “Nurse, how is that little girl doing who
// swallowed 10 quarters last night?”
// Nurse: “No change yet.”

Q: My son swallowed 10 quarters last night.

A: Any _change_ yet?


// My doctor told me that jogging could add years
// to my life.
// He was right — I feel 10 years older already.

Q:
You said _jogging_ would add years to my life.

A:
I did.

Q:
You were right - I feel 10 years older already.

// A man goes to the doctors and says, “Doctor, I think
// I’m going deaf!”
// And the doctor says, “Can you describe the symptoms?”
// The man responds, “Yes, Homer is fat and Marge has
// blue hair.”

Q:
Doctor, I think I'm going deaf!

A:
Can you describe the _symptoms_?

Q:
Yes, Homer is fat and Marge has blue hair.

// Patient: “Whenever I drink coffee, I have this sharp, excruciating pain.”
// Doctor: “Try to remember to remove the spoon from the cup before drinking.”

// Me: “Aren’t you going to treat me?”
// Doctor: “I am treating you.”
// Me: “You’re just staring at me.”
// Doc: “It’s called silent treatment.”

Q:
...

A:
...

Q:
... Well, aren't you doing to treat me?

A:
I am _treating_ you.

Q:
But you're just staring at me.

A:
It's called _silent_ treatment.

*/
