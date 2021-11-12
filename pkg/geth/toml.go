package geth

import (
	"fmt"
	"io"
	"reflect"
	"unicode"

	"github.com/naoina/toml"
)

// These settings ensure that TOML keys use the same names as Go struct fields.
// source: https://github.com/ethereum/go-ethereum/blob/v1.10.11/cmd/geth/config.go#L63
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		id := fmt.Sprintf("%s.%s", rt.String(), field)
		if deprecated(id) {
			// log.Warn("Config field is deprecated and won't have an effect", "name", id)
			return nil
		}
		var link string
		if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
			link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
		}
		return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
	},
}

// deprecated returns whether the given config field is deprecated.
// source: https://github.com/ethereum/go-ethereum/blob/v1.10.11/cmd/geth/config.go#L259
func deprecated(field string) bool {
	switch field {
	case "ethconfig.Config.EVMInterpreter":
		return true
	case "ethconfig.Config.EWASMInterpreter":
		return true
	default:
		return false
	}
}

// dumpConfig dumps the configuration to a file.
// example usage: dumpConfig(cfg, os.Stdout)
func dumpConfig(cfg interface{}, w io.Writer) error {
	out, err := tomlSettings.Marshal(&cfg)
	if err != nil {
		return err
	}

	// if _, err = io.WriteString(w, comment); err != nil {
	// 	return err
	// }

	if _, err = w.Write(out); err != nil {
		return err
	}

	return nil
}
