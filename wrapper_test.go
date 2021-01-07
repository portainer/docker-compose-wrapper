package wrapper

import (
	"log"
	"testing"
)

func TestCommand(t *testing.T) {
	w := NewComposeWrapper()

	file := "docker-compose-test.yml"
	_, err := w.Up(file)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = w.Down(file)
	if err != nil {
		log.Fatalln(err)
	}

}
