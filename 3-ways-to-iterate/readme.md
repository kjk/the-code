Code samples for <a href="https://blog.kowalczyk.info/book/go-cookbook.html">Go Cookbook</a> chapter on
<a href="">3 ways to iterate</a>.

Legend:

* <a href="inlined.go">inlined.go</a> : iteration and processing commingled
* <a href="callback.go">callback.go</a> : iteration code calls a processing callback for each item
* <a href="channel.go">channel.go</a> : iterator sends item over a channel
* <a href="next.go">next.go</a> : iterator is a struct that implements `Next` method to advance to the next item
