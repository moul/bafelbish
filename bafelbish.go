package bafelbish

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"camlistore.org/third_party/labix.org/v2/mgo/bson"

	"howett.net/plist"

	"github.com/BurntSushi/toml"
	"gopkg.in/vmihailenco/msgpack.v2"
	"gopkg.in/yaml.v2"
)

type format int

const (
	formatUnknown = iota
	formatYAML
	formatJSON
	formatTOML
	formatXML
	formatMsgpack
	formatPlist
	formatBson
	// FIXME: handle other plist formats: binary, openstep, gnustep
	// FIXME: add form-urlencoded
	// FIXME: SDL
	// FIXME: go format
	// FIXME: php format
	// FIXME: bson
	// FIXME: xdr
	// FIXME: add automatic mode
)

type Fish struct {
	InputFormat  format
	OutputFormat format
}

func NewFish() Fish {
	return Fish{
		InputFormat:  formatUnknown,
		OutputFormat: formatUnknown,
	}
}

func formatFromString(name string) (format, error) {
	formatMapping := map[string]format{
		"json":    formatJSON,
		"yaml":    formatYAML,
		"toml":    formatTOML,
		"xml":     formatXML,
		"msgpack": formatMsgpack,
		"plist":   formatPlist,
		"bson":    formatBson,
	}
	if match, found := formatMapping[strings.ToLower(name)]; found {
		return match, nil
	}
	return formatUnknown, fmt.Errorf("unsupported format: %q", name)
}

func (f *Fish) SetInputFormat(format string) (err error) {
	f.InputFormat, err = formatFromString(format)
	return
}

func (f *Fish) SetOutputFormat(format string) (err error) {
	f.OutputFormat, err = formatFromString(format)
	return
}

func Unmarshal(input []byte, inputFormat format) (interface{}, error) {
	var data interface{}
	var err error

	switch inputFormat {
	case formatJSON:
		decoder := json.NewDecoder(bytes.NewReader(input))
		decoder.UseNumber()
		err = decoder.Decode(&data)
		// FIXME: convert numbers to int64
	case formatTOML:
		_, err = toml.Decode(string(input), &data)
		// FIXME: use effective bytes to string instead whole copy
	case formatXML:
		err = xml.Unmarshal(input, &data)
	case formatMsgpack:
	case formatBson:
		err = bson.Unmarshal(input, &data)
	case formatPlist:
		input := bytes.NewReader(input)
		decoder := plist.NewDecoder(input)
		err = decoder.Decode(data)
	case formatYAML:
		err = yaml.Unmarshal(input, &data)
		if err == nil {
			data, err = convertMapsToStringMaps(data)
		}
	default:
		err = fmt.Errorf("unsupported input format")
	}

	return data, err
}

func Marshal(data interface{}, outputFormat format) ([]byte, error) {
	var result []byte
	var err error

	switch outputFormat {
	case formatJSON:
		result, err = json.Marshal(&data)
		// FIXME: option to indent json
	case formatXML:
		result, err = xml.Marshal(&data)
	case formatYAML:
		result, err = yaml.Marshal(&data)
	case formatMsgpack:
		result, err = msgpack.Marshal(&data)
	case formatBson:
		result, err = bson.Marshal(&data)
	case formatPlist:
		// result, err = plist.Marshal(&data, plist.XMLFormat)
		output := new(bytes.Buffer)
		encoder := plist.NewEncoder(output)
		err = encoder.Encode(data)
		result = output.Bytes()
	case formatTOML:
		buf := new(bytes.Buffer)
		err = toml.NewEncoder(buf).Encode(data)
		result = buf.Bytes()
	default:
		err = fmt.Errorf("unsupported output format")
	}

	return result, err
}

func (f *Fish) Parse(input io.Reader, output io.Writer) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(input)

	data, err := Unmarshal(buf.Bytes(), f.InputFormat)
	if err != nil {
		return err
	}

	result, err := Marshal(data, f.OutputFormat)
	if err != nil {
		return err
	}
	output.Write(result)
	return nil
}
