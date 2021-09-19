# redkstats <sup>🟥 🔑 📊</sup>
`redkstats` is a result of mixing parts of words "redis keys stats". This utility is intended for gathering keys and their idleTime to build basic statistics for further analysis.

## The real use-case
Imagine that you have plenty of unnecessary data. And keys are not set to be "expirable" or those data should be expired at some point in the distant future. But the problem is - you don't really know what are the keys or even their prefixes (namespaces).

So this is the moment where this utility comes in action. Using `redkstats` you could collect all the keys and their idle time (see redis command [OBJECT IDLETIME](https://redis.io/commands/OBJECT)), group them by their prefix and get some
basic stats. Summary that you got might help you to decide which keys should be deleted from the storage.

## License
This project is licensed under the GNU GPLv3 - see the [LICENSE](LICENSE) file for details

