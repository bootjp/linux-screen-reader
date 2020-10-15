package linux_screen_reader

import (
	"fmt"
	"log"
	"os"
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

var unSupported bool

func init() {
	_, err := os.Stat("/usr/local/bin/clipnotify")
	unSupported = err != nil
}

func (o *observe) Close() error {

	err := o.webserver.Close()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

const ClipboardPrefix = "--screen_reader "

func (o *observe) clipboard() {

	if unSupported {
		fmt.Println("clipnotify is not found. please install")
		fmt.Println("disabling clipboard observe.")
		return
	}

	lastText := ""
	for {

		c := exec.Command("/usr/local/bin/clipnotify")
		err := c.Start()
		if err != nil {
			log.Fatal(err)
		}
		err = c.Wait()
		if err != nil {
			log.Fatal(err)
		}

		text, err := clipboard.ReadAll()
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
		}
	}
}

func (o observe) rest() {

	o.webserver.HideBanner = true
	o.webserver.Use(middleware.Recover())
	o.webserver.POST("/tts/speech", func(c echo.Context) error {
		text := c.FormValue("text")

		if err := o.ttsHandle.play(text); err != nil {
			fmt.Println(err)
		}

		return nil
	})

	if err := o.webserver.Start(":50500"); err != nil {
		log.Fatal(err)
	}
}

func (o *observe) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go o.rest()
	go o.clipboard()
	wg.Wait()
}
