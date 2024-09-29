<!--
 * @Author: 7erry
 * @Date: 2024-09-29 17:40:53
 * @LastEditTime: 2024-09-29 17:50:14
 * @Description: 
-->
# Doc

// Todo. Detailed api info need to be filled in

## Collector

package collector // import "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"

func BriefClone(u *url.RepoURL) (*gogit.Repository, error)
func BriefCollect(u *url.RepoURL) (*gogit.Repository, error)
func Clone(u *url.RepoURL) (*gogit.Repository, error)
func Collect(u *url.RepoURL) (*gogit.Repository, error)
func EzCollect(u *url.RepoURL) (*gogit.Repository, error)
func MemClone(u *url.RepoURL) (*gogit.Repository, error)
func Open(path string) (*gogit.Repository, error)
func Update(u *url.RepoURL) (*gogit.Repository, error)

## Database

package database // import "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database"

const SQLITE_MAX_IDLE_CONNS int = 10 ...
type Metrics struct{ ... }
    func NewMetrics(Name, Owner, Source, URL string, CreatedSince, UpdatedSince time.Time, ...) Metrics

### Database/Psql

    package psql // import "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database/psql"

* @Date: 2024-09-07 16:30:21 * @LastEditTime: 2024-09-29 14:29:34 *
@Description:

func BatchInsertMetrics(db *gorm.DB, metrics [database.BATCH_SIZE]database.Metrics) error
func CreateTable(db *gorm.DB)
func InitDB() *gorm.DB
func InsertTable(db *gorm.DB, metrics *database.Metrics)

## Parser

package parser // import "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser"

const UNKNOWN_URL string = "Unknown URL" ...
var BEGIN_TIME = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC) ...
var LANGUAGE_EXTENSIONS = map[string]string{ ... }
var LANGUAGE_FILENAMES = map[string]string{ ... }

### Parser/git

package git // import "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/git"

func GetBlobs(r *git.Repository) *[]*object.Blob
func GetBranches(r *git.Repository) *[]*plumbing.Reference
func GetCommits(r *git.Repository) *[]*object.Commit
func GetConfig(r *git.Repository) *gitconfig.Config
func GetLanguages(r *git.Repository) *[]string
func GetObjects(r *git.Repository) *[]*object.Object
func GetReferences(r *git.Repository) *[]*plumbing.Reference
func GetRemotes(r *git.Repository) *[]*git.Remote
func GetTagRefs(r *git.Repository) *[]*plumbing.Reference
func GetTags(r *git.Repository) *[]*object.Tag
func GetTrees(r *git.Repository) *[]*object.Tree
func GetURL(r *git.Repository) string
func GetWorkTree(r *git.Repository) *git.Worktree
type Repo struct{ ... }
    func NewRepo() Repo
    func ParseGitRepo(r *git.Repository) *Repo
type RepoMetrics struct{ ... }
    func GetMetrics(r *git.Repository) *RepoMetrics
    func NewRepoMetrics() RepoMetrics

### Parser/url

package url // import "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"

func IsSsh(input string) bool
func Protocols(input string) []string
type RepoURL struct{ ... }
    func NewRepoURL(url string) RepoURL
    func ParseURL(input string) RepoURL

## Utils

package utils // import "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"

func CheckArgs(arg ...string)
func CheckIfError(err error)
func GetStdInput() string
func HandleErr(err error, u string) error
func Info(format string, args ...interface{})
func Warning(format string, args ...interface{})
