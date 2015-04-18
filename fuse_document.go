package main

import (
	"encoding/json"
	"log"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	"labix.org/v2/mgo/bson"
)

// DocumentFile implements both Node and Handle for a document from a collection.
type DocumentFile struct {
	coll string
	Id   bson.ObjectId `bson:"_id"`

	Dirent fuse.Dirent
	Fattr  fuse.Attr
}

func (d DocumentFile) Attr(a *fuse.Attr) {
	log.Println("DocumentFile.Attr() for: %+v", d)

	now := time.Now()
	a.Mode = 0400
	a.Size = 42
	a.Ctime = now
	a.Atime = now
	a.Mtime = now
}

func (d DocumentFile) Lookup(ctx context.Context, fname string) (fs.Node, error) {
	log.Println("DocumentFile[%s].Lookup(): %s\n", d.coll, fname)

	return nil, fuse.ENOENT
}

func (d DocumentFile) ReadAll(ctx context.Context) ([]byte, error) {
	log.Println("DocumentFile[%s].ReadAll(): %s\n", d.coll, d.Id)

	db, s := getDb()
	defer s.Close()

	var f interface{}
	err := db.C(d.coll).FindId(d.Id).One(&f)
	if err != nil {
		log.Fatal(err)
		return nil, fuse.EIO
	}

	buf, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatal(err)
		return nil, fuse.EIO
	}

	return buf, nil
}
