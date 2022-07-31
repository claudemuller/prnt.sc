package main

import (
	"flag"
	"image"
	"image/draw"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang-ui/nuklear/nk"
	"github.com/xlab/closer"
	"prnt.sc/internal/pkg"
)

const (
	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024
)

func main() {
	const (
		prntscURL = "https://prnt.sc/"
		idLen     = 6
	)

	maxRetries := flag.Int("retries", 3, "the number of retries if an image URL can't be found")

	flag.Parse()
	log.SetPrefix("prnt.sc >> ")

	img, _ := pkg.GetNewImage(prntscURL, idLen, maxRetries)
	// pkg.ShowImage(img)

	const (
		winWidth  = 200
		winHeight = 100
	)

	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		closer.Fatalln(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(winWidth, winHeight, "Nuklear Demo", nil, nil)
	if err != nil {
		closer.Fatalln(err)
	}

	win.MakeContextCurrent()

	width, height := win.GetSize()
	state := &State{
		bgColor:   nk.NkRgba(28, 48, 62, 255),
		winWidth:  float32(width),
		winHeight: float32(height),
	}

	log.Printf("glfw: created window %dx%d", width, height)

	if err := gl.Init(); err != nil {
		closer.Fatalln("opengl: init failed:", err)
	}

	gl.Viewport(0, 0, int32(width), int32(height))

	ctx := nk.NkPlatformInit(win, nk.PlatformInstallCallbacks)

	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	nk.NkFontStashEnd()

	exitC := make(chan struct{}, 1)
	doneC := make(chan struct{}, 1)

	closer.Bind(func() {
		close(exitC)
		<-doneC
	})

	nk.NkTexteditInitDefault(&state.text)

	fpsTicker := time.NewTicker(time.Second / 30)

	for {
		select {
		case <-exitC:
			nk.NkPlatformShutdown()
			glfw.Terminate()
			fpsTicker.Stop()
			close(doneC)

			return
		case <-fpsTicker.C:
			if win.ShouldClose() {
				close(exitC)

				continue
			}

			glfw.PollEvents()
			gfxMain(win, ctx, state, img)
		}
	}
}

type State struct {
	bgColor   nk.Color
	text      nk.TextEdit
	winWidth  float32
	winHeight float32
}

func NkImageFromRgba(tex *uint32, rgba *image.RGBA) nk.Image {
	gl.Enable(gl.TEXTURE_2D)

	if *tex == 0 {
		gl.GenTextures(1, tex)
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *tex)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,                         // Level of detail, 0 is base image level
		gl.RGBA8,                  // Format. COuld ge RGB8 or RGB16UI
		int32(rgba.Bounds().Dx()), // Width
		int32(rgba.Bounds().Dy()), // Height
		0,                         // Must be 0
		gl.RGBA,                   // Pixel data format of last parameter rgba,Pix, could be RGB
		gl.UNSIGNED_BYTE,          // Data type for of last parameter rgba,Pix, could be UNSIGNED_SHORT
		gl.Ptr(rgba.Pix))          // Pixel data
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return nk.NkImageId(int32(*tex))
}

func gfxMain(win *glfw.Window, ctx *nk.Context, state *State, img image.Image) {
	nk.NkPlatformNewFrame()

	// Layout
	bounds := nk.NkRect(0, 0, state.winWidth, state.winHeight)
	update := nk.NkBegin(ctx, "prnt.sc", bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)

	var tex uint32

	gl.GenTextures(1, &tex)

	if update > 0 {
		nk.NkLayoutRowStatic(ctx, 30, 120, 1)
		{
			b := img.Bounds()
			m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
			i := NkImageFromRgba(&tex, m)
			nk.NkImage(ctx, i)

			if nk.NkButtonLabel(ctx, "get another pic") > 0 {
				log.Println("[INFO] button pressed!")
				win.SetSize(100, 100)
			}
		}
	}

	nk.NkEnd(ctx)

	// Render
	bg := make([]float32, 4)
	nk.NkColorFv(bg, state.bgColor)

	width, height := win.GetSize()

	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])

	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)

	win.SwapBuffers()
}
