package kfuapi

import (
	"fmt"
	"net/http"
	"net/url"
)

func (e *Employee) GetFullName() string {
	if e.FullName != "" {
		return e.FullName
	}
	return fmt.Sprintf("%s %s %s", e.LastName, e.FirstName, e.MiddleName)
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

func (c *Client) SearchEmployees(fio string) ([]Employee, error) {
	encodedQuery := url.QueryEscape(fio)
	apiUrl := fmt.Sprintf("%s/employees?q=%s", baseURL, encodedQuery)

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	var resp EmployeesResponse
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("запрос не удался")
	}

	return filterTeachers(resp.Employees), nil
}
