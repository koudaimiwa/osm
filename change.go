package osm

import (
	"github.com/golang/protobuf/proto"
	"github.com/paulmach/go.osm/internal/osmpb"
)

// Change is the structure of a changeset to be
// uploaded or downloaded from the server.
// See: http://wiki.openstreetmap.org/wiki/OsmChange
type Change struct {
	Create *OSM `xml:"create"`
	Modify *OSM `xml:"modify"`
	Delete *OSM `xml:"delete"`
}

// Marshal encodes the osm data using protocol buffers.
func (c *Change) Marshal() ([]byte, error) {
	ss := &stringSet{}
	encoded := marshalChange(c, ss, true)
	encoded.Strings = ss.Strings()

	return proto.Marshal(encoded)
}

// UnmarshalChange will unmarshal the data into a Change object.
func UnmarshalChange(data []byte) (*Change, error) {

	pbf := &osmpb.Change{}
	err := proto.Unmarshal(data, pbf)
	if err != nil {
		return nil, err
	}

	return unmarshalChange(pbf, pbf.GetStrings(), nil)
}

func marshalChange(c *Change, ss *stringSet, includeChangeset bool) *osmpb.Change {
	if c == nil {
		return nil
	}

	return &osmpb.Change{
		Create: marshalOSM(c.Create, ss, includeChangeset),
		Modify: marshalOSM(c.Modify, ss, includeChangeset),
		Delete: marshalOSM(c.Delete, ss, includeChangeset),
	}
}

func unmarshalChange(encoded *osmpb.Change, ss []string, cs *Changeset) (*Change, error) {
	var err error
	c := &Change{}

	c.Create, err = unmarshalOSM(encoded.Create, ss, cs)
	if err != nil {
		return nil, err
	}

	c.Modify, err = unmarshalOSM(encoded.Modify, ss, cs)
	if err != nil {
		return nil, err
	}

	c.Delete, err = unmarshalOSM(encoded.Delete, ss, cs)
	if err != nil {
		return nil, err
	}

	return c, nil
}