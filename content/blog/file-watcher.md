---
title: "A file watcher"
date: 2024-08-25
slug: file-watcher
---

I wrote [`reload`](https://github.com/vegarsti/reload), my own file watcher!

The last year and a half I've been making my way through [CSPrimer](https://csprimer.com/).
It consists of hundreds of computer science problems, so working through these means I'm writing hundreds of small programs.
Like most programming, it's a very iterative process where I tweak my program and rerun it, typically switching from my editor to run a command.

```
gcc main.c -o main && ./main
```

It would be nice if this just reran automatically whenever I needed it to, so I didn't have to switch between my editor and the terminal!

Of course, there exist tools like this. They're sometimes called file watchers or hot reloaders.
I tried a few existing tools such as [entr](https://github.com/eradman/entr) and [wach](https://github.com/quackingduck/wach) but they weren't _just right_.

So I figured I'd write my own.

After a few hours the result was [`reload`](https://github.com/vegarsti/reload). It's one `main.go` file with less than 200 lines. It does exactly what I need, nothing more.

### Requirements

I wanted to be able to just put `reload` in front of the command I'd usually run.
To do that, what `reload` needed to do was

- figure out which file(s) to watch
- figure out which command to run
- run the command
- listen to file changes
- reload the command on file changes

And it also needed to

- run command pipelines like `gcc main.c -o main && ./main` sequentially
- cancel the command if it's still running when there's a file change
- deduplicate file events -- so that a change to `main.c` and the `main` binary file don't trigger two subsequent reloads

There needs to be some concurrent code here, but Go excels at this with channels and goroutines.

In the rest of the post, let's look at how `reload`

1. Figures out which files to watch
2. Detects file changes
3. Runs the command

### Figuring out which files to watch

Some of these tools have you explicitly list the files that should trigger a reload.
For my use case, my command usually has the filename already.
Why not just grab that?

So `reload` looks at each part of the command, and checks if it's the filename of an existing file.
If it is, we add it to the list of files to watch.
If we didn't find any file, we fall back to watching the working directory.

Kinda neat!

```go
// Find files in command that we will watch
toWatch := make([]string, 0)
for _, part := range input {
	// Check if there's a file to watch
	info, err := os.Stat(part)
	if os.IsNotExist(err) {
		continue
	}
	check(err)
	if !slices.Contains(toWatch, info.Name()) {
		toWatch = append(toWatch, info.Name())
	}
}

// Fall back to watching the working directory
if len(toWatch) == 0 {
	wd, err := os.Getwd()
	check(err)
	toWatch = append(toWatch, wd)
}
```

### Detecting file changes

The most unfamiliar thing to me here was watching for file changes.

I used the [fsnotify](https://github.com/fsnotify/fsnotify) package to handle file notifications.
It supports all the usual platforms like macOS, Windows, and Linux.
On macOS it uses the [kqueue](https://en.wikipedia.org/wiki/Kqueue) notification interface, whereas on Linux it uses [inotify](https://en.wikipedia.org/wiki/Inotify).

Here's the code for that.

```go
watcher, err := fsnotify.NewWatcher()
check(err)
defer watcher.Close()

for _, file := range toWatch {
	err = watcher.Add(file)
	check(err)
}

fileChanges := make(chan string)
wg.Add(1)
go func() {
	defer wg.Done()
	lastChange := time.Now()
	dedupWindow := 100 * time.Millisecond
	for event := range watcher.Events {
		if event.Has(fsnotify.Write) {
			// Treat multiple events at same time as one
			if time.Since(lastChange) < dedupWindow {
				continue
			}
			lastChange = time.Now()
			fileChanges <- event.Name
		}
	}
}()
```

And we listen on the `fileChanges` channel in the code to run the command, which we will see now.

### Running commands

To create a command that we can run, and cancel, we use [`exec.CommandContext`](https://pkg.go.dev/os/exec#CommandContext) from Go's standard library.
We can specify its standard in, standard out, and standard error.
In our case, we want standard out and standard error to be the operating system's stdout and stderr, since we want the command to behave as if we ran it directly in the terminal.

[`sh -c '<command>'`](https://man7.org/linux/man-pages/man1/sh.1p.html) allows us to pass the full provided command. All use of pipes and operators like `&&` and `||` will just work, as it's handed over to the shell.

```go
// In `main()`
// ...

// First run the command
runCommand(ctx, command, fileChanges)

// Then rerun it on file changes
for name := range fileChanges {
	fmt.Fprintf(os.Stderr, "--- Changed: %s\n", name)
	fmt.Fprintf(os.Stderr, "--- Running: %s\n", command)
	runCommand(ctx, command, fileChanges)
}

// ...

func runCommand(ctx context.Context, command string, fileChanges chan string) {
	// Create child context so we can cancel this command
	// without cancelling the entire program
	commandCtx, commandCancel := context.WithCancel(ctx)
	defer commandCancel()

	// Cancel and rerun the command if the file changes
	go func() {
		name, ok := <-fileChanges
		// The channel was closed, shut down
		if !ok {
			return
		}
		commandCancel()
		fileChanges <- name
	}()

	// Run the command using `sh -c <command>` to allow for
	// shell syntax such as pipes and boolean operators
	cmd := exec.CommandContext(commandCtx, "sh", []string{"-c", command}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // It's fine if the command fails!
}
```

I hope you'll try it out!
