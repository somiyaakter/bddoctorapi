package doctor

import "time"

type Doctor struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	BMDCRegNo       string    `json:"bmdc_reg_no"`
	Degrees         string    `json:"degrees"`
	ExperienceYears int       `json:"experience_years"`
	Specialties     string    `json:"specialties"`
	Designation     string    `json:"designation"`
	Workplace       string    `json:"workplace"`
	ImageURL        string    `json:"image_url"`
	ProfileURL      string    `json:"profile_url"`
	Chambers        []Chamber `json:"chambers"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Chamber struct {
	ID               int64  `json:"id"`
	DoctorID         int64  `json:"doctor_id"`
	Name             string `json:"name"`
	Address          string `json:"address"`
	VisitingHour     string `json:"visiting_hour"`
	AppointmentPhone string `json:"appointment_phone"`

	Doctor Doctor `json:"-"`
}
