package js

import "local/src/util"

type foldersAndFiles struct {
	folders      []string
	files        []util.FileContent
	dependencies []util.Dependencies
	suggestions  []*util.Suggestion
	espath       string
}

func initFoldersAndFiles(folders []string, files []util.FileContent) foldersAndFiles {
	var ff foldersAndFiles
	ff.folders = folders
	ff.files = files

	return ff
}

func (ff *foldersAndFiles) _addFiles(files ...util.FileContent) {
	ff.files = append(ff.files, files...)
}

func (ff *foldersAndFiles) _addFolders(folders ...string) {
	ff.folders = append(ff.folders, folders...)
}

func (ff *foldersAndFiles) _addSuggestion(sug ...*util.Suggestion) {
	ff.suggestions = append(ff.suggestions, sug...)
}

func (ff *foldersAndFiles) _addDependencies(dependencies ...util.Dependencies) {
	ff.dependencies = append(ff.dependencies, dependencies...)
}
