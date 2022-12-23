package apis

import "testing"

func TestLog_MarshalJSON(t *testing.T) {
	log := &Log{
		Level:      "info",
		Format:     "json",
		Source:     "test",
		MaxSize:    10,
		RawMessage: []byte(`xxx`),
	}
	json, err := log.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(json))
}
