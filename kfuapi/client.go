package kfuapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Employee struct {
	ID         int    `json:"employee_id"`
	LastName   string `json:"lastname"`
	FirstName  string `json:"firstname"`
	MiddleName string `json:"middlename"`
	IsTeacher  bool   `json:"is_teacher"`

	FullName string `json:"fullname,omitempty"`
}

func filterTeachers(employees []Employee) []Employee {
	var teachers []Employee
	for _, emp := range employees {
		if emp.IsTeacher {
			teachers = append(teachers, emp)
		}
	}
	return teachers
}

func (e *Employee) GetFullName() string {
	if e.FullName != "" {
		return e.FullName
	}
	return fmt.Sprintf("%s %s %s", e.LastName, e.FirstName, e.MiddleName)
}

// Структура для ответа API
type APIResponse struct {
	Success   bool       `json:"success"`
	Employees []Employee `json:"employees"`
}

func SearchEmployees(fio string) ([]Employee, error) {
	encodedQuery := url.QueryEscape(fio)
	apiUrl := fmt.Sprintf("https://auth.kpfu.tyuop.ru/api/v1/employees?q=%s", encodedQuery)

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("User-Agent", "KPFU-Teacher-Search/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус: %d", resp.StatusCode)
	}

	// Используем io.ReadAll вместо ioutil.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("запрос не удался")
	}

	teachers := filterTeachers(apiResponse.Employees)

	return teachers, nil

}
