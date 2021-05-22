# URL vector

URL parser based on [Vector API](https://github.com/koykov/vector) with minimum memory consumption.
This project is a part of policy to reducing memory consumption and amount of pointers.

Similar to [JSON vector](https://github.com/koykov/jsonvector) this package doesn't make a copy of the parsed URL.
It just stores an array with indexes of each part of the url (schema, hostname, path, ...) and doesn't contain any pointer.

## Usage

```go
	src := "http://x.com/x/y/z?arr[]=1&arr[]=2&arr[]=3&b=x&arr1[]=a&arr1[]=b&arr1[]=c"
vec := urlvector.Acquire()
defer urlvector.Release(vec)
_ = vec.ParseStr(src)
fmt.Println("host", vec.HostString())
fmt.Println("path", vec.PathString())
fmt.Println("query:")
vec.Query().Get("arr[]").Each(func(_ int, node *vector.Node){
    fmt.Println("arr[] ->", node.String())
})
fmt.Println("b ->", vec.Query().Get("b"))
vec.Query().Get("arr1[]").Each(func(_ int, node *vector.Node){
    fmt.Println("arr1[] ->", node.String())
})
```

Output:
```
host x.com
path /x/y/z
query:
arr[] -> 1
arr[] -> 2
arr[] -> 3
b -> x
arr1[] -> a
arr1[] -> b
arr1[] -> c
```
