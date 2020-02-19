# dissect-tester

![Docker Pulls](https://img.shields.io/docker/pulls/jorgelbg/dissect-tester?style=plastic)

This project presents a simple web UI to test a collection of log line samples against a pattern
supported by the [Filebeat dissect processor](https://www.elastic.co/guide/en/beats/filebeat/master/dissect.html).

> Both [Logstash](https://www.elastic.co/guide/en/logstash/current/plugins-filters-dissect.html) and
> [Elasticsearch pipelines](https://www.elastic.co/guide/en/elasticsearch/reference/master/dissect-processor.html)
> have a similar filter/processor that uses the same configuration pattern. Therefore, this UI can be
> used to test a pattern that will be used in either Logstash or Elasticsearch pipelines.

A brief description of your project, what it is used for and how does life get
awesome when someone starts to use it.

## 🎮 Installing / Getting started

🔗 If you only want to test some samples you can go directly to the demo instance running in http://dissect-tester.jorgelbg.me/.

New releases are published to a public [Docker image](https://hub.docker.com/repository/docker/jorgelbg/dissect-tester). To run it you can use the following command:

```shell
docker run --rm -ti -p 8080:8080 jorgelbg/dissect-tester
```

The terminal should print a message similar to
```json
{"level":"info","ts":1582130507.8919635,"caller":"build/main.go:68","msg":"Server is running","port":8080}
```

Indicating the the server is running. Head your browser to http://localhost:8080/ and enjoy 🎉.

Your browser should show the following:

![Screenshot](http://screen.jorgelbg.me/jorgelbg-dropshare/Screen-Shot-2020-02-19-5-43-11.30-PM.png)

## 👨🏻‍💻 Developing

```shell
git clone https://github.com/jorgelbg/dissect-tester
cd dissect-tester/
make
```

This will build a binary placed in `bin/github.com/jorgelbg/dissect-tester` for your native platform.

If you want to build a new Docker image use the following command:

```shell
make docker
```

## Contributing

If you'd like to contribute, please fork the repository and use a feature
branch. Pull requests are warmly welcome.

## Links

- Project homepage: http://dissect-tester.jorgelbg.me/