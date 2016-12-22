package main

import (
	"archive/zip"
	"strings"
	"fmt"
)

type ClassFile struct {
	CanonicalName string
	Content       []byte
	Size          int64
}

func loadJar(filename string) ([]*ClassFile) {
	var closer, err = zip.OpenReader(filename)
	if err != nil {
		fmt.Printf("Error! %v", err)
	}

	var files = closer.File
	return loadClasses(&files)
}

func loadClasses(files *[]*zip.File) ([]*ClassFile) {
	arr := make([]*ClassFile, 0)
	for _, f := range *files {
		if (!f.FileInfo().IsDir() && strings.HasSuffix(f.Name, ".class")) {
			arr = append(arr, loadClassFile(f))
		}
	}
	return arr
}

func loadClassFile(classfile *zip.File) (*ClassFile) {
	reader, err := classfile.Open()
	if err != nil {
		fmt.Printf("Error! %v", err)
	}

	size := classfile.FileInfo().Size()
	content := make([]byte, size)
	_, err = reader.Read(content)

	class := new(ClassFile)
	class.CanonicalName = classNameFromPath(classfile.Name)
	class.Content = content
	class.Size = size
	return class
}

func classNameFromPath(filename string) (string) {
	name := strings.Replace(filename, "/", ".", -1)
	return strings.Replace(name, ".class", "", -1)
}