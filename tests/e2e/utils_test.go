package e2e_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
}

func NewE2E(t *testing.T) *E2E {
	return &E2E{
		t: t,
	}
}

func (e *E2E) Start(config string, timeout time.Duration) {
	var ctx context.Context
	ctx, e.cancel = context.WithTimeout(context.Background(), timeout)

	e.cmd = exec.CommandContext(ctx, "go", "run", RootPath("main.go"), "start", "-v", "debug", "-f", "-")

	e.cmd.Stdin = strings.NewReader(strings.Trim(config, "\t"))

	mw := io.MultiWriter(os.Stdout, &e.out)

	e.cmd.Stdout = mw
	e.cmd.Stderr = mw

	e.startedAt = time.Now()

	assert.NoError(e.t, e.cmd.Start())
}

func (e *E2E) Wait() {
	defer e.cancel()
	e.cmd.Wait()
	e.endedAt = time.Now()
}

func (e *E2E) AssertExecutaionDuration(min time.Duration, max time.Duration) {
	assert.GreaterOrEqual(e.t, e.endedAt.Sub(e.startedAt), min)
	assert.LessOrEqual(e.t, e.endedAt.Sub(e.startedAt), max)
}

func (e *E2E) AssertExitCode(code int) {
	assert.Equal(e.t, fmt.Sprintf("exit status %d", code), e.GetLastOutputLine())
}

func (e *E2E) GetLastOutputLine() string {
	lines := strings.Split(strings.Trim(e.out.String(), "\n"), "\n")
	return lines[len(lines)-1]
}
