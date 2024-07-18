package structs

type Organization struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	ShortName         string `json:"short_name"`
	UnitTypeName      string `json:"unit_type_name"`
	Kind              any    `json:"kind"`
	Ogrn              any    `json:"ogrn"`
	OgrnDate          any    `json:"ogrn_date"`
	Inn               string `json:"inn"`
	Kpp               string `json:"kpp"`
	Oktmo             any    `json:"oktmo"`
	Okpo              any    `json:"okpo"`
	Okato             string `json:"okato"`
	Okopf             any    `json:"okopf"`
	Okfs              any    `json:"okfs"`
	Okogu             any    `json:"okogu"`
	Director          any    `json:"director"`
	Telephone         any    `json:"telephone"`
	Fax               any    `json:"fax"`
	Email             any    `json:"email"`
	WebSite           any    `json:"web_site"`
	FAddress          any    `json:"f_address"`
	FAddressStreet    any    `json:"f_address_street"`
	FAddressHouse     any    `json:"f_address_house"`
	FAddressFlat      any    `json:"f_address_flat"`
	FAddressHouseGUID any    `json:"f_address_house_guid"`
	FZipcode          any    `json:"f_zipcode"`
	FAddressFull      any    `json:"f_address_full"`
	UAddress          any    `json:"u_address"`
	UAddressStreet    any    `json:"u_address_street"`
	UAddressHouse     any    `json:"u_address_house"`
	UAddressFlat      any    `json:"u_address_flat"`
	UZipcode          any    `json:"u_zipcode"`
	UAddressFull      any    `json:"u_address_full"`
	UAddressHouseGUID any    `json:"u_address_house_guid"`
	TerritoryID       any    `json:"territory_id"`
	TerritoryName     any    `json:"territory_name"`
	TerritoryCode     any    `json:"territory_code"`
	TerritoryOkato    any    `json:"territory_okato"`
	TerritoryTypeID   int    `json:"territory_type_id"`
	TerritoryTypeName string `json:"territory_type_name"`
	Founders          []any  `json:"founders"`
	Parent            any    `json:"parent"`
}
