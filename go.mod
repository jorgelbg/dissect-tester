module github.com/jorgelbg/dissect-tester

go 1.13

replace (
	github.com/Shopify/sarama => github.com/elastic/sarama v1.19.1-0.20200629123429-0e7b69039eec
	github.com/dop251/goja => github.com/andrewkroh/goja v0.0.0-20190128172624-dd2ac4456e20
	github.com/fsnotify/fsevents => github.com/elastic/fsevents v0.0.0-20181029231046-e1d381a4d270
)

require (
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/dlclark/regexp2 v1.2.0 // indirect
	github.com/dop251/goja_nodejs v0.0.0-20200128125109-2d688c7e0ac4 // indirect
	github.com/elastic/beats v7.6.2+incompatible
	github.com/elastic/beats/v7 v7.8.0
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/google/go-cmp v0.4.1
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/spf13/cobra v1.0.0 // indirect
	github.com/stretchr/testify v1.5.0
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/zap v1.15.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
