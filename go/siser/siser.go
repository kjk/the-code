package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"time"

	"github.com/kjk/common/siser"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func siserWriteSimple() []byte {
	var buf bytes.Buffer
	w := siser.NewWriter(&buf)
	_, err := w.Write([]byte("hello"), time.Now(), "log")
	panicIfErr(err)
	_, err = w.Write([]byte("world!"), time.Now(), "log")
	panicIfErr(err)
	return buf.Bytes()
}

func siserReadSimple(d []byte) {
	br := bufio.NewReader(bytes.NewReader(d))
	r := siser.NewReader(br)
	for {
		hasMore := r.ReadNextData()
		if !hasMore {
			break
		}
		fmt.Printf("Record: name: '%s', data: '%s'\n", r.Name, r.Data)
	}
}

func showSiserWriteSimple() {
	d := siserWriteSimple()
	fmt.Printf("Data written:\n%s\n", string(d))
	siserReadSimple(d)
}

func siserWriteRecords() []byte {
	var buf bytes.Buffer
	w := siser.NewWriter(&buf)
	var r siser.Record
	r.Name = "log"

	err := r.Write("url", "http://example.com")
	panicIfErr(err)
	_, err = w.WriteRecord(&r)
	panicIfErr(err)

	r.Reset()
	err = r.Write("url", "http://example2.com")
	panicIfErr(err)
	_, err = w.WriteRecord(&r)
	panicIfErr(err)

	return buf.Bytes()
}

func siserReadRecords(d []byte) {
	br := bufio.NewReader(bytes.NewReader(d))
	r := siser.NewReader(br)
	for {
		hasMore := r.ReadNextRecord()
		if !hasMore {
			break
		}
		rec := r.Record
		fmt.Printf("Record: name: '%s', time: %d\n", rec.Name, rec.Timestamp.UnixMilli())
		for _, e := range rec.Entries {
			fmt.Printf("  key: '%s', value: '%s'\n", e.Key, e.Value)
		}
	}
}

func showSiserWriteRecord() {
	d := siserWriteRecords()
	fmt.Printf("Data written:\n%s\n", string(d))
	siserReadRecords(d)
}

func showRecordMarshal() {
	var r siser.Record
	r.Name = "log"
	r.Write("url", "http://example.com")
	r.Write2("code", 200)
	r.Write("long", ";lksaflashfdlkahfdlk a;lfka sdflkj a;sldk f;las f;lkjsadf ;lajsdf;ljasd ;k ;lksa j;las jf;lsajfd ;sa ;lkasjf;lsakjd ;l a;ldsafkj ;laskj f;lkj ;lsakjfd ;laskj ;lfjds ;lkjdsaf;d")
	d := r.Marshal()
	fmt.Printf("Record marshalled:\n%s\n", string(d))
}

func main() {
	var (
		flgSimple        bool
		flgRecord        bool
		flgRecordMarshal bool
	)
	{
		flag.BoolVar(&flgSimple, "simple", false, "show simple siser write/read example")
		flag.BoolVar(&flgRecord, "record", false, "show siser record write/read example")
		flag.BoolVar(&flgRecordMarshal, "record-marshal", false, "show siser record marshal example")
		flag.Parse()
	}
	if flgSimple {
		showSiserWriteSimple()
		return
	}
	if flgRecord {
		showSiserWriteRecord()
		return
	}
	if flgRecordMarshal {
		showRecordMarshal()
		return
	}
	flag.Usage()
}
