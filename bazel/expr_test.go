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

	if !reflect.DeepEqual(depsList, []string{"\":hello-greet\"", "\":useless\"", "\"//lib:hello-time\""}) {
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

	depsList, _ := ExtractEntriesFromFile(newContent, "hello-world", "deps")

	if !reflect.DeepEqual(depsList, []string{"\":hello-greet\"", "\":useless\""}) {
		t.Fail()
	}
}

func TestIsTest(t *testing.T) {
	res, err := IsTest("//main:hello-world", "../test/cppexample/")

	if err != nil {
		panic(err)
	}

	if res {
		t.Fail()
	}
}
