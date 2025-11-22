package dto

type UserPRStats struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	TotalPRs  int64  `json:"total_prs"`
	OpenPRs   int64  `json:"open_prs"`
	MergedPRs int64  `json:"merged_prs"`
}
