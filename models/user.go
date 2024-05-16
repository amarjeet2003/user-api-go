package models

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	DOB       string `json:"dob"`
}

// // MarshalJSON custom JSON marshaling for User struct
// func (u *User) MarshalJSON() ([]byte, error) {
// 	type Alias User
// 	return json.Marshal(&struct {
// 		DOB string `json:"dob"`
// 		*Alias
// 	}{
// 		DOB:   u.DOB.Format("2006-01-02"),
// 		Alias: (*Alias)(u),
// 	})
// }

// // UnmarshalJSON custom JSON unmarshaling for User struct
// func (u *User) UnmarshalJSON(data []byte) error {
// 	type Alias User
// 	aux := &struct {
// 		DOB string `json:"dob"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(u),
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	t, err := time.Parse("2006-01-02", aux.DOB)
// 	if err != nil {
// 		return err
// 	}
// 	u.DOB = t
// 	return nil
// }
