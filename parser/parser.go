package parser

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

const (
	baseURL                = "https://kpfu.ru"
	userAgent              = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	timeout                = 10 * time.Second
	schedule_link_selector = "#ss_content > div > div > div:nth-child(2) > div.right_width.right_block > ul > li:nth-child(7) > a"
	schedule_selector      = "#ss_content > div > div > div:nth-child(2) > div.left_width > div"
)

func ParseSchedule(name string, secondname string) (string, error) {

	// 1. Загружаем начальную страницу
	initialPage, err := fetchPage(baseURL + "/" + name + "." + secondname)
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки начальной страницы: %v", err)
	}

	// 2. Ищем ссылку на расписание
	scheduleLink, err := findScheduleLink(initialPage)
	if err != nil {
		return "", fmt.Errorf("не найдена ссылка на расписание: %v", err)
	}

	// 3. Загружаем страницу с расписанием
	schedulePage, err := fetchPage(scheduleLink)
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки страницы расписания: %v", err)
	}

	// 4. Извлекаем нужный блок с расписанием
	scheduleHTML, err := extractSchedule(schedulePage)
	if err != nil {
		return "", fmt.Errorf("ошибка извлечения расписания: %v", err)
	}

	// 5. Парсим HTML в удобный для Telegram формат
	formattedSchedule, err := formatScheduleForTelegram(scheduleHTML)
	if err != nil {
		return "", fmt.Errorf("ошибка форматирования расписания: %v", err)
	}

	return formattedSchedule, nil
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
		builder.WriteString(fmt.Sprintf("📅 *Расписание для %s*\n\n", strings.TrimSpace(name)))
	}

	// Обрабатываем каждый день
	doc.Find("div[style*='background-image']").Each(func(i int, dayDiv *goquery.Selection) {
		day := strings.TrimSpace(dayDiv.Text())
		builder.WriteString(fmt.Sprintf("📌 *%s*\n", day))

		// Находим следующую таблицу после этого div
		table := dayDiv.NextFiltered("table")
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			time := strings.TrimSpace(row.Find("td").First().Text())
			subject := strings.TrimSpace(row.Find("td").Last().Text())

			// Упрощаем текст
			subject = cleanSubjectText(subject)

			builder.WriteString(fmt.Sprintf("⏰ %s - %s\n", time, subject))
		})
		builder.WriteString("\n")
	})

	return builder.String(), nil
}

func cleanSubjectText(text string) string {
	// Удаляем лишние повторы и упрощаем текст
	replacements := map[string]string{
		"Учебное здание №14": "",
		"(гр.":               "гр.",
		"неч. нед.":          "(неч.)",
		"чет. нед.":          "(чет.)",
		"  ":                 " ", // Удаляем двойные пробелы
	}

	for old, new := range replacements {
		text = strings.ReplaceAll(text, old, new)
	}

	return strings.TrimSpace(text)
}

func fetchPage(url string) (*goquery.Document, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Устанавливаем заголовки как у обычного браузера
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP ошибка: %d", resp.StatusCode)
	}

	// Автоматически определяем кодировку и конвертируем в UTF-8
	utf8Reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromReader(utf8Reader)
}

func findScheduleLink(doc *goquery.Document) (string, error) {
	// Ищем ВСЕ ссылки на странице и фильтруем по тексту "Расписание"
	var scheduleLink string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if strings.TrimSpace(s.Text()) == "Расписание" {
			href, exists := s.Attr("href")
			if exists {
				scheduleLink = href
			}
		}
	})

	if scheduleLink == "" {
		return "", fmt.Errorf("ссылка с текстом 'Расписание' не найдена")
	}

	// Обрабатываем относительные ссылки
	if strings.HasPrefix(scheduleLink, "/") {
		return baseURL + scheduleLink, nil
	}

	// Проверяем валидность URL
	_, err := url.ParseRequestURI(scheduleLink)
	if err != nil {
		return "", fmt.Errorf("невалидный URL: %v", err)
	}

	return scheduleLink, nil
}

func extractSchedule(doc *goquery.Document) (string, error) {
	content := doc.Find(schedule_selector)
	if content.Length() == 0 {
		return "", fmt.Errorf("блок расписания не найден")
	}

	html, err := content.Html()
	if err != nil {
		return "", fmt.Errorf("ошибка извлечения HTML: %v", err)
	}

	return html, nil
}
