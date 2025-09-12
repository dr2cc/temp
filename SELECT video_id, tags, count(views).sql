SELECT video_id, MAX(tags) as tags, MAX(views) as views
FROM videos
GROUP BY video_id
ORDER BY MAX(views) DESC
LIMIT 5;