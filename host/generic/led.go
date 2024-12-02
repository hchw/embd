// LED control support.

package generic

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kidoman/embd"
)

type led struct {
	id string

	brightness *os.File

	initialized bool
	path        string
}

func NewLED(id string) embd.LED {
	return &led{id: id}
}

func NewLEDWithPath(id, path string) embd.LED {
	return &led{id: id, path: path}
}

func (l *led) init() error {
	if l.initialized {
		return nil
	}

	var err error
	if l.brightness, err = l.brightnessFile(); err != nil {
		return err
	}

	l.initialized = true

	return nil
}

func (l *led) brightnessFilePath() string {
	if l.path != "" {
		return l.path
	}
	return fmt.Sprintf("/sys/class/leds/%v/brightness", l.id)
}

func (l *led) openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, os.ModeExclusive)
}

func (l *led) brightnessFile() (*os.File, error) {
	return l.openFile(l.brightnessFilePath())
}

func (l *led) On() error {
	if err := l.init(); err != nil {
		return err
	}

	_, err := l.brightness.WriteString("1")
	return err
}

func (l *led) Off() error {
	if err := l.init(); err != nil {
		return err
	}

	_, err := l.brightness.WriteString("0")
	return err
}

func (l *led) isOn() (bool, error) {
	l.brightness.Seek(0, 0)
	bytes := make([]byte, 1)
	_, err := io.ReadAtLeast(l.brightness, bytes, 1)
	if err != nil {
		return false, err
	}
	str := string(bytes)
	str = strings.TrimSpace(str)
	if str == "1" {
		return true, nil
	}
	return false, nil
}

func (l *led) Toggle() error {
	if err := l.init(); err != nil {
		return err
	}

	state, err := l.isOn()
	if err != nil {
		return err
	}

	if state {
		return l.Off()
	}
	return l.On()
}

func (l *led) Close() error {
	if !l.initialized {
		return nil
	}

	if err := l.brightness.Close(); err != nil {
		return err
	}

	l.initialized = false

	return nil
}
