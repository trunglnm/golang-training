package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Respone struct {
	Result Result
	Status interface{}
}

type Result struct {
	Product  Product `json:"data"`
	Error    interface{}
	MetaData interface{}
}

type Product struct {
	Id      int
	Name    string
	AdminId int
	Price   int `json:"final_price"`
}

func main() {
	http.HandleFunc("/order", order)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		panic(err)
	}
}
func order(w http.ResponseWriter, req *http.Request) {
	s := fmt.Sprintf("Tong so tien cua 3 san pham: %d", total())
	fmt.Fprintf(w, s)
}

func total() int {
	urls := []string{
		"https://www.sendo.vn/m/wap_v2/full/san-pham/ao-so-mi-jean-nam-dai-tay-cao-cap-hang-vnxk-31331127?platform=web", // context cancel
		"https://www.sendo.vn/m/wap_v2/full/san-pham/ao-dui-nam-cao-cap-30157047",                                       // context cancel
		"https://www.sendo.vn/m/wap_v2/full/san-pham/ao-so-mi-nam-hang-hop-10036141"}                                    //404 ngay lap tuc => no se dc xu ly nhanh hon 2 thang o tren

	total := 0
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	one := sync.Once{}
	for _, url := range urls {
		wg.Add(1) // wait
		go func(url string) {
			defer wg.Done()
			product, err := getProduct(url, ctx)
			if err != nil {
				one.Do(cancel)
				return
			}
			total = total + product.Price
		}(url)
	}
	wg.Wait() //wait
	return total
}

func getProduct(url string, ctx context.Context) (*Product, error) {
	// Xu ly timeout cho URL
	httpClient := http.Client{
		Timeout: time.Duration(60 * time.Second), // net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("NewRequest:", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second) // context deadline exceeded
	defer cancel()
	req = req.WithContext(ctx)
	// url van bi goi
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("httpClient.Do:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	res := Respone{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("json.Unmarshal:", err)
		return nil, err
	}
	return &res.Result.Product, nil
}
