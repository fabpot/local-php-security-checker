package security

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLock(t *testing.T) {
	file, err := os.Open("fixtures/no_vulns.lock")
	if err != nil {
		panic(err)
	}
	lock, err := NewLock(bufio.NewReader(file))
	assert.Nil(t, err)
	assert.Equal(t, len(lock.DevPackages), 0)
	assert.Equal(t, len(lock.Packages), 1)
	assert.Equal(t, lock.Packages[0].Name, "symfony/apache-pack")
}

func TestIntegerAsVersionLock(t *testing.T) {
	file, err := os.Open("fixtures/integer_as_version.lock")
	if err != nil {
		panic(err)
	}
	lock, err := NewLock(bufio.NewReader(file))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(lock.DevPackages))
	assert.Equal(t, 1, len(lock.Packages))
	assert.Equal(t, "symfony/apache-pack", lock.Packages[0].Name)
	assert.Equal(t, "7", string(lock.Packages[0].Version))
}

func TestNotALock(t *testing.T) {
	file, err := os.Open("fixtures/not_a_lock.lock")
	if err != nil {
		panic(err)
	}
	_, err = NewLock(bufio.NewReader(file))
	assert.Equal(t, "lock file is not valid (no packages and no dev packages)", err.Error())
}

func TestLocateLock(t *testing.T) {
	for _, path := range []string{"fixtures/locate", "fixtures/locate/composer.json", "fixtures/locate/composer.lock"} {
		_, err := LocateLock(path)
		assert.Nil(t, err)
	}
}

func TestPrereleaseWithoutDot(t *testing.T) {
	file, err := os.Open("fixtures/prerelease_without_dot.lock")
	if err != nil {
		panic(err)
	}
	lock, err := NewLock(bufio.NewReader(file))
	if err != nil {
		panic(err)
	}
	assert.Equal(t, lock.Packages[0].Version, Version("v1.0.0-alpha.10"))
	assert.Equal(t, lock.Packages[1].Version, Version("2.0-beta.3"))
	assert.Equal(t, lock.Packages[2].Version, Version("2.0-RC.1"))
}
