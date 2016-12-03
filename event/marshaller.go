package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func (m MetaData)MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	length := len(m)
	count := 0
	for key, value := range m {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("\"%d\":%s", key, string(jsonValue)))
		count++
		if count < length {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

func (c *Container) UnmarshalJSON(def MetaDataDefinition, event interface{}, b []byte) error {
	var temp  struct {
		Event      Event
		MetaData   map[string]interface{}
		Identifier Identifier
		Version    Version
	}
	temp.Event = event

	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	metaData := MetaData{}
	for key, value := range temp.MetaData {
		id, err := strconv.Atoi(key)
		if err != nil {
			return err
		}
		metaData[def.Generator(id)] = value
	}

	c.Event = temp.Event
	c.MetaData = metaData
	c.Identifier = temp.Identifier
	c.Version = temp.Version

	return nil
}

func (m MetaData) UnmarshalJSON(def MetaDataDefinition, b []byte) error {
	var temp map[string]interface{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	for key, value := range temp {
		id, err := strconv.Atoi(key)
		if err != nil {
			return err
		}
		m[def.Generator(id)] = value
	}
	return nil
}