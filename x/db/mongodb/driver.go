package mongodb

var MaxPool int
var PATH string
var DBNAME string
var UNAME string
var PASSWORD string

func CheckAndInitServiceConnection() {
	if service.baseSession == nil {
		service.URL = PATH
		service.Uname = UNAME
		service.Password = PASSWORD
		err := service.New()
		if err != nil {
			panic(err)
		}
	}
}
