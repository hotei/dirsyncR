# Makefile for dirsyncR program

PROG = dirsyncR
VERSION = 0.0.1
TARDIR = 	$(HOME)/Desktop/TarPit/
DATE = 	`date "+%Y-%m-%d.%H_%M_%S"`
DOCOUT = README-$(PROG)-godoc.md

all:
	go build -v

# change cp to echo if you really don't want to install the program

install:
	go build -v
	go tool vet .
	go tool vet -shadow .
	gofmt -w *.go
#	killall -q $(PROG)dest
	cp $(PROG) $(PROG)dest
	cp $(PROG) $(PROG)source
	rm $(PROG)
#	go install
#	cp $(PROG) $(HOME)/bin

# note that godepgraph can be used to derive .travis.yml install: section
docs:
	godoc2md . > $(DOCOUT)
	godepgraph -md -p . >> $(DOCOUT)
	deadcode -md >> $(DOCOUT)
	cp README-$(PROG).md README.md
	cat $(DOCOUT) >> README.md

neat:
	go fmt ./...

dead:
	deadcode > problems.dead

index:
	cindex .

clean:
	go clean ./...
	rm -f *~ problems.dead count.out dirsyncRsource dirsyncRdest
	rm -f $(DOCOUT) README2.md

tar:
	echo $(TARDIR)$(PROG)_$(VERSION)_$(DATE).tar
	tar -ncvf $(TARDIR)$(PROG)_$(VERSION)_$(DATE).tar .

# for test only we allow both source and dest on same program run
test:
	./dirsyncRdest -send="/home/mdr/bin" -dest="/home/mdr/Public/binTwin" -verbose
