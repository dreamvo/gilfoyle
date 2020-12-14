package testutils

import (
	"github.com/dreamvo/gilfoyle/worker"
	"testing"
)

func CloseWorker(t *testing.T, w *worker.Worker) {
	err := w.Close()
	if err != nil {
		t.Error(err)
	}
}
