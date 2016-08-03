package nut

import (
	"encoding/json"
	"errors"

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

// Delete is a wrapper around (*bolt.Bucket).Delete() which takes a string for
// the key name in place of a []byte.
func (b *Bucket) Delete(key string) error {
	return b.raw.Delete([]byte(key))
}

// Get is a wrapper around (*bolt.Bucket).Get()
//
// Get takes a key and a struct to unmarshal into and returns an error if the
// key does not exist or the Unmarshaling failed.
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
// func (b *Bucket) CreateBucketIfNotExists(key []byte) (*Bucket, error)
// func (b *Bucket) Cursor() *Cursor
// func (b *Bucket) DeleteBucket(key []byte) error
// func (b *Bucket) ForEach(fn func(k, v []byte) error) error
// func (b *Bucket) NextSequence() (uint64, error)
// func (b *Bucket) Root() pgid
// func (b *Bucket) Stats() BucketStats
// func (b *Bucket) Tx() *Tx
// func (b *Bucket) Writable() bool

// Raw will return a reference to the backing *bolt.Bucket
func (b *Bucket) Raw() *bolt.Bucket {
	return b.raw
}
