package shutdowner

import (
	"io"
	"os"
	"os/signal"
	"syscall"
)

type shutdowner struct {
	list []io.Closer
	done chan struct{}
}

// New returns a new Shutdowner instance.
func New(list ...io.Closer) *shutdowner {
	sd := shutdowner{
		done: make(chan struct{}),
		list: list,
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(
		sigs,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	go func() {
		<-sigs
		sd.Shutdown()
	}()

	return &sd
}

// Add saves items to close on shutdown.
func (sd *shutdowner) Add(c ...io.Closer) {
	sd.list = append(sd.list, c...)
}

// Shutdown closes all items.
func (sd shutdowner) Shutdown() {
	for _, c := range sd.list {
		c.Close()
	}

	close(sd.done)
}

// Wait for the shutdown to complete.
func (sd shutdowner) Wait() {
	<-sd.done
}
