package main

import (
	"fmt"
	"os"

	"github.com/lucheng0127/license/pkg/license"
	"github.com/urfave/cli/v2"
)

func generate(cCtx *cli.Context) error {
	licMgr, err := license.NewLicenseManager(cCtx.String("key"), "./")
	if err != nil {
		return err
	}

	lic, err := licMgr.GenerateLicense(cCtx.String("dmi"), cCtx.Int("life"))
	if err != nil {
		return err
	}

	fmt.Printf("=== License - days %d ===\n%s\n", cCtx.Int("life"), lic)

	return nil
}

func importLic(cCtx *cli.Context) error {
	licMgr, err := license.NewLicenseManager("0123456789abcdef", cCtx.String("dir"))
	if err != nil {
		return err
	}

	err = licMgr.Import(cCtx.String("lic"))
	if err != nil {
		return err
	}

	fmt.Println("license import succeed")
	return nil
}

func showLife(cCtx *cli.Context) error {
	licMgr, err := license.NewLicenseManager(cCtx.String("key"), cCtx.String("dir"))
	if err != nil {
		return err
	}

	lifetime, err := licMgr.LifeTime()
	if err != nil {
		return err
	}

	fmt.Printf("License expired in %s\n", lifetime)

	return nil
}

func generateCmd() *cli.Command {
	return &cli.Command{
		Name:   "generate",
		Action: generate,
		Usage:  "generate a license",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dmi",
				Required: true,
				Usage:    "dmi code get by toolkit",
			},
			&cli.StringFlag{
				Name:     "key",
				Required: true,
				Usage:    "encrypt key, length 16, 24 or 32",
			},
			&cli.IntFlag{
				Name:     "life",
				Required: true,
				Usage:    "lifetime of license (days)",
			},
		},
	}
}

func importCmd() *cli.Command {
	return &cli.Command{
		Name:   "import",
		Action: importLic,
		Usage:  "import a license",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dir",
				Value: "/var/run/lic",
				Usage: "license directory default /var/run/lic",
			},
			&cli.StringFlag{
				Name:     "lic",
				Required: true,
				Usage:    "license string",
			},
		},
	}
}

func showCmd() *cli.Command {
	return &cli.Command{
		Name:   "show",
		Action: showLife,
		Usage:  "show lifetime of license",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "key",
				Required: true,
				Usage:    "encrypt key, length 16, 24 or 32",
			},
			&cli.StringFlag{
				Name:  "dir",
				Value: "/var/run/lic",
				Usage: "license directory default /var/run/lic",
			},
		},
	}
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			generateCmd(),
			importCmd(),
			showCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
