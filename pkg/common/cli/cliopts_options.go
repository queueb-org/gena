package cli

import "text/template"

// IOption keeps [IO] options.
type IOption struct {
	FuncMap template.FuncMap
}

// mergeIOptions merges [IOption] into a single option object.
func mergeIOptions(opts ...*IOption) *IOption {
	def := &IOption{
		FuncMap: nil,
	}

	for _, o := range opts {
		if o.FuncMap != nil {
			def.FuncMap = o.FuncMap
		}
	}

	return def
}
