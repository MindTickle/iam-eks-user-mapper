package main

import (
	"fmt"
	"strings"
)

func extractIAMK8sFromString(str string) (string, string) {
	splits := strings.Split(str, "::")
	if len(splits) != 2 {
		panic(fmt.Sprintf("Invalid flag value %s", str))
	}
	iam := splits[0]
	k8s := splits[1]
	return iam, k8s
}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
