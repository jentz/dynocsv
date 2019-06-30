package main

import (
	"bufio"
	"fmt"
	"github.com/zshamrock/dynocsv/aws/dynamodb"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"strings"
)

const (
	tableFlagName   = "table"
	columnsFlagName = "columns"
	limitFlagName   = "limit"
	profileFlagName = "profile"
	outputFlagName  = "output"
)

const appName = "dynocsv"

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = `Export DynamoDB table into CSV file`
	app.Version = "1.0.0"
	app.Author = "(c) Aliaksandr Kazlou"
	app.Metadata = map[string]interface{}{"GitHub": "https://github.com/zshamrock/dynocsv"}
	app.UsageText = fmt.Sprintf(`%s		 
		--table/-t <table> 
		[--columns/-c <comma separated columns>] 
		[--limit/-l <number>]
		[--profile/-p <AWS profile>]
		[--output/-o <output file name>]`,
		appName)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  fmt.Sprintf("%s, t", tableFlagName),
			Usage: "table to export",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s, c", columnsFlagName),
			Usage: "optional columns to export from the table, if skipped, all columns will be exported",
		},
		cli.UintFlag{
			Name:  fmt.Sprintf("%s, l", limitFlagName),
			Usage: "limit number of records returned, if not set all items are fetched",
		},
		cli.StringFlag{
			Name: fmt.Sprintf("%s, p", profileFlagName),
			Usage: "AWS profile to use to connect to DynamoDB, otherwise the value from AWS_PROFILE env var is used " +
				"if available, or then \"default\" if it is not set or empty",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s, o", outputFlagName),
			Usage: "output file, or the default <table name>.csv will be used",
		},
	}
	app.Action = action

	err := app.Run(os.Args)
	if err != nil {
		log.Panicf("error encountered while running the app %v", err)
	}
}

func action(c *cli.Context) error {
	table := mustFlag(c, tableFlagName)
	columns := c.String(columnsFlagName)
	filename := c.String(outputFlagName)
	if filename == "" {
		filename = fmt.Sprintf("%s.csv", table)
	}
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	limit := c.Uint(limitFlagName)
	profile := c.String(profileFlagName)
	headers := dynamodb.ExportToCSV(profile, table, columns, limit, bufio.NewWriter(file))
	if columns == "" {
		fmt.Println(strings.Join(headers, ","))
	}
	return file.Close()
}

func mustFlag(c *cli.Context, name string) string {
	value := c.String(name)
	if value == "" {
		log.Panic(fmt.Sprintf("%s is required", name))
	}
	return value
}
