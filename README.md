## go based HTTP simple bittorrent tracker

You can use this by import the `server` package that implements `http.Handler` to use in your own 
applications or use the binary `bttracker` to create a standalone bittorrent tracker.


You have the option to use an in-memory registry for keeping peer data or using a redis server for storing the 
peer data.  You can always implement your own registry as well.

### btracker

```bash
bttracker -h
Usage of bttracker:
    -addr=":9090": address of the tracker
    -debug=false: enable debug mode for logging
    -interval=120: interval for when Peers should poll for new peers
    -min-interval=30: min poll interval for new peers
    -redis-addr="": address to a redis server for persistent peer data
    -redis-pass="": password to use to connect to the redis server
```

### License MIT
Copyright (c) 2014 Michael Crosby. michael@crosbymichael.com

Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation 
files (the "Software"), to deal in the Software without 
restriction, including without limitation the rights to use, copy, 
modify, merge, publish, distribute, sublicense, and/or sell copies 
of the Software, and to permit persons to whom the Software is 
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be 
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED,
INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, 
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. 
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT 
HOLDERS BE LIABLE FOR ANY CLAIM, 
DAMAGES OR OTHER LIABILITY, 
WHETHER IN AN ACTION OF CONTRACT, 
TORT OR OTHERWISE, 
ARISING FROM, OUT OF OR IN CONNECTION WITH 
THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
