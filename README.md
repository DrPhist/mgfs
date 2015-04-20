mgfs
========

Allows mounting a MongoDb database as a file system via [FUSE](https://bazil.org/fuse/).


## Working
* `ls mountpoint/` lists all collections in that database
* `ls mountpoint/collection` lists all documents with their `_id` as filename
* `cat mountpoint/collection/document` returns the document content with `json.MarshalIndent`


## Todo
- [x] Collections are only read on startup
- [ ] Documents in `xxx.index` collections don't have `_id` fields, so they aren't listed yet
- [ ] Support writes on [GridFS](http://www.mongodb.org/display/DOCS/GridFS) (http://godoc.org/labix.org/v2/mgo#GridFS)
- [ ] Show file name instead of ID for GridFS files

## Credits
* Uses [bazil.org/fuse](http://bazil.org/fuse)
* Uses [labix.org/mgo](http://labix.org/mgo)
