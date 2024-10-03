package model

type Response struct {
	Status       string          `json:"status"`
	Error        string          `json:"error,omitempty"`
	MusikLibrary SongsConnection `json:"musikLibrary,omitempty"`
	SongText     string
}

type SongsConnection struct {
	Edges    []*SongEdge `json:"edge"`
	PageInfo *PageInfo   `json:"pageInfo"`
}

type PageInfo struct {
	EndCursor   *string `json:"endCursor,omitempty"`
	HasNextPage bool    `json:"hasNextPage"`
}

type SongEdge struct {
	Node   *Song  `json:"node"`
	Cursor string `json:"cursor"`
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
