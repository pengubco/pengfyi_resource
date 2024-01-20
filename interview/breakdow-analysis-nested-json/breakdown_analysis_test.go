package breakdowanalysisnestedjson

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBreakdown1(t *testing.T) {
	events := []Event{
		{Name: "SignUp", Properties: map[string]interface{}{
			"Device":  "Android",
			"Country": "US",
			"City":    "NYC",
		}},
		{Name: "SignUp", Properties: map[string]interface{}{
			"Name":    "SignUp",
			"Device":  "Android",
			"Country": "France",
		}},
		{Name: "SignUp", Properties: map[string]interface{}{
			"Name":    "SignUp",
			"Device":  "Desktop",
			"Country": "US",
			"City":    "Seattle",
			"Browser": "Chrome",
		}},
		{Name: "SignUp", Properties: map[string]interface{}{
			"Name":    "SignUp",
			"Device":  "iOS",
			"Country": "France",
			"Browser": "Safari",
		}},
	}

	cases := []struct {
		events    []Event
		eventName string
		props     []string
		result    Result
	}{
		{
			events:    events,
			eventName: "SignUp",
			props:     []string{"Device", "Country"},
			result: map[string]interface{}{
				"Android": map[string]interface{}{
					"US":     1,
					"France": 1,
				},
				"Desktop": map[string]interface{}{
					"US": 1,
				},
				"iOS": map[string]interface{}{
					"France": 1,
				},
			},
		},
		{
			events:    events,
			eventName: "SignUp",
			props:     []string{"Country"},
			result: map[string]interface{}{
				"US":     2,
				"France": 2,
			},
		},
		{
			events:    events,
			eventName: "SignUp",
			props:     []string{"Country", "City", "Browser"},
			result: map[string]interface{}{
				"France": map[string]interface{}{
					"UnKnown": map[string]interface{}{
						"UnKnown": 1,
						"Safari":  1,
					},
				},
				"US": map[string]interface{}{
					"NYC": map[string]interface{}{
						"UnKnown": 1,
					},
					"Seattle": map[string]interface{}{
						"Chrome": 1,
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			gotResult := breakdownAnalysis(tc.events, tc.eventName, tc.props)
			json1, _ := json.Marshal(tc.result)
			json2, _ := json.Marshal(gotResult)
			b, _ := areJSONsEqual(json1, json2)
			assert.True(t, b)
		})
	}
}

func areJSONsEqual(json1, json2 []byte) (bool, error) {
	var data1, data2 interface{}

	if err := json.Unmarshal(json1, &data1); err != nil {
		return false, err
	}

	if err := json.Unmarshal(json2, &data2); err != nil {
		return false, err
	}

	return reflect.DeepEqual(data1, data2), nil
}
