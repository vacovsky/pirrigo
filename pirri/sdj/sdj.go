package sdj

import (
	"fmt"
	"io"

	"github.com/coreos/go-systemd/sdjournal"
)

func doit() {
	jconf := sdjournal.JournalReaderConfig{
		Path: "/var/log/journal/",
		Matches: []sdjournal.Match{
			{
				Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
				Value: "pirri", // ${APPNAME}.service
			},
		},
	}

	jr, err := sdjournal.NewJournalReader(jconf)
	if err != nil {
		panic(err)
	}

	//jr.Follow(nil, os.Stdout)
	b := make([]byte, 64*1<<(10)) // 64KB.
	for {
		c, err := jr.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		fmt.Print(string(b[:c]))
	}
}
