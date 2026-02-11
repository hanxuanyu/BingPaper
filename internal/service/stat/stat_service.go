package stat

import (
	"BingPaper/internal/model"
	"BingPaper/internal/repo"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RecordStat 记录 API 调用统计
func RecordStat(endpoint, mkt string) {
	if mkt == "" {
		mkt = "default"
	}
	date := time.Now().Format("2006-01-02")

	// 异步记录
	go func() {
		stat := model.ApiStat{
			Date:     date,
			Endpoint: endpoint,
			Mkt:      mkt,
			Count:    1,
		}
		// 使用 OnConflict 实现聚合累加
		repo.DB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "date"}, {Name: "endpoint"}, {Name: "mkt"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"count":      gorm.Expr("count + ?", 1),
				"updated_at": time.Now(),
			}),
		}).Create(&stat)
	}()
}

// GetSummary 获取统计概览
func GetSummary() (map[string]any, error) {
	var totalCalls int64
	var todayCalls int64
	var yesterdayCalls int64

	repo.DB.Model(&model.ApiStat{}).Select("COALESCE(SUM(count), 0)").Scan(&totalCalls)

	today := time.Now().Format("2006-01-02")
	repo.DB.Model(&model.ApiStat{}).Where("date = ?", today).Select("COALESCE(SUM(count), 0)").Scan(&todayCalls)

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	repo.DB.Model(&model.ApiStat{}).Where("date = ?", yesterday).Select("COALESCE(SUM(count), 0)").Scan(&yesterdayCalls)

	return map[string]any{
		"total":     totalCalls,
		"today":     todayCalls,
		"yesterday": yesterdayCalls,
	}, nil
}

// GetTrend 获取调用趋势
func GetTrend(days int, endpoint string) ([]map[string]any, error) {
	if days <= 0 {
		days = 7
	}

	startDate := time.Now().AddDate(0, 0, -days+1).Format("2006-01-02")

	var results []struct {
		Date  string
		Count int64
	}

	query := repo.DB.Model(&model.ApiStat{}).Select("date, COALESCE(SUM(count), 0) as count").Where("date >= ?", startDate)
	if endpoint != "" {
		query = query.Where("endpoint = ?", endpoint)
	}
	err := query.Group("date").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Date] = r.Count
	}

	trend := make([]map[string]any, 0, days)
	for i := range days {
		d := time.Now().AddDate(0, 0, -days+1+i).Format("2006-01-02")
		trend = append(trend, map[string]any{
			"date":  d,
			"count": counts[d],
		})
	}

	return trend, nil
}

// GetEndpointDist 获取接口分布
func GetEndpointDist() ([]map[string]any, error) {
	var results []struct {
		Endpoint string
		Count    int64
	}
	err := repo.DB.Model(&model.ApiStat{}).Select("endpoint, COALESCE(SUM(count), 0) as count").Group("endpoint").Order("count desc").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	dist := make([]map[string]any, 0, len(results))
	for _, r := range results {
		dist = append(dist, map[string]any{
			"endpoint": r.Endpoint,
			"count":    r.Count,
		})
	}
	return dist, nil
}

// GetRegionDist 获取地区分布
func GetRegionDist() ([]map[string]any, error) {
	var results []struct {
		Mkt   string
		Count int64
	}
	err := repo.DB.Model(&model.ApiStat{}).Select("mkt, COALESCE(SUM(count), 0) as count").Group("mkt").Order("count desc").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	dist := make([]map[string]any, 0, len(results))
	for _, r := range results {
		dist = append(dist, map[string]any{
			"mkt":   r.Mkt,
			"count": r.Count,
		})
	}
	return dist, nil
}
