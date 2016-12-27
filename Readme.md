Twitter Streaming Sentiment
==

This is a simple demo of using the twitter streaming api to watch live tweets matching certain keywords and applying sentiment analysis to them. Originally written on the night of UK's EU Referendum in an attempt to measure live sentiment of brexit and other terms.

## Usage
* Run `go get` to download the dependencies.
* Update the `blank_config.json` with your twitter api keys/secrets, then copy/rename to `config.json`.
* Run with `go run twitterstreamingsentiment.go --terms` followed by comma separated search terms, i.e. `go run twitterstreamingsentiment.go --terms brexit,trump`
