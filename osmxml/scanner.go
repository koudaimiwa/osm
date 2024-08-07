package osmxml

import (
	"context"
	"encoding/xml"
	"io"
	"strings"

	"github.com/koudaimiwa/osm"
)

var _ osm.Scanner = &Scanner{}

// Scanner provides a convenient interface reading a stream of osm data
// from a file or url. Successive calls to the Scan method will step through the data.
//
// Scanning stops unrecoverably at EOF, the first I/O error, the first xml error or
// the context being cancelled. When a scan stops, the reader may have advanced
// arbitrarily far past the last token.
//
// The Scanner API is based on bufio.Scanner
// https://golang.org/pkg/bufio/#Scanner
type Scanner struct {
	ctx    context.Context
	done   context.CancelFunc
	closed bool

	decoder *xml.Decoder
	next    osm.Object
	err     error
}

// New returns a new Scanner to read from r.
func New(ctx context.Context, r io.Reader) *Scanner {
	if ctx == nil {
		ctx = context.Background()
	}

	s := &Scanner{
		decoder: xml.NewDecoder(r),
	}

	s.ctx, s.done = context.WithCancel(ctx)
	return s
}

// Close causes all future calls to Scan to return false.
// Does not close the underlying reader.
func (s *Scanner) Close() error {
	s.closed = true
	s.done()

	return nil
}

// Scan advances the Scanner to the next element, which will then be available
// through the Object method. It returns false when the scan stops, either
// by reaching the end of the input, an io error, an xml error or the context
// being cancelled. After Scan returns false, the Err method will return any
// error that occurred during scanning, except if it was io.EOF, Err will
// return nil.
func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}

Loop:
	for {
		if s.ctx.Err() != nil {
			return false
		}

		t, err := s.decoder.Token()
		if err != nil {
			s.err = err
			return false
		}

		se, ok := t.(xml.StartElement)
		if !ok {
			continue
		}

		s.next = nil
		switch strings.ToLower(se.Name.Local) {
		case "bounds":
			bounds := &osm.Bounds{}
			err = s.decoder.DecodeElement(&bounds, &se)
			s.next = bounds
		case "node":
			node := &osm.Node{}
			err = s.decoder.DecodeElement(&node, &se)
			s.next = node
		case "way":
			way := &osm.Way{}
			err = s.decoder.DecodeElement(&way, &se)
			s.next = way
		case "relation":
			relation := &osm.Relation{}
			err = s.decoder.DecodeElement(&relation, &se)
			s.next = relation
		case "changeset":
			cs := &osm.Changeset{}
			err = s.decoder.DecodeElement(&cs, &se)
			s.next = cs
		case "note":
			n := &osm.Note{}
			err = s.decoder.DecodeElement(&n, &se)
			s.next = n
		case "user":
			u := &osm.User{}
			err = s.decoder.DecodeElement(&u, &se)
			s.next = u
		default:
			continue Loop
		}

		if err != nil {
			s.err = err
			return false
		}

		return true
	}
}

// Object returns the most recent token generated by a call to Scan
// as a new osm.Object. This interface is implemented by:
//
//	*osm.Bounds
//	*osm.Node
//	*osm.Way
//	*osm.Relation
//	*osm.Changeset
//	*osm.Note
//	*osm.User
func (s *Scanner) Object() osm.Object {
	return s.next
}

// Err returns the first non-EOF error that was encountered by the Scanner.
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}

	if s.err != nil {
		return s.err
	}

	if s.closed {
		return osm.ErrScannerClosed
	}

	return s.ctx.Err()
}
