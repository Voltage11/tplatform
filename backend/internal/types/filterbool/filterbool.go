package filterbool

type FilterBool string

const (
	FilterBoolTrue  FilterBool = "true"
	FilterBoolFalse FilterBool = "false"
	FilterBoolAll   FilterBool = "all"
)

func (f FilterBool) IsTrue() bool {
	return f == FilterBoolTrue
}

func (f FilterBool) IsFalse() bool {
	return f == FilterBoolFalse
}

func (f FilterBool) IsAll() bool {
	return f == FilterBoolAll || f == ""
}

func (f FilterBool) String() string {
	return string(f)
}

func (f FilterBool) IsSet() bool {
	return f != FilterBoolAll && f != ""
}

func (f FilterBool) GetBool() *bool {
	switch f {
	case FilterBoolTrue:
		b := true
		return &b
	case FilterBoolFalse:
		b := false
		return &b
	default:
		return nil
	}
}

func NewFilterBool(s string) FilterBool {
	switch s {
	case "true":
		return FilterBoolTrue
	case "false":
		return FilterBoolFalse
	default:
		return FilterBoolAll
	}
}
