// Copyright 2010 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package mongo

import (
	"bytes"
	"testing"
	"reflect"
)

func testMap(value interface{}) map[string]interface{} {
	return map[string]interface{}{"test": value}
}

type stEmpty struct{}

type stFloat64 struct {
	Test float
}

type stString struct {
	Test string
}

type stDoc struct {
	Test map[string]interface{}
}

type stBinary struct {
	Test []byte
}

type stObjectId struct {
	Test ObjectId
}

type stBool struct {
	Test bool
}

type stRegexp struct {
	Test Regexp
}

type stSymbol struct {
	Test Symbol
}

type stInt32 struct {
	Test int32
}

type stInt64 struct {
	Test int64
}

type stDateTime struct {
	Test DateTime
}

type stTimestamp struct {
	Test Timestamp
}

type stKey struct {
	Test Key
}

type stCodeWithScope struct {
	Test CodeWithScope
}

type stAny struct {
	Test interface{}
}

var bsonTests = []struct {
	psv  interface{}
	sv   interface{}
	mv   map[string]interface{}
	data []byte
}{
	{new(stEmpty), stEmpty{}, map[string]interface{}{},
		[]byte("\x05\x00\x00\x00\x00")},
	{new(stFloat64), stFloat64{1.5}, testMap(float64(1.5)),
		[]byte("\x13\x00\x00\x00\x01test\x00\x00\x00\x00\x00\x00\x00\xf8?\x00")},
	{new(stString), stString{"world"}, testMap("world"),
		[]byte("\x15\x00\x00\x00\x02test\x00\x06\x00\x00\x00world\x00\x00")},
	{new(stAny), stAny{"world"}, testMap("world"),
		[]byte("\x15\x00\x00\x00\x02test\x00\x06\x00\x00\x00world\x00\x00")},
	{new(stDoc), stDoc{make(map[string]interface{})}, testMap(make(map[string]interface{})),
		[]byte("\x10\x00\x00\x00\x03test\x00\x05\x00\x00\x00\x00\x00")},
	{new(stBinary), stBinary{[]byte("test")}, testMap([]byte("test")),
		[]byte("\x14\x00\x00\x00\x05\x74\x65\x73\x74\x00\x04\x00\x00\x00\x00\x74\x65\x73\x74\x00")},
	{new(stObjectId), stObjectId{ObjectId{0x4C, 0x9B, 0x8F, 0xB4, 0xA3, 0x82, 0xAA, 0xFE, 0x17, 0xC8, 0x6E, 0x63}},
		testMap(ObjectId{0x4C, 0x9B, 0x8F, 0xB4, 0xA3, 0x82, 0xAA, 0xFE, 0x17, 0xC8, 0x6E, 0x63}),
		[]byte("\x17\x00\x00\x00\x07test\x00\x4C\x9B\x8F\xB4\xA3\x82\xAA\xFE\x17\xC8\x6E\x63\x00")},
	{new(stBool), stBool{true}, testMap(true),
		[]byte("\x0C\x00\x00\x00\x08test\x00\x01\x00")},
	{new(stBool), stBool{false}, testMap(false),
		[]byte("\x0C\x00\x00\x00\x08test\x00\x00\x00")},
	{new(stSymbol), stSymbol{Symbol("aSymbol")}, testMap(Symbol("aSymbol")),
		[]byte("\x17\x00\x00\x00\x0Etest\x00\x08\x00\x00\x00aSymbol\x00\x00")},
	{new(stInt32), stInt32{10}, testMap(int32(10)),
		[]byte("\x0F\x00\x00\x00\x10test\x00\x0A\x00\x00\x00\x00")},
	{new(stInt64), stInt64{256}, testMap(int64(256)),
		[]byte("\x13\x00\x00\x00\x12test\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00")},
	{new(stKey), stKey{MaxKey}, testMap(MaxKey),
		[]byte("\x0B\x00\x00\x00\x7Ftest\x00\x00")},
	{new(stKey), stKey{MinKey}, testMap(MinKey),
		[]byte("\x0B\x00\x00\x00\xFFtest\x00\x00")},
	{nil, stRegexp{Regexp{"a*b", "i"}}, testMap(Regexp{"a*b", "i"}),
		[]byte("\x11\x00\x00\x00\vtest\x00a*b\x00i\x00\x00")},
	{nil, stCodeWithScope{CodeWithScope{"test", nil}}, testMap(CodeWithScope{"test", nil}),
		[]byte("\x1d\x00\x00\x00\x0ftest\x00\x12\x00\x00\x00\x05\x00\x00\x00test\x00\x05\x00\x00\x00\x00\x00")},
	{new(stTimestamp), stTimestamp{1168216211000}, testMap(Timestamp(1168216211000)),
		[]byte("\x13\x00\x00\x00\x11test\x008\xbe\x1c\xff\x0f\x01\x00\x00\x00")},
	{new(stDateTime), stDateTime{1168216211000}, testMap(DateTime(1168216211000)),
		[]byte("\x13\x00\x00\x00\ttest\x008\xbe\x1c\xff\x0f\x01\x00\x00\x00")},
}

