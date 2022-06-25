package roadBuillder

import (
	"JwtAuth/src/authorization"
	"errors"
	"fmt"
	"strings"
)

type RoadBuilder[T any] struct {
	head       *RoadModel[T]
	levelCount int
}

func NewRoadBuilder[T any]() *RoadBuilder[T] {
	return &RoadBuilder[T]{head: newRoadModel[T](0), levelCount: 0}
}

func (r *RoadBuilder[T]) Put(road string, data *T) error {
	r.levelCount++
	formattedEndpoint := authorization.FormatEndpoint(road)
	splits := strings.Split(formattedEndpoint, "/")
	splits = splits[1 : len(splits)-1]

	current := r.head
	for i, split := range splits {
		if split == "**" {
			isLast := (len(splits) - 1) == i
			if !isLast {
				return fmt.Errorf("uknown road defination, road: %s", road)
			}
		}

		road := current.roads[split]
		if road == nil {
			newModel := newRoadModel[T](r.levelCount)
			current.roads[split] = newModel
			current = newModel
		} else {
			current = current.roads[split]
		}
	}

	if current.data != nil {
		return errors.New("duplicated RoadBuilder definition")
	}
	// ekleme sırasına gore yetki bulma
	current.data = data
	return nil
}

func (r RoadBuilder[T]) Get(road string) (*T, error) {
	formattedEndpoint := authorization.FormatEndpoint(road)
	splits := strings.Split(formattedEndpoint, "/")
	splits = splits[1 : len(splits)-1]

	current := r.head
	var result *RoadModel[T]

	for _, split := range splits {
		applyAll := current.roads["**"]
		if applyAll != nil && (result == nil || (current.level < result.level)) {
			result = applyAll
		}
		road := current.roads[split]
		if road != nil {
			current = road
		} else {
			if result != nil {
				return result.data, nil
			}
			return nil, errors.New("not found")
		}
	}

	if current.level < result.level {
		result = current
	}

	return result.data, nil
}

func newRoadModel[T any](level int) *RoadModel[T] {
	return &RoadModel[T]{roads: make(map[string]*RoadModel[T]), data: nil, level: level}
}
