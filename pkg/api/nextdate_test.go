package api

import (
	"testing"
	"time"
)

func TestNextDate_M31FebruaryDoesNotHang(t *testing.T) {
	now, err := time.Parse(dateFormat, "20260101")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		_, _ = NextDate(now, "20260101", "m 31 2")
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("NextDate hung on m 31 2")
	}
}
