package cartridge

import (
	"image"

	"github.com/TheMightyGit/marv/marvtypes"
)

type InputType struct { // prob should be in MarvBench
	MousePos            image.Point
	MousePosDelta       image.Point
	MouseWheelDelta     image.Point
	MousePressed        bool
	MouseHeld           bool
	MouseReleased       bool
	InputChars          []rune
	GamepadButtonStates [4]marvtypes.GamepadState // bitmask per gamepad (up to 4)
}
