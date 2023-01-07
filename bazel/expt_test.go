package bazel

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestExtractEntriesFromFile(t *testing.T) {
	content, err := ioutil.ReadFile("../test/cppexample/main/BUILD")

	if err != nil {
		panic(err)
	}

	depsList, err := ExtractEntriesFromFile(content, "hello-world", "deps")

	if !reflect.DeepEqual(depsList, []string{"\":hello-greet\"", "\"//lib:hello-time\""}) {
		t.Fail()
	}
}
