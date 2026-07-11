package customs

// TrendingTopicsQuery aggregates the most-used hashtags across all posts,
// ordered by frequency then tag. Bind the LIMIT (?) at call time.
func TrendingTopicsQuery() string {
	return `
SELECT tag, COUNT(*) AS posts
FROM post_hashtags
GROUP BY tag
ORDER BY posts DESC, tag ASC
LIMIT ?`
}
