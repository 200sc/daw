// Package oak is a game engine. It provides scene control, control over windows
// and what is drawn to them, propagates regular events to evaluate game logic,
// and so on.
//
// A minimal oak app follows:
//
// 	func main() {
//		oak.AddScene("myApp", scene.Scene{Start: func(ctx *scene.Context) {
//			// ... ctx.Draw(...), event.Bind(ctx, ...)
//		}})
//		oak.Init("myApp")
//	}
package oak

import (
	"context"
	"image"
	"io"
	"sort"
	"sync/atomic"
	"time"

	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/debugstream"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
	"github.com/oakmound/oak/v4/shiny/driver"
	"github.com/oakmound/oak/v4/shiny/screen"
	"github.com/oakmound/oak/v4/window"
)

var _ window.App = &Window{}

func (w *Window) windowController(s screen.Screen, x, y, width, height int) (*driver.Window, error) {
	dwin, err := s.NewWindow(screen.NewWindowGenerator(
		screen.Dimensions(width, height),
		screen.Title(w.config.Title),
		screen.Position(x, y),
		screen.Fullscreen(w.config.Fullscreen),
		screen.Borderless(w.config.Borderless),
		screen.TopMost(w.config.TopMost),
	))
	return dwin.(*driver.Window), err
}

// the number of rgba buffers oak's draw loop swaps between
const bufferCount = 2

type Window struct {
	// The keyboard state this window is aware of.
	key.State

	// the driver.Window embedded in this window exposes at compile time the OS level
	// options one has to manipulate this.
	*driver.Window

	// TODO: most of these channels are not closed cleanly
	transitionCh chan struct{}

	// The skip scene channel receives a debug
	// signal to forcibly go to the next
	// scene.
	skipSceneCh chan string

	// The quit channel receives a signal when
	// oak should stop active workers and return from Init.
	quitCh chan struct{}

	// The draw channel receives a signal when
	// drawing should cease (or resume)
	drawCh chan struct{}

	// The between draw channel receives a signal when
	// a function is provided to Window.DoBetweenDraws.
	betweenDrawCh chan func()

	// ScreenWidth is the width of the screen
	ScreenWidth int
	// ScreenHeight is the height of the screen
	ScreenHeight int

	// FrameRate is the current logical frame rate.
	// Changing this won't directly effect frame rate, that
	// requires changing the LogicTicker, but it will take
	// effect next scene
	FrameRate int

	// DrawFrameRate is the equivalent to FrameRate for
	// the rate at which the screen is drawn.
	DrawFrameRate int

	// IdleDrawFrameRate is how often the screen will be redrawn
	// when the window is out of focus.
	IdleDrawFrameRate int

	// The window buffer represents the subsection of the world which is available to
	// be shown in a window.
	winBuffers    [bufferCount]screen.Image
	screenControl screen.Screen

	windowTextures [bufferCount]screen.Texture
	bufferIdx      uint8

	windowRect image.Rectangle

	// DrawTicker is the parallel to LogicTicker to set the draw framerate
	DrawTicker *time.Ticker
	// animationFrame is used by the javascript driver instead of DrawTicker
	animationFrame chan struct{}

	bkgFn func() image.Image

	// SceneMap is a global map of scenes referred to when scenes advance to
	// determine what the next scene should be.
	// It can be replaced or modified so long as these modifications happen
	// during a scene or before the controller has started.
	SceneMap *scene.Map

	// viewPos represents the point in the world which the viewport is anchored at.
	viewPos    intgeom.Point2
	viewBounds intgeom.Rect2

	aspectRatio float64

	// Driver is the driver oak will call during initialization
	Driver Driver

	// prePublish is a function called each draw frame prior to publishing frames to the OS
	prePublish func(*image.RGBA)

	// LoadingR is a renderable that is displayed during loading screens.
	LoadingR render.Renderable

	firstScene string
	// ErrorScene is a scene string that will be entered if the scene handler
	// fails to enter some other scene, for example, because it's name was
	// undefined in the scene map. If the scene map does not have ErrorScene
	// as well, it will fall back to panicking.
	ErrorScene string

	eventHandler  event.Handler
	CallerMap     *event.CallerMap
	MouseTree     *collision.Tree
	CollisionTree *collision.Tree
	DrawStack     *render.DrawStack

	// LastMouseEvent is the last triggered mouse event,
	// tracked for continuous mouse responsiveness on events
	// that don't take in a mouse event
	LastMouseEvent         mouse.Event
	LastRelativeMouseEvent mouse.Event
	lastRelativePress      mouse.Event
	// LastPress is the last triggered mouse event,
	// where the mouse event was a press.
	// If TrackMouseClicks is set to false then this will not be tracked
	LastMousePress mouse.Event

	FirstSceneInput interface{}

	ControllerID int32

	config Config

	mostRecentInput int32

	exitError     error
	ParentContext context.Context

	useViewBounds bool
	// UseAspectRatio determines whether new window changes will distort or
	// maintain the relative width to height ratio of the screen buffer.
	UseAspectRatio bool

	inFocus bool
}

var (
	nextControllerID = new(int32)
)

