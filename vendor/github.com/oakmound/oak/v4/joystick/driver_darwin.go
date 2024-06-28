package joystick

import "github.com/oakmound/oak/v4/oakerr"

func osinit() error {
	return nil
}

func newOsJoystick() osJoystick {
	return osJoystick{}
}

type osJoystick struct {
}

func (j *Joystick) prepare() error {
	return oakerr.UnsupportedPlatform{Operation: "joystick"}
}

func (j *Joystick) getState() (*State, error) {
	return nil, oakerr.UnsupportedPlatform{Operation: "joystick"}
}

func (j *Joystick) vibrate(left, right uint16) error {
	return oakerr.UnsupportedPlatform{Operation: "joystick"}
}

func (j *Joystick) close() error {
	return oakerr.UnsupportedPlatform{Operation: "joystick"}
}

func getJoysticks() []*Joystick {
	return nil
}
