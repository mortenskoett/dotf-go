package terminalio

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mortenskoett/dotf-go/pkg/test"
)

func Test_backupFile_saves_file(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	fileToBackup := env.UserspaceDir.AddTempFile().Path
	expectedBackupPath := filepath.Join("/tmp/dotf-go/backups", fileToBackup)

	actual, err := backupFile(fileToBackup)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(expectedBackupPath)

	if expectedBackupPath != actual {
		test.Fail(actual, expectedBackupPath, t)
	}
}

func Test_copyFile_copies_file(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	todir := env.BackupDir.Path
	fileToMove := env.UserspaceDir.AddTempFile().Path
	dstFilename := "dstFileName"

	expectedPath := filepath.Join(todir, dstFilename)

	actualpath, err := copyFile(fileToMove, expectedPath)
	if err != nil {
		t.Fatal(err)
	}

	// compare returned path
	if expectedPath != actualpath {
		test.Fail(actualpath, expectedPath, t)
	}

	// check if file exists
	actualStat, err := os.Stat(actualpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			test.Fail(err, expectedPath, t)
		}
	}

	expectedStat, err := os.Stat(expectedPath)
	if err != nil {
		test.Fail(err, expectedPath, t)
	}

	// check if permissions are the same
	if actualStat.Mode() != expectedStat.Mode() {
		test.FailMsg("file mode not identical", actualStat, expectedStat, t)
	}
}

func Test_replacePrefixPath_replaces_prefix_of_path(t *testing.T) {
	file := "/dir1/dir2/file.txt"
	from := "/userdir"
	to := "/dotfiles"

	result, err := replacePrefixPath(from+file, from, to)
	if err != nil {
		test.Fail(err, "Shouldn't fail here", t)
	}

	expected := filepath.Join(to, file)
	if result != expected {
		test.Fail(result, expected, t)
	}
}

func Test_replacePrefixPath_using_test_env(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	fromdir := env.DotfilesDir
	subfolderPath := "/bla1/bla2/"
	subfolders := fromdir.AddTempDir(subfolderPath)
	fp := subfolders.AddTempFile()
	todir := env.UserspaceDir
	expected := filepath.Join(todir.Path, subfolderPath, filepath.Base(fp.Path))

	result, err := replacePrefixPath(fp.Path, fromdir.Path, todir.Path)
	if err != nil {
		test.Fail(result, err, t)
	}

	if result != expected {
		test.Fail(result, expected, t)
	}
}

func Test_trimBasePath_removes_shared_prefix(t *testing.T) {
	df := "/dotfiles/d1/d2/d3/"
	bp := "/d1/d2/d3/"
	fp := bp + "file.txt"

	p, err := trimBasePath(fp, df)

	if err != nil {
		test.Fail(err, "Shouldn't fail here", t)
	}

	if p != fp {
		test.Fail(p, bp, t)
	}
}

func Test_trimBasePath_using_test_env(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	// example
	// dotfiles/d1/d2/file.txt
	// detach(dotfiles/, dotfiles/d1/d2/file.txt)
	// returns d1/d2/file.txt

	somedir := env.DotfilesDir.AddTempDir("/dotfiles/")
	basepath := somedir.AddTempDir("/bla1/bla2/")
	f := basepath.AddTempFile()

	p, err := trimBasePath(f.Path, basepath.Path)
	if err != nil {
		test.Fail(err, "Should not fail here", t)
	}

	// Because result has leading slash
	expected := "/" + filepath.Base(f.Path)

	// Check filename
	if p != expected {
		test.Fail(p, expected, t)
	}
}

func Test_getAbsolutePath_returns_equal_abs_path(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	f := env.UserspaceDir.AddTempFile()

	actual, err := getAbsolutePath(f.Path)

	// Check error
	if err != nil {
		test.Fail(err, "Should not fail here", t)
	}

	expected := filepath.Join(f.Path)

	// Check path -- should return the same path
	if actual != expected {
		test.Fail(actual, expected, t)
	}
}

func Test_GetAndValidateAbsolutePath_fails_if_path_not_exist(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	f := "myrandomfile"

	f, err := GetAndValidateAbsolutePath(f)
	if err == nil {
		test.Fail(err, "Should fail here as file does not exist.", t)
	}
}

func Test_checkIfFileExists_determines_files_exist(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	file := env.UserspaceDir.AddTempFile()

	if exists, _ := checkIfFileExists(file.Path); !exists {
		test.Fail(exists, "Should not fail as file exists.", t)
	}
}

func Test_checkIfFileExists_determines_dirs_exist(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	dir := env.UserspaceDir.AddTempDir("mytestdir")

	if exists, _ := checkIfFileExists(dir.Path); !exists {
		test.Fail(exists, "Should not fail as file exists.", t)
	}
}

