# kvas

`kvas` is a minimal overhead key value store backed by the filesystem. 

It comes with the following features: 
- writes only happen for new data (verified with a SHA-256 hash); 
- modifications are timestamped, allowing querying last modified by a timestamp; 
- an in-memory index, allowing fast operations like confirming whether value set contains a key, enumerations;
- specific values in a set are stored individually as files;
- reductions - a slice of property values per key for fast in-memory operations like finding all keys that match specific value.
- list of reduction - a slice of reductions with additional fabric that establishes data relationships for transparent common operations (aliasing, joins).

## Using `kvas` in your app

- Run go get `github.com/boggydigital/kvas`
- Import `github.com/boggydigital/kvas`
- For a key value set - connect to a local store with `kvas.ConnectLocal(directory, extension)`
  - Extensions supported: `kvas.JsonExt (.json)`, `kvas.GobExt (.gob)`
- Use this key value set client to `Get`, `Set`, `Cut` values, as well as filter using `CreatedAfter`, `ModifiedAfter`, etc.

## Key types provided by `kvas`

`kvas` comes with the following types:
- `KeyValues` - key value store backed by local filesystem
- `ReduxValues` - key reductions store (backed by `kvas.KeyValues`)
- `ReduxAssets` - collection of `ReduxValues` plus additional data relationship fabric

## Example usage of `kvas.KeyValues`

NOTE: Error handling omitted for brevity.

```go
lkv, _ := kvas.ConnectLocal(os.TempDir(), kvas.JsonExt)

key := "value1"
start := time.Now().Unix()

_ = lkv.Set(key, strings.NewReader(key))

if lkv.Has(key) {
    readCloser, _ := lkv.Get(key)
    defer readCloser.Close()
    // use readCloser to read data stored under "value1" key
}

fmt.Println(lkv.ModifiedAfter(start, false)) // prints: [value1]
```

## Example usage of `kvas.ReduxValues`

NOTE: Error handling omitted for brevity.

```go
rdx, _ := kvas.ConnectRedux(os.TempDir(), "titles")

for _, key := range rdx.Keys() {
    allKeyTitles, _ := rdx.GetAllValues(key)
    // process all titles for a key
}

// find all keys that have values containing "specific_title" (any case) 
matches := rdx.Match([]string{"specific_title"}, nil, true, true)
```

## Example usage of `kvas.ReduxAssets`

NOTE: Error handling omitted for brevity.

```go
rxa, _ := kvas.ConnectReduxAssets(os.TempDir(), nil, "titles", "country")

query := map[string][]string{
    "title": {"title1", "title2"},
    "country": {"states"},
}

// find all keys across "titles" and "countries" reductions that match query
matches := rxa.Match(query, true)
```