package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/JohanVong/fiscal_v3_sample/internal/services/graylog"
	migrate "github.com/rubenv/sql-migrate"
)

var mssqlDbUser = os.Getenv("FISCAL_DB_USER")
var mssqlDbPassword = os.Getenv("FISCAL_DB_PASSWORD")
var mssqlDbHost = os.Getenv("FISCAL_DB_HOST")
var mssqlDbPort = os.Getenv("FISCAL_DB_PORT")
var mssqlDbName = os.Getenv("FISCAL_DB_NAME")

func parseMigrations(path string) (err error) {
	source := migrate.FileMigrationSource{Dir: path}
	migrations, err := source.FindMigrations()
	m := graylog.NewGELFMessage("Test")
	builder := strings.Builder{}
	for _, mig := range migrations {
		builder.WriteString(fmt.Sprintf("%#v", mig))
	}
	m.SendMessage = builder.String()
	graylog.SendMessage(m)
	return
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

type migrationInfo struct {
	version   int
	dir       string
	migration migrate.MigrationSource
}

type migrationInfoCollection []migrationInfo

func (c migrationInfoCollection) Len() int { return len(c) }
func (c migrationInfoCollection) Less(i, j int) bool {
	if c[i].version < c[j].version {
		return true
	}
	return false
}

func (c migrationInfoCollection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func LaunchMigration() (err error) {
	query := url.Values{}
	query.Add("database", mssqlDbName)
	dbUrl := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(mssqlDbUser, mssqlDbPassword),
		Host:     fmt.Sprintf("%s:%s", mssqlDbHost, mssqlDbPort),
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("mssql", dbUrl.String())
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	migrationsDir := filepath.Join("migrations")
	version := 0

	row := db.QueryRow("select [value] from system_config with(nolock) where [option]='db_version'")

	err = row.Scan(&version)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("DB version before migrations are applied: %d", version))

	infos, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return err
	}
	var versions []int
	for _, i := range infos {
		if i.IsDir() == false {
			tokens := strings.Split(i.Name(), "_")
			if len(tokens) > 0 {
				if tokens[0] == ".DS" { // для запуска на макоси
					continue
				}
				v, err := strconv.Atoi(tokens[0])
				if err != nil {
					return err
				}
				versions = append(versions, v)
			} else {
				return errors.New(fmt.Sprintf("migration file's name is in wrong format: %s, has to be in the form <version>_<description>.sql", i.Name()))
			}
		}

	}
	sort.Ints(versions)
	maxVersion := 0
	if len(versions) > 0 {
		maxVersion = versions[len(versions)-1]
	}

	log.Println(fmt.Sprintf("Target DB version after migrations are applied: %d", maxVersion))
	rec, err := migrate.GetMigrationRecords(db, "mssql")
	err = parseMigrations(migrationsDir)
	if err != nil {
		return err
	}

	source := migrate.FileMigrationSource{Dir: "migrations"}
	timeStart := time.Now()
	_, offset := timeStart.Zone()
	n, err := migrate.Exec(db, "mssql", source, migrate.Up)
	if err != nil {
		return err
	}

	rec, err = migrate.GetMigrationRecords(db, "mssql")
	if n > 0 {
		log.Println("Following migrations have been applied:")
		for _, r := range rec {
			after := r.AppliedAt.After(timeStart.Add(time.Second * time.Duration(offset)))
			if after {
				log.Println(r.Id)
			}

		}
	}
	if len(rec) > 0 {
		lastRecord := rec[len(rec)-1]
		tokens := strings.Split(lastRecord.Id, "_")
		var newVersion int
		if len(tokens) > 0 {
			newVersion, err = strconv.Atoi(tokens[0])
			if err != nil {
				return err
			}
		}
		log.Println(newVersion)
		_, err = db.Query("update system_config set value=? where [option]=?", newVersion, "db_version")
		if err != nil {
			return err
		}
	}

	log.Println("Total number of migrations applied: ", n)
	return nil
}
