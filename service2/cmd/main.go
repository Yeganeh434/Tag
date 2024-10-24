package main
// module tag_project
import(
	"service2/internal/adapters/controllers/http"
	"log"
)
func main(){
	err := http.RunWebServer()
	if err != nil {
		log.Printf("could not start server:%v", err)
	}
}
