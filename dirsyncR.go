// dirsyncR (c) 2015 David Rook - all rights reserved
//
// directory replication suitable for use with cron
//
// Usage
//
// dirsyncR -dest="/home/mdr/bin" on the "receiver" system
// dirsyncR -send="/home/mdr/bin"	on the "sending" system
//
// options:
//	-port=8283
//	-verbose
//	usage:	dirsyncR -send="/home/mdr/bin" -port=8283
//			dirsyncR -dest="/home/koor/bin" -port=8283
//
//	send (source dir) must be a full path ie. begins with slash
//	dest (receive dir) must be a full path
//	OPTIONS:
//		-send="/home/mdr/bin         send this directory
//		-dest="/home/koor/bin         receive to this directory
//		-port=12345                  send/receive on this port
//		-version                     print version info and exit
//		-v                           print version info and exit
//		-verbose                     use extra output detail
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
	//
	"github.com/hotei/mdr" // for hash helpers
)

const (
	ServerPort = ":8283"
	Version    = "dirsyncR 0.0.1 (c) 2015 David Rook"
)

type CmdT int

const (
	LookupDigest = iota
	RefreshDigest
)

type SyncRequest struct {
	Id       int
	CmdType  int
	FilePath string
	CRC64    uint64
}

type SyncReply struct {
	Id       int // Id of request to which we're replying
	Digest   string
	FilePath string
	CRC64    uint64
}

type DigestFullPathT struct {
	Digest   string
	FullPath string
}

var (
	flagDest    string
	flagPortNum int
	flagSource  string
	flagVersion bool
	flagVerbose bool

	hasSource        bool // implies you want to run as server
	hasPort          bool
	hasDest          bool
	invokedByName    string // ./dirsyncR --> dirsyncR
	serverMapIsReady bool
	niceValue        = 1 * time.Millisecond
	hashUpdate       sync.Mutex

	// server stuff
	gDigests    []*DigestFullPathT //
	gDigestMap  map[string]string  // map[filename]digest
	requestChan chan SyncRequest
	replyChan   chan SyncReply
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	flag.StringVar(&flagDest, "dest", "/home/mdr/bin", "receiving directory")
	flag.StringVar(&flagSource, "send", "/home/mdr/bin", "directory to push")
	flag.IntVar(&flagPortNum, "port", 8283, "port from which to serve")
	flag.BoolVar(&flagVersion, "version", false, "show version info and stop")
	flag.BoolVar(&flagVerbose, "v", false, "use verbose mode")
	flag.BoolVar(&flagVerbose, "verbose", false, "use verbose mode")
	gDigestMap = make(map[string]string)
}

func usage() {
	s := `usage: dirsyncR -send="/home/mdr/bin" -port=8283
	send (source dir) must be a full path ie. begins with slash
		-send="/home/mdr/bin         send this directory
		-dest="/home/mdr/bin         receive to this directory
		-port=12345                  send/receive on this port
		-version                     print version info and exit
		-v                           print version info and exit
		-verbose                     use extra output detail
	`
	fmt.Printf("%s\n", s)
	os.Exit(0)
}

// buildSHA256List walks pathname tree and fills out
// gDigests [DigestFullPathT] as side effect.
func buildSHA256List(pathname string, info os.FileInfo, err error) error {
	time.Sleep(niceValue) // don't hog CPU
	if info == nil {
		fmt.Printf("WARNING --->  no stat info: %s\n", pathname)
		os.Exit(1)
	}
	if info.IsDir() {
		// do nothing
	} else {
		// TODO(mdr) test for regular file? skip pipes and devices etc
		digest, err := mdr.FileSHA256(pathname)
		if err != nil {
			fmt.Printf("!Err --->  SHA256 failed on %s\n", pathname)
			os.Exit(1)
		}
		var digPath DigestFullPathT = DigestFullPathT{
			Digest:   digest,
			FullPath: pathname,
		}
		gDigests = append(gDigests, &digPath)
		Verbose.Printf("adding %d %v\n", len(gDigests), digPath)
	}
	return nil
}

