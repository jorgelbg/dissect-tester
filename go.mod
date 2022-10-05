module github.com/jorgelbg/dissect-tester

go 1.19

replace (
	github.com/Shopify/sarama => github.com/elastic/sarama v1.19.1-0.20200629123429-0e7b69039eec
	github.com/dop251/goja => github.com/andrewkroh/goja v0.0.0-20190128172624-dd2ac4456e20
	github.com/dop251/goja_nodejs => github.com/dop251/goja_nodejs v0.0.0-20171011081505-adff31b136e6
	github.com/fsnotify/fsevents => github.com/elastic/fsevents v0.0.0-20181029231046-e1d381a4d270
)

require (
	github.com/elastic/beats/v7 v7.15.0
	github.com/google/go-cmp v0.5.2
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/zap v1.14.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
