package nut

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/boltdb/bolt"
)

// Bucket is a wrapper around *bolt.Bucket
type Bucket struct {
	raw *bolt.Bucket
}

// Bucket is a wrapper around (*bolt.Bucket).Bucket() which takes a string for
// the bucket name in place of a []byte.
func (b *Bucket) Bucket(name string) *Bucket {
	ret := &Bucket{
		raw: b.raw.Bucket([]byte(name)),
	}

	return ret
}

// CreateBucketIfNotExists is a wrapper around
// (*bolt.Bucket).CreateBucketIfNotExists which takes a string for the key in
// place of a []byte.
func (b *Bucket) CreateBucketIfNotExists(name string) (*Bucket, error) {
	rawBucket, err := b.raw.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, err
	}

	return &Bucket{
		raw: rawBucket,
	}, nil
}

// Cursor is a wrapper around (*bolt.Bucket).Cursor()
func (b *Bucket) Cursor() *Cursor {
	ret := &Cursor{
		raw: b.raw.Cursor(),
	}

	return ret
}

// Delete is a wrapper around (*bolt.Bucket).Delete() which takes a string for
// the key name in place of a []byte.
func (b *Bucket) Delete(key string) error {
	return b.raw.Delete([]byte(key))
}

// Get is a wrapper around (*bolt.Bucket).Get() which takes a key and
// a struct to unmarshal into and returns an error if the key does not
// exist or the Unmarshaling failed.
func (b *Bucket) Get(key string, out interface{}) error {
	data := b.raw.Get([]byte(key))
	if data == nil {
		return errors.New("Empty key, or a bucket")
	}

	return json.Unmarshal(data, out)
}

// Put is a wrapper around (*bolt.Bucket).Put()
//
// Put takes a key and a struct to store and returns an error if the Marshaling
// or the underlying Put failed.
func (b *Bucket) Put(key string, in interface{}) error {
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return b.raw.Put([]byte(key), data)
}

// Unimplemented Bucket methods
// func (b *Bucket) CreateBucket(key []byte) (*Bucket, error)
// func (b *Bucket) DeleteBucket(key []byte) error
// func (b *Bucket) ForEach(fn func(k, v []byte) error) error
// func (b *Bucket) Root() pgid
// func (b *Bucket) Stats() BucketStats
// func (b *Bucket) Tx() *Tx
// func (b *Bucket) Writable() bool

// NextID is a loose wrapper around (*bolt.Bucket).NextSequence()
// which will return the next sequence from the DB in base32. This may
// be changed later, but if this happens, it will be ensured that any
// new IDs will not conflict with lower IDs.
func (b *Bucket) NextID() (string, error) {
	i, err := b.raw.NextSequence()
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(i, 32), nil
}

// Raw will return a reference to the backing *bolt.Bucket
func (b *Bucket) Raw() *bolt.Bucket {
	return b.raw
}
