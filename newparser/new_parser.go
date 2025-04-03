package newparser

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

func ParseSchedule(id int) (string, error) {
	url := fmt.Sprintf("https://kpfu.ru/main?p_id=%d&p_type=8&p_from=1", id)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при запросе:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Сервер вернул ошибку:", resp.Status)
	}

	utf8Reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println("Ошибка кодировки UTF8")
	}

	html_schedule, err := io.ReadAll(utf8Reader)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	formatted, err := formatScheduleForTelegram(string(html_schedule))
	if err != nil {
		return "", fmt.Errorf("ошибка при форматировании: %v", err)
	}

	return formatted, nil

}

func formatScheduleForTelegram(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	// Извлекаем имя преподавателя (если есть)
	name := doc.Find(".menu_header").Text()
	if name != "" {
		builder.WriteString(fmt.Sprintf("📅 %s\n\n", strings.TrimSpace(name)))
	}

	doc.Find("div[style*='background-image']").Each(func(i int, dayDiv *goquery.Selection) {
		day := strings.TrimSpace(dayDiv.Text())
		builder.WriteString(fmt.Sprintf("<b>📌 %s</b>\n", day))

		table := dayDiv.NextFiltered("table")
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			time := strings.TrimSpace(row.Find("td").First().Text())
			subject := strings.TrimSpace(row.Find("td").Last().Text())

			builder.WriteString(fmt.Sprintf("⏰ %s - %s\n", time, subject))
		})
		builder.WriteString("\n")
	})

	return builder.String(), nil
}
