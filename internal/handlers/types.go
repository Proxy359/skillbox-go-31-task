package handlers

type MakeFriendsRequest struct {
	SourceId int `json:"source_id"`
	TargetId int `json:"target_id"`
}

type DelieteUserRequest struct {
	TargetId int `json:"target_id"`
}

type NewAgeRequest struct {
	NewAge int `json:"new age"`
}

type CrateResponse struct {
	ID int `json:"ID"`
}
