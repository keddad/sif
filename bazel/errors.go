package bazel

import "errors"

var ErrNoSuchParam = errors.New("no such param")
var ErrNoSuchRule = errors.New("no such rule")
var ErrNoSuchListElem = errors.New("no such list element")
var ErrNotListAssingment = errors.New("param is not a list")
var ErrInvalidInvariant = errors.New("wtf") // What a Terrible Failure! Probably can't happen with valid file
