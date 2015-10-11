package bafelbish

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var inputs map[string][]byte = map[string][]byte{
	"json": []byte(`{"hello":["world","43"],"toto":true}`),
	"toml": []byte("hello = [\"world\", \"43\"]\ntoto = true\n"),
	"yaml": []byte("hello:\n- world\n- \"43\"\ntoto: true\n"),
	// "xmp":
	// "msgpack": []byte{130, 165, 104, 101, 108, 108, 111, 146, 165, 119, 111, 114, 108, 100, 162, 52, 51, 164, 116, 111, 116, 111, 195},
}

func TestFish_Parse(t *testing.T) {
	Convey("Testing Fish.Parse()", t, func() {
		for inputFormat, inputBuf := range inputs {
			for outputFormat, outputBuf := range inputs {
				Convey(fmt.Sprintf("%s -> %s", inputFormat, outputFormat), func() {
					fish := NewFish()

					err := fish.SetInputFormat(inputFormat)
					So(err, ShouldBeNil)

					err = fish.SetOutputFormat(outputFormat)
					So(err, ShouldBeNil)

					input := bytes.NewBufferString(string(inputBuf))
					output := new(bytes.Buffer)
					err = fish.Parse(input, output)
					// fmt.Printf("----\n%s -> %s\n%s\n", inputFormat, outputFormat, string(inputBuf))
					So(output.String(), ShouldEqual, string(outputBuf))
				})
			}
		}
	})
}