// buildDirMap creates gDigestMap as a side effect, maps from filename to digest
// for later lookups by server, sets a flag when done. server should check this
// flag before going into full service mode (send "please retry" if needed until
// fully ready.
func buildDirMap(pathName string) {
	stats, err := os.Stat(pathName)
	if err != nil {
		fmt.Printf("Can't get fileinfo for %s\n", pathName)
		os.Exit(1)
	}
	if stats.IsDir() {
		Verbose.Printf("%s is a directory, walking starts now\n", pathName)
		filepath.Walk(pathName, buildSHA256List)
	} else {
		fmt.Printf("this argument must be a directory (but %s isn't)\n", pathName)
		os.Exit(-1)
	}
	Verbose.Printf("map[digest]filename built with %d items\n", len(gDigests))
	for i := 0; i < len(gDigests); i++ {
		x := gDigests[i]
		gDigestMap[x.FullPath] = x.Digest
	}
	serverMapIsReady = true
	// TODO(mdr): print hashrate here
}

// doServer runs on destination system, serves up digests on request.
func doServer(requestChan chan SyncRequest, replyChan chan SyncReply) {
	// loops forever
	for {
		for {
			if serverMapIsReady {
				break
			}
			time.Sleep(time.Second)
		}
		// select should block till a request arrives
		select {
		case request := <-requestChan:
			Verbose.Printf("Rcvd a request: Id(%d) CmdType(%d) File %s\n", request.Id,
				request.CmdType, request.FilePath)
			requestedFile := filepath.Join(flagDest, request.FilePath)
			switch request.CmdType {
			case RefreshDigest:
				hashUpdate.Lock()
				log.Printf("Refreshing digest for %s\n", requestedFile)
				info, err := os.Stat(requestedFile)
				if err != nil {
					log.Printf("!Err --->  stat failed on %s\n", requestedFile)
					os.Exit(1)
				}
				if !info.Mode().IsRegular() {
					log.Printf("!Err ---> %s is not a regular file\n", requestedFile)
					return
				}
				digest, err := mdr.FileSHA256(requestedFile)
				if err != nil {
					fmt.Printf("!Err --->  SHA256 failed on %s\n", requestedFile)
					os.Exit(1)
				}
				//var digPath DigestFullPathT = DigestFullPathT{
				//	Digest:   digest,
				//	FullPath: requestedFile,
				//}
				gDigestMap[requestedFile] = digest
				hashUpdate.Unlock()
			case LookupDigest:
				hashUpdate.Lock()
				val, ok := gDigestMap[requestedFile]
				if ok {
					Verbose.Printf("Digest for %s is : %s\n", requestedFile, val)
					var reply SyncReply = SyncReply{
						Id:       request.Id,
						Digest:   val,
						FilePath: request.FilePath,
					}
					replyChan <- reply
				} else {
					fmt.Printf("Digest for %s is unknown\n", requestedFile)
					var reply SyncReply = SyncReply{
						Id:       request.Id,
						Digest:   "",
						FilePath: request.FilePath,
					}
					replyChan <- reply
				}
				hashUpdate.Unlock()
			}
		}
	}
}

// safeCopy runs on sender and uses os.Exec of "scp" for the heavy lifting.
func safeCopy(sourceName, destName string) error {
	hashUpdate.Lock()
	defer hashUpdate.Unlock()
	serverMapIsReady = false
	defer func() {
		serverMapIsReady = true
	}()
	fmt.Printf("ready to exec command: scp -pBC %s %s\n", sourceName, destName)
	if true {
		cmd := exec.Command("/usr/bin/scp", "-pBC", sourceName, destName)
		// can we leave Stdin unset?
		//cmd.Stdin = strings.NewReader("Some Input")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Printf("exec error %v\n", err)
			return err
		}
		fmt.Printf("command output = %v\n", out)
		Verbose.Printf("command ran with no errors\n")
	}
	return nil
}

