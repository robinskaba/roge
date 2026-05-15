package conversion

import (
	"bytes"
	"fmt"

	"github.com/robloxapi/rbxfile"
	"github.com/robloxapi/rbxfile/rbxl"
)

func BuildRbxm(instance *rbxfile.Instance) ([]byte, error) {
	model := &rbxfile.Root{Instances: []*rbxfile.Instance{instance}}

	var buf bytes.Buffer
	_, err := rbxl.Encoder{Mode: rbxl.Model}.Encode(&buf, model)
	if err != nil {
		return nil, fmt.Errorf("encoding rbxm: %w", err)
	}

	return buf.Bytes(), nil
}

func DecodeRbxm(rbxm []byte) (*rbxfile.Instance, error) {
	root, _, err := rbxl.Decoder{Mode: rbxl.Model}.Decode(bytes.NewReader(rbxm))
	if err != nil {
		return nil, fmt.Errorf("parsing rbxm: %w", err)
	}
	if len(root.Instances) == 0 {
		return nil, fmt.Errorf("rbxm contains no instances")
	}
	return root.Instances[0], nil
}
