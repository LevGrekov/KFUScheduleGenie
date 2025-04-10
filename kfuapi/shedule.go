package kfuapi

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (c *Client) GetSchedule(employeeID int) (string, error) {
	subjects, err := c.getTeacherSubjects(employeeID)
	if err != nil {
		return "", err
	}
	return formatTeacherSchedule(subjects), nil
}

func (c *Client) getTeacherSubjects(employeeID int) ([]Subject, error) {
	apiUrl := fmt.Sprintf("%s/employees/%d/schedule", baseURL, employeeID)

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	var resp ScheduleResponse
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("запрос не удался")
	}

	return resp.Subjects, nil
}

func formatTeacherSchedule(subjects []Subject) string {
	if len(subjects) == 0 {
		return "Расписание не найдено"
	}

	// Получаем имя преподавателя
	teacher := subjects[0]
	teacherName := fmt.Sprintf("%s %s %s ",
		teacher.TeacherLastname,
		teacher.TeacherFirstname,
		teacher.TeacherMiddlename)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📅 Расписание преподавателя: %s\n\n", teacherName))

	// Группируем по дням недели и времени
	schedule := make(map[int]map[string][]Subject)
	for _, subj := range subjects {
		if schedule[subj.DayWeekSchedule] == nil {
			schedule[subj.DayWeekSchedule] = make(map[string][]Subject)
		}
		timeSlot := fmt.Sprintf("%s-%s", subj.BeginTimeSchedule, subj.EndTimeSchedule)
		schedule[subj.DayWeekSchedule][timeSlot] = append(schedule[subj.DayWeekSchedule][timeSlot], subj)
	}

	// Выводим дни недели по порядку
	for day := 1; day <= 7; day++ {
		if timeSlots, ok := schedule[day]; ok {
			sb.WriteString(fmt.Sprintf("📌 %s:\n", weekdayName(day)))

			// Сортируем временные слоты
			var times []string
			for time := range timeSlots {
				times = append(times, time)
			}
			sort.Strings(times)

			for _, time := range times {
				subs := timeSlots[time]
				if len(subs) == 0 {
					continue
				}

				// Берем первый предмет из группы (они одинаковые по времени)
				subj := subs[0]

				// Формируем блоки информации
				timeInfo := fmt.Sprintf("🕒 %s", time)
				weekType := ""
				if subj.TypeWeekSchedule == 1 {
					weekType = " (неч.)"
				} else if subj.TypeWeekSchedule == 2 {
					weekType = " (чет.)"
				}
				period := formatPeriod(subj.StartDaySchedule, subj.FinishDaySchedule)
				sb.WriteString(timeInfo + weekType + period + " ")

				// Название предмета
				subjectInfo := fmt.Sprintf("📚 %s", subj.SubjectName)
				if subj.SubjectKindName != "" {
					subjectInfo += fmt.Sprintf(" (%s)", subj.SubjectKindName)
				}
				sb.WriteString(subjectInfo + " ")

				// Группы
				groupsInfo := fmt.Sprintf("👥 %s", formatGroupList(subj.GroupList))
				sb.WriteString(groupsInfo + " ;")

				// Аудитория (если указана)
				auditoryInfo := formatAuditory(subj.BuildingName, subj.NumAuditorium)
				if auditoryInfo != "" {
					sb.WriteString(auditoryInfo + " ")
				}

				if subj.NoteSchedule != "" {
					noteInfo := fmt.Sprintf("📝 %s", subj.NoteSchedule)
					sb.WriteString(noteInfo + " ")
				}
				sb.WriteString("\n")
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func weekdayName(day int) string {
	days := []string{"", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}
	if day >= 1 && day <= 7 {
		return days[day]
	}
	return "Неизвестный день"
}

func formatPeriod(start, end string) string {
	if start == "" || end == "" {
		return ""
	}
	startDate, err1 := time.Parse("02.01.06", start)
	endDate, err2 := time.Parse("02.01.06", end)
	if err1 != nil || err2 != nil {
		return ""
	}
	return fmt.Sprintf(" (%s-%s)", startDate.Format("02.01"), endDate.Format("02.01"))
}

func formatAuditory(building, room string) string {
	result := building
	if room != "" {
		result += fmt.Sprintf(", ауд. %s", room)
	}
	return result
}

func formatGroupList(groups string) string {
	if groups == "" {
		return ""
	}

	// Разбиваем строку на отдельные группы
	groupItems := strings.Split(groups, ", ")
	if len(groupItems) == 0 {
		return groups
	}

	// Сортируем группы для правильной обработки
	sort.Strings(groupItems)

	var result []string
	var currentPrefix string
	var startNum, prevNum int

	for i, group := range groupItems {
		// Разделяем префикс и номер (например "11" и "400")
		parts := strings.Split(group, "-")
		if len(parts) != 2 {
			// Если группа не в формате "XX-XXX", оставляем как есть
			result = append(result, group)
			continue
		}

		prefix := parts[0]
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			result = append(result, group)
			continue
		}

		if i == 0 {
			// Первая группа в списке
			currentPrefix = prefix
			startNum = num
			prevNum = num
			continue
		}

		if prefix == currentPrefix && num == prevNum+1 {
			// Продолжаем последовательность
			prevNum = num
		} else {
			// Завершаем текущую последовательность
			if startNum == prevNum {
				result = append(result, fmt.Sprintf("%s-%d", currentPrefix, startNum))
			} else {
				result = append(result, fmt.Sprintf("%s-%d..%d", currentPrefix, startNum, prevNum))
			}
			// Начинаем новую последовательность
			currentPrefix = prefix
			startNum = num
			prevNum = num
		}
	}

	// Добавляем последнюю последовательность
	if currentPrefix != "" {
		if startNum == prevNum {
			result = append(result, fmt.Sprintf("%s-%d", currentPrefix, startNum))
		} else {
			result = append(result, fmt.Sprintf("%s-%d..%d", currentPrefix, startNum, prevNum))
		}
	}

	return strings.Join(result, ", ")
}