// flagSetup
func flagSetup() error {
	if flagVerbose {
		Verbose = VerboseType(true)
	}
	if flagVersion {
		fmt.Printf("Version %s\n", Version)
		os.Exit(0)
	}
	invokedByName = os.Args[0]
	invokedByName = filepath.Base(invokedByName)

	if len(flagSource) > 0 {
		hasSource = true
		Verbose.Printf("Has source is true, %s\n", flagSource)
		if flagSource[0] != '/' {
			return fmt.Errorf("Source directory path must begin with slash\n")
		}
	}
	if len(flagDest) > 0 {
		hasDest = true
		Verbose.Printf("Has dest is true, %s\n", flagDest)
		if flagDest[0] != '/' {
			return fmt.Errorf("Destination directory path must begin with slash\n")
		}
	}
	// this is ok during testing
	//	if hasSource && hasDest {
	//		return fmt.Errorf("cant have both source and dest flags")
	//	}

	if flagPortNum > 1024 {
		hasPort = true
	}
	Verbose.Printf("Will run on port %d\n", flagPortNum)
	if hasSource && hasPort {
		return nil
	}
	if hasDest && hasPort {
		return nil
	}
	return fmt.Errorf("flag setup failed")
}

func main() {
	// fmt.Printf("%s will run for 10 minutes then quit\n", Version)
	flag.Parse()
	err := flagSetup()
	if err != nil {
		usage()
		log.Panicf("flagSetup failed\n")
	}

	requestChan = make(chan SyncRequest)
	replyChan = make(chan SyncReply)
	// this is the receiving side
	if hasDest && hasPort {
		Verbose.Printf("Destination setup\n")
		buildDirMap(flagDest)
		go doServer(requestChan, replyChan)
		fmt.Printf("%s finished normally\n", invokedByName)
	}
	time.Sleep(200 * time.Millisecond) // allow a bit for receiver to finish setup

	// this is the sending side - not recursive for test phase
	if hasSource && hasPort {
		Verbose.Printf("Source setup\n")
		fileList, err := filepath.Glob(flagSource + "/*")
		if err != nil {
			log.Panic("here")
		}
		fileList2, err := filepath.Glob(flagSource + "/moreBins/*")
		if err != nil {
			log.Panic("here")
		}
		for i := 0; i < len(fileList2); i++ {
			fileList = append(fileList, fileList2[i])
		}
		for i := 0; i < len(fileList); i++ {
			fmt.Printf("sending %d %s\n", i, fileList[i])
		}
		for i := 0; i < len(fileList); i++ {
			// skip directories if any
			fname := fileList[i]
			info, err := os.Stat(fname)
			if err != nil {
				log.Panicf("!Err --->  stat failed on %s\n", fname)
			}
			if !info.Mode().IsRegular() {
				log.Printf("!Warning ---> %s is not a regular file\n", fname)
				continue
			}

			fmt.Printf("\n\n")
			// create name relative to top of source tree
			relativeName := fileList[i]
			relativeName = relativeName[len(flagSource)+1:] // remove sending directory
			log.Printf("Requesting dest's hash of %s\n", relativeName)

			var req SyncRequest = SyncRequest{
				Id:       i,
				CmdType:  LookupDigest,
				FilePath: relativeName,
			}
			myHash, err := mdr.FileSHA256(fileList[i]) // use full path here
			if err != nil {
				log.Panic("here")
			}
			log.Printf("sending hash is %s\n", myHash)
			requestChan <- req
			select {
			case reply := <-replyChan:
				Verbose.Printf("reply is: Id(%d) Digest(%s) File(%s)\n", reply.Id,
					reply.Digest, reply.FilePath)
				if myHash != reply.Digest {
					sourceName := filepath.Join(flagSource, req.FilePath)
					destName := filepath.Join(flagDest, req.FilePath)
					log.Printf("Need to copy %s to %s\n", sourceName, destName)
					err := safeCopy(sourceName, destName)
					if err != nil {
						log.Printf("!Err ---> safeCopy failed with %v\n", err)
						return
					}
					if true {
						var req SyncRequest = SyncRequest{
							Id:       i,
							CmdType:  RefreshDigest,
							FilePath: reply.FilePath,
						}
						log.Printf("requesting hash refresh for %s\n", req.FilePath)
						requestChan <- req
					}
				} else {
					//fmt.Printf("%s has same digest in both source and dest\n", req.FilePath)
				}
			}
		}

		fmt.Printf("%s finished normally\n", invokedByName)
	}
	//time.Sleep(3 * time.Second) // server will shutdown on wakeup
	// should send quit signal to doServer instead or use waitgroup
	fmt.Printf("%s finished normally\n", invokedByName)
}
