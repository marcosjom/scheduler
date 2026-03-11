package client

import (
	"path/filepath"
	"testing"
)

// Testing some invalid examples
func Test_Client_LoadTasksFromFolder(t *testing.T) {
	folderpath := ".." + string(filepath.Separator) + ".." + string(filepath.Separator) + "test" + string(filepath.Separator) + "tasks"
	client := New()
	if err := client.ReloadTasks(folderpath); err != nil {
		t.Errorf("ReloadTasks failed for '%s': %s.", folderpath, err.Error())
	}
}
