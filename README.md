# indexer

This is a simple indexer that reads a directory and indexes all the go files in it. It makes to us easier to search for a function or a variable in a project.

## Usage

```bash
git clone github.com/snowmerak/indexer
cd indexer

podman compose up
# docker compose up

go install

# initialize the database
indexer init
```

And move to the project you want to index and run:

```bash
indexer index . # or the path you want to index
```

After that, you can search for a function or a variable in the project:

```bash
indexer search <query> <count>
```

When you want to remove the index, you can run:

```bash
indexer clean
```
