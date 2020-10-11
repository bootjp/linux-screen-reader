package google_tts_screenreader

import (
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/atotto/clipboard"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type observe struct {
	webserver *echo.Echo
	ttsHandle *ttsHandle
}

type observer interface {
	Close() error
	Run()
}

func NewObserve() *observe {
	return &observe{
		webserver: echo.New(),
		ttsHandle: &ttsHandle{},
	}
}

func (o *observe) Close() error {

	err := o.webserver.Close()
	if err != nil {
		log.Println(err)
	}
	return err
}

const ClipboardPrefix = "--screen_reader "

func (o *observe) clipboard() {
	lastText := ""
	for {

		c := exec.Command("./bin/clipnotify")
		err := c.Start()
		if err != nil {
			log.Println(err)
		}
		err = c.Wait()
		if err != nil {
			log.Println(err)
		}

		text, err := clipboard.ReadAll()
		if err != nil {
			log.Println(err)
		}

		if !strings.HasPrefix(text, ClipboardPrefix) {
			continue
		}

		text = strings.Replace(text, ClipboardPrefix, "", 1)

		if lastText == text {
			continue
		}
		lastText = text

		err = o.ttsHandle.play(text)
		if err != nil {
			log.Println(err)
		}
	}
}

func (o observe) rest() {

	o.webserver.Use(middleware.Recover())
	o.webserver.POST("/tts/speech", func(c echo.Context) error {
		text := c.FormValue("text")

		if err := o.ttsHandle.play(text); err != nil {
			log.Println(err)
		}

		return nil
	})

	if err := o.webserver.Start(":50500"); err != nil {
		log.Println(err)
	}
}

func (o *observe) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go o.rest()
	go o.clipboard()
	wg.Wait()
}
