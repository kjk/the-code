package main

// To run:
// go run main.g

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

const (
	dockerStatusExited  = "exited"
	dockerStatusRunning = "running"
)

var (
	// using https://hub.docker.com/_/mysql/
	// to use the latest mysql, use mysql:8
	dockerImageName = "mysql:5.6"
	// name must be unique across containers runing on this computer
	dockerContainerName = "mysql-db-multi"
	// where mysql stores databases. Must be on local disk so that
	// database outlives the container
	dockerDbDir = "~/data/db-multi"
	// 3306 is standard MySQL port, I use a unique port to be able
	// to run multiple mysql instances for different projects
	dockerDbLocalPort = "7200"
)

type containerInfo struct {
	id       string
	name     string
	mappings string
	status   string
}

func quoteIfNeeded(s string) string {
	if strings.Contains(s, " ") || strings.Contains(s, "\"") {
		s = strings.Replace(s, `"`, `\"`, -1)
		return `"` + s + `"`
	}
	return s
}

func cmdString(cmd *exec.Cmd) string {
	n := len(cmd.Args)
	a := make([]string, n, n)
	for i := 0; i < n; i++ {
		a[i] = quoteIfNeeded(cmd.Args[i])
	}
	return strings.Join(a, " ")
}

func runCmdWithLogging(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func decodeContainerStaus(status string) string {
	// Convert "Exited (0) 2 days ago" into statusExited
	if strings.HasPrefix(status, "Exited") {
		return dockerStatusExited
	}
	// convert "Up <time>" into statusRunning
	if strings.HasPrefix(status, "Up ") {
		return dockerStatusRunning
	}
	return strings.ToLower(status)
}

// given:
// 0.0.0.0:7200->3306/tcp
// return (0.0.0.0, 7200) or None if doesn't match
func decodeIPPortMust(mappings string) (string, string) {
	parts := strings.Split(mappings, "->")
	panicIf(len(parts) != 2, "invalid mappings string: '%s'", mappings)
	parts = strings.Split(parts[0], ":")
	panicIf(len(parts) != 2, "invalid mappints string: '%s'", mappings)
	return parts[0], parts[1]
}

func dockerContainerInfoMust(containerName string) *containerInfo {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}|{{.Status}}|{{.Ports}}|{{.Names}}")
	outBytes, err := cmd.CombinedOutput()
	panicIfErr(err, "cmd.CombinedOutput() for '%s' failed with %s", cmdString(cmd), err)
	s := string(outBytes)
	// this returns a line like:
	// 6c5a934e00fb|Exited (0) 3 months ago|0.0.0.0:7200->3306/tcp|mysql-db-multi
	s = strings.TrimSpace(s)
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		parts := strings.Split(line, "|")
		panicIf(len(parts) != 4, "Unexpected output from docker ps:\n%s\n. Expected 4 parts, got %d (%v)\n", line, len(parts), parts)
		id, status, mappings, name := parts[0], parts[1], parts[2], parts[3]
		if containerName == name {
			return &containerInfo{
				id:       id,
				status:   decodeContainerStaus(status),
				mappings: mappings,
				name:     name,
			}
		}
	}
	return nil
}

// returns host and port on which database accepts connection
func startLocalDockerDbMust() (string, string) {
	// docker must be running
	cmd := exec.Command("docker", "ps")
	err := cmd.Run()
	panicIfErr(err, "docker must be running! Error: %s", err)
	// ensure directory for database files exists
	dbDir := expandTildeInPath(dockerDbDir)
	err = os.MkdirAll(dbDir, 0755)
	panicIfErr(err, "failed to create dir '%s'. Error: %s", err)
	info := dockerContainerInfoMust(dockerContainerName)
	if info != nil && info.status == dockerStatusRunning {
		return decodeIPPortMust(info.mappings)
	}
	// start or resume container
	if info == nil {
		// start new container
		volumeMapping := dockerDbDir + "s:/var/lib/mysql"
		dockerPortMapping := dockerDbLocalPort + ":3306"
		cmd = exec.Command("docker", "run", "-d", "--name"+dockerContainerName, "-p", dockerPortMapping, "-v", volumeMapping, "-e", "MYSQL_ALLOW_EMPTY_PASSWORD=yes", "-e", "MYSQL_INITDB_SKIP_TZINFO=yes", dockerImageName)
	} else {
		// start stopped container
		cmd = exec.Command("docker", "start", info.id)
	}
	runCmdWithLogging(cmd)

	// wait max 8 seconds for the container to start
	for i := 0; i < 8; i++ {
		info := dockerContainerInfoMust(dockerContainerName)
		if info != nil && info.status == dockerStatusRunning {
			return decodeIPPortMust(info.mappings)
		}
		time.Sleep(time.Second)
	}

	panicIf(true, "docker container '%s' didn't start in time", dockerContainerName)
	return "", ""
}

// helper functions

func fmtArgs(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	format := args[0].(string)
	if len(args) == 1 {
		return format
	}
	return fmt.Sprintf(format, args[1:]...)
}

func panicIfErr(err error, args ...interface{}) {
	if err == nil {
		return
	}
	s := fmtArgs(args...)
	if s == "" {
		s = err.Error()
	}
	panic(s)
}

func panicIf(cond bool, args ...interface{}) {
	if !cond {
		return
	}
	s := fmtArgs(args...)
	if s == "" {
		s = "fatalIf: cond is false"
	}
	panic(s)
}

// userHomeDir returns $HOME diretory of the user
func userHomeDir() string {
	// user.Current() returns nil if cross-compiled e.g. on mac for linux
	if usr, _ := user.Current(); usr != nil {
		return usr.HomeDir
	}
	return os.Getenv("HOME")
}

// expandTildeInPath converts ~ to $HOME
func expandTildeInPath(s string) string {
	if strings.HasPrefix(s, "~") {
		return userHomeDir() + s[1:]
	}
	return s
}

func main() {
	ipAddr, port := startLocalDockerDbMust()
	fmt.Printf("mysql is running insider docker, connect to ip: %s, port: %s\n", ipAddr, port)
	// now connect to the database
}
