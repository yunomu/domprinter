package domprinter

import (
	"context"
	"io"
	"io/ioutil"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

type DomPrinter struct {
	cmdPath string
}

func New() *DomPrinter {
	return &DomPrinter{
		cmdPath: "chromium-browser",
	}
}

func (p *DomPrinter) Print(ctx context.Context, url string, out io.Writer) error {
	return p.PrintWithError(ctx, url, out, ioutil.Discard)
}

func (p *DomPrinter) PrintWithError(ctx context.Context, url string, out io.Writer, errOut io.Writer) error {
	cmd := exec.CommandContext(ctx, p.cmdPath, "--headless", "--disable-gpu", "--dump-dom", "--no-sandbox", url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer stdout.Close()

		_, err := io.Copy(out, stdout)
		return err
	})

	g.Go(func() error {
		defer stderr.Close()

		io.Copy(errOut, stderr)
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return cmd.Wait()
}
