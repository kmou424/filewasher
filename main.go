package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/joho/godotenv"
	"github.com/kmou424/filewasher/consts"
	"github.com/kmou424/filewasher/logger"
	"github.com/kmou424/filewasher/types"
	"github.com/kmou424/filewasher/utils"
	"github.com/ohanakogo/exceptiongo"
	"github.com/ohanakogo/exceptiongo/pkg/etype"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	initial bool
)

func parse() {
	flag.BoolVar(&initial, "init", false, "initial .env for app")
	flag.Parse()
}

//go:embed resource
var resource embed.FS

func initEnv() {
	bytes, err := resource.ReadFile("resource/filewasher.default.env")
	exceptiongo.QuickThrow[types.IOException](err)

	path := filepath.Join(sysutil.Workdir(), "filewasher.env")
	err = fsutil.WriteFile(path, bytes, 0644)
	exceptiongo.QuickThrow[types.IOException](err)

	logger.Logf("%s has been created", path)
}

func run() {
	err := godotenv.Load("filewasher.env")
	if err != nil {
		logger.Logf(`you must execute "-init" at first`)
		return
	}

	scanExtensionsStr := utils.RequireEnv(consts.EnvScanExtensions)
	if strutil.Trim(scanExtensionsStr) == "" {
		exceptiongo.QuickThrowMsg[types.EnvVarNotCompatibleException](fmt.Sprintf(
			`env "%s" can't be empty'`,
			consts.EnvScanExtensions),
		)
	}
	scanExtensions := strutil.Split(scanExtensionsStr, "|")
	for i := 0; i < len(scanExtensions); i++ {
		for {
			if strutil.HasPrefix(scanExtensions[i], ".") {
				scanExtensions[i] = scanExtensions[i][1:]
				continue
			}
			break
		}
		scanExtensions[i] = fmt.Sprintf(".%s", scanExtensions[i])
	}
	patchMode := utils.RequireEnv(consts.EnvPatchMode)

	logger.Logf("scanning for: %s", strutil.JoinList(", ", scanExtensions))

	err = fsutil.WalkDir(sysutil.Workdir(), func(path string, d fs.DirEntry, err error) error {
		if !strutil.HasOneSuffix(path, scanExtensions) {
			return nil
		}
		relPath, err := filepath.Rel(sysutil.Workdir(), path)
		exceptiongo.QuickThrow[types.IOException](err)

		fmt.Println()
		logger.Logf("%s", relPath)
		if !utils.CheckFileWasher(path) {
			patchFooter := utils.GenFileFooter()

			newFilePath := utils.GetNewFilePath(path)
			newRelPath, err := filepath.Rel(sysutil.Workdir(), newFilePath)
			exceptiongo.QuickThrow[types.IOException](err)

			if fsutil.FileExists(newFilePath) {
				if utils.CheckFileWasher(newFilePath) {
					logger.Logf("detected patched file: %s, skipping", newRelPath)
					return nil
				} else {
					fsutil.MustRemove(newFilePath)
				}
			}

			switch patchMode {
			case "patch":
				err := os.Rename(path, newFilePath)
				exceptiongo.QuickThrow[types.IOException](err)
			case "copy":
				err = fsutil.CopyFile(path, newFilePath)
				exceptiongo.QuickThrow[types.IOException](err)
			default:
				exceptiongo.QuickThrowMsg[types.EnvVarNotCompatibleException](fmt.Sprintf(
					`%s="%s" is not compatible`,
					consts.EnvPatchMode,
					patchMode,
				))
			}

			logger.Logf("patching file %s", relPath)
			utils.PatchFile(newFilePath, patchFooter)

			logger.Logf("patched file has been saved to %s", newRelPath)
		} else {
			logger.Logf("file has patched, skipping")
		}
		return nil
	})
	exceptiongo.QuickThrow[types.IOException](err)
}

func main() {
	defer exceptiongo.NewExceptionHandler(func(e *etype.Exception) {
		e.PrintStackTrace()
	}).Deploy()

	parse()

	switch {
	case initial:
		initEnv()
	default:
		run()
	}
}
