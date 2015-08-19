<center>
dirsyncr
===
</center>.

<h3>   <a href="http://godoc.org/github.com/hotei/dirsyncr">
<img src="https://godoc.org/github.com/hotei/dirsyncr?status.png" alt="GoDoc" />
</a></h1>

Travis tags are good but require registration at http://travis-ci.org and then
a bit of setup to put in a tarvis.yaml file.  See their homepage 
http://godoc.org and http://docs.travis-ci.com/user/languages/go/ for specifics.

<!---
<a href="http://travis-ci.org/hotei/dirsyncr">
<img src="https://secure.travis-ci.org/hotei/dirsyncr.png" alt="Build Status" /></a>
--->

This "README" document is (c) 2015 David Rook. 

dirsyncr is (c) 2015 David Rook - all rights reserved. The program
and related files in this repository are released under BSD 2-Clause License.
License details are at the end of this document. Bugs/issues can be reported on github.
Comments can be sent to <hotei1352@gmail.com>. "Thank you"s, constructive 
criticism and job offers are always welcome.  If you are so inclined you may
donate money to help me continue the project by sending PayPal contributions to
my wife at diane@celticpapers.com. 


### Description

dirsyncr keeps two directories upto date so the contents are identical (within
the resolution of the update interval).

### Why should I use it instead of rsync?

The main reason of course is that it's written in go and open source so you can 
modify it to suit your needs.  While rsync is powerful, the list of options is
several pages long and intimidating to new users.    

### How does it work?

### <font color=red> >>> Please Read This First <<< </font>

### Background Info on the problem

### Installation

If you have a working go installation on a Unix-like OS:

> ```go get github.com/hotei/dirsyncr```

If you don't have go installed you can always download the git repository as
a zipfile and extract it or just browse the godoc.org if you are just curious.

### Configuration

### Features

### Usage (incl option flags)

### Notes

### To do

NOTE:  "higher" relative priority is at top of list

1.  Added as needed (none active at the moment)

### Limitations


### Issues (please file on github.com)

1.  Added as needed (none active at the moment)
`

### Development Environment
	Mint 17.1 Linux on i7/2500 mhz 8 'core' (4+4HT) HP Envy Laptop
	X11/R6
	gnu g++ compiler gcc version 4.8.2 (Ubuntu 4.8.2-19ubuntu1)
	go 1.5rc1

	
### Change Log

* 2015-08-18 built with go 1.5rc1

### References

* [go language reference] [1] 
* [go standard library package docs] [2]
* [Source for dirsyncr on github] [3]
* [Go projects list(s)] [7]
* [Excellent godoc howto by Nate Finch] [8]

[1]: http://golang.org/ref/spec/ "go reference spec"
[2]: http://golang.org/pkg/ "go package docs"
[3]: http://github.com/hotei/dirsyncr "github.com/hotei/dirsyncr"
[4]: http://golang.org/doc/go1compat.html "Go 1.x API contract"
[5]: http://blog.golang.org/2011/06/profiling-go-programs.html "Profiling go code"
[6]: http://golang.org/doc/articles/godoc_documenting_go_code.html "GoDoc HowTo"
[7]: https://github.com/golang/go/wiki/Projects "go project list"
[8]: https://github.com/natefinch/godocgo "Nate Finch's Tutorial for GoDoc"

Comments can be sent to David Rook  <hotei1352@gmail.com>

### Disclaimer
Any trademarks mentioned herein are the property of their respective owners.

### License

The 'dirsyncr' go program/package and demo programs are distributed under the Simplified BSD License:

> Copyright (c) 2015 David Rook. All rights reserved.
> 
> Redistribution and use in source and binary forms, with or without modification, are
> permitted provided that the following conditions are met:
> 
>    1. Redistributions of source code must retain the above copyright notice, this list of
>       conditions and the following disclaimer.
> 
>    2. Redistributions in binary form must reproduce the above copyright notice, this list
>       of conditions and the following disclaimer in the documentation and/or other materials
>       provided with the distribution.
> 
> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDER ``AS IS'' AND ANY EXPRESS OR IMPLIED
> WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
> FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> OR
> CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
> CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
> SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
> ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
> NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
> ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

----

<center>
# Automatically Generated Documentation Follows
</center>


