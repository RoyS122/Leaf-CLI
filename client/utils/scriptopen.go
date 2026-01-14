package utils

import "os/exec"

func OpenLuaExternal(editor, filePath string) {
	exec.Command(editor, filePath).Start()
}
