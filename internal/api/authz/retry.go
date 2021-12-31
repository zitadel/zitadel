package authz

import "time"

//TODO: workaround if org projection is not yet up-to-date
func retry(retriable func() error) (err error) {
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		err = retriable()
		if err == nil {
			return nil
		}
	}
	return err
}
