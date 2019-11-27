package util

import (
	"strings"
	"io/ioutil"
)

func gitHeadFile(branch string) string {
	return ".git/refs/heads/" + branch
}

func GitRepoHead(dir string) (string, error) {
	val, err := ioutil.ReadFile(dir + "/" + gitHeadFile("master"))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(val)), nil
}

