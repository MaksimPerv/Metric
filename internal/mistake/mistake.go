package mistake

import (
	"fmt"
	"time"
)

func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if i < attempts-1 {
			time.Sleep(delay)
			delay += 2
		}
	}
	return fmt.Errorf("after %d attempts, last error: %w", attempts, err)
}
