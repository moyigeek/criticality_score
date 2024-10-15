/*
 * @Date: 2024-08-31 03:57:05
 * @LastEditTime: 2024-09-29 15:55:23
 * @Description:
 */
package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/plumbing"
	"github.com/go-git/go-git/plumbing/format/pktline"
	"github.com/go-git/go-git/plumbing/transport"
	gogit "github.com/go-git/go-git/v5"
)

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		// logger.Warnf("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}
	// logger.Warnf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	// logger.Infof("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	// logger.Warnf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
	//fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
	fmt.Fprintf(os.Stderr, "\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func GetStdInput() string {
	var input string
	fmt.Scanln(&input)
	return input
}

func HandleErr(err error, u string) error {
	if err == gogit.ErrRemoteNotFound {
		Warning("[!] %s Not Found!", u)
		err = nil
	}
	if err == pktline.ErrInvalidPktLen {
		Warning("[!] Wrong URL: %s!", u)
		err = nil
	}
	if err == transport.ErrAuthorizationFailed {
		Warning("[!] %s Authorization Failed!", u)
		err = nil
	}
	if err == gogit.ErrNonFastForwardUpdate {
		Warning("[!] %s non-fast-forward Update!", u)
		err = nil
	}
	if err == plumbing.ErrObjectNotFound {
		Warning("[!] %s Object Not Found!", u)
		err = nil
	}
	if err == transport.ErrEmptyRemoteRepository {
		Warning("[!] Repo %s is empty!", u)
		err = nil
	}
	if err == gogit.ErrUnstagedChanges {
		Warning("[!] %s Work tree contains unstaged changes!", u)
		err = nil
	}
	return err
}
