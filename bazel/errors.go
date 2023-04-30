package bazel

import "errors"

var ErrNoSuchParam = errors.New("no such param")
var ErrNoSuchRule = errors.New("no such rule")
