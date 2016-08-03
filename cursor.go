package nut

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

var (
	// ErrCursorEOF is a marker error signifying that the end of the
	// cursor has been reached.
	ErrCursorEOF = errors.New("End of Cursor")

	// ErrCursorBucket is a marker error signifying that the value is
	// actually a nested bucket.
	ErrCursorBucket = errors.New("Cursor is at bucket")
)

// Cursor is a wrapper around *bolt.Cursor
type Cursor struct {
	raw *bolt.Cursor
}

// Delete is a wrapper around (*bolt.Cursor).Delete()
func (c *Cursor) Delete() error {
	return c.raw.Delete()
}

// First is a wrapper around (*bolt.Cursor).First()
//
// First takes a struct to unmarshal into and returns the key along
// with any errors that occured.
//
// ErrCursorEOF will be returned if the bucket is empty.
func (c *Cursor) First(out interface{}) (string, error) {
	k, v := c.raw.First()
	if k == nil {
		return "", ErrCursorEOF
	} else if v == nil {
		return string(k), ErrCursorBucket
	}

	return string(k), json.Unmarshal(v, out)
}

// Last is a wrapper around (*bolt.Cursor).Last()
//
// Last takes a struct to unmarshal into and returns the key along
// with any errors that occured.
//
// ErrCursorEOF will be returned if the bucket is empty.
func (c *Cursor) Last(out interface{}) (string, error) {
	k, v := c.raw.Last()
	if k == nil {
		return "", ErrCursorEOF
	} else if v == nil {
		return string(k), ErrCursorBucket
	}

	return string(k), json.Unmarshal(v, out)
}

// Next is a wrapper around (*bolt.Cursor).Next()
//
// Next takes a struct to unmarshal into and returns the key along
// with any errors that occured.
//
// ErrCursorEOF will be returned if the cursor is at the end of the
// bucket.
func (c *Cursor) Next(out interface{}) (string, error) {
	k, v := c.raw.Next()
	if k == nil {
		return "", ErrCursorEOF
	} else if v == nil {
		return string(k), ErrCursorBucket
	}

	return string(k), json.Unmarshal(v, out)
}

// Prev is a wrapper around (*bolt.Cursor).Prev()
//
// Prev takes a struct to unmarshal into and returns the key along
// with any errors that occured.
//
// ErrCursorEOF will be returned if the cursor is at the beginning of
// the bucket.
func (c *Cursor) Prev(out interface{}) (string, error) {
	k, v := c.raw.Prev()
	if k == nil {
		return "", ErrCursorEOF
	} else if v == nil {
		return string(k), ErrCursorBucket
	}

	return string(k), json.Unmarshal(v, out)
}

// Seek is a wrapper around (*bolt.Cursor).Seek()
//
// Prev takes the target key and a struct to unmarshal into and
// returns the found key along with any errors that occured.
//
// Because this is a wrapper around (*bolt.Cursor).Seek(), if a key
// does not exist, the next key (and value) will be used.
//
// ErrCursorEOF will be returned if the key does not exist and no keys
// follow.
func (c *Cursor) Seek(key string, out interface{}) (string, error) {
	k, v := c.raw.Seek([]byte(key))
	if k == nil {
		return "", ErrCursorEOF
	} else if v == nil {
		return string(k), ErrCursorBucket
	}

	return string(k), json.Unmarshal(v, out)
}

// Unimplemented Cursor methods
// func (c *Cursor) Bucket() *Bucket

// Raw will return a reference to the backing *bolt.Cursor
func (c *Cursor) Raw() *bolt.Cursor {
	return c.raw
}
