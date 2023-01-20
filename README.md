# ip.zxq.co

An ipinfo.io clone, without the rate limiting.

## Why?

In my projects, I sometimes need easy ways to get the user's GeoIP data, without
having to set up my own code to update GeoIP databases and so on. Other times, I
want to get the GeoIP information through JavaScript. In both cases, ipinfo.io
would work well, but it does far more things than I need and it sets a limit of
50k requests per month. For most projects, this is unlikely to be hit, but I
still thought it could be good to write my own clone where I know this cannot
happen.

## Usage

Querying https://ip.zxq.co will return the current IP address' information. To
query a specific IP address, pass it in the path:

```sh
$ curl https://ip.zxq.co
{"ip":"[stripped]","city":"Gerbido","region":"Piedmont","country":"IT","country_full":"Italy","continent":"EU","continent_full":"Europe","loc":"45.0445,7.6141","postal":""}
$ curl https://ip.zxq.co/1.1.1.1
{"ip":"1.1.1.1","city":"Sydney","region":"New South Wales","country":"AU","country_full":"Australia","continent":"OC","continent_full":"Oceania","loc":"-33.8688,151.2090","postal":""}
```

For each field in the JSON object, you may request it specifically to be given
to you as simple text:

```sh
$ curl https://ip.zxq.co/1.1.1.1/loc
-33.8688,151.2090
$ curl https://ip.zxq.co/self/city # To use the current user's IP
Gerbido
```

Please remember the GeoIP information is not precise and your users could very
easily be in a different city, and sometimes a different country altogether
(unlikely, but has happened).

For convenience, the `pretty=1` query string flag is provided to format the JSON
body:

```json
$ curl "https://ip.zxq.co/8.8.8.8?pretty=1"
{
  "ip": "8.8.8.8",
  "city": "Mountain View",
  "region": "California",
  "country": "US",
  "country_full": "United States",
  "continent": "NA",
  "continent_full": "North America",
  "loc": "37.4223,-122.0850",
  "postal": ""
}
```

## Notes about features, and contributions

This project is only meant to provide simple lookups using a GeoIP database.
It does not aim to provide perfect backwards-compatibility with ipinfo.io, nor
does it want to provide all of the information returned by ipinfo.io. I consider
this project to be feature-complete, thus no major features will be added or
more information from other datasources; bug reports and fixes, however, are
welcome provided they are backwards-compatible with the current version of the
software.

The project is run for free and has been online without major disruptions since
2015, and you are welcome to use it for any purpose (see below for info about
licensing), however remember that it is provided without any guarantees of
uptime or SLAs, so if that's what you're after you can host it yourself or use a
paid service :)

## Deployment and development

For your convenience, a docker-compose file is available, as well as a images on
docker hub: https://hub.docker.com/repository/docker/howl/ipzxqco/general .

`docker compose up` with the docker-compose file will set up ip.zxq.co on your
machine, listening on port 8123, with attempts to update every week.

## License

The code is under the MIT license. Check LICENSE file for the full text.

The software, including its hosted version on https://ip.zxq.co, make use of
[db-ip.com's GeoIP database](https://db-ip.com/db/download/ip-to-city-lite),
specifically the IP to City lite database. It is provided under the
<a href="https://creativecommons.org/licenses/by/4.0/">Creative Commons Attribution 4.0 International License.</a>

Quoting from the page:

>You are free to use this IP to City Lite database in your application, provided you give attribution to DB-IP.com for the data.
>
>In the case of a web application, you must include a link back to DB-IP.com
>on pages that display or use results from the database.
>You may do it by pasting the HTML code snippet below into your code:
>
>```
<a href='https://db-ip.com'>IP Geolocation by DB-IP</a>
```
