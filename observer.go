package google_tts_screenreader

import (
	"fmt"
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
	ttsHandle *tts_handle
}

type observer interface {
	Close() error
	Run()
}

func NewObserve() *observe {
	return &observe{
		webserver: echo.New(),
		ttsHandle: &tts_handle{},
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
		fmt.Println("last: "+lastText, "current: "+text)
		lastText = text

		fmt.Println("clipboard ", text)

		err = o.ttsHandle.play(text)
		if err != nil {
			log.Println(err)
		}
	}
}

func (o observe) rest() {

	o.webserver.Use(middleware.Recover())
	o.webserver.POST("/", func(c echo.Context) error {
		text := c.FormValue("text")

		fmt.Println(text)

		if err := o.ttsHandle.play(text); err != nil {
			log.Println(err)
		}

		fmt.Println("played")
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
