package model

type Response struct {
	Status       string           `json:"status"`
	Error        string           `json:"error,omitempty"`
	SongsLibrary *SongsConnection `json:"songLibrary,omitempty"`
	SongText     *TextConnection  `json:"songText,omitempty"`
}

type SongsConnection struct {
	Edges    []*SongEdge      `json:"edges"`
	PageInfo *LibraryPageInfo `json:"pageInfo"`
}

type LibraryPageInfo struct {
	EndCursor   *int64 `json:"endCursor,omitempty"`
	HasNextPage bool   `json:"hasNextPage"`
}

type SongEdge struct {
	Node   *Song `json:"node"`
	Cursor int64 `json:"cursor"`
}

type TextConnection struct {
	Edges    []*CoupletEdge `json:"edges"`
	PageInfo *TextPageInfo  `json:"pageInfo"`
}

type TextPageInfo struct {
	EndCursor   *int `json:"endCursor,omitempty"`
	HasNextPage bool `json:"hasNextPage"`
}

type CoupletEdge struct {
	Node   *string `json:"node"`
	Cursor int     `json:"cursor"`
}

func OK() Response {
	return Response{
		Status: "OK",
	}
}

func StatusError(msg string) Response {
	return Response{
		Status: "Error",
		Error:  msg,
	}
}
