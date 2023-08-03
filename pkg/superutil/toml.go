package superutil

import (
	"os"

	"github.com/BurntSushi/toml"
)

// ReadTOML read value from toml file
func ReadTOML(filename string, v interface{}) error {
	_, err := toml.DecodeFile(filename, v)
	return err
}

// WriteTOML write value to toml file
func WriteTOML(filename string, v interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	enc.Indent = "\t"
	return enc.Encode(v)
}
