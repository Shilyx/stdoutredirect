package stdoutredirect

import (
	"os"
	"sync"
)

var (
	once      sync.Once
	newStdout chan string
)

// GetNewStdout get new stdout channel
func GetNewStdout() chan string {
	once.Do(func() {
		newStdout = make(chan string)
		r, w, _ := os.Pipe()

		os.Stdout = w

		go func() {
			for {
				buf := make([]byte, 10240)
				n, err := r.Read(buf)

				if err != nil {
					break
				}

				if n == 0 {
					continue
				}

				select {
				case newStdout <- string(buf[0:n]):
				default:
				}
			}
		}()
	})

	return newStdout
}
