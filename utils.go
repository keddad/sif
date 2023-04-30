package main

import "strings"

func splitArgs(params string) []string {
	var res []string

	if params != "" {
		res = strings.Split(params, ",")
	} else {
		res = make([]string, 0)
	}

	return res
}
