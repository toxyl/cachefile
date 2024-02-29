package cachefile

import (
	"fmt"
	"os"
	"time"
)

type CacheFile struct {
	file        string
	permissions os.FileMode
	maxAge      time.Duration
	fnLoad      func([]byte) ([]byte, error)
	fnRetrieve  func() ([]byte, error)
	data        []byte
}

func (cf *CacheFile) File() string {
	return cf.file
}

func (cf *CacheFile) Permissions() os.FileMode {
	return cf.permissions
}

func (cf *CacheFile) MaxAge() time.Duration {
	return cf.maxAge
}

func (cf *CacheFile) Expired() bool {
	fileInfo, err := os.Stat(cf.file)
	return err != nil || time.Since(fileInfo.ModTime()) > cf.maxAge
}

func (cf *CacheFile) Data() ([]byte, error) {
	if !cf.Expired() {
		// the file exists and hasn't expired
		if cf.data == nil {
			// looks like we haven't read its contents yet, let's do it then
			data, err := os.ReadFile(cf.file)
			if err != nil {
				return nil, fmt.Errorf("error reading file: %s", err)
			}

			// check if we have a function to process the data with
			if cf.fnLoad != nil {
				data, err = cf.fnLoad(data)
				if err != nil {
					return nil, err
				}
			}

			// store the content in memory
			cf.data = data
		}
		// return the content from memory
		return cf.data, nil
	}

	// the file is missing or has expired, retrieve the data
	bytes, err := cf.fnRetrieve()
	if err != nil {
		return bytes, err
	}
	if err = os.WriteFile(cf.file, bytes, cf.permissions); err != nil {
		return bytes, err
	}
	if cf.fnLoad != nil {
		bytes, err = cf.fnLoad(bytes)
	}
	cf.data = bytes
	return cf.data, err
}

func New(
	file string, permissions os.FileMode, maxAge time.Duration,
	fnRetrieve func() ([]byte, error),
	fnLoad func([]byte) ([]byte, error),
) *CacheFile {
	if fnRetrieve == nil {
		panic("CacheFile needs a retrieval function")
	}

	cf := &CacheFile{
		file:        file,
		permissions: permissions,
		maxAge:      maxAge,
		fnLoad:      fnLoad,
		fnRetrieve:  fnRetrieve,
		data:        nil,
	}

	return cf
}
