package queue

import (
	def "config"
	"log"
	//"time"
	"encoding/json"
	"io/ioutil"
	"os"
)

/*/--------------Dette legges der vi skal bruke det----------------
	var ordrs = []int{-2, 2, -3, 2}
	var listlist = []int{-1, 1, -1, 1}
	var backupOrd queue.BackupQueue
	var kokoko queue.BackupQueue
	backupOrd.List = ordrs
	kokoko.List = listlist

	backupOrd.Save_to_file("orderBackup")
	kokoko.Load_from_disk("orderBackup")
	fmt.Printf("%v\n", kokoko.List)
*/


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

// loadFromDisk checks if a file of the given name is available on disk, and
// saves its contents to a queue if the file is present.
func (q *backupQueue) load_from_file(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		log.Println(def.ColG, "Backup file found, processing...", def.ColN)

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

