package daw

import (
	"context"
	"image/color"
	"image/draw"
	"io"
	"math"
	"sync"

	oak "github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/audio"
	"github.com/oakmound/oak/v4/audio/pcm"
	"github.com/oakmound/oak/v4/dlog"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

func init() {
	err := audio.InitDefault()
	if err != nil {
		panic(err)
	}
}

func Play(ctx context.Context, r pcm.Reader) error {
	return audio.Play(ctx, r)
}

func Main(fn func(io.Writer)) {
	w := NewWriter()
	ioW := NotAnIOWriter{w}
	fn(ioW)
}

func VisualMain(fn func(io.Writer)) {
	viz := VisualWriter(DefaultFormat)
	ioW := NotAnIOWriter{viz}
	fn(ioW)
}

type NotAnIOWriter struct {
	Writer
}

func (w NotAnIOWriter) Write(b []byte) (n int, err error) {
	return w.WritePCM(b)
}

func Write(data []byte) (int, error) {
	return NewWriter().WritePCM(data)
}

func WritePCM(data []byte) (int, error) {
	return NewWriter().WritePCM(data)
}

type Writer interface {
	io.Closer
	PCMFormat() pcm.Format
	WritePCM([]byte) (n int, err error)
}

func NewWriter() Writer {
	return audio.MustNewWriter(DefaultFormat)
}

func NewPCMWriter(format pcm.Format) Writer {
	return audio.MustNewWriter(format)
}

var DefaultFormat = pcm.Format{
	SampleRate: 44100,
	Channels:   2,
	Bits:       32,
}

func BufferLength(format pcm.Format) uint32 {
	return uint32(float64(format.BytesPerSecond()) * audio.WriterBufferLengthInSeconds)
}

func VisualWriter(format pcm.Format) Writer {
	var monitor *pcmMonitor

	var wg sync.WaitGroup
	wg.Add(1)
	oak.AddScene("visualizer", scene.Scene{
		Start: func(ctx *scene.Context) {
			speaker := audio.MustNewWriter(format)
			monitor = newPCMMonitor(ctx, speaker)
			monitor.SetPos(0, 0)
			render.Draw(monitor)
			wg.Done()
		},
	})
	go oak.Init("visualizer", func(c oak.Config) (oak.Config, error) {
		c.Screen.Height = 240
		c.Title = "Audio Visualizer"
		c.Debug.Level = dlog.INFO.String()
		return c, nil
	})
	wg.Wait()
	return monitor
}

func PlayTo(dst pcm.Writer, src pcm.Reader) {
	audio.Play(context.Background(), src, func(po *audio.PlayOptions) {
		po.Destination = dst
	})
}

func Loop(dst pcm.Writer, src pcm.Reader) {
	audio.Play(context.Background(), audio.LoopReader(src), func(po *audio.PlayOptions) {
		po.Destination = dst
	})
}

func LoopContext(ctx context.Context, dst pcm.Writer, src pcm.Reader) {
	audio.Play(ctx, audio.LoopReader(src), func(po *audio.PlayOptions) {
		po.Destination = dst
	})
}

type pcmMonitor struct {
	event.CallerID
	render.LayeredPoint
	pcm.Writer
	pcm.Format
	written []byte
	at      int
}

var globalMagnification float64 = 1

func newPCMMonitor(ctx *scene.Context, w pcm.Writer) *pcmMonitor {
	fmt := w.PCMFormat()
	pm := &pcmMonitor{
		Writer:       w,
		Format:       w.PCMFormat(),
		LayeredPoint: render.NewLayeredPoint(0, 0, 0),
		written:      make([]byte, int(float64(fmt.BytesPerSecond())*audio.WriterBufferLengthInSeconds)),
	}
	event.GlobalBind(ctx, mouse.ScrollDown, func(_ *mouse.Event) event.Response {
		mag := globalMagnification - 0.5
		if mag < 1 {
			mag = 1
		}
		globalMagnification = mag
		return 0
	})
	event.GlobalBind(ctx, mouse.ScrollUp, func(_ *mouse.Event) event.Response {
		globalMagnification += 0.5
		return 0
	})
	return pm
}

func (pm *pcmMonitor) CID() event.CallerID {
	return pm.CallerID
}

func (pm *pcmMonitor) PCMFormat() pcm.Format {
	return pm.Format
}

func (pm *pcmMonitor) WritePCM(b []byte) (n int, err error) {
	copy(pm.written[pm.at:], b)
	if len(b) > len(pm.written[pm.at:]) {
		copy(pm.written[0:], b[len(pm.written[pm.at:]):])
	}
	pm.at += len(b)
	pm.at %= len(pm.written)
	return pm.Writer.WritePCM(b)
}

func (pm *pcmMonitor) Draw(buf draw.Image, xOff, yOff float64) {
	const width = 640
	const height = 200.0
	xJump := len(pm.written) / width
	xJump = int(float64(xJump) / globalMagnification)
	c := color.RGBA{255, 255, 255, 255}
	for x := 0.0; x < width; x++ {
		wIndex := int(x) * xJump

		var val int16
		switch pm.Format.Bits {
		case 8:
			val8 := pm.written[wIndex]
			val = int16(val8) << 8
		case 16:
			wIndex -= wIndex % 2
			val = int16(pm.written[wIndex+1])<<8 +
				int16(pm.written[wIndex])
		case 32:
			wIndex = wIndex - wIndex%4
			val32 := int32(pm.written[wIndex+3])<<24 +
				int32(pm.written[wIndex+2])<<16 +
				int32(pm.written[wIndex+1])<<8 +
				int32(pm.written[wIndex])
			val = int16(val32 / int32(math.Pow(2, 16)))
		}

		// -32768 -> 200
		// 0 -> 100
		// 32768 -> 0
		var y float64
		if val < 0 {
			y = height/2 + float64(val)*float64(height/2/-32768.0)
		} else {
			y = height/2 + -(float64(val) * float64(height/2/32768.0))
		}
		buf.Set(int(x+xOff+pm.X()), int(y+yOff+pm.Y()), c)
	}
}
