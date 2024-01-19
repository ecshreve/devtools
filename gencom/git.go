package main

import (
	"os/exec"
)

type GitInterface interface {
	GetDiff() (string, error)
}

type GitCommand struct{}

func (g GitCommand) GetDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--unified=0")
	output, err := cmd.Output()
	return string(output), err
}

type MockGitCommand struct{}

func (mg MockGitCommand) GetDiff() (string, error) {
	const gitDiff = `diff --git a/internal/cache/tinycache.go b/internal/cache/tinycache.go
index b271459..2d4d8eb 100644
--- a/internal/cache/tinycache.go
+++ b/internal/cache/tinycache.go
@@ -24,0 +25,3 @@ type TinyCache struct {
+// TODO:
+// - Add a better way to seed the cache.
+// - Persist to disk
@@ -63 +66 @@ func (c *TinyCache) Set(key string, value interface{}) error {
-func (c *TinyCache) Reset(ctx context.Context) error {
+func (c *TinyCache) Reset(ctx context.Context) {
@@ -68,2 +70,0 @@ func (c *TinyCache) Reset(ctx context.Context) error {
-
-	return nil
@@ -114 +115 @@ func (c *TinyCache) CacheSummary(url string, summary string) error {
-func (c *TinyCache) Seed(data map[string]string) error {
+func (c *TinyCache) Seed(data map[string]string) {
@@ -121 +121,0 @@ func (c *TinyCache) Seed(data map[string]string) error {
-	return nil
diff --git a/internal/cache/tinycache_test.go b/internal/cache/tinycache_test.go
index c807705..749f5a9 100644
--- a/internal/cache/tinycache_test.go
+++ b/internal/cache/tinycache_test.go
@@ -83 +83,4 @@ func TestTinyCache_ConcurrentAccess(t *testing.T) {
-	c := cache.NewTinyCache()
+	env := utils.DefaultTestEnv(t)
+	defer env.Clean()
+
+	c := env.Server.Cache
@@ -99,4 +101,0 @@ func TestTinyCache_ConcurrentAccess(t *testing.T) {
-
-			// Reset the cache.
-			err = c.Reset(context.Background())
-			assert.NoError(t, err)
@@ -118,2 +117 @@ func TestTinyCache_SeedCache(t *testing.T) {
-	err := tc.Seed(seedData)
-	assert.NoError(t, err)
+	tc.Seed(seedData)
diff --git a/internal/server/feed.go b/internal/server/feed.go
index aa4629b..ee4ffd9 100644
--- a/internal/server/feed.go
+++ b/internal/server/feed.go
@@ -5,0 +6,2 @@ import (
+	"html"
+	"regexp"
@@ -55,0 +58,2 @@ func (s *Server) GetSummary(ctx context.Context, url string) (string, error) {
+		log.Debug("Summary not found in cache", "url", url)
+		log.Error(err)
@@ -58 +62 @@ func (s *Server) GetSummary(ctx context.Context, url string) (string, error) {
-		log.Info("Cached summary", "data", cachedSummary)
+		log.Debug("Cached summary", "data", cachedSummary)
@@ -80,14 +83,0 @@ func (s *Server) GetSummary(ctx context.Context, url string) (string, error) {
-// buildStringToSummarize builds a string that summarizes the feed by concatenating
-// the title and description of each item.
-// It takes a pointer to a Server struct and a pointer to a gofeed.Feed struct as parameters.
-// It returns a string that contains the conatenated title and description of each item.
-func (s *Server) buildStringToSummarize(feed *gofeed.Feed) string {
-	items := make([]string, 0)
-	for _, item := range feed.Items {
-		itemStr := fmt.Sprintf("Title: %s\nDescription: %s\n\n", item.Title,
-			item.Description)
-		items = append(items, itemStr)
-	}
-	return strings.Join(items, "\n")
-}
-
@@ -102 +92 @@ func (s *Server) SummarizeFeed(ctx context.Context, feed *gofeed.Feed) (string,
-	data := s.buildStringToSummarize(feed)
+	data := FormatFeedItems(feed.Items)
@@ -109,0 +100,28 @@ func (s *Server) SummarizeFeed(ctx context.Context, feed *gofeed.Feed) (string,
+
+// stripHTML takes a string and returns a string with all HTML tags removed.
+// It first unescapes HTML entities, then compiles a regex to match HTML tags.
+// It then replaces all HTML tags with an empty string and returns the result.
+// If any error occurs during the process, it returns an empty string.
+func stripHTML(input string) string {
+	unescaped := html.UnescapeString(input)
+
+	tagRegex := regexp.MustCompile(<[^>]+>)
+
+	return tagRegex.ReplaceAllString(unescaped, "")
+}
+
+// FormatFeedItems takes a slice of gofeed.Item pointers and returns a formatted string.
+// The string is formatted specifically for the summarizer prompt.
+// It iterates over the items and constructs a string for each item.
+// The string contains the item title and the description, with HTML tags removed.
+// It returns the formatted string.
+func FormatFeedItems(items []*gofeed.Item) string {
+	var sb strings.Builder
+
+	for i, item := range items {
+		description := stripHTML(item.Description)
+		sb.WriteString(fmt.Sprintf("%d. Title: %s\n   - Description: %s\n", i+1, item.Title, description))
+	}
+
+	return sb.String()
+}`

	return gitDiff, nil
}