// NewWindow creates a window with default settings.
func NewWindow() *Window {
	return &Window{
		State:         key.NewState(),
		transitionCh:  make(chan struct{}),
		skipSceneCh:   make(chan string),
		quitCh:        make(chan struct{}),
		drawCh:        make(chan struct{}),
		betweenDrawCh: make(chan func()),
		SceneMap:      scene.NewMap(),
		Driver:        driver.Main,
		prePublish:    func(*image.RGBA) {},
		bkgFn: func() image.Image {
			return image.Black
		},
		eventHandler:  event.DefaultBus,
		MouseTree:     mouse.DefaultTree,
		CollisionTree: collision.DefaultTree,
		CallerMap:     event.DefaultCallerMap,
		DrawStack:     render.GlobalDrawStack,
		ControllerID:  atomic.AddInt32(nextControllerID, 1),
		ParentContext: context.Background(),
	}
}

// Propagate triggers direct mouse events on entities which are clicked
func (w *Window) Propagate(ev event.EventID[*mouse.Event], me mouse.Event) {
	hits := w.MouseTree.SearchIntersect(me.ToSpace().Bounds())
	sort.Slice(hits, func(i, j int) bool {
		return hits[i].Location.Min.Z() > hits[j].Location.Max.Z()
	})
	for _, sp := range hits {
		<-event.TriggerForCallerOn(w.eventHandler, sp.CID, ev, &me)
		if me.StopPropagation {
			break
		}
	}
	me.StopPropagation = false

	if ev == mouse.RelativePressOn {
		w.lastRelativePress = me
	} else if ev == mouse.PressOn {
		w.LastMousePress = me
	} else if ev == mouse.ReleaseOn {
		if me.Button == w.LastMousePress.Button {
			event.TriggerOn(w.eventHandler, mouse.Click, &me)

			pressHits := w.MouseTree.SearchIntersect(w.LastMousePress.ToSpace().Bounds())
			sort.Slice(pressHits, func(i, j int) bool {
				return pressHits[i].Location.Min.Z() > pressHits[j].Location.Max.Z()
			})
			for _, sp1 := range pressHits {
				for _, sp2 := range hits {
					if sp1.CID == sp2.CID {
						<-event.TriggerForCallerOn(w.eventHandler, sp1.CID, mouse.ClickOn, &me)
						if me.StopPropagation {
							return
						}
					}
				}
			}
		}
	} else if ev == mouse.RelativeReleaseOn {
		if me.Button == w.lastRelativePress.Button {
			pressHits := w.MouseTree.SearchIntersect(w.lastRelativePress.ToSpace().Bounds())
			sort.Slice(pressHits, func(i, j int) bool {
				return pressHits[i].Location.Min.Z() > pressHits[j].Location.Max.Z()
			})
			for _, sp1 := range pressHits {
				for _, sp2 := range hits {
					if sp1.CID == sp2.CID {
						<-event.TriggerForCallerOn(w.eventHandler, sp1.CID, mouse.RelativeClickOn, &me)
						if me.StopPropagation {
							return
						}
					}
				}
			}
		}
	}
}

// Width returns the absolute bounds of a window in pixels. It does not include window elements outside
// of the client area (OS provided title bars).
func (w *Window) Bounds() intgeom.Point2 {
	return intgeom.Point2{w.ScreenWidth, w.ScreenHeight}
}

// SetLoadingRenderable sets what renderable should display between scenes
// during loading phases.
func (w *Window) SetLoadingRenderable(r render.Renderable) {
	w.LoadingR = r
}

// SetBackground sets this window's background.
func (w *Window) SetBackground(b Background) {
	w.bkgFn = func() image.Image {
		return b.GetRGBA()
	}
}

// SetColorBackground sets this window's background to be a standard image.Image,
// commonly a uniform color.
func (w *Window) SetColorBackground(img image.Image) {
	w.bkgFn = func() image.Image {
		return img
	}
}

// GetBackgroundImage returns the image this window will display as its background
func (w *Window) GetBackgroundImage() image.Image {
	return w.bkgFn()
}

// SetLogicHandler swaps the logic system of the engine with some other
// implementation. If this is never called, it will use event.DefaultBus
func (w *Window) SetLogicHandler(h event.Handler) {
	w.eventHandler = h
}

// NextScene  causes this window to immediately end the current scene.
func (w *Window) NextScene() {
	w.GoToScene("")
}

// GoToScene causes this window to skip directly to the given scene.
func (w *Window) GoToScene(nextScene string) {
	go func() {
		w.skipSceneCh <- nextScene
	}()
}

// InFocus returns whether this window is currently in focus.
func (w *Window) InFocus() bool {
	return w.inFocus
}

// EventHandler returns this window's event handler.
func (w *Window) EventHandler() event.Handler {
	return w.eventHandler
}

// MostRecentInput returns the most recent input type (e.g keyboard/mouse or joystick)
// recognized by the window. This value will only change if the window is
// set to TrackInputChanges
func (w *Window) MostRecentInput() InputType {
	return InputType(w.mostRecentInput)
}

func (w *Window) exitWithError(err error) {
	w.exitError = err
	w.Quit()
}

func (w *Window) debugConsole(input io.Reader, output io.Writer) {
	debugstream.AttachToStream(w.ParentContext, input, output)
	debugstream.AddDefaultsForScope(w.ControllerID, w)
}
