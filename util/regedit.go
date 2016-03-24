package util

import (
	// go
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	Add = iota
	Query
	Delete
	Copy
	Save
	Restore
	Load
	Unload
	Compare
	Export
	Import
)

const (
	HKCR = iota
	HKLM
	HKU
	HKCU
	HKCC
)

const (
	SZ = iota
	BINARY
	MULTI
	EXPAND
)

type (
	Regedit struct {
		Action string
		Field  string
		Reg
	}
	Reg struct {
		Key   string
		Type  string
		Value string
	}
	RegCmd struct {
		cmd *exec.Cmd
		reg *Regedit
	}
)

var (
	Actions = map[int]string{
		0:  "add",
		1:  "query",
		2:  "delete",
		3:  "copy",
		4:  "save",
		5:  "restore",
		6:  "load",
		7:  "unload",
		8:  "compare",
		9:  "export",
		10: "import",
	}
	Fields = map[int]string{
		0: "HKEY_CLASSES_ROOT",
		1: "HKEY_LOCAL_MACHINE",
		2: "HKEY_USERS",
		3: "HKEY_CURRENT_USER",
		4: "HKEY_CURRENT_CONFIG",
	}
	Types = map[int]string{
		0: "REG_SZ",
		1: "REG_BINARY",
		2: "REG_MULTI_SZ",
		3: "REG_EXPAND_SZ",
	}
)

func (this *Regedit) Add(reg Reg) RegCmd {
	(*this).Reg = reg
	return RegCmd{exec.Command("cmd", "/c", "reg", this.Action, this.Field, "/v", this.Key, "/t", this.Type, "/d", this.Value), this}
}

func (this *Regedit) Search(reg Reg) RegCmd {
	(*this).Reg = reg
	return RegCmd{exec.Command("cmd", "/c", "reg", this.Action, this.Field, "/s"), this}
}

func (this RegCmd) Exec() ([]Reg, error) {
	if this.reg.Action == "add" {
		this.cmd.Stdout = os.Stdout
		this.cmd.Stderr = os.Stderr
		this.cmd.Stdin = os.Stdin
		if err := this.cmd.Run(); err != nil {
			return nil, err
		}
	} else if this.reg.Action == "query" {
		if out, err := this.cmd.Output(); err != nil {
			return nil, err
		} else {
			buff := bytes.NewBuffer(out)
			for {
				content, err := buff.ReadString('\n')
				content = strings.TrimSpace(content)
				if err != nil || err == io.EOF {
					break
				}
				if strings.Index(content, this.reg.Key) != -1 {
					regList := make([]Reg, 0)
					if arr := strings.Fields(content); len(arr) == 3 {
						regList = append(regList, Reg{arr[0], arr[1], arr[2]})
						return regList, nil
					}
				}
			}
		}

	}
	return nil, nil
}
