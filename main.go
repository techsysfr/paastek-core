package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	//"github.com/gorilla/mux"
	//"github.com/urfave/negroni"
)

var productsDB map[string]*Product

func init() {
	productsDB = make(map[string]*Product, 0)
	productsDB = map[string]*Product{
		"a": &Product{
			Family: "mafamille",
			Size:   "micro",
			Hardware: Hardware{
				CPU: "i7",
				RAM: 2,
			},
		},
		"b": &Product{
			Family: "mafamille",
			Size:   "medium",
			Hardware: Hardware{
				CPU: "i7",
				RAM: 4,
			},
		},
		"c": &Product{
			Family: "mafamille",
			Size:   "large",
			Hardware: Hardware{
				CPU: "i5",
				RAM: 4,
			},
		},
	}
}

type Product struct {
	ID       string
	Family   string // `json: "Family, omitempty"`
	Size     string
	Hardware Hardware
}

type Hardware struct {
	CPU string
	RAM int
}

func RootHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello root")
}

func ProductsHandler(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	if param != nil {

		prd, err := getProduct(param.ByName("ID"))
		if err != nil {
			log.Fatal(err)
		}
		if prd == nil {
			http.NotFound(w, r)
			return
		}
		aff, err := json.MarshalIndent((*prd), "", "   ")
		if err != nil {
			log.Fatal(err)
		}
		log.Println(prd)
		fmt.Fprintf(w, "Hello product, %s!\n", string(aff))

	} else {
		prd, err := getProducts()
		if err != nil {
			log.Fatal(err)
		}
		aff, err := json.MarshalIndent((prd), "", "   ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Hello product, %s!\n", string(aff))
	}

}

func InsertProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var prd Product
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&prd)

	fmt.Fprintf(w, "%v\n", err)
	fmt.Fprintf(w, "%v\n", prd)
}

func ArticlesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello articles")
}

func getProducts() ([]*Product, error) {
	var prds []*Product
	for k, v := range productsDB {
		v.ID = k
		prds = append(prds, v)
	}
	return prds, nil
}
func getProduct(id string) (*Product, error) {
	return productsDB[id], nil
}

func addProduct(p *Product, id string) error {
	productsDB[id] = p
	return nil
}

func main() {

	/*	prd := Product{
			ID:     "monID",
			Family: "t2",
			Size:   "micro",
			Hardware: Hardware{
				CPU: "i7",
				RAM: 4,
			},
		}
		//aff,err := json.Marshal(prd)
		aff, err := json.MarshalIndent(prd, "", "   ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(aff))
	*/

	router := httprouter.New()
	router.GET("/", RootHandler)
	router.GET("/products", ProductsHandler)
	router.GET("/products/:ID", ProductsHandler)
	router.POST("/products", InsertProduct)
	router.GET("/articles", ArticlesHandler)

	log.Fatal(http.ListenAndServe(":8080", router))

	/* With Negroni
	r := mux.NewRouter()
	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/products", ProductsHandler)
	r.HandleFunc("/articles", ArticlesHandler)
	//http.Handle("/", r)

	//log.Fatal(http.ListenAndServe(":8080", nil))

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(r)
	log.Fatal(http.ListenAndServe(":3000", n))
	*/

}
