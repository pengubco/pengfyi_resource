package breakdowanalysisnestedjson

type Event struct {
	Name       string
	Properties map[string]interface{}
}

type Result map[string]interface{}

func breakdownAnalysis(events []Event, eventName string, props []string) Result {
	result := make(Result)
	n := len(props)
	propValues := make([]string, n)
	for _, e := range events {
		if e.Name != eventName {
			continue
		}
		// use "UnKnown" for prop that does not exist
		for i := 0; i < n; i++ {
			v, ok := e.Properties[props[i]]
			if !ok {
				propValues[i] = "UnKnown"
			} else {
				propValues[i] = v.(string)
			}
		}
		insertToResult(result, propValues)
	}
	return result
}

// Recursively insert prop values to the result.
func insertToResult(r Result, propValues []string) {
	if len(propValues) == 0 {
		return
	}
	curValue := propValues[0]
	v, ok := r[curValue]
	if ok {
		if len(propValues) == 1 {
			r[curValue] = v.(int) + 1
		} else {
			insertToResult(v.(Result), propValues[1:])
		}
		return
	}

	if len(propValues) == 1 {
		r[curValue] = 1
	} else {
		r[curValue] = make(Result)
		insertToResult(r[curValue].(Result), propValues[1:])
	}
}
