package kfuapi

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/LevGrekov/KFUScheduleGenie/utils"
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

	// Получаем имя преподавателя из первого предмета
	teacher := subjects[0]
	teacherName := fmt.Sprintf("%s %s.%s.",
		teacher.TeacherLastname,
		teacher.TeacherFirstname,
		teacher.TeacherMiddlename)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📅 %s\n\n", teacherName))

	// Группируем по дням недели
	days := make(map[int][]Subject)
	for _, subj := range subjects {
		days[subj.DayWeekSchedule] = append(days[subj.DayWeekSchedule], subj)
	}

	// Выводим дни недели по порядку (1-7)
	for day := 1; day <= 7; day++ {
		if subs, ok := days[day]; ok {
			sb.WriteString(fmt.Sprintf("📌 %s:\n", utils.WeekdayName(day)))

			for _, subj := range subs {
				// Форматируем информацию о предмете
				weekType := ""
				if subj.TypeWeekSchedule == 1 {
					weekType = " неч."
				} else if subj.TypeWeekSchedule == 2 {
					weekType = " чет."
				}

				// Период проведения
				period := ""
				if subj.StartDaySchedule != "" && subj.FinishDaySchedule != "" {
					startDate, err1 := time.Parse("2006-01-02", subj.StartDaySchedule)
					endDate, err2 := time.Parse("2006-01-02", subj.FinishDaySchedule)
					if err1 == nil && err2 == nil {
						period = fmt.Sprintf(" (%s - %s)",
							startDate.Format("02.01"),
							endDate.Format("02.01"))
					}
				}

				timeInfo := fmt.Sprintf("🕒 %s-%s%s%s",
					subj.BeginTimeSchedule,
					subj.EndTimeSchedule,
					weekType,
					period)
				sb.WriteString(timeInfo + "\n")

				subjectInfo := "📚 " + subj.SubjectName
				if subj.SubjectKindName != "" {
					subjectInfo += fmt.Sprintf(" (%s)", subj.SubjectKindName)
				}
				sb.WriteString(subjectInfo + "\n")

				if subj.GroupList != "" {
					groupsInfo := fmt.Sprintf("👥 Группы: %s", subj.GroupList)
					sb.WriteString(groupsInfo + "\n")
				}

				auditoryInfo := "🏫" + subj.BuildingName
				if subj.NumAuditorium != "" {
					auditoryInfo += fmt.Sprintf(", ауд. %s", subj.NumAuditorium)
				}
				sb.WriteString(auditoryInfo + "\n")

				if subj.NoteSchedule != "" {
					noteInfo := fmt.Sprintf("📝 %s", subj.NoteSchedule)
					sb.WriteString(noteInfo + "\n")
				}

				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
}
