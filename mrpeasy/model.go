package mrpeasy

import (
	"encoding/json"
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

func (p *Timestamp) UnmarshalJSON(bytes []byte) error {
	var raw string
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	i, err := strconv.Atoi(raw)
	if err != nil {
		return err
	}

	p.Time = time.Unix(int64(i), 0)
	return nil
}

type Product struct {
	ArticleID      string `json:"article_id"`
	ProductID      string `json:"product_id"`
	ItemCode       string `json:"item_code"`
	ItemTitle      string `json:"item_title"`
	LotID          string `json:"lot_id"`
	LotCode        string `json:"lot_code"`
	LotStatus      string `json:"lot_status"`
	SiteID         string `json:"site_id"`
	Site           string `json:"site"`
	LocationID     string `json:"location_id"`
	Location       string `json:"location"`
	QuantityPicked int    `json:"quantity_picked"`
	QuantityBooked int    `json:"quantity_booked"`
	UnitID         string `json:"unit_id"`
	Unit           string `json:"unit"`
	ExpiryDate     string `json:"expiry_date"`
	LotStatusTxt   string `json:"lot_status_txt"`
}

type Order struct {
	CustomerOrderID   string `json:"customer_order_id"`
	CustomerOrderCode string `json:"customer_order_code"`
}

type Shipment struct {
	ShipmentID        string     `json:"shipment_id"`
	Code              string     `json:"code"`
	Created           Timestamp  `json:"created"`
	Status            string     `json:"status"`
	CustomerOrderID   string     `json:"customer_order_id"`
	CustomerOrderCode string     `json:"customer_order_code"`
	RmaOrderID        string     `json:"rma_order_id"`
	RmaOrderCode      string     `json:"rma_order_code"`
	PurchaseOrderID   string     `json:"purchase_order_id"`
	PurchaseOrderCode string     `json:"purchase_order_code"`
	WaybillNotes      string     `json:"waybill_notes"`
	PackingNotes      string     `json:"packing_notes"`
	TrackingNumber    string     `json:"tracking_number"`
	ShippingAddress   string     `json:"shipping_address"`
	StatusTxt         string     `json:"status_txt"`
	Products          []*Product `json:"products"`
	Orders            []*Order   `json:"orders"`
}

type PurchaseTerms struct {
	VendorID          string  `json:"vendor_id"`
	VendorCode        string  `json:"vendor_code"`
	VendorTitle       string  `json:"vendor_title"`
	VendorProductCode string  `json:"vendor_product_code"`
	Priority          float64 `json:"priority"`
	LeadTime          string  `json:"lead_time"`
	Unit              string  `json:"unit"`
	UnitRate          float64 `json:"unit_rate"`
	MinQuantity       float64 `json:"min_quantity"`
	VendorMinQuantity float64 `json:"vendor_min_quantity"`
	Price             float64 `json:"price"`
	CurrencyPrice     float64 `json:"currency_price"`
	Currency          string  `json:"currency"`
}

type Parameter struct {
	ParameterID string `json:"parameter_id"`
	Ord         string `json:"ord"`
	Title       string `json:"title"`
}

type StockItem struct {
	ArticleID  string  `json:"article_id"`
	ProductID  string  `json:"product_id"`
	Code       string  `json:"code"`
	Title      string  `json:"title"`
	UnitID     *string `json:"unit_id,omitempty"`
	GroupID    *string `json:"group_id,omitempty"`
	GroupCode  string  `json:"group_code"`
	GroupTitle string  `json:"group_title"`
	IsRaw      bool    `json:"is_raw"`
	// SellingPrice -> how to model? can be decimal or an array or null
	AvgCost           *float64         `json:"avg_cost,omitempty"`
	InStock           float64          `json:"in_stock"`
	Available         float64          `json:"available"`
	Booked            float64          `json:"booked"`
	ExpectedTotal     float64          `json:"expected_total"`
	ExpectedAvailable float64          `json:"expected_available"`
	ExpectedBooked    float64          `json:"expected_booked"`
	MinQuantity       string           `json:"min_quantity"`
	Icon              string           `json:"icon"`
	Deleted           bool             `json:"deleted"`
	PurchaseTerms     []*PurchaseTerms `json:"purchase_terms"`
	Parameters        []*Parameter     `json:"parameters"`
}

// type CustomerOrderStatus int

// const (
// 	CustomerOrderStatusQuotation              CustomerOrderStatus = 10
// 	CustomerOrderStatusWaitingForConfirmation CustomerOrderStatus = 20
// 	CustomerOrderStatusConfirmed              CustomerOrderStatus = 30
// 	CustomerOrderStatusWaitingForProduction   CustomerOrderStatus = 40
// 	CustomerOrderStatusInProduction           CustomerOrderStatus = 50
// 	CustomerOrderStatusReadyForShipment       CustomerOrderStatus = 60
// 	CustomerOrderStatusShippedAndInvoiced     CustomerOrderStatus = 70
// 	CustomerOrderStatusDelivered              CustomerOrderStatus = 80
// 	CustomerOrderStatusArchived               CustomerOrderStatus = 85
// 	CustomerOrderStatusCancelled              CustomerOrderStatus = 90
// )

// type CustomerOrderPartStatus int

// const (
// 	CustomerOrderPartStatusNotBooked        CustomerOrderPartStatus = 10
// 	CustomerOrderPartStatusRequested        CustomerOrderPartStatus = 15
// 	CustomerOrderPartStatusDelayed          CustomerOrderPartStatus = 20
// 	CustomerOrderPartStatusPossiblyDelayed  CustomerOrderPartStatus = 25
// 	CustomerOrderPartStatusExpectedOnTime   CustomerOrderPartStatus = 30
// 	CustomerOrderPartStatusReadyForShipment CustomerOrderPartStatus = 40
// 	CustomerOrderPartStatusDelivered        CustomerOrderPartStatus = 50
// )

type StockLot struct {
	LotID                       string `json:"lot_id"`
	LotCode                     string `json:"lot_code"`
	LotStatus                   string `json:"lot_status"`
	LotStatusTxt                string `json:"lot_status_txt"`
	ManOrdID                    string `json:"man_ord_id"`
	ManufacturingOrderCode      string `json:"manufacturing_order_code"`
	ManufacturingOrderStatus    string `json:"manufacturing_order_status"`
	ManufacturingOrderStatusTxt string `json:"manufacturing_order_status_txt"`
	PurOrdID                    string `json:"pur_ord_id"`
	PurchaseOrderCode           string `json:"purchase_order_code"`
	PurchaseOrderStatus         string `json:"purchase_order_status"`
	PurchaseOrderStatusTxt      string `json:"purchase_order_status_txt"`
}

type CustomerOrderProduct struct {
	LineID          string      `json:"line_id"`
	Ord             string      `json:"ord"`
	ArticleID       string      `json:"article_id"`
	Description     string      `json:"description"`
	Quantity        float64     `json:"quantity"`
	Shipped         float64     `json:"shipped"`
	DelvieryDate    Timestamp   `json:"delviery_date"`
	ItemPrice       float64     `json:"item_price"`
	ItemPriceCur    float64     `json:"item_price_cur"`
	TotalPrice      float64     `json:"total_price"`
	TotalPriceCur   float64     `json:"total_price_cur"`
	DiscountRate    float64     `json:"discount_rate"`
	DiscountRateCur float64     `json:"discount_rate_cur"`
	Cost            float64     `json:"cost"`
	Profit          float64     `json:"profit"`
	PartStatus      string      `json:"part_status"`
	PartStatusTxt   string      `json:"part_status_txt"`
	Source          []*StockLot `json:"source"`
}

type CustomerOrder struct {
	CustomerOrderID   string `json:"cust_ord_id"`
	Code              string `json:"code"`
	Reference         string `json:"reference"`
	CustomerID        string `json:"customer_id"`
	CustomerCode      string `json:"customer_code"`
	CustomerName      string `json:"customer_name"`
	PricelistID       int    `json:"pricelist_id"`
	PricelistCode     string `json:"pricelist_code"`
	PricelistTitle    string `json:"pricelist_title"`
	ShippingAddressID string `json:"shipping_address_id"`
	// ShippingAddress string or array -> How to model?
	Status             string                  `json:"status"`
	StatusTxt          string                  `json:"status_txt"`
	PartStatus         string                  `json:"part_status"`
	PartStatusTxt      string                  `json:"part_status_txt"`
	InvoiceStatus      string                  `json:"invoice_status"`
	InvoiceStatusTxt   string                  `json:"invoice_status_txt"`
	PaymentStatus      string                  `json:"payment_status"`
	PaymentStatusTxt   string                  `json:"payment_status_txt"`
	Created            Timestamp               `json:"created"`
	DeliveryDate       *Timestamp              `json:"delivery_date"`
	ActualDeliveryDate *Timestamp              `json:"actual_delivery_date"`
	Currency           string                  `json:"currency"`
	CurrencyRate       float64                 `json:"currency_rate"`
	TotalPrice         float64                 `json:"total_price"`
	TotalPriceCur      float64                 `json:"total_price_cur"`
	TotalCost          float64                 `json:"total_cost"`
	Profit             float64                 `json:"profit"`
	DiscountRate       float64                 `json:"discount_rate"`
	DiscountSum        float64                 `json:"discount_sum"`
	DiscountSumCur     float64                 `json:"discount_sum_cur"`
	Notes              string                  `json:"notes"`
	DeliveryTerms      string                  `json:"delivery_terms"`
	Products           []*CustomerOrderProduct `json:"products"`
}

type ContactDetails struct {
	LineID string          `json:"line_id"`
	Type   string          `json:"type"`
	Value  json.RawMessage `json:"value"` // this can be either a string or an object
}

type Customer struct {
	CustomerID    string            `json:"customer_id"`
	Code          string            `json:"code"`
	Title         string            `json:"title"`
	RegNr         string            `json:"reg_nr"`
	TaxNr         string            `json:"tax_nr"`
	Created       Timestamp         `json:"created"`
	NextContact   *Timestamp        `json:"next_contact"`
	Status        string            `json:"status"`
	PaymentPeriod int               `json:"payment_period"`
	UserID        string            `json:"user_id"`
	LanguageID    string            `json:"language_id"`
	ContactData   []*ContactDetails `json:"contact_data"`
}
