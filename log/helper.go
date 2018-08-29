package log

import (
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"
)

func Flags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewStringFlag(
			&cli.StringFlag{
				Name:    "log_level",
				Aliases: []string{"log.level"},
				Value:   "info",
				EnvVars: []string{"LOG_LEVEL"},
			}),
		altsrc.NewStringFlag(
			&cli.StringFlag{
				Name:    "log_format",
				Aliases: []string{"log.format"},
				Value:   "json",
				EnvVars: []string{"LOG_FORMAT"},
			}),
	}
}

// Setup helps with setting up logger
func Setup(c *cli.Context) error {
	if c.String("env") == "prod" || c.String("env") == "int" {
		if err := New(c.String("log_format"), c.String("log_level")); err != nil {
			return err
		}
	} else if c.String("env") == "dev" {
		if err := NewDevelopment(c.String("log_level")); err != nil {
			return err
		}
	} else if c.String("env") == "test" {
		if _, err := NewTest(); err != nil {
			return err
		}
	} else {
		return cli.Exit("failed to create logger", 1)
	}
	return nil
}
