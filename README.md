# Kvas

`kvas` is a key value store with the following features: writes only happen for new data (verified with a SHA-256); modifications are timestamped, allowing querying last modified by a timestamp. Kvas builds an index and stores index in memory, allowing fast operations like count and contains. Values are stored individually as files.

# Implementation plans

## Connecting 

Clients need to Connect, specifying the base directory (value set, e.g. "products") that would be used to transform relative data set path to an absolute filepath.

## Keys

- There is a single key used for data access.
- The key is a string constructed as a file path, e.g. `kvasProducts.Get("movies/12345")`, the path components are:
    - Values set initialized at connection that specifies base directory (`kvasProducts` in this example, initialized like `kvasProducts := kvas.Connect("products", ".json")`)
    - Values subset that specifies folders where file is stored (`movies` in this example)
    - ID - the last path component, that identifies an individual file (`12345`)
- There is no limit on the number of "folders" in a values set portion
- IDs are not globally unique, only unique to the values set
- Values are stored in the files that have ID as a name and values set as the full directory structure (`products/movies` are the relative directories and `12345.json` in this example)
- All values are contained within initial connection entry point (`products` in this example).
- Index file is stored at the root of the values set

## Index

Index holds the key metadata for a set of entries in a data set, for each entry the following is stored:
- ID
- Hash of that currently stored value
- Timestamp when value was Created
- Timestamp when value was Modified

## Operations

There are basic operations for creating (Create), reading (Get), updating (Update), deleting (Delete) records.

### Additional operations

- Contains(key): Check if the specific value exists.
- All(dataSet): Return all IDs for a data set.
- CreatedAfter(timestamp): Return all IDs for a values set created after a timestamp.
- ModifiedAfter(timestamp): Return all IDs for a values set modified after a timestamp.

## Prevent write operation if data is the same

For each entry index holds the hash of a last known written state. Value writes happen only if the hash of that entry is different from currently stored value hash.

## Clients, SDKs

The implementation is a golang module that can be used in other golang apps. In addition to that there are two more ways to use `kvas` is with the CLI client, and a service that exposes REST API endpoints.