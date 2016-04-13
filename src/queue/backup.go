package queue

/*
import (
	def "config"
	"log"
	//"time"
	"encoding/json"
	"io/ioutil"
	"os"
)

//--------------Dette legges der vi skal bruke det----------------
	var ordrs = []int{-2, 2, -3, 2}
	var backupOrd queue.OrdQueue
	var kokoko queue.OrdQueue
	backupOrd.List = ordrs

	backupOrd.SaveToDisk("orderBackup")
	kokoko.LoadFromDisk("orderBackup")
	fmt.Printf("%v\n", kokoko.List)
//-----------------------------------------


type OrdQueue struct {
	List []int
}

func (q *OrdQueue) SaveToDisk(filename string) error {

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
func (q *OrdQueue) LoadFromDisk(filename string) error {
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
*/
