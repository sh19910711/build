package worker_test

import (
	"testing"
)

func TestAttach(t *testing.T) {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	w, err := createFakeWorker(ctx)
	if err != nil {
		t.Fatal(err)
	}

	r, err := w.Attach(ctx)
	if err != nil {
		t.Fatal(err)
	}

	finished := make(chan bool) // true if the log contains expected output
	go func() {
		s := readerToString(r)
		finished <- (contains(s, "hello 1") && contains(s, "hello 2") && contains(s, "hello 3"))
	}()

	if err := w.Start(ctx); err != nil {
		t.Fatal(err)
	}

	if exitCode, err := w.Wait(ctx); err != nil {
		t.Fatal(err)
	} else if exitCode != 0 {
		t.Fatal("worker exited with status ", exitCode)
	}

	if !<-finished {
		t.Fatal("something went wrong")
	}
}
