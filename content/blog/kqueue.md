---
title: Detecting file changes on macOS with kqueue
slug: kqueue
date: 2026-03-24
collapsible_code: false
---

A while ago I wrote [a small file watcher in Go for my own use](https://github.com/vegarsti/reload) with an [accompanying blog post](/file-watcher).
I needed a tool that I could just plop in front of the command I was running as part of my iteration loop.
I use it for recompiling C files when I modify them,

```bash
reload gcc main.c -o main && ./main
```

and for rebuilding and reloading my static site on file changes.

```bash
reload make
```

It has two modes.
1. If one or more explicit files are mentioned in the command, it will watch those files.
2. If it does not find any filenames, it will watch all files in the working directory.

The only thing `reload` needs to know is whether any file it is watching has changed.
If a file has changed, it reruns the command.

It works great! But I copped out on the part that was the most unfamiliar to me, namely [detecting file changes](/file-watcher/#detecting-file-changes).
I used the [fsnotify](https://github.com/fsnotify/fsnotify) package which is a nice cross-platform Go library for this.
It supports macOS as well as Linux but since it's for my own use, I don't care about Linux support.
More importantly, I wanted to understand what fsnotify did under the hood.

On macOS it uses the [kqueue](https://en.wikipedia.org/wiki/Kqueue) event notification interface.
Let's take a look at how this works, writing some C code to test it, and then finally implement it in the reload Go program.

## `kqueue` data structures
The `kqueue()` [function call](https://man.freebsd.org/cgi/man.cgi?kqueue) creates a new kernel event queue (a kqueue) and returns a file descriptor.
We register and wait for system events using the `kevent()` function call, which uses the `kevent` data structure.
This has five fields that we need to care about.

- `ident` The source of the event. In this case, it will be a file descriptor for the file we want to watch.
- `filter` The kernel filter used to process the event.
- `flags` Actions to perform on the event.
- `fflags` Filter-specific flags.
- `udata` Opaque user data identifier

Well, which kernel filter do we use if we want to watch a file for changes?
There are 9 possible filters, but the one we're looking for is `EVFILT_VNODE`.

```
EVFILT_VNODE   Takes a file descriptor as the identifier and the events to watch for in fflags,
               and returns when one or more of the requested events occurs on the descriptor.
```

It goes on to list 10 possible events that can be used in `fflags`, but the only flag we need is `NOTE_WRITE`.

```
NOTE_WRITE     A write occurred on the file referenced by the descriptor.
```

Finally, `flags` defines the actions to perform on the event.
There are 10 possible flags, but we only need `EV_ADD` to register the event in the `kqueue`, and `EV_CLEAR` to reset the event state after delivery.
Without `EV_CLEAR`, we'd get the first file change again and again - try it!

If there are multiple writes to this file before we retrieve the event from the kqueue, these will be collapsed into one event.
However, if multiple files are written, these will be distinct events.
So in practice, in `reload`, we'll likely want a window such that we don't rerun the command on every single event if they are close in time.

We now have all we need to initialize a kevent structure.

## Watching named files
Let's write a program to watch all files passed as arguments.
We omit error handling for rare errors for brevity, but keep for stuff that might happen such as trying to open a file that doesn't exist.

First we open the files to watch with `O_EVTONLY`.
The `open` man pages say the event-only mode is only intended for monitoring a file for changes, such as with `kqueue`.
We'll need the following headers:

```c
#include <sys/types.h>
#include <sys/event.h>
#include <sys/time.h>
```

We create an array of kevent structs representing the changes we care about,
namely writes to files provided on the command-line.
We use the `EV_SET` macro to initialize the kevent struct.

```c
// Open all files and set up change events
int nfiles = argc - 1;
int *fds = malloc(nfiles * sizeof(int));
struct kevent *changes = malloc(nfiles * sizeof(struct kevent));
for (int i = 0; i < nfiles; i++) {
    fds[i] = open(argv[i + 1], O_EVTONLY);
    if (fds[i] == -1) {
        fprintf(stderr, "open(%s): ", argv[i + 1]);
        exit(1);
    }

    EV_SET(
        &changes[i],
        fds[i],
        EVFILT_VNODE,
        EV_ADD | EV_CLEAR,
        NOTE_WRITE,
        0,
        (void *)argv[i + 1] // udata: opaque user data. Store filename here.
    );
}
```

We register the events with a call to `kevent()`.

```c
// Register all events at once
int kq = kqueue();
kevent(
    kq, // the queue
    changes, // array of kevent events to register
    nfiles, // length of array
    NULL, // struct to populate with event (not used here)
    0, // number of events to wait for
    NULL // timeout if waiting for event (irrelevant here)
);
```

We're now ready to start an event loop and listen to file changes.
We use `kevent()` for this as well.
After we get an event back, we can look at `event.fflags` to see which file event was emitted.
In our case we're only listening for `NOTE_WRITE` so that should always be true.

```c
struct kevent event;
while (1) {
    kevent(
        kq, // queue
        NULL, // array of events to register
        0, // no events to register
        &event, // struct to populate with event
        1, // number of events to wait for
        NULL // no timeout; wait forever
    );

    if (event.fflags & NOTE_WRITE) {
        const char *name = (const char *)event.udata;
        printf("[%s] written\n", name);
    }
}
```

The full code for this is in [this GitHub Gist](https://gist.github.com/vegarsti/2fca3bc234839b1ccb3bc33ff03e535d).

## Watching a directory
In the second mode of reload, we watch the current working directory for any file changes.

```bash
reload make
```

Let's first look at how to watch a directory.
First we open the directory itself, and we watch it like we did above for a single file.
```c
int fd = open(directory, O_EVTONLY);
```
This emits events for new files added to the directory and for file deletions.
Such changes involve writing to the file on disk that represents the directory.

However, this does not emit events when there are changes to an existing file, so it's not sufficient for our use.
We need to open all files within the directory and watch these for changes individually!
This means that when a file is created, we need to add this file to be watched.

I wrote the reload program in Go, so we now switch over to looking at Go code.

## Implementing in Go

We need a reference to a kqueue, and we need to keep track of which file descriptors we open.
We also want to be able to refer to files by path.

```go
type watcher struct {
	kq      int            // kqueue file descriptor
	fds     map[string]int // path -> file descriptor mapping
	fdPaths map[int]string // file descriptor -> path mapping (reverse lookup)
}
```

Notice the CloseOnExec call when we create the kqueue below.
When `reload` re-runs your command, Go's exec package uses the fork + exec pattern.
A fork clones the parent's open file descriptors (fd) into the child process,
including our kqueue fd and every file descriptor we have opened for watching.
The child process doesn't need these, and leaking them can cause subtle problems.
The watched files can't be fully deleted by the operating system since their reference count never hits zero, and the child holds a kqueue it never drains.

The `O_CLOEXEC` flag tells the kernel to automatically close these file descriptors when exec runs.
We set it on the kqueue itself here, and on each watched file in `Add` below.

```go
func newWatcher() (*watcher, error) {
	kq, err := unix.Kqueue()
	if err != nil {
		return nil, err
	}

	// Set close-on-exec flag so child processes don't inherit the kqueue fd
	unix.CloseOnExec(kq)
	return &watcher{
		kq:      kq,
		fds:     make(map[string]int),
		fdPaths: make(map[int]string),
	}, nil
}
```

When adding a file to be watched, we open the file and register for file writes on the kqueue.

```go
func (w *watcher) Add(path string) error {
	// Skip if already watching
	if _, exists := w.fds[path]; exists {
		return nil
	}

	// Open file with event-only flag to only get events,
	// and close on exec, so that exec'ing the command to reload
	// does not copy the file descriptors.
	fd, err := unix.Open(path, unix.O_EVTONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return err
	}

	// Register file/directory
	w.fds[path] = fd
	w.fdPaths[fd] = path
	_, err = unix.Kevent(
		w.kq,
		// changes to look for
		[]unix.Kevent_t{{
			Ident:  uint64(fd),
			Filter: unix.EVFILT_VNODE,
			Flags:  unix.EV_ADD | unix.EV_CLEAR,
			Fflags: uint32(unix.NOTE_WRITE),
		}},
		nil, // events to populate: none here
		nil, // no timeout; not populating events anyway
	)

	if err != nil {
		unix.Close(fd)
		return err
	}

	return nil
}
```

To add all files in a directory, we walk the directory tree and call `Add` for each entry, be it a file or nested directory.
```go
func (w *watcher) addRecursive(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// If it's a file, watch it directly
	if !info.IsDir() {
		return w.Add(path)
	}

	// Walk the directory tree and add all directories and files
	return filepath.WalkDir(path, func(walkPath string, d os.DirEntry, err error) error {
		// non-nil when WalkDir fails to access walkPath
		if err != nil {
			return err
		}

		// Watch the directory for new files and deletions
		err = w.Add(walkPath)
		if err != nil {
			return err
		}
		return nil
	})
}
```

Finally, we set up the event loop where we listen on the kqueue until we get a file event.
Once we do, we return back so the main program can reload the command.

```go
// Wait until any file changes, and return the filename of the changed file
func (w *watcher) Wait() (string, error) {
	events := make([]unix.Kevent_t, 1)
	for {
		n, err := unix.Kevent(
			w.kq,   // the queue
			nil,    // changes
			events, // events to populate
			nil,    // no timeout
		)

		if err != nil {
			if err == unix.EINTR {
				continue
			}
			return "", err
		}

		if n > 0 {
			path, ok := w.fdPaths[int(events[0].Ident)]
			if ok {
				return path, nil
			}
		}
	}
}
```

As I mentioned, if a file is created in a directory, we need to add that file to be watched as well.
We solve that by re-walking the directory if the file event comes from a directory and not a file.

```go
// If the file event was from a directory, it may have
// created a new file. If so, we need to add it!
info, _ := os.Stat(path)
if info != nil && info.IsDir() {
	// Re-walk to pick up new files
	watcher.addRecursive(path)
}
```

Note that we don't remove a file if it was deleted.
This means we are leaking file descriptors when files are deleted, and causing the OS to be unable to clean these up.
But for my use case it's so rare that I both delete a file and need to keep the program running for a long time that we don't cover this here.

## Final thoughts

This was fun!
I did not know about kqueue before, but I learned that it's easy to work with, and can be used for many inter process communication use cases on macOS.

Another way to solve the file watching problem would be the naive way, that is, to poll for changes.
Since kqueue requires an open fd per watched file, kqueue won't scale to very large directory trees.
The third way this could be done on macOS is to use [FSEvents](https://en.wikipedia.org/wiki/FSEvents), which does not have the problem of fd exhaustion.

I would be very happy if you checked out [the reload code on GitHub](https://github.com/vegarsti/reload), tried it out on your Mac, or reached out to me with any thoughts or experiences doing something similar.

Finally, if you're interested in reading more about kqueue,
you can read [this post](https://nima101.github.io/kqueue_server) about implementing a bidirectional streaming server using kqueue in C,
and [the original paper](https://people.freebsd.org/~jlemon/papers/kqueue.pdf).
