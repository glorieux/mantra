package todo

import "github.com/sirupsen/logrus"

// NotImplemented prints a mesage specifing
// that a given function is not yet implemented
func NotImplemented(name string) {
	logrus.Warnf("%s is not implemented!", name)
}
