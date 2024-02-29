# CacheFile
A simple lib to cache data locally using a file. When creating an instance you have to provide a function for data retrieval which will be used when CacheFile data is requested. If the local file has expired or is missing, the retrieval function will be triggered to generate the local file. Otherwise the cached data will be returned. 

## Installation
```bash
go get github.com/toxyl/cachefile
```

## Example
```bash
go run example/main.go
```
