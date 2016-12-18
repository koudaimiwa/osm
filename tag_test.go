package osm

import (
	"bytes"
	"reflect"
	"testing"
)

func TestTagsMarshalJSON(t *testing.T) {
	data, err := Tags{}.MarshalJSON()
	if err != nil {
		t.Errorf("marshal error: %v", err)
	}

	if !bytes.Equal(data, []byte(`{}`)) {
		t.Errorf("incorrect data, got: %v", string(data))
	}

	t2 := Tags{
		Tag{Key: "highway 🏤 ", Value: "crossing"},
		Tag{Key: "source", Value: "Bind 🏤 "},
	}

	data, err = t2.MarshalJSON()
	if err != nil {
		t.Errorf("marshal error: %v", err)
	}
	if !bytes.Equal(data, []byte(`{"highway 🏤 ":"crossing","source":"Bind 🏤 "}`)) {
		t.Errorf("incorrect data, got: %v", string(data))
	}
}

func TestTagsUnmarshalJSON(t *testing.T) {
	tags := Tags{}
	data := []byte(`{"highway 🏤 ":"crossing","source":"Bind 🏤 "}`)

	err := tags.UnmarshalJSON(data)
	if err != nil {
		t.Errorf("unmarshal error: %v", err)
	}

	tags.SortByKeyValue()
	t2 := Tags{
		Tag{Key: "highway 🏤 ", Value: "crossing"},
		Tag{Key: "source", Value: "Bind 🏤 "},
	}

	if !reflect.DeepEqual(tags, t2) {
		t.Errorf("incorrect tags: %v", tags)
	}
}

func TestTagsSortByKeyValue(t *testing.T) {
	tags := Tags{
		Tag{Key: "highway", Value: "crossing"},
		Tag{Key: "source", Value: "Bind"},
	}

	tags.SortByKeyValue()
	if v := tags[0].Key; v != "highway" {
		t.Errorf("incorrect sort got %v", v)
	}

	if v := tags[1].Key; v != "source" {
		t.Errorf("incorrect sort got %v", v)
	}

	tags = Tags{
		Tag{Key: "source", Value: "Bind"},
		Tag{Key: "highway", Value: "crossing"},
	}

	tags.SortByKeyValue()
	if v := tags[0].Key; v != "highway" {
		t.Errorf("incorrect sort got %v", v)
	}

	if v := tags[1].Key; v != "source" {
		t.Errorf("incorrect sort got %v", v)
	}
}