func TestEncodeMap(t *testing.T) {
	for i, bt := range bsonTests {
		var buf bytes.Buffer
		err := Encode(&buf, bt.mv)
		data := buf.Bytes()
		if err != nil {
			t.Errorf("%d: error encoding %s: %s", i, bt.mv, err)
		} else if !bytes.Equal(bt.data, buf.Bytes()) {
			t.Errorf("%d: doc=%s,\n  expected %q\n  actual   %q", i, bt.mv, bt.data, data)
		}
	}
}

func TestEncodeStruct(t *testing.T) {
	for i, bt := range bsonTests {
		var buf bytes.Buffer
		err := Encode(&buf, bt.sv)
		data := buf.Bytes()
		if err != nil {
			t.Errorf("%d: error encoding %s: %s", i, bt.sv, err)
		} else if !bytes.Equal(bt.data, buf.Bytes()) {
			t.Errorf("%d: doc=%s,\n  expected %q\n  actual   %q", i, bt.mv, bt.data, data)
		}
	}
}

func TestDecodeMap(t *testing.T) {
	for i, bt := range bsonTests {
		if bt.psv == nil {
			continue
		}
		m := map[string]interface{}{}
		err := Decode(bt.data, &m)
		if err != nil {
			t.Errorf("%d: error decoding %q: %s", i, bt.data, err)
		} else if !reflect.DeepEqual(bt.mv, m) {
			t.Errorf("%d: data=%s,\n  expected %q\n  actual   %q", i, bt.data, bt.mv, m)
		}
	}
}

func TestDecodeStruct(t *testing.T) {
	for i, bt := range bsonTests {
		if bt.psv == nil {
			continue
		}
		pt := reflect.NewValue(bt.psv).Type().(*reflect.PtrType)
		psv := reflect.MakeZero(pt).(*reflect.PtrValue)
		psv.PointTo(reflect.MakeZero(pt.Elem()))
		err := Decode(bt.data, psv.Interface())
		sv := psv.Elem().Interface()
		if err != nil {
			t.Errorf("%d: error decoding %q: %s", i, bt.data, err)
		} else if !reflect.DeepEqual(sv, bt.sv) {
			t.Errorf("%d: data=%q,\n\texpected %q\n\tactual  %q", i, bt.data, bt.sv, sv)
		}
	}
}

func TestEncodeOrderedMap(t *testing.T) {
	m := OrderedMap{{"test", "hello world"}}
	data := []byte("\x1b\x00\x00\x00\x02test\x00\f\x00\x00\x00hello world\x00\x00")
	var buf bytes.Buffer
	err := Encode(&buf, m)
	if err != nil {
		t.Error("error encoding map %s", err)
	} else if !bytes.Equal(data, buf.Bytes()) {
		t.Errorf("  expected %q\n  actual   %q", data, buf.Bytes())
	}
}
