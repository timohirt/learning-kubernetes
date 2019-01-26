# Setup go dev environment

All the services I use while learning kubernetes are written in Go. I'm new to
Go and part of this project is to also learn how to write Go code. So, if you
find anything that is not proper Go or agains best practices, I would be really
happy for feedback, being it a pull request or via Twitter 
([@timohirt](https://twitter.com/TimoHirt)).

## Dependency Management

`dep` is used for dependency management. For more details take a look at [the
project GitHub](https://github.com/golang/dep).

On OSX you can easily install it with brew:

```bash
brew install dep
```

## GOPATH

All go code projects are located in `services/src`. Set the `$GOPATH` to
`/your/path/to/learning-kubernetes/services`.

If you are using tools like autoenv, you can just put the following into your
`.env`.

```bash
export GOPATH=/Users/name/learning-kubernetes/services 
```

