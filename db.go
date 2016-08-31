package nut

import (
	"os"

	"github.com/boltdb/bolt"
)

// DB is a wrapper around *bolt.DB
type DB struct {
	// We want to keep this mostly internal, but let the user get to
	// it if they absolutely have to.
	raw *bolt.DB
}

// Open is a wrapper around bolt.Open which will create a *DB from the resulting
// *bolt.DB.
//
// NOTE: Options are not currently supported. If you need fancy options set on
// the underlying *bolt.DB, please use NewDB and pass in the created underlying
// *bolt.DB yourself.
func Open(path string, mode os.FileMode) (*DB, error) {
	raw, err := bolt.Open(path, mode, nil)
	if err != nil {
		return nil, err
	}

	return NewDB(raw), nil
}

// NewDB will create a new nutdb given a boltdb.
func NewDB(rawDB *bolt.DB) *DB {
	db := &DB{
		raw: rawDB,
	}

	return db
}

// Close wraps (*bolt.DB).Close()
func (db *DB) Close() error {
	return db.raw.Close()
}

// Update wraps (*bolt.DB).Update()
func (db *DB) Update(fn func(*Tx) error) error {
	return db.raw.Update(func(rawTx *bolt.Tx) error {
		tx := &Tx{
			raw: rawTx,
		}

		return fn(tx)
	})
}

// View wraps (*bolt.DB).View()
func (db *DB) View(fn func(*Tx) error) error {
	return db.raw.View(func(rawTx *bolt.Tx) error {
		tx := &Tx{
			raw: rawTx,
		}

		return fn(tx)
	})
}

// Unimplemented DB methods
// func (db *DB) Batch(fn func(*Tx) error) error
// func (db *DB) Begin(writable bool) (*Tx, error)
// func (db *DB) GoString() string
// func (db *DB) Info() *Info
// func (db *DB) IsReadOnly() bool
// func (db *DB) Path() string
// func (db *DB) Stats() Stats
// func (db *DB) String() string
// func (db *DB) Sync() error

// EnsureBucket is an addition to the API which will ensure a bucket exists and
// return an error if it doesn't exist and fails to be created.
func (db *DB) EnsureBucket(name string) error {
	return db.Update(func(tx *Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})
}

// Raw will return a reference to the backing *bolt.DB
func (db *DB) Raw() *bolt.DB {
	return db.raw
}
