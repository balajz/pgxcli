package cliio

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/charmbracelet/x/term"
	"github.com/fatih/color"
	"github.com/google/shlex"
)

var (
	printErr  = color.New(color.FgHiRed).FprintfFunc()
	printInfo = color.New(color.FgWhite).FprintfFunc()
	printTime = color.New(color.FgHiCyan).FprintfFunc()
)

type Printer interface {
	SetOut(out io.Writer)
	SetErrOut(errOut io.Writer)
	Print(str string)
	PrintError(err error)
	PrintTime(time time.Duration)
	PrintViaPager(str string)
}

type PgxPrinter struct {
	out    io.Writer
	errOut io.Writer
}

func NewPgxPrinter(out io.Writer, errOut io.Writer) *PgxPrinter {
	return &PgxPrinter{out: out, errOut: errOut}
}

func (p *PgxPrinter) SetOut(out io.Writer) {
	p.out = out
}

func (p *PgxPrinter) SetErrOut(errOut io.Writer) {
	p.errOut = errOut
}

func (p *PgxPrinter) Print(str string) {
	printInfo(p.out, str)
}

func (p *PgxPrinter) PrintError(err error) {
	printErr(p.errOut, "%v\n", err)
}

func (p *PgxPrinter) PrintTime(time time.Duration) {
	printTime(p.out, "Time: %.3fs\n", time.Seconds())
}

func (p *PgxPrinter) PrintViaPager(str string) {
	err := EchoViaPager(func(w io.Writer) error {
		_, err := io.WriteString(w, str)
		return err
	})
	if err != nil {
		p.PrintError(err)
	}
}

func EchoViaPager(writeFn func(io.Writer) error) error {
	stdout := os.Stdout
	stdin := os.Stdin

	if !term.IsTerminal(stdin.Fd()) || !term.IsTerminal(stdout.Fd()) {
		return writeFn(stdout)
	}

	pagerCmd := getPager()

	if tryPipePager(pagerCmd, writeFn) {
		return nil
	}

	if tryTempfilePager(pagerCmd, writeFn) {
		return nil
	}

	return writeFn(stdout)
}

func tryPipePager(pagerCmd []string, writeFn func(io.Writer) error) bool {
	cmdPath, err := exec.LookPath(pagerCmd[0])
	if err != nil {
		return false
	}

	cmd := exec.Command(cmdPath, pagerCmd[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return false
	}
	if err := cmd.Start(); err != nil {
		return false
	}

	writeErr := writeFn(stdin)
	_ = stdin.Close()

	waiterr := waitIgnoringInterrupt(cmd)

	return writeErr == nil && waiterr == nil
}

func tryTempfilePager(pagerCmd []string, writerFn func(io.Writer) error) bool {
	cmdPath, err := exec.LookPath(pagerCmd[0])
	if err != nil {
		return false
	}
	tmp, err := os.CreateTemp("", "pager-*")
	if err != nil {
		return false
	}
	defer func() {
		_ = os.Remove(tmp.Name())
	}()

	buf := &bytes.Buffer{}
	if err := writerFn(buf); err != nil {
		return false
	}

	if _, err := tmp.Write(buf.Bytes()); err != nil {
		return false
	}
	_ = tmp.Close()

	cmd := exec.Command(cmdPath, tmp.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run() == nil
}

type waiter interface {
	Wait() error
}

func waitIgnoringInterrupt(w waiter) error {
	for {
		err := w.Wait()
		if err == nil {
			return nil
		}
		if errors.Is(err, syscall.EINTR) {
			continue
		}
		return err
	}
}

func getPager() []string {
	if pager := os.Getenv("PAGER"); pager != "" {
		parts, err := shlex.Split(pager)
		if err == nil && len(parts) > 0 {
			return parts
		}
	}

	if runtime.GOOS == "windows" {
		return []string{"more"}
	}

	if _, okay := os.LookupEnv("LESS"); !okay {
		_ = os.Setenv("LESS", "-SRFX")
	}
	return []string{"less"}
}
