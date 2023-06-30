package utils

import (
	"fmt"
	"github.com/gookit/goutil/fsutil"
	"github.com/kmou424/filewasher/types"
	"github.com/ohanakogo/exceptiongo"
	"path/filepath"
)

func GetNewFilePath(filename string) string {
	extName := filepath.Ext(filename)
	newFilename := fmt.Sprintf(
		"%s-filewasher%s",
		filename[:len(filename)-len(extName)],
		extName,
	)
	return newFilename
}

func PatchFile(filename string, data []byte) {
	file, err := fsutil.OpenAppendFile(filename)
	exceptiongo.QuickThrow[types.IOException](err)

	writeLen, err := fsutil.WriteOSFile(file, data)
	exceptiongo.QuickThrow[types.IOException](err)

	if len(data) != writeLen {
		exceptiongo.QuickThrowMsg[types.IOException](fmt.Sprintf(
			"expected to write %d bytes but actual count is %d",
			len(data),
			writeLen,
		))
	}
}
