package main

import (
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"golang.org/x/net/context"
)

func mount(point string) {

	// startup mount
	c, err := fuse.Mount(point)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer c.Close()

	log.Println("Mounted: ", point)
	err = fs.Serve(c, mgoFS{})

	// check if the mount process has an error to report
	<-c.Ready
}

// mgoFS implements my mgo fuse filesystem
type mgoFS struct{}

func (mgoFS) Root() (fs.Node, error) {
	log.Println("returning root node")
	return Dir{"Root"}, nil
}

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	name string
}

func (d Dir) Attr(a *fuse.Attr) {
	log.Println("Dir.Attr() for ", d.name)
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
}

func (Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Println("Dir.Lookup():", name)

	db, s := getDb()
	defer s.Close()

	names, err := db.CollectionNames()
	if err != nil {
		log.Panic(err)
		return nil, fuse.EIO
	}

	for _, collName := range names {
		if collName == name {
			return CollFile{Name: name}, nil
		}
	}

	return nil, fuse.ENOENT
}

func (d Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("Dir.ReadDirAll():", d.name)

	db, s := getDb()
	defer s.Close()

	names, err := db.CollectionNames()
	if err != nil {
		log.Panic(err)
		return nil, fuse.EIO
	}

	ents := make([]fuse.Dirent, 0, len(names))

	for _, name := range names {
		ents = append(ents, fuse.Dirent{Name: name, Type: fuse.DT_Dir})
	}
	return ents, nil
}
