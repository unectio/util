package util

import (
	"io"
	"bufio"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type YAMLRaw struct {
	uf func(interface{}) error
}

func (raw *YAMLRaw)UnmarshalYAML(uf func(interface{}) error) error {
	raw.uf = uf
	return nil
}

func (raw *YAMLRaw)Unmarshal(v interface{}) error {
	return raw.uf(v)
}

func LoadYAML(path string, into interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.UnmarshalStrict(data, into)
	}
	return err
}

func SplitYAML(in io.Reader) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		defer close(ch)

		sc := bufio.NewScanner(in)

		var data []byte
		for sc.Scan() {
			ln := sc.Text()
			if ln != "---" {
				ln += "\n"
				data = append(data, []byte(ln)...)
				continue
			}

			if len(data) != 0 {
				ch <- data
				data = []byte{}
			}
		}

		if len(data) != 0 {
			ch <- data
		}
	}()

	return ch
}
