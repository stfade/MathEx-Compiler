package errs

import "os"

func ErrCheck(err error) {
	if err != nil {
		println("[!] ERROR:", err.Error())
		os.Exit(1)
	}
}
