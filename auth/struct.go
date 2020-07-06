package auth

//Token struct
type Token struct {
	UserID int    `json:"user_id"`
	TTL    int64  `json:"ttl"`
	Email  string `json:"email"`
}

//LogoPassData struct for login/password sending
type LogoPassData struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

//JSONAnswer struct for mock test
type JSONAnswer struct {
	Token    string `json:"token"`
	Accepted bool   `json:"accepted"`
	Reason   string `json:"reason"`
	Data     string `json:"data"`
}
