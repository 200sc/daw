package render

import (
	"time"

	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/timing"
)

// LogicFPS is a Stackable that will draw the logical fps onto the screen when a part
// of the draw stack.
type LogicFPS struct {
	event.CallerID
	*Text
	fps       int
	lastTime  time.Time
	Smoothing float64
}

func (lf LogicFPS) CID() event.CallerID {
	return lf.CallerID
}

// NewLogicFPS returns a LogicFPS, which will render a counter of how fast it receives event.Enter events.
// If font is not provided, DefaultFont is used. If smoothing is 0, a reasonable default is used.
func NewLogicFPS(smoothing float64, font *Font, x, y float64) *LogicFPS {
	if smoothing == 0.0 {
		smoothing = defaultFpsSmoothing
	}
	if font == nil {
		font = DefaultFont().Copy()
	}
	lf := &LogicFPS{
		Smoothing: smoothing,
		lastTime:  time.Now(),
	}
	lf.Text = font.NewIntText(&lf.fps, x, y)
	lf.CallerID = event.DefaultCallerMap.Register(lf)
	// TODO: not default bus
	event.Bind(event.DefaultBus, event.Enter, lf, logicFPSBind)

	return lf
}

func logicFPSBind(lf *LogicFPS, _ event.EnterPayload) event.Response {
	t := time.Now()
	lf.fps = int((timing.FPS(lf.lastTime, t) * lf.Smoothing) + (float64(lf.fps) * (1 - lf.Smoothing)))
	lf.lastTime = t
	return 0
}
