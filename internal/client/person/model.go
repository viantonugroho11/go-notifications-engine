package person

import "time"

type PersonProfileAddress struct {
	FullAddress string `json:"full_address"`
	CityID      int    `json:"city_id"`
	ProvinceID  int    `json:"province_id"`
}

type PersonProfile struct {
	DateOfBirth  string               `json:"date_of_birth"`
	PlaceOfBirth string               `json:"place_of_birth"`
	Address      PersonProfileAddress `json:"address"`
	Disability   bool                 `json:"disability"`
}

type PersonClass struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Grade        int    `json:"grade"`
	AcademicYear string `json:"academic_year"`
}

type PersonParent struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Relationship string         `json:"relationship"`
	Phone        string         `json:"phone"`
	Email        string         `json:"email"`
	Devices      []PersonDevice `json:"devices"`
}

type PersonDevice struct {
	DeviceType   string    `json:"device_type"`
	DeviceToken  string    `json:"device_token"`
	LastActiveAt time.Time `json:"last_active_at"`
}

type Person struct {
	ID        string         `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	FullName  string         `json:"full_name"`
	Email     string         `json:"email"`
	Phone     string         `json:"phone"`
	Gender    string         `json:"gender"`
	Type      string         `json:"type"`
	Profile   PersonProfile  `json:"profile"`
	Class     PersonClass    `json:"class"`
	Parents   []PersonParent `json:"parents"`
	Devices   []PersonDevice `json:"devices"`
}
