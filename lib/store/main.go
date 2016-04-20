package store

import "github.com/boltdb/bolt"

func OpenDB() (*bolt.DB, error) {

	db, err := bolt.Open(".store/main.db", 0600, nil)
	return db, err

}

func Get(key string, bucket string, cb func(string)) error {

	var Error error

	db, Error := OpenDB()

	if Error == nil {

		defer db.Close()

		Error = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(bucket))

			v := string(b.Get([]byte(key)))

			cb(v)

			return err
		})

	}

	return Error
}
