package command

import (
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func Test_IsStart(t *testing.T) {
	c.Convey("Is start bytes", t, func() {
		c.So(IsStart(startBytes()), c.ShouldEqual, true)
	})
	c.Convey("Is end bytes", t, func() {
		c.So(IsEnd(endBytes()), c.ShouldEqual, true)
	})
}
