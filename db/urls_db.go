package db

import (
	"ohurlshortener/core"
	"ohurlshortener/utils"
)

var Max_Insert_Count = 5000

func FindShortUrl(url string) (core.ShortUrl, error) {
	found := core.ShortUrl{}
	query := `SELECT * FROM public.short_urls WHERE short_url = $1`
	err := Get(query, &found, url)
	return found, err
}

func FindAllShortUrls() ([]core.ShortUrl, error) {
	found := []core.ShortUrl{}
	query := `SELECT * FROM public.short_urls ORDER BY created_at DESC`
	err := Select(query, &found)
	return found, err
}

func FindPagedShortUrls(url string, page int, size int) ([]core.ShortUrl, error) {
	found := []core.ShortUrl{}
	offset := (page - 1) * size
	query := "SELECT * FROM public.short_urls u ORDER BY u.id DESC LIMIT $1 OFFSET $2"
	if !utils.EemptyString(url) {
		query := "SELECT * FROM public.short_urls u WHERE u.short_url = $1 ORDER BY u.id DESC LIMIT $2 OFFSET $3"
		var foundUrl core.ShortUrl
		err := Get(query, &foundUrl, url, size, offset)
		if !foundUrl.IsEmpty() {
			found = append(found, foundUrl)
		}
		return found, err
	}
	return found, Select(query, &found, size, offset)
}

func InsertShortUrl(url core.ShortUrl) error {
	query := `INSERT INTO public.short_urls (short_url, dest_url, created_at, is_valid, memo)
	 VALUES(:short_url,:dest_url,:created_at,:is_valid,:memo)`
	return NamedExec(query, url)
}

func GetUrlStats(url string) (core.ShortUrlStats, error) {
	found := core.ShortUrlStats{}
	query := `select * from public.url_ip_count_stats WHERE short_url = $1`
	err := Get(query, &found, url)
	return found, err
}

func splitLogsArray(array []core.AccessLog, size int) [][]core.AccessLog {
	var chunks [][]core.AccessLog
	for {
		if len(array) <= 0 {
			break
		}
		if len(array) < size {
			size = len(array)
		}
		chunks = append(chunks, array[0:size])
		array = array[size:]
	}
	return chunks
}