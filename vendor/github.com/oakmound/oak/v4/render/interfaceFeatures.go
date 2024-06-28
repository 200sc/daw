package render

import "github.com/oakmound/oak/v4/event"

// NonStatic types are not always static. If something is not NonStatic,
// it is equivalent to having IsStatic always return true.
type NonStatic interface {
	IsStatic() bool
}

// Triggerable types can have an ID set so when their animations finish,
// they trigger AnimationEnd on that ID.
type Triggerable interface {
	SetTriggerID(event.CallerID)
}

type updates interface {
	update()
}
