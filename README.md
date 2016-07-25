HttpMux
==========

HttpMux is a lightweight high performance HTTP request router (also called multiplexer or just mux for short) for Go.

Installation
------------

```
go get github.com/niflheims-io/httpmux
```

Usage
------------

```go
    mux := httpmux.New()
    // set context
    mux.Ctx().Set("some key", "some interface")

	mux.Get("/1", func(r *httpmux.Request) {
	    // v, ok := r.Ctx().Get("some_key")
		r.Response().Status(http.StatusOK).String(time.Now().String())
		// also support json, xml, bytes, html/template..
	})

	mux.Get("/1/:k", func(r *httpmux.Request) {
		k := r.Query("k")
		r.Response().Status(http.StatusOK).String("[" + k + "]" + time.Now().String())
	})

    // also support http2
    srv := http.Server{
		Addr:        ":8080",
		Handler:     mux,
		ReadTimeout: 120 * time.Second,
	}

    fmt.Println(srv.ListenAndServe())

```

Status
------

* Golang >= 1.6.2


License
-------

GNU GENERAL PUBLIC LICENSE

Copyright (C) 2015-2016 niflheims-io 