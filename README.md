# sif

Golang app which removes unused dependencies in your Bazel projects. For example, say you have dependency graph like this:

```
         --- //app:useful_dep_a
        /
//app:app --- //app:useful_dep_b --- //app:useless_dep
                                  \- //...


```

Sif can find that //app:useless_dep is redundant dependency in //app:useful_dep_b and remove it, while ensuring that //app:app still compiles. It is done with some smart brtuteforce made fast by Bazel caches. It does parse BUILD files without actual Starlark interpreter, so it can fail to build a full build graph (if your deps are generated using a function, or if there is a weird macro generating a target, or ...), but it is generally quite useful. It is still in development.

TBD:
* Optimize by multuple params at the same time (say, deps and sources)
* Recursive optimization

Building:
```
go build .
```

Usage example:
```
./sif --workspace test/cppexample --label //main:hello-world --param deps,hdrs --recparams deps
```

Params:
* --workspace: Path to workspace where targets are, if not current folder
* --label: Label of target to optimize
* --params: Params to optimize. Split with ","
* --check: List of additional targets (tests) to ensure still compile (pass) while optimizing. Split with ","
* -v: Verbose mode
* --recparams: Params to recursivly optimize dependency graph
* --recblacklist: Regexp to filter out unwanted recursive optimization targets