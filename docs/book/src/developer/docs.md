# Docs

The markdown files are in the directory `docs/book/src`.

You can build the docs like this:

```
make -C docs/book/ build
```

Then open your web browser:
```
open ./docs/book/book/index.html
```

Via Github Actions, the docs get uploaded to: [hivelocity.github.io/cluster-api-provider-hivelocity/](https://hivelocity.github.io/cluster-api-provider-hivelocity/)

You can check links by using [nektos/act](https://github.com/nektos/act):

```
act -j markdown-link-check
```

