# URL vector

URL parser based on [Vector API](https://github.com/koykov/vector) with minimum memory consumption.
This project is a part of policy to reducing memory consumption and amount of pointers.

Similar to [JSON vector](https://github.com/koykov/jsonvector) this package doesn't make a copy of the parsed URL.
It just stores an array with indexes of each part of the url (schema, hostname, path, ...) and doesn't contain any pointer.
