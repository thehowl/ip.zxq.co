# Benchmarks

We can't benchmark ipinfo.io, because it's not open source so we cannot test it in localhost. Although, we did do a test for ip.zxq.co. So, this is what [boom](https://github.com/rakyll/boom) tells us after doing 7000 requests to localhost (please note that this test was done on a Windows machine, so it probably runs slightly faster on Linux, and thus the server):

```
$ boom -n 7000 -c 200 http://localhost:61430 # (lookup on ::1)
7000 / 7000  100.00 % s

Summary:
  Total:        7.6446 secs.
  Slowest:      0.4261 secs.
  Fastest:      0.0350 secs.
  Average:      0.2176 secs.
  Requests/sec: 915.6815
  Total Data Received:  980000 bytes.
  Response Size per Request:    140 bytes.

Status code distribution:
  [200] 7000 responses

Response time histogram:
  0.035 [3]     |
  0.074 [26]    |
  0.113 [83]    |∎
  0.152 [521]   |∎∎∎∎∎∎∎∎∎∎
  0.191 [1844]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.231 [1983]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.270 [1238]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.309 [828]   |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.348 [350]   |∎∎∎∎∎∎∎
  0.387 [96]    |∎
  0.426 [28]    |

Latency distribution:
  10% in 0.1541 secs.
  25% in 0.1761 secs.
  50% in 0.2091 secs.
  75% in 0.2501 secs.
  90% in 0.2941 secs.
  95% in 0.3191 secs.
  99% in 0.3661 secs.
$ boom -n 7000 -c 200 http://localhost:61430/8.8.8.8
7000 / 7000  100.00 % s

Summary:
  Total:        7.7096 secs.
  Slowest:      0.4682 secs.
  Fastest:      0.0450 secs.
  Average:      0.2192 secs.
  Requests/sec: 907.9589
  Total Data Received:  1400000 bytes.
  Response Size per Request:    200 bytes.

Status code distribution:
  [200] 7000 responses

Response time histogram:
  0.045 [2]     |
  0.087 [27]    |
  0.130 [314]   |∎∎∎∎∎∎
  0.172 [1157]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.214 [2051]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.257 [1683]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.299 [1126]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.341 [419]   |∎∎∎∎∎∎∎∎
  0.384 [157]   |∎∎∎
  0.426 [55]    |∎
  0.468 [9]     |

Latency distribution:
  10% in 0.1511 secs.
  25% in 0.1781 secs.
  50% in 0.2131 secs.
  75% in 0.2581 secs.
  90% in 0.2951 secs.
  95% in 0.3231 secs.
  99% in 0.3811 secs.
```
