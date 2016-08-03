package nuts

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

type Bucket struct {
	raw *bolt.Bucket
}

func (b *Bucket) Bucket(name string) *Bucket {
	ret := &Bucket{
		raw: b.raw.Bucket([]byte(name)),
	}

	return ret
}

func (b *Bucket) Delete(key string) error {
	return b.raw.Delete([]byte(key))
}

func (b *Bucket) Get(key string, out interface{}) error {
	data := b.raw.Get([]byte(key))
	if data == nil {
		return errors.New("Empty key, or a bucket")
	}

	return json.Unmarshal(data, out)
}

func (b *Bucket) Put(key string, in interface{}) error {
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return b.raw.Put([]byte(key), data)
}

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
