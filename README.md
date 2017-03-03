# An example Go web application

You have been hacking on Go and want to use it for a web application. You want a simple, opinionated, idiomatic (mostly!), and easy to understand Go application that has batteries included: Authentication, sessions, databases, middleware, etc. You are tired of piecing together things from various snippets and poring over other Go code bases. Look no further!

This is a simple Go web application. It is a skeleton and example that can quickly be customized. It is not a framework. It is not a library. It is meant to be cloned and hacked on, getting an app up and running very quickly. Why? Mostly so I can use this for my own purposes, but it is a really nice starting point. 

Inspired by common Go patterns. Some code borrowed from GoPhish (in particular, the Context package, simple middleware, context handling, and bits of pieces of the Login code). Boiled down to its essence and simplified. This project is ready for hacking. 

No fancy HTML. Everything in one place, easy to see with minimal abstraction. 

You write code. You build it. You get a deploy anywhere on almost any platform web application in one binary. 

## Demonstrated simply in this app

 * Context passing
 * Cookie based session store
 * Native templates
 * Storm (an ORM based on BoltDB)
 * Gorilla packages (mux, session)
 * Middleware
 * Sessions (flashing messages, authentication)
 * Authenticated routes
 * Simple old school CRUD interface for Users

All less than 375 lines of Go.

## About Storm and BoltDB

Bolt is a really nice database. Storm is one of the easiest to use "ORMs" you will find. If you want to get hacking quickly give it a try. If you don't like it throw it out and use something else.

## Quick Tour

There are two commands `gowe` and `gowe-user`. 

Use `gowe-user` to add a new user.

Use `gowe` to launch gowebexample. 

To install try something like:

`go get github.com/bitexploder/gowebexample/...`

Cd to the `gowebexample` directory (where `static` and `template` dirs are)

Then run:

`gowe-user -username admin -password admin -name First Last`
`gowe`

Now visit:

`http://127.0.0.1:8080`

And login with username admin password admin. 

That's it!

## Implementation Notes

Everything but the user model and context package lives in: `cmd/gowe/web.go`

There is a lot left as an exercise to the reader (including tightening up security settings, such as login error messages, cookie settings, and other low hanging fruit.)

Bcrypt as used is a nice way to store password hashes. The authentication mechanism and authenticated routes should be reasonable sound. 

If you need TLS (you probably do!), just generate a cert and key and modify ListenAndServer to use ListenAndServeTLS.

If you want to edit and update passwords checkout gowe-user for how to use the bcrypt package and add it to the `EditUserHandler`


## Thoughts on running in production

All logging goes to stdout. Use tee and output redirection to put it in a file.

# License

~~~~~
Gowebexample - An example Go application

The MIT License (MIT)

Copyright (c) 2017 Jeremy Allen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software ("Gophish Community Edition") and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.



(Context and bits and pieces here and there)
Copyright (c) 2013 - 2017 Jordan Wright

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software ("Gophish Community Edition") and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions...
~~~~~

 


