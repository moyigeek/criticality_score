package enumerator

import "testing"

func Test_enumerateBitbucket(t *testing.T) {
	t.Run("Bitbucket", func(t *testing.T) {
		c := NewEnumerator()
		c.enumerateBitbucket()
	})
}

func Test_enumerateGitlab(t *testing.T) {
	t.Run("Gitlab", func(t *testing.T) {
		c := NewEnumerator()
		c.enumerateGitlab()
	})
}
