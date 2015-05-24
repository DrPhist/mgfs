package main

import (
	"log"
	"os"
	"strings"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	"labix.org/v2/mgo/bson"
)

// CollFile implements both Node and Handle for a collection.
type CollFile struct {
	Name string

	Dirent fuse.Dirent
	Fattr  fuse.Attr
}

func (f CollFile) Attr(a *fuse.Attr) {
	log.Printf("CollFile.Attr() for: %+v", f)
	a.Mode = os.ModeDir | 0700
	a.Uid = uint32(os.Getuid())
	a.Gid = uint32(os.Getgid())
}

func (c CollFile) Lookup(ctx context.Context, fname string) (fs.Node, error) {
	log.Printf("CollFile[%s].Lookup(): %s\n", c.Name, fname)

	if !strings.HasSuffix(fname, ".json") {
		return nil, fuse.ENOENT
	}
	fname = fname[0 : len(fname)-5]

	if !bson.IsObjectIdHex(fname) {
		return nil, fuse.ENOENT
	}

	db, s := getDb()
	defer s.Close()

	var f DocumentFile
	err := db.C(c.Name).FindId(bson.ObjectIdHex(fname)).One(&f)
	if err != nil {
		log.Printf("Error while looking up %s: %s \n", fname, err.Error())
		return nil, fuse.EIO
	}

	f.coll = c.Name

	return f, nil
}

func (c CollFile) ReadDirAll(ctx context.Context) (ents []fuse.Dirent, ferr error) {
	log.Println("CollFile.ReadDirAll(): ", c.Name)

	db, s := getDb()
	defer s.Close()

	iter := db.C(c.Name).Find(nil).Select(bson.M{"text": 0}).Iter()

	var f DocumentFile
	for iter.Next(&f) {
		ents = append(ents, fuse.Dirent{Name: f.Id.Hex() + ".json", Type: fuse.DT_File})
	}

	if err := iter.Err(); err != nil {
		log.Panic(err)
		return nil, fuse.EIO
	}

	return ents, nil
}
