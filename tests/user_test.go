package tests

import (
	"testing"

	"github.com/LeeReindeer/lightblog/controllers"
)

func TestHash(t *testing.T) {
	h := controllers.PasswordHash("0000")
	t.Log(h)
	t.Log(len(h))
}
