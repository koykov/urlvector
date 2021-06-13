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

## Performance

See [versus](https://github.com/koykov/versus) project for performance comparison between urlvector and [net/url](https://golang.org/pkg/net/url/) packages:
```
BenchmarkNetUrl_ParseUrl-8        	 1027659	      1139 ns/op	     192 B/op	       2 allocs/op
BenchmarkNetUrl_ParseQuery-8      	  153732	      7690 ns/op	    2919 B/op	      22 allocs/op
BenchmarkNetUrl_ModHost-8         	  704062	      1799 ns/op	     840 B/op	       7 allocs/op
BenchmarkNetUrl_ModQuery-8        	   86022	     14279 ns/op	    5649 B/op	      42 allocs/op
BenchmarkUrlvector_ParseUrl-8     	 1000000	      1107 ns/op	       0 B/op	       0 allocs/op
BenchmarkUrlvector_ParseQuery-8   	  337450	      3489 ns/op	       0 B/op	       0 allocs/op
BenchmarkUrlvector_ModHost-8      	  796296	      1317 ns/op	       0 B/op	       0 allocs/op
BenchmarkUrlvector_ModQuery-8     	  290557	      4110 ns/op	       0 B/op	       0 allocs/op
```
