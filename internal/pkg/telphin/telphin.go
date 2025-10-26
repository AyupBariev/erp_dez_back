package telphin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TelphinClient struct {
	Token   string
	BaseURL string
}

func NewTelphinClient(token string) (*TelphinClient, error) {
	client := &TelphinClient{
		Token:   token,
		BaseURL: "https://api.telphin.ru/v2",
	}

	//// Если у тебя нет контекста, нужно создать его
	//ctx := context.Background()
	//
	//if err := client.; err != nil {
	//	return nil, fmt.Errorf("telphin ping failed: %w", err)
	//}

	return client, nil
}

// Call создает звонок через Telphin API
func (c *TelphinClient) Call(from, to string) (string, error) {
	payload := map[string]interface{}{
		"from": from,
		"to":   to,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/calls", c.BaseURL), bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		ID string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&res)

	return res.ID, nil
}

//
//func (h *OrderHandler) scheduleVoiceCalls(order *models.Order, eng *models.Engineer) {
//	ctx := context.Background()
//	virtualNumber := os.Getenv("TELPHIN_VIRTUAL_NUMBER")
//
//	now := time.Now()
//
//	// 1️⃣ Сразу после назначения — 3 звонка с интервалом 5 минут
//	for i := 1; i <= 3; i++ {
//		scheduledAt := now.Add(time.Duration((i-1)*5) * time.Minute)
//		key := fmt.Sprintf("call:%d:%d:%d", order.ID, eng.ID, i)
//		val, _ := json.Marshal(map[string]interface{}{
//			"from":         virtualNumber,
//			"to":           eng.Phone,
//			"scheduled_at": scheduledAt,
//			"status":       "queued",
//		})
//		h.Redis.Set(ctx, key, val, scheduledAt.Sub(now))
//
//		// Можно запускать worker, который будет проверять Redis и делать звонки
//	}
//
//	// 2️⃣ За час до ScheduledAt, если больше часа
//	if order.ScheduledAt.Sub(now) > time.Hour {
//		reminderTime := order.ScheduledAt.Add(-1 * time.Hour)
//		for i := 1; i <= 3; i++ {
//			scheduledAt := reminderTime.Add(time.Duration((i-1)*5) * time.Minute)
//			key := fmt.Sprintf("call:%d:%d:reminder%d", order.ID, eng.ID, i)
//			val, _ := json.Marshal(map[string]interface{}{
//				"from":         virtualNumber,
//				"to":           eng.Phone,
//				"scheduled_at": scheduledAt,
//				"status":       "queued",
//			})
//			h.Redis.Set(ctx, key, val, scheduledAt.Sub(now))
//		}
//	}
//}
