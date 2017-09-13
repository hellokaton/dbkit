package main

import (
	"fmt"
	"os/exec"
	"os"
	"time"
)

// export mysql data
func ExportMySql(x DbKitConfig) *ExportResult {

	result := &ExportResult{MIME: "application/x-tar"}
	// local storage path, default is .
	if x.Storage.Location == "" {
		x.Storage.Location = "."
	} else {
		if b, _ := PathExists(x.Storage.Location); !b {
			// mkdir storage location
			os.MkdirAll(x.Storage.Location, 0700)
		}
	}

	dumpPath := fmt.Sprintf(`%v/dbkit_%v_%v.sql`, x.Storage.Location, x.Db.Database, time.Now().Unix())

	options := append(dumpOptions(x), fmt.Sprintf(`-r%v`, dumpPath))

	out, err := exec.Command(MysqlDumpCmd, options...).Output()
	if err != nil {
		result.Error = MakeErr(err, string(out))
		return result
	}

	result.Path = dumpPath + ".tar.gz"
	_, err = exec.Command(TarCmd, "-czf", result.Path, dumpPath).Output()
	if err != nil {
		result.Error = MakeErr(err, string(out))
		return result
	}
	os.Remove(dumpPath)

	return result
}

// append mysql dump options
func dumpOptions(x DbKitConfig) []string {
	options := x.Db.Options
	options = append(options, fmt.Sprintf(`-h%v`, x.Db.Host))
	options = append(options, fmt.Sprintf(`-P%v`, x.Db.Port))
	options = append(options, fmt.Sprintf(`-u%v`, x.Db.Username))
	if x.Db.Password != "" {
		options = append(options, fmt.Sprintf(`-p%v`, x.Db.Password))
	}
	options = append(options, x.Db.Database)
	return options
}
