# Kvas

`kvas` is a key value store with the following features: writes only happen for new data (verified with a SHA-256 hash); modifications are timestamped, allowing querying last modified by a timestamp. Kvas builds an in-memory index, allowing fast operations like confirming whether value set contains a key and enumerations. Values in a value set are stored individually as files.

## Using `kvas` in your app

- Run go get `github.com/boggydigital/kvas`
- In your app import `github.com/boggydgital/kvas`
- For a value set - create new JSON Client with `kvas.NewJsonClient(value_set_location)`
- Use this value set client as appropriate for your app

## Example usage of `kvas`

NOTE: Error handling omitted for brevity.

```
  vs, _ := kvas.NewJsonClient("test")

  key := "value1"
  ts := time.Now().Unix()

  _ = vs.Set(key, []byte(key))

  if vs.Contains(key) {
    bytes, _ := vs.Get(key)
	fmt.Println(string(bytes)) // prints: value1
  } 

  fmt.Println(vs.ModifiedAfter(ts)) // prints: [value1]
```

## 'kvas' operations

- `NewJsonClient`
- `Get`
- `Set`
- `Remove`
- `Contains`
- `All`
- `CreatedAfter`
- `ModifiedAfter`

## Frequently asked questions

- Q: Is `kvas` suitable for concurrent read/write operations?
- A: Yes with a note. `Kvas` itself makes no effort to protect indexes with locks, instead you're encouraged to leverage Go language features like channels (along with `select`).

