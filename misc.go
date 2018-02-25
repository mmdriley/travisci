package travisci

import "encoding/json"

func mustJSONMarshal(v interface{}) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return bytes
}
