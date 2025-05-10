package domain

type ProfileStats struct {
	TotalFollowersCount int `json:"total_followers_count"`
	TotalFollowingCount int `json:"total_following_count"`

	TotalReviewsCount   int `json:"total_reviews_count"`
	TotalPlaylistsCount int `json:"total_playlists_count"`

	TotalLikesCount    int `json:"total_likes_count"`
	TotalDislikesCount int `json:"total_dislikes_count"`

	TotalCommentsCount int `json:"total_comments_count"`
	//TotalPhotosCount    int `json:"total_photos_count"`
	//TotalSongsCount int `json:"total_songs_count"`
}

type TrackStats struct {
	TotalLikesCount    int `json:"total_likes_count"`
	TotalDislikesCount int `json:"total_dislikes_count"`

	TotalReviewsCount int `json:"total_reviews_count"`

	TotalNotesCount int `json:"total_comments_count"`
}
