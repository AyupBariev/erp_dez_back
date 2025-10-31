package repositories

import (
	"database/sql"
	"encoding/json"
	"erp/internal/app/models"
	"erp/internal/utils"
	"fmt"
	"strings"
	"time"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type QueryOptions struct {
	Fields  []string
	Where   string
	Args    []interface{}
	Joins   []string
	OrderBy string
	GroupBy string
	Having  string
	Limit   int
	Offset  int
}

func (r *OrderRepository) queryOrders(opts QueryOptions) ([]*models.Order, error) {
	if len(opts.Fields) == 0 {
		opts.Fields = []string{
			"id", "erp_number", "engineer_id", "status", "client_name",
			"phones", "address", "note", "scheduled_at",
		}
	}

	query := "SELECT " + strings.Join(opts.Fields, ", ") + " FROM orders o"

	for _, join := range opts.Joins {
		query += " " + join
	}

	if strings.TrimSpace(opts.Where) != "" {
		query += " WHERE " + opts.Where
	}

	if opts.GroupBy != "" {
		query += " GROUP BY " + opts.GroupBy
	}

	if opts.OrderBy != "" {
		query += " ORDER BY " + opts.OrderBy
	}

	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", opts.Offset)
	}

	fmt.Println("SQL:", query, "ARGS:", opts.Args)

	rows, err := r.db.Query(query, opts.Args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	engineerIDs := make(map[int64]struct{}) // соберём ID инженеров для последующей выборки

	for rows.Next() {
		values := make([]interface{}, len(opts.Fields))
		ptrs := make([]interface{}, len(opts.Fields))
		for i := range opts.Fields {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		data := make(map[string]interface{}, len(opts.Fields))
		for i, f := range opts.Fields {
			fieldName := f
			if strings.Contains(strings.ToLower(fieldName), " as ") {
				parts := strings.Split(strings.ToLower(fieldName), " as ")
				fieldName = strings.TrimSpace(parts[1])
			} else if parts := strings.Split(fieldName, "."); len(parts) > 1 {
				fieldName = parts[1]
			}
			data[fieldName] = values[i]
		}

		order := &models.Order{}
		if err := utils.MapToStruct(data, order); err != nil {
			return nil, err
		}

		if order.EngineerID.Valid {
			engineerIDs[order.EngineerID.Int64] = struct{}{}
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 🔹 Автоматическая подгрузка инженеров
	if len(engineerIDs) > 0 {
		if err := r.attachEngineers(orders, engineerIDs); err != nil {
			return nil, err
		}
	}

	if err := r.attachProblems(orders); err != nil {
		return nil, err
	}

	if err := r.attachAggregators(orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) attachEngineers(orders []*models.Order, ids map[int64]struct{}) error {
	if len(ids) == 0 {
		return nil
	}

	// Преобразуем map → slice для IN (...)
	engineerIDs := make([]interface{}, 0, len(ids))
	for id := range ids {
		engineerIDs = append(engineerIDs, id)
	}

	// Соберём плейсхолдеры (?, ?, ?)
	placeholders := strings.Repeat("?,", len(engineerIDs))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`
		SELECT id, first_name, second_name, username, phone, telegram_id, is_approved
		FROM engineers
		WHERE id IN (%s)
	`, placeholders)

	rows, err := r.db.Query(query, engineerIDs...)
	if err != nil {
		return err
	}
	defer rows.Close()

	engineerMap := make(map[int64]*models.Engineer)

	for rows.Next() {
		var e models.Engineer
		err := rows.Scan(&e.ID, &e.FirstName, &e.SecondName, &e.Username, &e.Phone, &e.TelegramID, &e.IsApproved)
		if err != nil {
			return err
		}
		engineerMap[int64(e.ID)] = &e
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Привяжем инженеров к заказам
	for _, order := range orders {
		if order.EngineerID.Valid {
			if eng, ok := engineerMap[order.EngineerID.Int64]; ok {
				order.Engineer = eng
			}
		}
	}

	return nil
}
func (r *OrderRepository) attachProblems(orders []*models.Order) error {
	rows, err := r.db.Query(`SELECT id, name FROM problems`)
	if err != nil {
		return err
	}
	defer rows.Close()

	problems := map[int64]*models.BaseDictionary{}
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		problems[id] = &models.BaseDictionary{ID: int(id), Name: name}
	}

	for _, o := range orders {
		if o.ProblemID.Valid {
			o.Problem = problems[o.ProblemID.Int64]
		}
	}
	return nil
}

func (r *OrderRepository) attachAggregators(orders []*models.Order) error {
	rows, err := r.db.Query(`SELECT id, name FROM aggregators`)
	if err != nil {
		return err
	}
	defer rows.Close()

	aggs := map[int64]*models.BaseDictionary{}
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		aggs[id] = &models.BaseDictionary{ID: int(id), Name: name}
	}

	for _, o := range orders {
		if o.AggregatorID != 0 { // теперь просто проверяем на 0
			o.Aggregator = aggs[o.AggregatorID]
		}
	}

	return nil
}

func (r *OrderRepository) GetOrders(date *string) ([]*models.Order, error) {
	opts := QueryOptions{
		Fields: []string{
			"o.id",
			"o.erp_number",
			"o.client_name",
			"o.address",
			"o.our_percent",
			"o.price",
			"o.address",
			"o.work_volume",
			"o.scheduled_at",
			"o.status",
			"o.engineer_id",
			"o.admin_id",
			"o.phones",
			"o.problem_id",
			"o.aggregator_id",
		},
	}

	if date != nil && *date != "" {
		opts.Where = "DATE(o.scheduled_at) = ?"
		opts.Args = []interface{}{*date}
	}

	return r.queryOrders(opts)
}

// Получить максимальный ERP номер
func (r *OrderRepository) GetNextERPNumber() (int64, error) {
	var nextNumber int64
	err := r.db.QueryRow("SELECT COALESCE(MAX(erp_number), 0) + 1 FROM orders FOR UPDATE").Scan(&nextNumber)
	return nextNumber, err
}

// Создание нового заказа
func (r *OrderRepository) Create(order *models.Order) error {
	phonesJSON, _ := json.Marshal(order.Phones)

	res, err := r.db.Exec(`
    INSERT INTO orders (
        erp_number, aggregator_id, price, our_percent, client_name, phones,
        address, work_volume, note, scheduled_at, status, engineer_id, admin_id, problem_id
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		order.ERPNumber, order.AggregatorID, order.Price, order.OurPercent, order.ClientName,
		phonesJSON, order.Address, order.WorkVolume, order.Note,
		order.ScheduledAt, order.Status, order.EngineerID, order.AdminID, order.ProblemID,
	)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	order.ID = int(id)

	return r.db.QueryRow(`SELECT created_at, updated_at FROM orders WHERE id = ?`, id).
		Scan(&order.CreatedAt, &order.UpdatedAt)
}

// Получение заказа по ERP-номеру
func (r *OrderRepository) GetOrderByErpNumber(erpNumber int64) (*models.Order, error) {
	opts := QueryOptions{
		Fields: []string{
			"id", "erp_number", "engineer_id", "status", "client_name",
			"phones", "address", "note", "scheduled_at",
			"problem_id", "aggregator_id", "admin_id",
		},
		Where: "erp_number = ?",
		Args:  []interface{}{erpNumber},
	}

	orders, err := r.queryOrders(opts)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, sql.ErrNoRows
	}
	return orders[0], nil
}

// Заказы инженера за сегодня
func (r *OrderRepository) GetTodayOrders(chatID int64) ([]models.Order, error) {
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	opts := QueryOptions{
		Fields: []string{"id", "erp_number", "note", "status", "scheduled_at",
			"problem_id", "aggregator_id", "admin_id"},
		Where: "engineer_id IN (SELECT id FROM engineers WHERE telegram_id = ?) AND scheduled_at >= ? AND scheduled_at < ? AND status = 'working'",
		Args:  []interface{}{chatID, startOfDay, endOfDay},
	}

	orders, err := r.queryOrders(opts)
	if err != nil {
		return nil, err
	}

	var result []models.Order
	for _, o := range orders {
		result = append(result, *o)
	}
	return result, nil
}

// Повторные заказы инженера
func (r *OrderRepository) GetRepeatOrders(chatID int64) ([]models.Order, error) {
	opts := QueryOptions{
		Fields: []string{"id", "erp_number", "note", "status", "scheduled_at",
			"problem_id", "aggregator_id", "admin_id"},
		Where: "engineer_id IN (SELECT id FROM engineers WHERE telegram_id = ?) AND is_repeat = TRUE AND status = 'working'",
		Args:  []interface{}{chatID},
	}

	orders, err := r.queryOrders(opts)
	if err != nil {
		return nil, err
	}

	var result []models.Order
	for _, o := range orders {
		result = append(result, *o)
	}
	return result, nil
}

// Наличные заказы инженера
func (r *OrderRepository) GetCashOrders(chatID int64) ([]models.Order, error) {
	opts := QueryOptions{
		Fields: []string{"id", "erp_number", "note", "status", "scheduled_at",
			"problem_id", "aggregator_id", "admin_id"},
		Where: "engineer_id IN (SELECT id FROM engineers WHERE telegram_id = ?) AND payment_type = 'cash'",
		Args:  []interface{}{chatID},
	}

	orders, err := r.queryOrders(opts)
	if err != nil {
		return nil, err
	}

	var result []models.Order
	for _, o := range orders {
		result = append(result, *o)
	}
	return result, nil
}

// Обновление заказа (частичное)
func (r *OrderRepository) Update(order *models.Order) error {
	query := `
		UPDATE orders 
		SET engineer_id = ?, status = ?, confirmed_at = ?, updated_at = NOW()
		WHERE id = ?`
	_, err := r.db.Exec(query, order.EngineerID, order.Status, order.ConfirmedAt, order.ID)
	return err
}

func (r *OrderRepository) UpdateStatus(erpNumber int64, status string) error {
	query := `
		UPDATE orders 
		SET status = ?, updated_at = NOW()
		WHERE erp_number = ?`
	_, err := r.db.Exec(query, status, erpNumber)
	return err
}

// Полное обновление заказа
func (r *OrderRepository) UpdateFull(order *models.Order) error {
	query := `
		UPDATE orders 
		SET aggregator_id = ?, our_percent = ?, client_name = ?, phones = ?, 
		    address = ?, work_volume = ?, note = ?, scheduled_at = ?, status = ?, engineer_id = ?, updated_at = NOW()
		WHERE id = ?`
	_, err := r.db.Exec(query,
		order.AggregatorID, order.OurPercent, order.ClientName, order.Phones,
		order.Address, order.WorkVolume, order.Note, order.ScheduledAt,
		order.Status, order.EngineerID, order.ID)
	return err
}
