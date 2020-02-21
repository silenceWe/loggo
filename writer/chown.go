// +build !linux

package writer

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
