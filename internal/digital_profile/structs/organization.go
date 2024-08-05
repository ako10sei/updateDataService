package structs

import (
	"fmt"
)

type Organization struct {
	ID            int    `json:"id"`                // Идентификатор организации
	Name          string `json:"name"`              // Наименование организации
	ShortName     string `json:"short_name"`        // Краткое наименование организации
	Ogrn          any    `json:"ogrn"`              // ОГРН организации
	Director      any    `json:"director"`          // Руководитель организации
	Telephone     any    `json:"telephone"`         // Телефон
	Fax           any    `json:"fax"`               // Факс
	Email         any    `json:"email"`             // Email
	WebSite       any    `json:"web_site"`          // Web Site организации
	FAddressFull  any    `json:"f_address_full"`    // Полный фактический адрес организации
	UAddressFull  any    `json:"u_address_full"`    // Полный юридический адрес организации
	TerritoryName string `json:"territory_name"`    // Наименование района
	Parent        int    `json:"parent"`            // Родительская организация
	MaxOccupancy  int    `json:"maximum_occupancy"` // Проектная мощность организации
}

// GetValueByField возвращает карту, содержащую определенные поля структуры Organization
// в форматированном виде для дальнейшей обработки. Ключи карты представляют имена столбцов в таблице Visiology,
// а значения соответствуют полям структуры Organization.
func (o *Organization) GetValueByField() map[string]any {
	return map[string]any{
		"ОГРН":               o.Ogrn,
		"Адрес юридический":  o.UAddressFull,
		"Адрес фактический":  o.FAddressFull,
		"Руководитель":       fmt.Sprintf("%s, %s, %s", o.Director, o.Email, o.Telephone),
		"Район":              GetAreaIDByName(o.TerritoryName),
		"Проектная мощность": o.MaxOccupancy,
	}
}

// GetAreaIDByName возвращает идентификатор района на основе его названия.
// Функция использует словарь для хранения названий районов и их соответствующих идентификаторов.
func GetAreaIDByName(areaName string) int {
	areaMap := map[string]int{
		"Александровский район":  1,
		"Вязниковский район":     2,
		"Гороховецкий район":     3,
		"Гусь-Хрустальный район": 4,
		"Камешковский район":     5,
		"Киржачский район":       6,
		"Ковровский район":       7,
		"Кольчугинский район":    8,
		"Меленковский район":     9,
		"Муром":                  10,
		"Петушинский район":      11,
		"Селивановский район":    12,
		"Собинский район":        13,
		"Судогодский район":      14,
		"Суздальский район":      15,
		"Юрьев-Польский район":   16,
		"Владимир":               17,
		"Гусь-Хрустальный":       18,
		"Ковров":                 19,
		"Радужный":               20,
		"Муромский район":        21,
	}
	// Если указанное название района не найдено в словаре, функция возвращает -1.
	id, ok := areaMap[areaName]
	if !ok {
		return -1
	}

	return id
}
