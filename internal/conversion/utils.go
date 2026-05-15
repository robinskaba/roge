package conversion

import "github.com/robloxapi/rbxfile"

func sourceOf(inst *rbxfile.Instance) string {
	raw, ok := inst.Properties["Source"]
	if !ok {
		return ""
	}
	switch v := raw.(type) {
	case rbxfile.ValueProtectedString:
		return string(v)
	case rbxfile.ValueString:
		return string(v)
	default:
		return ""
	}
}
