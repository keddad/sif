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

func TestDeleteEntryFromFile(t *testing.T) {
	content, err := ioutil.ReadFile("../test/cppexample/main/BUILD")

	if err != nil {
		panic(err)
	}

	newContent, err := DeleteEntryFromFile(content, "hello-world", "deps", "\"//lib:hello-time\"")

	if err != nil {
		panic(err)
	}

	depsList, err := ExtractEntriesFromFile(newContent, "hello-world", "deps")

	if !reflect.DeepEqual(depsList, []string{"\":hello-greet\""}) {
		t.Fail()
	}
}
