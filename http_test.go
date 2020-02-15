package hyper

import "testing"

func TestCommand(t *testing.T) {
	want := "Hello, World!"

	c := MakeCommand()
	c.Arguments["foo"] = "data:text/plain;base64,SGVsbG8sIFdvcmxkIQ%3D%3D"

	bs := c.Arguments.Bytes("foo")
	if want != string(bs) {
		t.Errorf("want: %s, got: %s", want, string(bs))
	}
}
