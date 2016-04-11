package queue
/*
import (
	def "config"
	"encoding/binary"
	"encoding/json"
	"log"
	"fmt"
	"os"
	"io/ioutil"
)

type backupQueue []int

// Saves this elevs queue to file on disk
func Local_save(){
	file, err := os.OpenFile("file.bin", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error : os.OpenFile : " )
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, &def.Orders)
	if err != nil {
		fmt.Println("error : binary.Write : ")
	}
}
/*
	if err := ioutil.WriteFile(filename, int_to_byte(def.Orders), 0644); err != nil {
		log.Println(def.ColR, "ioutil.WriteFile() error: Failed to backup.", def.ColN)
		return err
	}
	return nil
}

//Checks for file "filename" in folder and loads content to a queue if present 
func local_load() []int{
	file, err := os.OpenFile("file.bin", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error : os.OpenFile : " )
	}
	defer file.Close()


}
*/
