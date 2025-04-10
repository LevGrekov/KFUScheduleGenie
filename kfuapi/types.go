package kfuapi

type EmployeesResponse struct {
	Success   bool       `json:"success"`
	Employees []Employee `json:"employees"`
}

type ScheduleResponse struct {
	Success  bool      `json:"success"`
	Subjects []Subject `json:"subjects"`
}

type Subject struct {
	ID                string `json:"id"`
	Semester          int    `json:"semester"`
	Year              int    `json:"year"`
	SubjectName       string `json:"subject_name"`
	SubjectID         int    `json:"subject_id"`
	StartDaySchedule  string `json:"start_day_schedule"`
	FinishDaySchedule string `json:"finish_day_schedule"`
	DayWeekSchedule   int    `json:"day_week_schedule"`
	TypeWeekSchedule  int    `json:"type_week_schedule"`
	NoteSchedule      string `json:"note_schedule"`
	TotalTimeSchedule string `json:"total_time_schedule"`
	BeginTimeSchedule string `json:"begin_time_schedule"`
	EndTimeSchedule   string `json:"end_time_schedule"`
	TeacherID         int    `json:"teacher_id"`
	TeacherLastname   string `json:"teacher_lastname"`
	TeacherFirstname  string `json:"teacher_firstname"`
	TeacherMiddlename string `json:"teacher_middlename"`
	NumAuditorium     string `json:"num_auditorium_schedule"`
	BuildingName      string `json:"building_name"`
	BuildingID        string `json:"building_id"`
	GroupList         string `json:"group_list"`
	SubjectKindName   string `json:"subject_kind_name"`
}

type Employee struct {
	ID         int    `json:"employee_id"`
	LastName   string `json:"lastname"`
	FirstName  string `json:"firstname"`
	MiddleName string `json:"middlename"`
	IsTeacher  bool   `json:"is_teacher"`

	FullName string `json:"fullname,omitempty"`
}
