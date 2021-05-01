module github.com/stmcginnis/ctlfish

go 1.15

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/lithammer/dedent v1.1.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stmcginnis/gofish v0.9.1-0.20210501200310-7ddb6f97be05
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.3 // indirect
)

replace github.com/stmcginnis/ctlfish => ./
