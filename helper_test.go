package gooik

import (
  "reflect"
  "testing"
)

func expect(title string, t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("[%v] Expected %v (type %v) - Got %v (type %v)", title, b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(title string, t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("[%v] Did not expect %v (type %v) - Got %v (type %v)", title, b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
