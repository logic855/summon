package provider

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

func TestResolve(t *testing.T) {
	Convey("Passing no provider should return an error", t, func() {
		// Point to a tempdir to avoid pollution from dev env
		tempDir, _ := ioutil.TempDir("", "summontest")
		defer os.RemoveAll(tempDir)
		DefaultPath = tempDir

		_, err := Resolve("")

		So(err, ShouldNotBeNil)
	})

	Convey("Passing the provider via CLI should return it without error", t, func() {
		expected := "/bin/bash"
		provider, err := Resolve(expected)

		So(err, ShouldBeNil)
		So(provider, ShouldEqual, expected)

	})

	Convey("Setting the provider via environment variable works", t, func() {
		expected := "/bin/bash"
		os.Setenv("SUMMON_PROVIDER", expected)
		provider, err := Resolve("")
		os.Unsetenv("SUMMON_PROVIDER")

		So(err, ShouldBeNil)
		So(provider, ShouldEqual, expected)

	})

	Convey("Given a provider path", t, func() {
		tempDir, _ := ioutil.TempDir("", "summontest")
		defer os.RemoveAll(tempDir)
		DefaultPath = tempDir

		Convey("If there is 1 executable, return it as the provider", func() {
			f, err := ioutil.TempFile(DefaultPath, "")
			f.Chmod(755)
			provider, err := Resolve("")

			So(err, ShouldBeNil)
			So(provider, ShouldEqual, f.Name())

		})

		Convey("If there are > 1 executables, return an error to user", func() {
			// Create 2 exes in provider path
			ioutil.TempFile(DefaultPath, "")
			ioutil.TempFile(DefaultPath, "")
			_, err := Resolve("")

			So(err, ShouldNotBeNil)
		})
	})
}

func TestCall(t *testing.T) {
	Convey("When I call a provider", t, func() {
		Convey("If it returns exit code 0, return stdout", func() {
			arg := "provider.go"
			out, err := Call("ls", arg)

			So(out, ShouldEqual, arg)
			So(err, ShouldBeNil)
		})
		Convey("If it returns exit code > 0, return stderr", func() {
			out, err := Call("ls", "README.notafile")

			So(out, ShouldBeBlank)
			So(err.Error(), ShouldContainSubstring, "No such file or directory")
		})
	})
}
