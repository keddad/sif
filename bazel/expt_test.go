package bazel

import (
	"github.com/bazelbuild/buildtools/build"
	"github.com/bazelbuild/buildtools/edit"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestExtractEntries(t *testing.T) {
	content, err := ioutil.ReadFile("../test/cppexample/main/BUILD")

	if err != nil {
		panic(err)
	}

	origBuildFile, err := build.ParseBuild("BUILD", content)

	if err != nil {
		panic(err)
	}

	targetRule := edit.FindRuleByName(origBuildFile, "hello-world")

	depsList, err := ExtractEntries(targetRule, "deps")

	if !reflect.DeepEqual(depsList, []string{"\":hello-greet\"", "\"//lib:hello-time\""}) {
		t.Fail()
	}
}
