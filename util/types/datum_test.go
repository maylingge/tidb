// Copyright 2016 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	. "github.com/pingcap/check"
	"github.com/pingcap/tidb/mysql"
)

var _ = Suite(&testDatumSuite{})

type testDatumSuite struct {
}

func (ts *testDatumSuite) TestDatum(c *C) {
	values := []interface{}{
		int64(1),
		uint64(1),
		1.1,
		"abc",
		[]byte("abc"),
		[]int{1},
	}
	for _, val := range values {
		var d Datum
		d.SetValue(val)
		x := d.GetValue()
		c.Assert(x, DeepEquals, val)
	}
}

func testDatumToBool(c *C, in interface{}, res int) {
	datum := NewDatum(in)
	res64 := int64(res)
	b, err := datum.ToBool()
	c.Assert(err, IsNil)
	c.Assert(b, Equals, res64)
}

func (ts *testDatumSuite) TestToBool(c *C) {
	testDatumToBool(c, int(0), 0)
	testDatumToBool(c, int64(0), 0)
	testDatumToBool(c, uint64(0), 0)
	testDatumToBool(c, float32(0), 0)
	testDatumToBool(c, float64(0), 0)
	testDatumToBool(c, "", 0)
	testDatumToBool(c, "0", 0)
	testDatumToBool(c, []byte{}, 0)
	testDatumToBool(c, []byte("0"), 0)
	testDatumToBool(c, mysql.Hex{Value: 0}, 0)
	testDatumToBool(c, mysql.Bit{Value: 0, Width: 8}, 0)
	testDatumToBool(c, mysql.Enum{Name: "a", Value: 1}, 1)
	testDatumToBool(c, mysql.Set{Name: "a", Value: 1}, 1)

	t, err := mysql.ParseTime("2011-11-10 11:11:11.999999", mysql.TypeTimestamp, 6)
	c.Assert(err, IsNil)
	testDatumToBool(c, t, 1)

	td, err := mysql.ParseDuration("11:11:11.999999", 6)
	c.Assert(err, IsNil)
	testDatumToBool(c, td, 1)

	ft := NewFieldType(mysql.TypeNewDecimal)
	ft.Decimal = 5
	v, err := Convert(3.1415926, ft)
	c.Assert(err, IsNil)
	testDatumToBool(c, v, 1)
	d := NewDatum(&invalidMockType{})
	_, err = d.ToBool()
	c.Assert(err, NotNil)
}

func (ts *testDatumSuite) TestEqualDatums(c *C) {
	testCases := []struct {
		a    []interface{}
		b    []interface{}
		same bool
	}{
		// Positive cases
		{[]interface{}{1}, []interface{}{1}, true},
		{[]interface{}{1, "aa"}, []interface{}{1, "aa"}, true},
		{[]interface{}{1, "aa", 1}, []interface{}{1, "aa", 1}, true},

		// Negative cases
		{[]interface{}{1}, []interface{}{2}, false},
		{[]interface{}{1, "a"}, []interface{}{1, "aaaaaa"}, false},
		{[]interface{}{1, "aa", 3}, []interface{}{1, "aa", 2}, false},

		// Corner cases
		{[]interface{}{}, []interface{}{}, true},
		{[]interface{}{nil}, []interface{}{nil}, true},
		{[]interface{}{}, []interface{}{1}, false},
		{[]interface{}{1}, []interface{}{1, 1}, false},
		{[]interface{}{nil}, []interface{}{1}, false},
	}
	for _, t := range testCases {
		testEqualDatums(c, t.a, t.b, t.same)
	}
}

func testEqualDatums(c *C, a []interface{}, b []interface{}, same bool) {
	res, err := EqualDatums(MakeDatums(a), MakeDatums(b))
	c.Assert(err, IsNil)
	c.Assert(res, Equals, same, Commentf("a: %v, b: %v", a, b))
}

func testDatumToInt64(c *C, val interface{}, expect int64) {
	d := NewDatum(val)
	b, err := d.ToInt64()
	c.Assert(err, IsNil)
	c.Assert(b, Equals, expect)
}

func (ts *testTypeConvertSuite) TestToInt64(c *C) {
	testDatumToInt64(c, "0", int64(0))
	testDatumToInt64(c, int(0), int64(0))
	testDatumToInt64(c, int64(0), int64(0))
	testDatumToInt64(c, uint64(0), int64(0))
	testDatumToInt64(c, float32(3.1), int64(3))
	testDatumToInt64(c, float64(3.1), int64(3))
	testDatumToInt64(c, mysql.Hex{Value: 100}, int64(100))
	testDatumToInt64(c, mysql.Bit{Value: 100, Width: 8}, int64(100))
	testDatumToInt64(c, mysql.Enum{Name: "a", Value: 1}, int64(1))
	testDatumToInt64(c, mysql.Set{Name: "a", Value: 1}, int64(1))

	t, err := mysql.ParseTime("2011-11-10 11:11:11.999999", mysql.TypeTimestamp, 0)
	c.Assert(err, IsNil)
	testDatumToInt64(c, t, int64(20111110111112))

	td, err := mysql.ParseDuration("11:11:11.999999", 6)
	c.Assert(err, IsNil)
	testDatumToInt64(c, td, int64(111112))

	ft := NewFieldType(mysql.TypeNewDecimal)
	ft.Decimal = 5
	v, err := Convert(3.1415926, ft)
	c.Assert(err, IsNil)
	testDatumToInt64(c, v, int64(3))

	_, err = ToInt64(&invalidMockType{})
	c.Assert(err, NotNil)
}
