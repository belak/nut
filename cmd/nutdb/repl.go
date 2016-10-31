package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/belak/nut"
)

type NutDBRepl struct {
	db            *nut.DB
	bucketName    string
	commandLookup map[string]func(string) error
}

func NewNutDBRepl(path string) (*NutDBRepl, error) {
	db, err := nut.Open(path, 0644)
	if err != nil {
		return nil, err
	}

	r := &NutDBRepl{
		db: db,
	}

	r.commandLookup = map[string]func(string) error{
		"bucket": r.Bucket,
		"exit":   r.Exit,
		"get":    r.Get,
		"help":   r.Help,
		"list":   r.List,
		"set":    r.Set,
	}

	return r, nil
}

func (r *NutDBRepl) Run() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(r.Prompt())
		rawLine, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if len(rawLine) == 0 {
			continue
		}

		line := strings.TrimSpace(string(rawLine))
		split := strings.SplitN(line, " ", 2)
		if len(split) == 1 {
			split = append(split, "")
		}

		cb := r.commandLookup[split[0]]
		if cb == nil {
			fmt.Println("Invalid command")
			continue
		}

		err = cb(split[1])
		if err != nil {
			return err
		}
	}
}

func (r *NutDBRepl) Prompt() string {
	out := bytes.NewBufferString("")
	if r.bucketName != "" {
		_, _ = out.WriteString(r.bucketName)
	}
	_, _ = out.WriteString("> ")
	return out.String()
}

// Commands begin here

func (r *NutDBRepl) Bucket(remainder string) error {
	return r.db.View(func(tx *nut.Tx) error {
		bucket := tx.Bucket(remainder)
		if bucket == nil {
			fmt.Printf("No bucket named %q\n", remainder)
		} else {
			r.bucketName = remainder
		}
		return nil
	})
}

func (r *NutDBRepl) Exit(remainder string) error {
	return io.EOF
}

func (r *NutDBRepl) Get(remainder string) error {
	return r.db.View(func(tx *nut.Tx) error {
		if r.bucketName == "" {
			fmt.Println("No bucket selected")
			return nil
		}

		// TODO: This isn't a valid assumption to make.
		out := make(map[string]interface{})

		bucket := tx.Bucket(r.bucketName)
		bucket.Get(remainder, &out)

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.Encode(out)

		return nil
	})
}

func (r *NutDBRepl) Help(remainder string) error {
	cmds := []string{}
	for cmd := range r.commandLookup {
		cmds = append(cmds, cmd)
	}

	sort.Strings(cmds)

	fmt.Println("Available commands:")
	for _, cmd := range cmds {
		fmt.Println(cmd)
	}

	return nil
}

func (r *NutDBRepl) List(remainder string) error {
	return r.db.View(func(tx *nut.Tx) error {
		if r.bucketName == "" {
			fmt.Println("No bucket selected")
			return nil
		}

		bucket := tx.Bucket(r.bucketName)

		out := &json.RawMessage{}

		cursor := bucket.Cursor()
		for key, err := cursor.First(&out); err != nut.ErrCursorEOF; key, err = cursor.Next(&out) {
			fmt.Println(key)
		}
		return nil
	})
}

func (r *NutDBRepl) Set(remainder string) error {
	fmt.Println("Not implemented")
	return nil
}
