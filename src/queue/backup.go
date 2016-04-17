package queue

import (
	def "config"
	"log"
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
)


type backupQueue struct {
	List []int
}

func Get_backup_from_file(){
	var backup backupQueue
	backup.load_from_file(def.BackupFileName)
	Set_Orders(backup.List)
}

func Save_backup_to_file(){
	var backup backupQueue
	backup.List = Get_Orders()
	backup.save_to_file(def.BackupFileName)
}

func (q *backupQueue) save_to_file(filename string) error {

	data, err := json.Marshal(&q)
	if err != nil {
		log.Println(def.ColR, "json.Marshal() error: Failed to backup.", def.ColN)
		return err
	}
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		log.Println(def.ColR, "ioutil.WriteFile() error: Failed to backup.", def.ColN)
		return err
	}
	return nil
}

func (q *backupQueue) load_from_file(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		fmt.Println(def.ColG,"Backup file found, content uploaded", def.ColN)

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println(def.ColR, "loadFromDisk() error: Failed to read file.", def.ColN)
		}
		if err := json.Unmarshal(data, q); err != nil {
			log.Println(def.ColR, "loadFromDisk() error: Failed to Unmarshal.", def.ColN)
		}
	}
	return nil
}

