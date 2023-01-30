package base

import "github.com/pkg/errors"

func InvalidMethod(method string) error {
	return errors.Errorf("invalid method: %s", method)
}
