package hyper

import (
	"io/ioutil"
	"testing"
)

func TestCommand(t *testing.T) {

	c := MakeCommand()
	c.Arguments["foo"] = "data:text/plain;base64,SGVsbG8sIFdvcmxkIQ%3D%3D"
	c.Arguments["bar"] = "data:image/gif;base64,R0lGODlhEAAQAMQAAORHHOVSKudfOulrSOp3WOyDZu6QdvCchPGolfO0o/XBs/fNwfjZ0frl3/zy7////wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACH5BAkAABAALAAAAAAQABAAAAVVICSOZGlCQAosJ6mu7fiyZeKqNKToQGDsM8hBADgUXoGAiqhSvp5QAnQKGIgUhwFUYLCVDFCrKUE1lBavAViFIDlTImbKC5Gm2hB0SlBCBMQiB0UjIQA7"

	want := "Hello, World!"
	bs := c.Arguments.Bytes("foo")
	if want != string(bs) {
		t.Errorf("want: %s, got: %s", want, string(bs))
	}

	bs = c.Arguments.Bytes("bar")
	if bs == nil {
		t.Errorf("error")
	}

	err := ioutil.WriteFile("testdata/bar.gif", bs, 0644)
	if err != nil {
		t.Error(err)
	}
}
