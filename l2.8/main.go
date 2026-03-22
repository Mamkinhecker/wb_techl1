package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	t, err := ntp.Time("pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось получить время от NTP: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(t.Format(time.RFC3339Nano))
}
