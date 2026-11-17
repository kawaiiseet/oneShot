package server

type PersonalData struct {
	ID        int64    `json:"id,omitempty"`
	Name      string   `json:"name"`
	Login     string   `json:"login"`
	Age       int      `json:"age"`
	Password  string   `json:"password"`
	Country   string   `json:"country"`
	Interests []string `json:"interests"`
	Likes     []string `json:"likes"`
}

// FindRequest — тело /find.
// Можно передать login (тогда критерии берутся из профиля),
// либо country + interests напрямую.
type FindRequest struct {
	Login     string   `json:"login,omitempty"`
	Country   string   `json:"country,omitempty"`
	Interests []string `json:"interests,omitempty"`
	Likes     []string `json:"likes,omitempty"`
}
