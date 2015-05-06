## Synopsis

An IRC bot for monitoring the [daily free book](https://www.packtpub.com/packt/offers/free-learning) provided by Packt. Type '!packt' and the bot will say the current free book's title.

I mostly wrote this to get practice using [Go](https://golang.org/). No warranty is expressed or implied.

## Installation

```
go get github.com/muraiki/packtbot
go build github.com/muraiki/packtbot
go install github.com/muraiki/packtbot
```

Create a `packtbot.yaml` (see sample in this repo) and configure the bot's name, the channel it will join, and the server it will join. To use a secure server, use field `secureserver` instead of `server`.

Once you have a `packtbot.yaml`, run `packtbot` from the same directory.

## Contributors

(c)2015 Erik Ferguson

## License

Do whatever you want with this. Just don't be mean. :)

