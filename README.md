# ip.zxq.co

An ipinfo.io clone, without the rate limiting. And some other thingies.

## Why?

The main reason for this is: rate limiting. ipinfo.io sets a rate limiting of 1000 requests per day. I understand it, although that's a bit of a bummer. So I said, whatever. Let's just write our own clone. So, there you have it.

## Using

This is really similiar to the way ipinfo.io does it. Every response will be identical from ipinfo.io's, almost. So, basic usage: you will just make a GET request to <http://ip.zxq.co/[ip]>, like http://ip.zxq.co/8.8.8.8. Need to get something specific? http://ip.zxq.co/8.8.8.8/country.

```
$ curl ip.zxq.co/8.8.8.8
{"city":"Mountain View","continent":"NA","continent_full":"North America","country":"US","country_full":"United States","ip":"8.8.8.8","loc":"37.3860,-122.0838","postal":"94040","region":"California"}
```

```
$ curl ip.zxq.co/8.8.8.8/country
US
```

That is a bit ugly, though. Especially for testing purposes. Let's make it pretty, shall we?

```
$ curl "ip.zxq.co/8.8.8.8?pretty=1"
{
  "city": "Mountain View",
  "continent": "NA",
  "continent_full": "North America",
  "country": "US",
  "country_full": "United States",
  "ip": "8.8.8.8",
  "loc": "37.3860,-122.0838",
  "postal": "94040",
  "region": "California"
}
```

Much better! :smile:

But what if we don't know the user's IP? In that case, then, we can call `/self`.

```
$ curl "ip.zxq.co/self?pretty=1"
{
  "city": "",
  "continent": "EU",
  "continent_full": "Europe",
  "country": "IT",
  "country_full": "Italy",
  "ip": "87.16.45.15",
  "loc": "42.8333,12.8333",
  "postal": "",
  "region": ""
}
```

Using /self works just like any normal request with an IP, which means you can still get a specific field. By default if you do not specify any IP it will assume it's /self, although it's recommended to always specify it, as if for any reason you have "webkit", "opera" or "mozilla" in your user agent, it will assume you're accessing from a browser and redirect you to this README.

We're aren't done just yet! You want to use JSONP. You guess it, we are using the same system as ipinfo.io's. Just provide a `callback` parameter to your GET request.

```
$ curl "ip.zxq.co/8.8.8.8?pretty=1&callback=myFancyFunction"
/**/ typeof myFancyFunction === 'function' && myFancyFunction({
  "city": "Mountain View",
  "continent": "NA",
  "continent_full": "North America",
  "country": "US",
  "country_full": "United States",
  "ip": "8.8.8.8",
  "loc": "37.3860,-122.0838",
  "postal": "94040",
  "region": "California"
});
```

```html
<script>
var myFancyFunction = function(data) {
  alert("The city of the IP address 8.8.8.8 is: " + data.city);
}
</script>
<script src="http://ip.zxq.co/8.8.8.8?callback=myFancyFunction"></script>
```

## Features that aren't here (and not going to be implemented)

* Hostname. We would have to pick that data from another data source, which is too much effort.
* Organization field. Need another datasource. Effort. Performance. I'll just stick with not having it.

## Features that aren't on ipinfo.io but are here

* JSON minified, so it gets to your server quicker.
* Full name for the country!
* We also got continent info, with the full name too.

## Some advantages:

* We are using Go and not nodejs like them. Go is a compiled language, and therefore is [amazingly fast. A response can be generated in a very short time.](Benchmarks.md)
* We get data only from one data source. Which means no lookups on other databases, which results in being faster overall.
* We are open source. Which means you can compile and put it on your own server!

## Running locally/development/contributing:

Feel free to open an issue or pull request for anything! If you want to run it locally for whatever reason, you can do so this way if you don't need to touch the code:

```sh
go get -d http://github.com/TheHowl/ip.zxq.co
cd $GOPATH/src/github.com/TheHowl/ip.zxq.co
go build
./ip.zxq.co # .exe if you're on windows
```

(the reason you can't just do `go get` and then execute it from the terminal is that the software requires `GeoLite2-City.mmdb` to be in the same folder)

If you want to hack in the future, this is a better way:

```sh
cd $GOPATH
mkdir -p src/github.com/TheHowl
cd src/github.com/TheHowl
git clone git@github.com:TheHowl/ip.zxq.co.git
cd ip.zxq.co
go build
./ip.zxq.co
# Or if you don't want to create the binary in the folder
go run main.go
```

## Data source

This product includes GeoLite2 data created by MaxMind, available from http://www.maxmind.com.

## License

MIT. Check LICENSE file.
