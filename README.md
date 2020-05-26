# Garbage Diving with GO

Code for talk of the same name by [@sleepypioneer](https://twitter.com/sleepypioneer) you will find the slides on [speaker deck](https://speakerdeck.com/sleepypioneer/garbage-diving-for-gophers)

What is Garbage collection and why do we need it? How does Garbage collection happen and when? How can we as engineers be more sympathetic to the machines our code runs on? If you've ever found yourself asking any of these questions or are curious to know more, this is the talk for you.

## Full talk description

Garbage collection in programming languages ranges from none to fully configurable (ie in Java) Go keeps things simple but that doesn't mean we should take it for granted. Understanding how our code runs under the hood allows us to build sympathy for the machine it runs on and in turn allows us to write informed code which optimizes where possible for these mechanics.
Last year I found myself digging into Go's garbage collector while researching for a book project which is currently on hold. I enjoyed learning about it so much I thought it'd be great to put my findings together and share them. This talk will cover what is garbage collection, how golangs's garbage collection is implemented and when it runs as well as exploring how armed with this knowledge, we as developers can write our code to provide better effectiveness of this mechanism.
If you enjoy talks which take a deep dive on one aspect of a language I hope you will enjoy this one and come out with some interesting insights!

## Speaker bio

Jessica Greene is a backend developer at Ecosia.org interested in cloud computing, devops, K8s, AI & IoT! She codes mainly in Go, Python & JavaScript. A career changer her previous activities have seen her on film sets and coffee farms. She is a self taught/community taught developer. Find her at local Meetups or at climate marches!

## About this repo

Please note this repo is intended to be a programatic aid to the talk and is not yet documented to a level to make it a useful independant learning tool.

### Memory management in GO

First we will create an allocation in our program which we know will escape to the heap.

```go
func myFunc() {
    _ = make([]int, 1000000)
}
```

### Heap or Stack

We can check this did indeed escape to the heap with the following command:

```sh
go run -gcflags -m main.go
```

In it's return we see that indeed our slice has escaped to the heap, and that the length has also escaped, this is because we passed it to the function fmt.Println():

```sh
.\main.go:15:11: make([]int, 1000000) escapes to heap
.\main.go:16:17: len(a) escapes to heap
```

## Tracing our program

Let's add tracing to our program so we can dig deeper. Uncomment the line `tracer.WithTrace(repeatXTimes, 10, myFunc)` and comment out first command in `main()` `repeatXTimes(10, myFunc)` then run the following commands:

```sh
# Now we restricted the progame to one process thread for simplicity.
# We also add the debugger option to trace our GC in the cli.
env GODEBUG=gctrace=1 GOMAXPROCS=1 go run -gcflags -m main.go
```

We see in the output from the gctrace that our GC runs 10 times, once each time our function that allocates runs.
We can also see here some more interesting information:

```sh
gc 4 @0.278s 0%: 0+5.0+0 ms clock, 0+0.99/0/0+0 ms cpu, 4->4->0 MB, 5 MB goal, 1 P
gc 5 @0.558s 0%: 0+6.0+0 ms clock, 0+0/0/0+0 ms cpu, 4->5->1 MB, 5 MB goal, 1 P
gc 6 @0.610s 0%: 0+8.9+0 ms clock, 0+0/1.9/0+0 ms cpu, 4->5->1 MB, 5 MB goal, 1 P
gc 7 @0.642s 0%: 0+0.99+0 ms clock, 0+0/0/0+0 ms cpu, 4->4->0 MB, 5 MB goal, 1 P
gc 8 @0.680s 0%: 0+2.9+0 ms clock, 0+0.99/0/0+0 ms cpu, 4->4->1 MB, 5 MB goal, 1 P
```

The format of this output changes with every version of Go, currently, it is:
`gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P`
where the fields are as follows:

```sh
gc #        the GC number, incremented at each GC
@#s         time in seconds since program start
#%          percentage of time spent in GC since program start
#+...+#     wall-clock/CPU times for the phases of the GC
#->#-># MB  heap size at GC start, at GC end, and live heap
# MB goal   goal heap size
# P         number of processors used
```

[Reference](https://golang.org/pkg/runtime/)

We can also see information about the internal state of the concurrent pacer adding the option `gcpacertrace=1` to the GODEBUG variable:

```sh
env GODEBUG=gctrace=1,gcpacertrace=1  GOMAXPROCS=1 go run -gcflags -m main.go
```

Next lets run the trace tool to a see a visualisation of this:

```go
go tool trace trace.out
```

This will open up a browser tab, lets navigate to view trace, the first option in the list, also found at endpoint `/trace`.

We can clearly see the steps of Garbage collection taking place, the phases are stop-the-world (STW) sweep termination, concurrent mark and scan, and STW mark termination. Additionally we can see the heap value go up and down as allocations are made and the GC cleans them up.

### The GOGC control

The only control point that we have for the garbage collector is GOGC, the default value is 100, which means garbage collection will not be triggered until the heap has grown by 100% (effectively doubled) since the previous collection.

So if after the first run of GC the heap has 2mb of in-use memory then the pacer will schedule the next collection just before the heap reaches 4mb. This is assesed at the end of each run.

Setting this value higher, say GOGC=400, will delay the start of a garbage collection cycle until the live heap has grown to 400% of the previous size.

```sh
env GODEBUG=gctrace=1 GOGC=400 GOMAXPROCS=1 go run -gcflags -m main.go
```

Next let's try setting the value lower, say GOGC=20 this  will cause the garbage collector to be triggered more often as less new data can be allocated on the heap before triggering a collection.

```sh
env GODEBUG=gctrace=1 GOGC=20 GOMAXPROCS=1 go run -gcflags -m main.go
```

Setting `GOGC=off` will disable garbage collection entirely.

We can also manually run GC in our program with the runtime method `runtime.GC()` though it's important to note this will block the caller.
