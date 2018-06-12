package hvue

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCloak(t *testing.T) {
	type testGoType struct {
		Name string
	}

	Convey("Cloak/Uncloak", t, func() {
		v1 := testGoType{Name: "asdf"}
		cl1 := Cloak(&v1)
		uc1 := Uncloak(cl1).(*testGoType)
		So(&v1, ShouldEqual, uc1)

		cl1_2 := Cloak(v1)
		uc1_2 := Uncloak(cl1_2).(testGoType)
		So(v1 == uc1_2, ShouldBeTrue)
		So(&v1 != &uc1_2, ShouldBeTrue)

		So(cl1 == Cloak(&v1), ShouldBeTrue)

		v2 := testGoType{Name: "qwer"}
		So(&v1, ShouldNotEqual, &v2)
		So(&v2, ShouldEqual, &v2)
		cl2_1 := Cloak(&v2)
		cl2_2 := Cloak(&v2)
		So(cl2_1 == cl2_2, ShouldBeTrue)

		v3 := 6
		cl3 := Cloak(v3)
		uc3 := Uncloak(cl3).(int)
		So(v3 == uc3, ShouldBeTrue)

	})
}
