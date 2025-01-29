# indexer

This is a simple indexer that reads a directory and indexes all the go files in it. It makes to us easier to search for a function or a variable in a project.

## Usage

```bash
git clone github.com/snowmerak/indexer
cd indexer

podman compose up
# docker compose up

go install
```

And move to the project you want to index and run:

```bash
indexer new
```

Then indexer makes a new index file(`config.yaml`) in the project directory.

And initialize the database and index the project:

```bash
# initialize the database
indexer init

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

If you want to use other implementation of the interfaces, you can edit the `main.go` file to use the implementation you want.  
You can check the implementations in the `pkg` package and the interfaces in the `lib` package.
