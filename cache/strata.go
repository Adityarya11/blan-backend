package cache

import (
	"log"

	stratakv "github.com/Adityarya11/StrataKV/engine"
)

var DB *stratakv.DB

const strataHealthKey = "__stratakv_health__"

func InitStrataKV(dataDir string) {
	var err error
	DB, err = stratakv.Open(dataDir)
	if err != nil {
		log.Fatalf("Strata failed to open: %v", err)
	}
	log.Println("StrataKV LSM-Tree Engine Initialised")
}

// AI Generated Health-Check
func CheckStrataKV() (bool, string) {
	if DB == nil {
		return false, "not initialized"
	}

	if err := DB.Put([]byte(strataHealthKey), []byte("ok")); err != nil {
		return false, "put failed: " + err.Error()
	}

	val, found := DB.Get([]byte(strataHealthKey))
	if !found {
		return false, "get failed"
	}
	if string(val) != "ok" {
		return false, "unexpected value"
	}

	return true, "ok"
}

func GetCachedOutput(hashKey string) (string, bool) {
	if DB == nil {
		log.Println("Strata is not Initiialised")
		return "", false
	}

	val, found := DB.Get([]byte(hashKey))
	if !found {
		return "", false
	}

	return string(val), true
}

func SaveCacheOutput(hashkey, output string) {
	if DB == nil {
		return
	}

	err := DB.Put([]byte(hashkey), []byte(output))
	if err != nil {
		log.Printf("StrataKV Write Error for %s: %v\n", hashkey[:8], err)
		return
	}

	log.Printf("StrataKV: Cached execution result for %s\n", hashkey[:8])
}

func CloseStrata() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("error is closing strata: %v", err)
		} else {
			log.Println("Strata closed.")
		}
	}
}
