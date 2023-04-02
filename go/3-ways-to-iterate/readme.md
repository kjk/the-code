Code samples from [Go Cookbook](https://blog.kowalczyk.info/book/go-cookbook.html).

Chapter: [3 ways to iterate](https://blog.kowalczyk.info/article/1Bkr/3-ways-to-iterate-in-go.html).

Legend:

* [inlined.go](inlined.go) : iteration and processing commingled
* [callback.go](callback.go) : iteration code calls a processing callback for each item
* [next.go](next.go) : iterator is a struct that implements `Next` method to advance to the next item
* [channel.go](channel.go) : generator sends items over a channel
* [channel-cancellable.go](channel-cancellable.go) : generator sends items over a channel, added ability to stop the generator by the client
