package e2e_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RootPath(relevantPath string) string {
	root, _ := os.Getwd()

	if strings.Contains(root, "/tests/e2e") {
		// Get the path before "/test/e2e"
		index := strings.LastIndex(root, "/tests/e2e")
		pathBefore := root[:index]

		root = pathBefore
	}

	return root + "/" + strings.Trim(relevantPath, "/")
}

type E2E struct {
	t         *testing.T
	cmd       *exec.Cmd
	startedAt time.Time
	endedAt   time.Time
	out       bytes.Buffer
	cancel    context.CancelFunc
	context   context.Context
	timedout  bool
	envs      []string
}

func NewE2E(t *testing.T) *E2E {
	return &E2E{
		t: t,
	}
}

func (e *E2E) WithEnv(env string) {
	e.envs = append(e.envs, env)
}

// Get exit code from command process if binary is used
// or parse the last output line to determine the exit code.
func (e *E2E) GetExitCode() (int, error) {
	if e.GetKermooBinaryPath() != "" {
		return e.cmd.ProcessState.ExitCode(), nil
	}

	re, err := regexp.Compile(`\d+`)
	if err != nil {
		return 0, fmt.Errorf("error compiling regex: %v", err)
	}

	// Find the first match in the string
	lastLine := e.GetLastOutputLine()
	match := re.FindString(lastLine)
	if match == "" {
		return 0, fmt.Errorf("no match found: %s", lastLine)
	}

	// Convert the matched string to an integer
	i, err := strconv.Atoi(match)
	if err != nil {
		return 0, fmt.Errorf("error converting exit code string to int: %v", err)
	}

	return i, nil
}

func (e *E2E) GetFileLogWriter() *os.File {
	logDir := RootPath("e2e-test-results")

	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	name := strings.ReplaceAll(e.t.Name(), "/", "_")

	file, err := os.OpenFile(logDir+"/"+name+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}

	file.Write([]byte("prepared.\n"))

	// io.Writer to append content to the file
	return file
}

func (e *E2E) GetKermooBinaryPath() string {
	return os.Getenv("KERMOO_BINARY")
}

func (e *E2E) Start(config string, timeout time.Duration) {
	if e.GetKermooBinaryPath() != "" {
		e.StartBinary(config, timeout)
	} else {
		e.StartBuild(config, timeout)
	}
}

func (e *E2E) StartBuild(config string, timeout time.Duration) {
	e.context, e.cancel = context.WithTimeout(context.Background(), timeout)

	e.cmd = exec.CommandContext(e.context, "go", "run", RootPath("main.go"), "start", "-v", "debug", "-f", "-")

	e.cmd.Env = os.Environ()
	e.cmd.Env = append(e.cmd.Env, e.envs...)

	// This is a workaround for terminating process with child
	// See: https://stackoverflow.com/a/71714364
	if e.cmd.SysProcAttr == nil {
		e.cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	e.cmd.SysProcAttr.Setpgid = true

	e.cmd.Stdin = strings.NewReader(strings.Trim(config, "\t"))

	mw := io.MultiWriter(os.Stdout, &e.out)

	e.cmd.Stdout = mw
	e.cmd.Stderr = mw

	e.startedAt = time.Now()

	assert.NoError(e.t, e.cmd.Start())
}

func (e *E2E) StartBinary(config string, timeout time.Duration) {
	e.context, e.cancel = context.WithTimeout(context.Background(), timeout)

	var mw io.Writer

	e.cmd = exec.CommandContext(e.context, e.GetKermooBinaryPath(), "start", "-v", "debug", "-f", "-")

	//mw = io.MultiWriter(os.Stdout, e.GetFileLogWriter(), &e.out)
	mw = io.MultiWriter(e.GetFileLogWriter(), &e.out)

	e.cmd.Env = os.Environ()
	e.cmd.Env = append(e.cmd.Env, e.envs...)

	e.cmd.Stdin = strings.NewReader(strings.Trim(config, "\t"))

	e.cmd.Stdout = mw
	e.cmd.Stderr = mw

	e.startedAt = time.Now()

	require.NoError(e.t, e.cmd.Start())
}

func (e *E2E) Wait() {
	defer e.cancel()

	go func() {
		<-e.context.Done()
		if e.cmd.Process == nil {
			return
		}
		// Kill by negative PID to kill the process group, which includes
		// the top-level process we spawned as well as any subprocesses
		// it spawned.
		_ = syscall.Kill(-e.cmd.Process.Pid, syscall.SIGKILL)
		e.timedout = true
	}()

	e.cmd.Wait()
	e.endedAt = time.Now()
}

func (e *E2E) AssertExecutaionDuration(min time.Duration, max time.Duration) {
	assert.GreaterOrEqual(e.t, e.endedAt.Sub(e.startedAt), min)
	assert.LessOrEqual(e.t, e.endedAt.Sub(e.startedAt), max)
}

func (e *E2E) RequireNotTimedOut() {
	require.False(e.t, e.timedout, "command is timedout")
}

func (e *E2E) RequireTimedOut() {
	require.True(e.t, e.timedout, "command is NOT timedout")
}

func (e *E2E) AssertExitCode(expectedCode int) {
	exitCode, err := e.GetExitCode()

	require.Nil(e.t, err)

	assert.Equal(e.t, expectedCode, exitCode)
}

func (e *E2E) GetOutput() string {
	return strings.Trim(e.out.String(), "\n")
}

func (e *E2E) GetLastOutputLine() string {
	lines := strings.Split(e.GetOutput(), "\n")
	return lines[len(lines)-1]
}
