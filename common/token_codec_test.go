package common

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestTokenCharEncode(t *testing.T) {
	cnt := 0
	timeout := time.NewTimer(0)
	<-timeout.C // ignore first timeout
	defer timeout.Stop()
	resetTimeout := func() {
		timeout.Reset(1 * time.Second)
		if cnt >= 3 {
			time.Sleep(5 * time.Second)
		}
	}

	for {
		cnt++
		resetTimeout()
		select {
		case now := <-timeout.C:
			fmt.Println("nnnn", now.Unix(), time.Now().Unix())

		}
	}
}

func TestRandomToken(t *testing.T) {
	count := 100000
	for index := 0; index < count; index++ {
		data := rand.Intn(int(TOKENIDMAX))

		deData, err := TokenIdDecode(uint64(data))
		if err != nil {
			fmt.Println("data", data)
			panic(err)
		}
		newData := TokenIDEncode(deData)
		if newData != uint64(data) {
			t.Fatalf("data:%v newData:%v", data, newData)
		}
	}
}
