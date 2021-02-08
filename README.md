# skewer is ...

It rebuilds the Go source.
"embed" was implemented from Go1.16, and I wanted to reload HTTP server template.

# install 

```
go install github.com/secondarykey/skewer/_cmd/skewer
```

You should able to "skewer" command in "GOBIN" or "GOPATH/bin".

# Command

Take an argument to build the Go source.

```
skewer main.go
```

Please read Help for details.

# Operation Explanation

It usually looks for "go.mod" and monitors and builds all directories under it.
If there is an updatei in the monitored directories,it will be rebuild and the process will be starter.

Build temporarily and create a binary file to make the build error easier to understand and process status easier to understand.

At the beginning of development, it was realized as a proxy server, but I realized that it was not necessary and deleted it.

# Issue

I will write a memo for development.

## Argument

-t Change the directory to monitor
-a Specifying arguments when stating a process.
-d Monitoring lap time(default 5s)

