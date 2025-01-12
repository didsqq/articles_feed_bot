package botkit

import "encoding/json"

func ParseJSON[T any](src string) (T, error) {
	var res T
	if err := json.Unmarshal([]byte(src), &res); err != nil {
		return *(new(T)), err
	}
	return res, nil
}
