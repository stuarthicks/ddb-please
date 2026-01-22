# ddb-please

My own take on an easy to use cli for simple read operations against AWS DynamoDB. For anything other than simple list/get operations, you're probably better off just using the aws-sdk-go-v2 SDK directly or the awscli.

## Install

Nix:

    nix run github:stuarthicks/ddb-please

Using Homebrew:

    brew install stuarthicks/brews/ddb-please

Using Go:

    go install github.com/stuarthicks/ddb-please@latest

## Usage

```
Usage of ddb-please:
  -get
        Get an item
  -info
        Get info about a table (e.g. secondary indexes)
  -list
        List tables
  -p string
        Partition (HASH) Key
  -pretty
        Pretty-print JSON
  -region string
        AWS Region (default "us-east-1")
  -s string
        Sort (RANGE) Key
  -t string
        Name of table to query
````

### Examples

#### List tables
```
$ ddb-please -list
```

#### Get info on a table
```
$ ddb-please -info -t foo-table
Key Schema:
        id (HASH)

Global Secondary Indexes:
        account_id-index
                - account_id (HASH), id (RANGE)
```

#### Get an item by ID

Use `-p {partition_key}` and optionally `-s {sort_key}`.

Add `-pretty` to pretty-print the JSON.

```
$ ddb-please -get -t foo-table -p foo
{"id": "foo"}
```

## Notes

- Currently only supports partion/sort keys of String type.
- Query/Scan operations are not implemented (and so cannot get items via secondary indexes).
