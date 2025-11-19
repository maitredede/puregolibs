package tools

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeval(t *testing.T) {

	sec := 2
	mSec := 5

	dur := time.Duration(sec)*time.Second + time.Duration(mSec)*time.Millisecond

	tv := DurationToTimeVal(dur)

	assert.EqualValues(t, sec, tv.Sec, "tv.Sec wrong value")
	assert.EqualValues(t, mSec*1000, tv.Usec, "tv.Usec wrong value")
}
