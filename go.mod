module github.com/jorgelbg/dissect-tester

go 1.13

replace (
	github.com/Shopify/sarama => github.com/elastic/sarama v1.19.1-0.20200629123429-0e7b69039eec
	github.com/dop251/goja => github.com/andrewkroh/goja v0.0.0-20190128172624-dd2ac4456e20
	github.com/dop251/goja_nodejs => github.com/dop251/goja_nodejs v0.0.0-20171011081505-adff31b136e6
	github.com/fsnotify/fsevents => github.com/elastic/fsevents v0.0.0-20181029231046-e1d381a4d270
)

require (
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/dop251/goja v0.0.0-20201022115936-e21ccf39bfce // indirect
	github.com/dop251/goja_nodejs v0.0.0-20200811150831-9bc458b4bbeb // indirect
	github.com/elastic/beats/v7 v7.9.3
	github.com/elastic/go-sysinfo v1.4.0 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/google/go-cmp v0.4.1
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/prometheus/procfs v0.2.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/spf13/cobra v1.0.0 // indirect
	go.elastic.co/apm v1.8.0 // indirect
	go.elastic.co/ecszap v0.2.0 // indirect
	go.elastic.co/fastjson v1.1.0 // indirect
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/net v0.0.0-20201022231255-08b38378de70 // indirect
	golang.org/x/sys v0.0.0-20201022201747-fb209a7c41cd // indirect
	golang.org/x/tools v0.0.0-20201023174141-c8cfbd0f21e6 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	howett.net/plist v0.0.0-20200419221736-3b63eb3a43b5 // indirect
)
