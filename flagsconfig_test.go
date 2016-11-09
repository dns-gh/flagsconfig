package flagsconfig

import (
	"flag"
	"github.com/dns-gh/gotest"
	"os"
	"path/filepath"
	"testing"
)

const (
	configFileName       = "test.config"
	flagTest             = "flag-test"
	flagTestValue        = "test"
	flagTestDescription  = "configuration test flag"
	flagTest2            = "flag-test2"
	flagTestValue2       = "test2"
	flagTestDescription2 = "configuration test flag 2"
	flagTest3            = "flag-test3"
	flagTestValue3       = "test3"
	flagTestDescription3 = "configuration test flag 3"
)

func makeConfigFile(dir string) string {
	return filepath.Join(dir, configFileName)
}

func checkFileInfo(t *testing.T, file string) {
	info, err := os.Stat(file)
	gotest.Assert(t, err)
	gotest.Check(t, info.Name() == configFileName)
	gotest.Check(t, info.Size() != 0)
	gotest.Check(t, !info.IsDir())
}

func checkFlag(t *testing.T, flagTest, flagTestValue string) {
	testFlag := flag.Lookup(flagTest)
	gotest.Check(t, testFlag.Value.String() == flagTestValue)
}

func TestFlagsConfig(t *testing.T) {
	defer gotest.RemoveTestFolder(t)
	flag.String(flagTest, "", flagTestDescription)

	dir := gotest.MakeUniqueTestFolder(t)
	file := makeConfigFile(dir)
	config, err := NewConfig(file)
	gotest.Assert(t, err)
	checkFileInfo(t, file)
	checkFlag(t, flagTest, "")

	err = config.Update(flagTest, flagTestValue)
	gotest.Assert(t, err)
	checkFlag(t, flagTest, "")

	err = config.Parse(file)
	gotest.Assert(t, err)
	checkFlag(t, flagTest, flagTestValue)
}

func TestFlagsConfigFiltered(t *testing.T) {
	defer gotest.RemoveTestFolder(t)
	flag.String(flagTest2, "", flagTestDescription2)
	flag.String(flagTest3, "", flagTestDescription3)

	dir := gotest.MakeUniqueTestFolder(t)
	file := makeConfigFile(dir)
	config, err := NewConfig(file, flagTest3)
	gotest.Assert(t, err)
	checkFileInfo(t, file)

	checkFlag(t, flagTest2, "")
	checkFlag(t, flagTest3, "")

	err = config.Update(flagTest2, flagTestValue2)
	gotest.Assert(t, err)

	checkFlag(t, flagTest2, "")
	checkFlag(t, flagTest3, "")

	err = config.Parse(file)
	gotest.Assert(t, err)

	checkFlag(t, flagTest2, flagTestValue2)
	checkFlag(t, flagTest3, "")
}
