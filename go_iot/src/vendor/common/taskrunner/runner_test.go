package taskrunner

import (
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	d := func(dc DataChan) error {
		log.Printf("Dispatcher sent: %d", 1)
		dc <- 1
		return nil
	}

	e := func(dc DataChan) error {
	forloop:
		for {
			select {
			case d := <-dc:
				log.Printf("Executor received: %v", d)
			default:
				break forloop
			}
		}
		return nil
	}

	runner := NewRunner(30, false, d, e)
	go runner.startAll()
	time.Sleep(3 * time.Second)

}
