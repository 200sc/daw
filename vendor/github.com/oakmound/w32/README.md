About w32
==========

w32 is a wrapper of windows apis for the Go Programming Language.

It wraps win32 apis to "Go style" to make them easier to use.


## About This Fork ##

This is forked from [AllenDang/w32](https://github.com/AllenDang/w32).

This fork is used internally by [oakmound/shiny](https://github.com/oakmound/shiny). 

It removes a few functions from w32 that relied on a C dependency.


Setup
=====

1. Install Go

2. go get -u github.com/oakmound/w32

3. go install github.com/oakmound/w32
