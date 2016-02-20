package pkg

import (
	pf "gitlab.com/abduld/wgx-pandoc/pkg/pandocfilter"
)

func HeaderFilter(k string, v interface{}, format string, meta interface{}) interface{} {

	if k == "Header" {
		level := v.([]interface{})[0].(float64)
		attrs := v.([]interface{})[1].([]interface{})
		inlines := v.([]interface{})[2].([]interface{})

		return pf.Header(int(level)-1, attrs, inlines)
	}

	return nil
}

func init() {
	AddFilter(HeaderFilter)
}
