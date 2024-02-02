package stream

import (
	"encoding/json"

	"github.com/r3labs/sse/v2"
	"github.com/szkiba/xk6-top/internal/digest"
)

func parse(msg *sse.Event) (*digest.Event, error) {
	var (
		etype digest.EventType
		edata interface{}
		err   error
	)

	if err = etype.UnmarshalText(msg.Event); err != nil {
		return nil, err
	}

	edata, err = unmarshalData(etype, msg.Data)
	if err != nil {
		return nil, err
	}

	return &digest.Event{Type: etype, Data: edata}, nil
}

func unmarshalData(etype digest.EventType, data []byte) (interface{}, error) {
	switch etype {
	case digest.EventTypeMetric:
		target := make(digest.Metrics)

		if err := json.Unmarshal(data, &target); err != nil {
			return nil, err
		}

		return target, nil

	case digest.EventTypeParam:
		target := new(digest.ParamData)

		if err := json.Unmarshal(data, target); err != nil {
			return nil, err
		}

		return target, nil

	case digest.EventTypeConfig:
		target := make(digest.ConfigData)

		if err := json.Unmarshal(data, &target); err != nil {
			return nil, err
		}

		return target, nil

	case digest.EventTypeStart,
		digest.EventTypeStop,
		digest.EventTypeSnapshot,
		digest.EventTypeCumulative:
		target := make(digest.Aggregates)

		if err := json.Unmarshal(data, &target); err != nil {
			return nil, err
		}

		return target, nil

	default:
		return nil, nil //nolint:nilnil
	}
}
