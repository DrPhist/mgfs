mgfs
========

Allows mounting a MongoDb database as a file system via [FUSE](https://bazil.org/fuse/).

## Installation
You need to have [Golang](https://golang.org/doc/install) installed.
Open your terminal, and run `go install github.com/amsa/mgfs`. Now you should be able to
run `mgfs`.

## How to use
First mount your MongoDb database: `mgfs test /path/to/mount/dir`. You may now go to the directory specified 
as the mount point, and see the collections (directories), and documents (json files). You may read, update, 
or delete the documents. You may also read GridFs files under the specified prefix (`fs` by default).

Don't forget to unmount the database when you are done (`umount /path/to/mount/dir`).

## Todo
- [x] Support GridFS read 
- [ ] Support GridFS write
- [ ] Show GridFS file names

## Credits
* Uses [bazil.org/fuse](http://bazil.org/fuse)
* Uses [labix.org/mgo](http://labix.org/mgo)
