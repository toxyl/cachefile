// This will read the data from source.txt in intervals of 30 seconds and store it to test.txt.
// Every 5 seconds the data will be printed. Try it out by changing source.txt. Changes will be reflected a couple of update cycles later.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/toxyl/cachefile"
)

func main() {
	cf := cachefile.New(
		"example/test.txt", 0644, 30*time.Second,
		func() ([]byte, error) {
			return os.ReadFile("example/source.txt")
		},
		nil,
	)
	for {
		data, err := cf.Data()
		fmt.Println(err, string(data))
		time.Sleep(5 * time.Second)
	}
}
