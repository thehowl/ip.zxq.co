# ip.zxq.co

An ipinfo.io clone, without the rate limiting. And some other thingies.

## Why?

The main reason for this is: rate limiting. ipinfo.io sets a rate limiting of 1000 requests per day. I understand it, although that's a bit of a bummer. So I said, whatever. Let's just write our own clone. So, there you have it.

## Features that aren't here (and not going to be implemented)

* Hostname. We would have to pick that data from another data source, which is too much effort.
* Organization field. Need another datasource. Effort. Performance. I'll just stick with not having it.
* JSONP. Might bother to implement one day.

## Features that aren't in ipinfo.io but are here

* JSON minified, so it gets to your server quicker.
* Full name for the country!
* We also got continent info, with the full name too.

## Using

This is really similiar to the way ipinfo.io does it. Every response will be identical from ipinfo.io's, almost. So, basic usage: you will just make a GET request to <http://ip.zxq.co/&lt;ip&gt;>, like http://ip.zxq.co/8.8.8.8. Need to get something specific? http://ip.zxq.co/8.8.8.8/country.

## Why us rather than them?

tl;dr: we are faster than ipinfo.io. Also open sauce!

Long version:

* We are using Go and not nodejs like them. Go is a compiled language, and therefore is amazingly fast. A response can be generated in a very short time.
* We get data only from one data source. Which means no lookups on other databases, which results in being faster overall.
* We are open source. Which means you can compile and put it on your own server!

## License

MIT. Check LICENSE file.
