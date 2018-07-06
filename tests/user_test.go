package tests

import (
	"github.com/LeeReindeer/lightblog/controllers"
	"testing"
)

func TestHash(t *testing.T) {
	h := controllers.PasswordHash("0000")
	t.Log(h)
	t.Log(len(h))
}
