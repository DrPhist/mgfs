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
	_, size, err := d.readDocument()

	if err != nil {
		return
	}

	now := time.Now()
	a.Mode = 0400
	a.Size = size
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

	strval, _, err := d.readDocument()
	if err != nil {
		return nil, err
	}

	return []byte(strval), nil
}

// Read a document and return it as a JSON string
func (d DocumentFile) readDocument() (string, uint64, error) {
	db, s := getDb()
	defer s.Close()

	var f interface{}
	err := db.C(d.coll).FindId(d.Id).One(&f)
	if err != nil {
		log.Fatal(err)
		return "", 0, fuse.EIO
	}

	buf, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatal(err)
		return "", 0, fuse.EIO
	}

	strval := string(buf) + "\n"
	return strval, uint64(len(buf)), nil
}
