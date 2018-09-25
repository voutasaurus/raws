# raws

Find all raw string literals in a Go project (searches all subdirectories).

# install

```
$ go get github.com/voutasaurus/raws
```

# usage

```
$ cd cbtest # or whatever your project is
$ raws
api/api.go
 158: `json:"err"`
database/database.go
  88: `json:"type"`
  89: `json:"x"`
 131: `SELECT {type,x} AS doc,
		  META(`
 132: `).id AS id
		  FROM `
 133: `
		  WHERE type == $type AND x LIKE $prefix`
```
