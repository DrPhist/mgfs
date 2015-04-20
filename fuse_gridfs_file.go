package main

import (
	//"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	"labix.org/v2/mgo/bson"
)

// GridFsFile implements both Node and Handle for a document from a collection.
type GridFsFile struct {
	Id     bson.ObjectId `bson:"_id"`
	Name   string
	Prefix string

	Dirent fuse.Dirent
	Fattr  fuse.Attr
}

func (g GridFsFile) Attr(a *fuse.Attr) {
	log.Printf("GridFsFile.Attr() for: %+v", g)

	db, s := getDb()
	defer s.Close()

	file, err := db.GridFS(g.Prefix).OpenId(g.Id)
	checkError(err)
	defer file.Close()

	now := time.Now()
	a.Mode = 0400
	a.Size = uint64(file.Size())
	a.Ctime = now
	a.Atime = now
	a.Mtime = now
}

func (g GridFsFile) Lookup(ctx context.Context, fname string) (fs.Node, error) {
	log.Printf("GridFsFile.Lookup(): %s\n", fname)

	return nil, fuse.ENOENT
}

func (g GridFsFile) ReadAll(ctx context.Context) ([]byte, error) {
	log.Printf("GridFsFile.ReadAll(): %s\n", g.Id)

	db, s := getDb()
	defer s.Close()

	file, err := db.GridFS(g.Prefix).OpenId(g.Id)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf []byte
	buf, err = ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return buf, nil

}
