package main

import "slices"

//go:generate stringer -type=CookType -trimprefix=Cook
type CookType int

const (
	CookFull CookType = iota
	CookIncremental
	CookIterative
)

var cookTypeMap = map[string]CookType{
	CookFull.String():        CookFull,
	CookIncremental.String(): CookIncremental,
	CookIterative.String():   CookIterative,
}

func ParseCookType(s string) (CookType, bool) {
	t, ok := cookTypeMap[s]
	return t, ok
}

func GetCookTypeStrings() []string {
	// Keys of the map contain all the string names
	names := make([]string, 0, len(cookTypeMap))
	for name := range cookTypeMap {
		names = append(names, name)
	}

	slices.Sort(names)

	return names
}
