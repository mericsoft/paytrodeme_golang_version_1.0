package paytrlibrary

import (
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PaytrOde(r *http.Request) []byte { //ödeme fonksiyonu
	rand.Seed(time.Now().UnixNano())        // bu yapılmazsa hep aynı random sayı döner
	Randoid := rand.Intn(9999999 - 1000000) // random sayı üretici merhcant_oid alanı içn
	MerchantKey := "blablabla"              // paytr panelden alınacak
	MerchantSalt := ""                      //paytr panelden alınacak
	Vals := map[string]string{
		"merchant_id":       "111111",                       //paytr kaydı sırasında alınmış olacak
		"user_ip":           "3.124.199.224",                // "176.227.88.243",
		"merchant_oid":      strconv.Itoa(Randoid),          //random sayıdan gelecek
		"email":             r.FormValue("email"),           // sistemdeki emailiniz
		"payment_amount":    r.FormValue("price"),           // tutar
		"currency":          "TL",                           // burası sabit kalsın
		"no_installment":    r.FormValue("no_installment"),  //0 yazarsanız taksit yok 1 yazadsanız var
		"max_installment":   r.FormValue("max_installment"), // maksimum taksit sayısı
		"user_name":         r.FormValue("user_name"),       // müşteri adı
		"user_address":      r.FormValue("user_address"),    //müşteri adresi
		"user_phone":        r.FormValue("user_phone"),      //müşteri telefon
		"merchant_ok_url":   "https://mericsoft.com/manage/payresult",
		"merchant_fail_url": "https://mericsoft.com/manage/payresult",
		"debug_on":          "1", // hata dönecek mi
		"test_mode":         "0", // panelde test modundan canlıya alınca 0 yapılmalı
	}
	Price, _ := strconv.ParseFloat(r.FormValue("price"), 64)                                                                                                                                                                                                   // ödemeyi bu formata çevirmek gerekiyor
	PriceN := (Price / 100)                                                                                                                                                                                                                                    // ödeme kuruş olarak gider sisteme
	PriceSt := strconv.FormatFloat(PriceN, 'f', 2, 64)                                                                                                                                                                                                         // ödemeyi 2 haneli yani kuruş yapıyor
	Bask := [][]string{{r.FormValue("user_basket"), PriceSt, "1"}}                                                                                                                                                                                             // burası ürün sepeti için
	Basket, _ := json.Marshal(Bask)                                                                                                                                                                                                                            // ürün sepetini json formatına çevirecek
	Vals["user_basket"] = b64.StdEncoding.EncodeToString(Basket)                                                                                                                                                                                               // json foratındakini string halinde encode etti
	Nsalt := hmac.New(sha256.New, []byte(Merchant_key))                                                                                                                                                                                                        // veriyi tuzlamak için token hazırlığı
	Nsalt.Write([]byte(Vals["merchant_id"] + Vals["user_ip"] + Vals["merchant_oid"] + Vals["email"] + Vals["payment_amount"] + Vals["user_basket"] + Vals["no_installment"] + Vals["max_installment"] + Vals["currency"] + Vals["test_mode"] + Merchant_salt)) // veriyi hashladı
	Vals["paytr_token"] = b64.StdEncoding.EncodeToString(Nsalt.Sum(nil))                                                                                                                                                                                       // tekrar string formatında enocad etti
	var gon string                                                                                                                                                                                                                                             //keylerin isimlerini buraya alacağız
	for key, val := range Vals {
		gon += key + "=" + val + "&" //post için kay=value& şeklinde birleştiriyoruz
	}
	gon = gon[:len(gon)-1] // son & karakterini sildik

	read := strings.NewReader(gon)                                                                                    // stringgs kütüphanesi ile post datayı okuduk ve post edilir formata çevirdik
	ResSend, err := http.Post("https://www.paytr.com/odeme/api/get-token", "application/x-www-form-urlencoded", read) // datayı yolladık
	if err != nil {
		log.Fatal(err) // hata varsa error yazdı, isterseniz log yerine return dönebilirsiniz
	}
	defer ResSend.Body.Close()                  // iş bitince cevabı kapat
	ResRet, err := ioutil.ReadAll(ResSend.Body) // cevabı oku
	if err != nil {
		log.Fatal(err) // hata varsa  bas
	}
	return ResRet

}

func PaytrResult(r *http.Request) string { // odeme sonucu, dönüş string olmak zorunda
	MerchantKey := "blablabla"
	MerchantSalt := "blablabla"
	Ressalt := hmac.New(sha256.New, []byte(MerchantKey))
	Ressalt.Write([]byte(r.FormValue("merchant_oid") + MerchantSalt + r.FormValue("status") + r.FormValue("total_amount")))
	ResHash := b64.StdEncoding.EncodeToString(Ressalt.Sum(nil))
	if ResHash != r.FormValue("hash") {
		ResultMessage = "hash problem"
		//w.Write([]byte()

	}
	if r.FormValue("status") == "success" {
		ResultMessage = " The payment is in your account now!"
	} else {
		ResultMessage = "There is a problem"
	}
	return "OK" //  dönüşün ok olması zorunlu

}
