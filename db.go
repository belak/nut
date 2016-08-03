package nuts

import (
	"os"

	"github.com/boltdb/bolt"
)

type DB struct {
	// We want to keep this mostly internal, but let the user get to
	// it if they absolutely have to.
	raw *bolt.DB
}

func Open(path string, mode os.FileMode) (*DB, error) {
	var err error

	db := &DB{}
	db.raw, err = bolt.Open(path, mode, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Raw() *bolt.DB {
	return db.raw
}

func (db *DB) Update(fn func(*Tx) error) error {
	return db.raw.Update(func(rawTx *bolt.Tx) error {
		tx := &Tx{
			raw: rawTx,
		}

		return fn(tx)
	})
}

func (db *DB) View(fn func(*Tx) error) error {
	return db.raw.View(func(rawTx *bolt.Tx) error {
		tx := &Tx{
			raw: rawTx,
		}

		return fn(tx)
	})
}

// func (db *DB) Batch(fn func(*Tx) error) error
// func (db *DB) Begin(writable bool) (*Tx, error)
// func (db *DB) Close() error
// func (db *DB) GoString() string
// func (db *DB) Info() *Info
// func (db *DB) IsReadOnly() bool
// func (db *DB) Path() string
// func (db *DB) Stats() Stats
// func (db *DB) String() string
// func (db *DB) Sync() error

func (db *DB) EnsureBucket(name string) error {
	return db.Update(func(tx *Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})
}
