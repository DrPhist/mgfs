package main

import (
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	"labix.org/v2/mgo/bson"
)

// GridFs implements both Node and Handle for a collection.
type GridFs struct {
	Name string // the GridFS prefix

	Dirent fuse.Dirent
	Fattr  fuse.Attr
}

func (g GridFs) Attr(a *fuse.Attr) {
	log.Printf("GridFs.Attr() for: %+v", g)
	a.Mode = os.ModeDir | 0400
	a.Uid = uint32(os.Getuid())
	a.Gid = uint32(os.Getgid())
}

func (g GridFs) Lookup(ctx context.Context, fname string) (fs.Node, error) {
	log.Printf("GridFs[%s].Lookup(): %s\n", g.Name, fname)

	if !bson.IsObjectIdHex(fname) {
		return nil, fuse.ENOENT
	}

	db, s := getDb()
	defer s.Close()

	id := bson.ObjectIdHex(fname)
	gf := &GridFsFile{Id: id}
	file, err := db.GridFS(g.Name).OpenId(id)
	if err != nil {
		log.Fatal(err)
		return nil, fuse.EIO
	}
	defer file.Close()

	gf.Name = file.Name()
	gf.Prefix = g.Name

	return gf, nil
}

func (g GridFs) ReadDirAll(ctx context.Context) (ents []fuse.Dirent, ferr error) {
	log.Println("GridFs.ReadDirAll(): ", g.Name)

	db, s := getDb()
	defer s.Close()

	iter := db.GridFS(g.Name).Find(nil).Iter()

	var f GridFsFile
	for iter.Next(&f) {
		ents = append(ents, fuse.Dirent{Name: f.Id.Hex(), Type: fuse.DT_File})
	}

	if err := iter.Err(); err != nil {
		log.Fatal(err)
		return nil, fuse.EIO
	}

	return ents, nil
}
