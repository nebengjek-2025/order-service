package token

type Claim struct {
	Iss      string   `json:"iss"`
	Metadata Metadata `json:"metadata"`
	Aud      string   `json:"aud"`
	Exp      string   `json:"exp"`
}

type Metadata struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
}
