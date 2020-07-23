package file

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func init() {
	if _, err := os.Stat("./data/"); os.IsNotExist(err) {
		_ = os.Mkdir("./data/", 0755)

	}
}

func Open(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
}

// function that gets the complete content of a specified file
func GetContent(id string, message bool) ([]byte, error) {
	var content []byte
	var err error
	// read the content of the file
	if message {
		// open the file
		content, err = ioutil.ReadFile("./data/messages.data")
	} else {
		// open the file
		content, err = ioutil.ReadFile("./data/errors.data")
	}
	// if there is an error return it
	if err != nil {
		return nil, err
	}
	//FILTER CONTENT FOR ID
	data := bytes.Split(content, []byte("\n"))
	content = []byte("")
	maxLines := 100
	//for _, str := range data {
	//	if strings.HasPrefix(string(str), id +",") {
	//		content = append(content, []byte(string(str))...)
	//	}
	//}
	for i := len(data)-1;i>=0;i-- {
		if strings.HasPrefix(string(data[i]), id +",") {
			content = append(content, []byte(string(data[i]))...)
			maxLines--
			if maxLines == 0 {
				break
			}
		}
	}
	// return the content
	return content, nil
}

// appends the given data to the end of the file, also checks amount of lines in the file to make sure that the
// files don't get too large
func WriteContent(id []byte, content []byte, message bool) {
	var f *os.File
	var err error
	if message {
		// open the file
		f, err = Open("./data/messages.data")
	} else {
		// open the file
		f, err = Open("./data/errors.data")
	}
	// if there is an error print it and return
	if err != nil {
		// log.Print(err)
		return
	}
	defer f.Close()
	CheckFileSize(f, 1000000, 50)
	if _, err := f.Write(append([]byte(string(id)+", "), []byte(string(content)+"\n")...)); err != nil {
		// log.Print(err)
	}
}

func CheckFileSize(f *os.File, lineLimit int, removeAmount int) {
	lines, _ := lineCounter(f)
	if lines >= lineLimit {
		removeLines(f, removeAmount)
	}
	_, _ = f.Seek(0, 2)
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func removeLines(f *os.File, amount int) {
	// get file stats
	fi, err := f.Stat()
	if err != nil {
		return
	}
	// create a buffer with the size of the file
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))
	// seek to the beginning of the file
	_, err = f.Seek(0, 0)
	if err != nil {
		return
	}
	// copy the content of the file into the buffer
	_, err = io.Copy(buf, f)
	if err != nil {
		return
	}
	// loop over the amount of lines to remove
	for i := 0; i < amount; i++ {
		// read from the buffer until the delimiter is hit for the first time
		_, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return
		}
	}
	// seek back to the top of the file
	_, err = f.Seek(0, 0)
	if err != nil {
		return
	}
	// copy the data of the buffer into the file (overwrites?)
	nw, err := io.Copy(f, buf)
	if err != nil {
		return
	}
	// truncate the size of the file to the size of the copied data
	err = f.Truncate(nw)
	if err != nil {
		return
	}
	// commit the changes (write from memory to disk)
	err = f.Sync()
	if err != nil {
		return
	}
	// seek back to the top of the file
	_, err = f.Seek(0, 0)
	if err != nil {
		return
	}
}
