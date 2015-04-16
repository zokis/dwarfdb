package dwarfdb

import (
    "gopkg.in/vmihailenco/msgpack.v2"
    "io/ioutil"
    "os"
    "errors"
)

func pathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

type DwarfDB struct {
    path string
    force bool
    db map[string]interface{}
}

func DwarfDBLoad(path string, force bool) DwarfDB {
    exist, err := pathExists(path)
    if err != nil {
        panic(err)
    }
    ddb := DwarfDB{path, force, make(map[string]interface{})}
    if exist {
        ddb.loaddb()
    }
    return ddb
}

func (ddb *DwarfDB) Dump() bool {
    ddb.dumpdb(true)
    return true
}

func (ddb *DwarfDB) loaddb() bool {
    input, err := ioutil.ReadFile(ddb.path)
    if err != nil {
        panic(err)
    }

    if err := msgpack.Unmarshal([]byte(input), &ddb.db); err != nil {
        panic(err)
    }

    return true
}

func (ddb *DwarfDB) dumpdb(force bool) bool {
    dump, _ := msgpack.Marshal(ddb.db)
    err := ioutil.WriteFile(ddb.path, []byte(string(dump)), 0644)
    if err != nil {
        panic(err)
    }
    return true
}

func (ddb *DwarfDB) Set(key string, value interface{}) bool {
    ddb.db[key] = value
    ddb.dumpdb(ddb.force)
    return true
}

func (ddb *DwarfDB) Get(key string) (interface{}, error) {
    value, ok := ddb.db[key]
    if ok {
        return value, nil
    } else {
        return nil, errors.New("Not Found")
    }
}

func (ddb *DwarfDB) GetAll() []string {
    keys := make([]string, 0, len(ddb.db))
    for k := range ddb.db {
        keys = append(keys, k)
    }
    return keys
}

func (ddb *DwarfDB) Len() int {
    return len(ddb.db)
}

func (ddb *DwarfDB) Rem(key string) bool{
    delete(ddb.db, key)
    ddb.dumpdb(ddb.force)
    return true
}

func (ddb *DwarfDB) DelDB() bool{
    ddb.db = make(map[string]interface{})
    ddb.dumpdb(ddb.force)
    return true
}