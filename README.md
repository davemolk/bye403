# bye403


## thanks
looked at these and more or less ported to Go
* [byp4xx](https://github.com/lobuhi/byp4xx)
* [403bypasser](https://github.com/yunemse48/403bypasser)

## header order not fixed :(
Go's http package doesn't preserve header order because http doesn't require it, so if the order of the request headers matters for what you're trying to do, well, maybe look elsewhere...