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
	log.Printf("DocumentFile.Attr() for: %+v", d)
	strval, err := d.readDocument()

	if err != nil {
		return
	}

	now := time.Now()
	a.Mode = 0400
	a.Size = uint64(len(strval))
	a.Ctime = now
	a.Atime = now
	a.Mtime = now
}

func (d DocumentFile) Lookup(ctx context.Context, fname string) (fs.Node, error) {
	log.Printf("DocumentFile[%s].Lookup(): %s\n", d.coll, fname)

	return nil, fuse.ENOENT
}

func (d DocumentFile) ReadAll(ctx context.Context) ([]byte, error) {
	log.Printf("DocumentFile[%s].ReadAll(): %s\n", d.coll, d.Id)

	strval, err := d.readDocument()
	if err != nil {
		return nil, err
	}

	return []byte(strval), nil
}

func (d DocumentFile) readDocument() (string, error) {
	db, s := getDb()
	defer s.Close()

	var f interface{}
	err := db.C(d.coll).FindId(d.Id).One(&f)
	if err != nil {
		log.Fatal(err)
		return "", fuse.EIO
	}

	buf, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatal(err)
		return "", fuse.EIO
	}

	strval := string(buf) + "\n"
	return strval, nil
}