func Test_copyDir_copies_recursively_files_folders(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	/* Dir structure example:
	/tmp/userspace-dir1205530807
	/tmp/userspace-dir1205530807/inside
	/tmp/userspace-dir1205530807/inside/file-1561196035
	/tmp/userspace-dir1205530807/inside/here
	/tmp/userspace-dir1205530807/inside/here/file-4056809082
	/tmp/userspace-dir1205530807/inside/here/is
	/tmp/userspace-dir1205530807/inside/here/is/nice
	*/
	src := env.UserspaceDir
	inside := src.AddTempDir("inside")
	insideFile := inside.AddTempFile()
	here := inside.AddTempDir("here")
	hereFile := here.AddTempFile()
	here.AddTempDir("is/nice")

	dst := env.BackupDir

	res, err := copyDir(src.Path, dst.Path)
	if err != nil {
		test.Fail(err, "Should not fail here", t)
	}

	dstInsideFile, err := replacePrefixPath(insideFile.Path, src.Path, dst.Path)
	if err != nil {
		test.Fail(err, "Should not fail here", t)
	}

	dstHereFile, err := replacePrefixPath(hereFile.Path, src.Path, dst.Path)
	if err != nil {
	}

	// Dst folder and content should exist
	if exists, _ := checkIfFileExists(res); !exists {
		test.Fail(res, "File does not exist", t)
	}
	if exists, _ := checkIfFileExists(dstInsideFile); !exists {
		test.Fail(dstInsideFile, "File does not exist", t)
	}
	if exists, _ := checkIfFileExists(dstHereFile); !exists {
		test.Fail(dstHereFile, "File does not exist", t)
	}
	// Src folder should NOT have been deleted
	if exists, _ := checkIfFileExists(src.Path); !exists {
		test.Fail(src.Path, "Dir should still exist", t)
	}
}

func Test_deleteDirectory_deletes_existing_directory(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	d := env.UserspaceDir.AddTempDir("dirtodelete").Path

	// check if dir exists
	if exists, _ := checkIfFileExists(d); !exists {
		test.Fail(exists, "dir should exist", t)
	}

	// check if what is created is considered dir
	ok, err := isDirectory(d)
	if err != nil {
		test.Fail(err, "shouldnt fail", t)
	}
	if !ok {
		test.Fail(ok, "this should be a directory", t)
	}

	// delete dir
	err = deleteDirectory(d)
	if err != nil {
		test.Fail(err, "shouldnt fail", t)
	}

	// check if deleted properly
	if exists, _ := checkIfFileExists(d); exists {
		test.Fail(exists, "dir should NOT exist now", t)
	}
}

func Test_deleteFile_deletes_existing_file(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	f := env.UserspaceDir.AddTempFile().Path

	if exists, _ := checkIfFileExists(f); !exists {
		test.Fail(exists, "file should exist", t)
	}

	err := deleteFile(f)
	if err != nil {
		test.Fail(err, "shouldnt fail", t)
	}

	if exists, _ := checkIfFileExists(f); exists {
		test.Fail(exists, "file should NOT exist now", t)
	}
}

func Test_getFileLocationInfo_returns_correct_fileinfo(t *testing.T) {
	env := test.NewTestEnvironment()
	defer env.Cleanup()

	type testcase struct {
		file   string
		uspace string
		dfiles string
		want   fileLocationInfo
	}

	uspace := env.UserspaceDir
	dfiles := env.DotfilesDir

	file1 := dfiles.AddTempFile()
	file2 := uspace.AddTempFile()

	cases := []testcase{
		{
			file:   file1.Path,
			want: fileLocationInfo{
				insideDotfiles: true,
				fileOrgPath:    file1.Path,
				userspaceFile:  filepath.Join(uspace.Path, file1.Name),
				dotfilesFile:   filepath.Join(dfiles.Path, file1.Name),
			},
		},
		{
			file:   file2.Path,
			want: fileLocationInfo{
				insideDotfiles: false,
				fileOrgPath:    file2.Path,
				userspaceFile:  filepath.Join(uspace.Path, file2.Name),
				dotfilesFile:   filepath.Join(dfiles.Path, file2.Name),
			},
		},
	}

	for _, tc := range cases {
		result, err := getFileLocationInfo(tc.file, uspace.Path, dfiles.Path)
		if err != nil {
			test.Fail(err, "shouldnt fail", t)
		}

		if result.fileOrgPath != tc.want.fileOrgPath {
			test.FailMsg("org file path", result.fileOrgPath, tc.want.fileOrgPath, t)
		}

		if result.dotfilesFile != tc.want.dotfilesFile {
			test.FailMsg("Path in dotfiles dir", result.dotfilesFile, tc.want.dotfilesFile, t)
		}

		if result.userspaceFile != tc.want.userspaceFile {
			test.FailMsg("Path in userspace", result.userspaceFile, tc.want.userspaceFile, t)
		}

		if result.insideDotfiles != tc.want.insideDotfiles {
			test.FailMsg("inside dotfiles", result.insideDotfiles, tc.want.insideDotfiles, t)
		}
	}
}
