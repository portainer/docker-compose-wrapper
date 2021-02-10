package wrapper

import (
	"testing"
)

func TestCommand(t *testing.T) {
	w, err := NewComposeWrapper("")
	if err != nil {
		t.Fatal(err)
	}

	file := "docker-compose-test.yml"
	_, err = w.Up(file, "", "", "")
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Down(file, "", "")
	if err != nil {
		t.Fatal(err)
	}

}
