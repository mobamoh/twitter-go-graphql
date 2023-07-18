package faker

import (
	"fmt"
	"github.com/mobamoh/twitter-go-graphql/uuid"
	"math/rand"
	"time"
)

var EncryptedPassword = "$2a$04$U9TeZV83pXGXdMYinZMTyOIGZ47wN76H0wCdY0Tl5UhXexUUQp7Va"
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func init() {
	rand.Seed(time.Now().UnixNano())
}
func RandStringRunes(n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(s)
}

func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
func Username() string {
	return RandStringRunes(RandInt(2, 10))
}

func Email() string {
	return fmt.Sprintf("%s@mail.com", RandStringRunes(RandInt(5, 10)))
}

func UUID() string {
	return uuid.Generate()
}
