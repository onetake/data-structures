package trie

import "testing"

import (
    "os"
    "math/rand"
    "fmt"
    "sort"
)

import (
    bs "file-structures/block/byteslice"
)

import (
    "github.com/timtadh/data-structures/types"
)


func init() {
    if urandom, err := os.Open("/dev/urandom"); err != nil {
        return
    } else {
        seed := make([]byte, 8)
        if _, err := urandom.Read(seed); err == nil {
            rand.Seed(int64(bs.ByteSlice(seed).Int64()))
        }
        urandom.Close()
    }
}

func randslice(length int) []byte {
    if urandom, err := os.Open("/dev/urandom"); err != nil {
        panic(err)
    } else {
        slice := make(bs.ByteSlice, length)
        if _, err := urandom.Read(slice); err != nil {
            panic(err)
        }
        urandom.Close()
        // return append([]byte("b"), slice...)
        return slice
    }
    panic("unreachable")
}

func has_zero(bytes []byte) bool {
    for _,ch := range bytes {
        if ch == 0 {
            return true
        }
    }
    return false
}

func randslice_nonzero(length int) []byte {
    slice := randslice(length)
    for ; has_zero(slice); slice = randslice(length) { }
    return slice
}

func write(name, contents string) {
    file, _ := os.Create(name)
    fmt.Fprintln(file, contents)
    file.Close()
}

type ByteSlices []types.ByteSlice

func (self ByteSlices) Len() int {
    return len(self)
}

func (self ByteSlices) Less(i, j int) bool {
    return self[i].Less(self[j])
}

func (self ByteSlices) Swap(i, j int) {
    self[i], self[j] = self[j], self[i]
}

func TestIteratorPrefixFindDotty(t *testing.T) {
    items := ByteSlices{
        types.ByteSlice("cat"),
        types.ByteSlice("catty"),
        types.ByteSlice("car"),
        types.ByteSlice("cow"),
        types.ByteSlice("candy"),
        types.ByteSlice("coo"),
        types.ByteSlice("coon"),
        types.ByteSlice("andy"),
        types.ByteSlice("alex"),
        types.ByteSlice("andrie"),
        types.ByteSlice("alexander"),
        types.ByteSlice("alexi"),
        types.ByteSlice("bob"),
        types.ByteSlice("bobcat"),
        types.ByteSlice("barnaby"),
        types.ByteSlice("baskin"),
        types.ByteSlice("balm"),
    }
    table := new(TST)
    for _, key := range items {
        if err := table.Put(key, nil); err != nil { t.Error(table, err) }
        if has := table.Has(key); !has { t.Error(table, "Missing key") }
    }
    write("TestDotty.dot", table.Dotty())
    sort.Sort(items)
    i := 0
    for k, _, next := table.Iterate()(); next != nil; k, _, next = next() {
        if !k.Equals(types.ByteSlice(items[i])) {
            t.Error(string(k.(types.ByteSlice)), "!=", string(items[i]))
        }
        i++
    }
    co_items := ByteSlices{
        types.ByteSlice("coo"),
        types.ByteSlice("coon"),
        types.ByteSlice("cow"),
    }
    i = 0
    for k, _, next := table.PrefixFind([]byte("co"))(); next != nil; k, _, next = next() {
        if !k.Equals(types.ByteSlice(co_items[i])) {
            t.Error(string(k.(types.ByteSlice)), "!=", string(co_items[i]))
        }
        i++
    }
}

func TestPutHasGet(t *testing.T) {

    type record struct {
        key bs.ByteSlice
        value bs.ByteSlice
    }

    ranrec := func() *record {
        return &record{ randslice_nonzero(3), randslice(3) }
    }

    test := func(table *TST) {
        records := make([]*record, 1000)
        for i := range records {
            r := ranrec()
            records[i] = r
            err := table.Put(r.key, "")
            if err != nil {
                t.Error(err)
            }
            err = table.Put(r.key, r.value)
            if err != nil {
                t.Error(err)
            }
        }

        for _, r := range records {
            if has := table.Has(r.key); !has {
                t.Error(table, "Missing key")
            }
            if has := table.Has(randslice(12)); has {
                t.Error("Table has extra key")
            }
            if val, err := table.Get(r.key); err != nil {
                t.Error(err)
            } else if !(val.(bs.ByteSlice)).Eq(r.value) {
                t.Error("wrong value")
            }
        }
    }

    test(new(TST))
}
