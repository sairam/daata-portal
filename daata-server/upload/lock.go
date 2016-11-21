package upload

import (
	"errors"
	"sync"
	"time"
)

type locky struct {
	info map[string]time.Time
	*sync.Mutex
}

var (
	// Locky takes care of locking during parallel uploads
	Locky = &locky{}
)

func getLock(path string) error {
	i := 0
	for _, ok := Locky.info[path]; ok; i++ {
		time.Sleep(5 * time.Millisecond)
		// wait for 10 times and timeout if lock is not released
		if i > 30 {
			return errors.New("Unable to unlock")
			// timeout with 504 and time taken to process the request
			// and the value of the path
		}
	}
	Locky.Lock()
	defer Locky.Unlock()

	// time.Sleep(500 * time.Millisecond) // test with this to confirm
	Locky.info[path] = time.Now()

	return nil
}

func releaseLock(path string) error {
	Locky.Lock()
	defer Locky.Unlock()

	delete(Locky.info, path)

	return nil
}
