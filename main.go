package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	// "golang.org/x/tools/go/analysis/passes/unmarshal"
)

func main() {
	proxyURL, _ := url.Parse("http://127.0.0.1:12334")
	httpClient := &http.Client{
    	Transport: &http.Transport{
        	Proxy: http.ProxyURL(proxyURL),
    },
}

	bot, err := tgbotapi.NewBotAPIWithClient("8971346126:AAEr4fkeyStMJugl0DZzLKRgltj8-ikC06s", tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			article := extractArticle(update.Message.Text)
			if article == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отправь ссылку на товар с goldapple.ru")
				bot.Send(msg)
			} else {
				price, err := getPrice(article)
				if err != nil {
					log.Println("ошибка getPrice:", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка получения цены")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Цена: %d руб.", price))
					bot.Send(msg)
				}
			}
		}
	}

}


func extractArticle(url string) string {
	
	if !strings.Contains(url, "goldapple.ru") {
        return ""
    } 
	article := strings.Split(strings.Split(url, "/")[3], "-")[0] 
    return article
}

type GAResponse struct {
    Data GAData `json:"data"`
}

type GAData struct {
    Variants []GAVariant `json:"variants"`
}

type GAVariant struct {
    Price GAPrice `json:"price"`
}

type GAPrice struct {
    Actual GAPriceValue `json:"actual"`
}

type GAPriceValue struct {
    Amount int `json:"amount"`
}

func getPrice(article string) (int, error) {
    url := "https://goldapple.ru/front/api/catalog/product-card/base/v3?locale=ru&cityId=0c5b2444-70a0-4932-980c-b4dc0d3f02b5&regionId=0c5b2444-70a0-4932-980c-b4dc0d3f02b5&itemId=" + article

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
    req.Header.Set("Accept", "application/json, text/plain, */*")
    req.Header.Set("Accept-Language", "ru,en-US;q=0.9,en;q=0.8")
    req.Header.Set("plaid-city-id", "0c5b2444-70a0-4932-980c-b4dc0d3f02b5")
    req.Header.Set("plaid-device-id", "ZKZF9tv5leF6t7bkt5Ryv")
    req.Header.Set("plaid-language-id", "ru_RU")
    req.Header.Set("plaid-platform", "web")
    req.Header.Set("plaid-region-id", "0c5b2444-70a0-4932-980c-b4dc0d3f02b5")
    req.Header.Set("plaid-store-id", "ru")
    req.Header.Set("plaid-version", "1.0.0")
    req.Header.Set("Referer", "https://goldapple.ru/")
	req.Header.Set("Cookie", "ga-gun-init=true; __zzatw-goldapple=MDA0dBA=Fz2+aQ==; ga-lang=ru; ga-device-id=ZKZF9tv5leF6t7bkt5Ryv; client-store-code=default; mage-cache-sessid=true; advcake_track_id=b25088d5-ff96-bd2f-377a-c67aa8f9fb17; advcake_session_id=fbb80f83-a148-9b97-fc90-711b69d1f440; _ym_uid=1781324515444689476; _ym_d=1781324515; _ym_isad=2; mindboxDeviceUUID=fe5dcd9d-1c27-4d5d-b7fe-ab0130f0fae0; directCrm-session=%7B%22deviceGuid%22%3A%22fe5dcd9d-1c27-4d5d-b7fe-ab0130f0fae0%22%7D; _gcl_au=1.1.378911014.1781324516; _ga=GA1.1.641025950.1781324516; gsscw-goldapple=tpElXqNoeo/tNGZ7GehI7w7HdzKaXzoca+zpCQxiZgIpBPhClgov0U689MBKfbX9aI4IRZX5gqr90ENFFiQY+gNBXAGK9V9rCNmYlctWP7E4388ieAyNrPCWK+u0sGg9Udc8pYFTx09vnLX1ymcgOUALY0hX8c+oyfLEqvkpPSjR0g+1iSJqAnL/t1j2MeLKREYkiN7h/PtcowwSdmf6qE/QxlF2FUfmn1DlRG2o7JwbxsZAX5WLIjRXxzMeGa74wrzoMCNJ; cfidsw-goldapple=D13SOjHJEIuA1lLRutPG2ItyssPz0peGIH5uAI7ESejq/I01mdyDTPegZQyx5dw1S46QS1vNih41uIkka+a9uUBgt3Y86QE9P+ImTvr1T3P8cAn29CVlDdItcX2kV90CEyUiMvEdrpU0f9q0mdYnhmY6Gv6fzNWhT/2GoSk=; fgsscw-goldapple=M9oxd39551598b25436e0829118b672f581bc637")
	req.Header.Set("x-gib-fgsscw-goldapple", "M9oxd39551598b25436e0829118b672f581bc637")
	req.Header.Set("x-gib-gsscw-goldapple", "tpElXqNoeo/tNGZ7GehI7w7HdzKaXzoca+zpCQxiZgIpBPhClgov0U689MBKfbX9aI4IRZX5gqr90ENFFiQY+gNBXAGK9V9rCNmYlctWP7E4388ieAyNrPCWK+u0sGg9Udc8pYFTx09vnLX1ymcgOUALY0hX8c+oyfLEqvkpPSjR0g+1iSJqAnL/t1j2MeLKREYkiN7h/PtcowwSdmf6qE/QxlF2FUfmn1DlRG2o7JwbxsZAX5WLIjRXxzMeGa74wrzoMCNJ")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return 0, err
    }

    log.Println("ответ ЗЯ:", string(body[:200]))

    var result GAResponse
    err = json.Unmarshal(body, &result)
    if err != nil {
        return 0, err
    }

    return result.Data.Variants[0].Price.Actual.Amount, nil
}


