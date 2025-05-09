package errorUtils

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

const baseFolderName = "/go/src/u9/3032-porsche-native-app-back-end/internal"

func Wrap(err error) (newErr error) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		_, fn, line, _ := runtime.Caller(1)
		newErr = errors.Wrap(err, fmt.Sprintf("[error] in [%s:%d]\n\n", strings.Replace(fn, baseFolderName, "", 1), line))
	}
	return
}
