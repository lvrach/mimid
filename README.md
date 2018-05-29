# mimid
HTTP Mocking go library  

Attemts to ease the creation of mocks for HTTP APIs.

## How it works 

### Proxy 
Inititally mimid acts as a proxy between your service and the extranal API, capturing all HTTP requests and responses. 

`mimid proxy http://api.3rd-pary.com`

note: You need to configure your service to use the mimid ip:port instead of the real api.

The requests/responses pairs are captures as json files.
You can inspect and modify this files to make sure the assertions are correct.

### Mock

Once you have the json files can run mimid as mock server.
Your request will not reach the real server. 
Mimid, instead will try to return the most sutable response, based on your request and the provided files.


`mimid mock`


### Statefull APIs

If for the same request we always get the same response, mimid will work without a problem.

However in most cases that is not the case, look at the following example:

```
  -> GET /books
  <- []
  
  POST /books { ... }
  <- OK
  
  -> GET /books
  <- [{ ... }]

```

In the example above, `GET /books` has multiple possible repsonses.

When mimid detects the same request has diffrent responses it examines all the request that happen in between.

In this case: `POST /books { ... }`

From this it assumes:

When a request `GET /books` takes place 
If a `POST /books { ... }` accured before it, 
  then answer `[{ ... }]`
  
Else 
  `[ ]`
  

This does not always works as expected. If you had:
```
  -> GET /books
  <- []
  
  POST /books { ... }
  <- OK
  
  POST /avatar {....}
  <- OK
  
  -> GET /books
  <- [{ ... }]

```

You will get a rule with `POST /books` and `POST /avatar`. Althouht `/avatar` should affect the `/books`.

This is why manual inspection is needed afterwards. You may have to add/remove some rules. 
