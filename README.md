# mimid
HTTP Mocking go library  

Mimid attempts to ease the creation of mocks for HTTP APIs.

## How it works 

### Proxy 
Initially mimid acts as a proxy between your service and the external API, capturing all HTTP requests and responses. 

`mimid proxy http://api.3rd-pary.com`

note: You need to configure your service to use the mimid ip:port instead of the real api.

The requests/responses pairs are captures as json files.
You can inspect and modify this files to make sure the assertions are correct.

### Mock

Once you have the json files can run mimid as mock server.
Your request will not reach the real server. 
Mimid, instead will try to return the most suitable response, based on your request and the provided files.


`mimid mock`


### Statefull APIs

Same request same response cases work excellent with mimid.

This will not be the cases for statefull APIs, look at the following example:

```
  -> GET /books
  <- []
  
  POST /books { ... }
  <- OK
  
  -> GET /books
  <- [{ ... }]

```

In the example above, `GET /books` has two possible responses: empty array or array with single object.

When mimid detects two request that are the same, but they have different responses, it examines all the request in between, in this case: `POST /books { ... }`

From this it assumes the following rule:

When a request `GET /books` takes place 
If a `POST /books { ... }` accrued before it, 
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

You will get a rule that will require both `POST /books` and `POST /avatar`. Although `/avatar` should affect the `/books`. Thankfully you can inspect the file and fix this manually.
