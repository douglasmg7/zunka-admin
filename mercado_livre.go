package main

import "time"

// User info
type MercadoLivreUserInfo struct {
	ID               int    `json:"id"`
	Nickname         string `json:"nickname"`
	RegistrationDate string `json:"registration_date"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Gender           string `json:"gender"`
	CountryID        string `json:"country_id"`
	Email            string `json:"email"`
	Identification   struct {
		Number string `json:"number"`
		Type   string `json:"type"`
	} `json:"identification"`
	InternalTags []string `json:"internal_tags"`
	Address      struct {
		Address string `json:"address"`
		City    string `json:"city"`
		State   string `json:"state"`
		ZipCode string `json:"zip_code"`
	} `json:"address"`
	Phone struct {
		AreaCode  interface{} `json:"area_code"`
		Extension string      `json:"extension"`
		Number    string      `json:"number"`
		Verified  bool        `json:"verified"`
	} `json:"phone"`
	AlternativePhone struct {
		AreaCode  string `json:"area_code"`
		Extension string `json:"extension"`
		Number    string `json:"number"`
	} `json:"alternative_phone"`
	UserType         string      `json:"user_type"`
	Tags             []string    `json:"tags"`
	Logo             interface{} `json:"logo"`
	Points           int         `json:"points"`
	SiteID           string      `json:"site_id"`
	Permalink        string      `json:"permalink"`
	ShippingModes    []string    `json:"shipping_modes"`
	SellerExperience string      `json:"seller_experience"`
	BillData         struct {
		AcceptCreditNote string `json:"accept_credit_note"`
	} `json:"bill_data"`
	SellerReputation struct {
		LevelID           string      `json:"level_id"`
		PowerSellerStatus interface{} `json:"power_seller_status"`
		Transactions      struct {
			Canceled  int    `json:"canceled"`
			Completed int    `json:"completed"`
			Period    string `json:"period"`
			Ratings   struct {
				Negative float64 `json:"negative"`
				Neutral  int     `json:"neutral"`
				Positive float64 `json:"positive"`
			} `json:"ratings"`
			Total int `json:"total"`
		} `json:"transactions"`
		Metrics struct {
			Sales struct {
				Period    string `json:"period"`
				Completed int    `json:"completed"`
			} `json:"sales"`
			Claims struct {
				Period string  `json:"period"`
				Rate   float64 `json:"rate"`
				Value  int     `json:"value"`
			} `json:"claims"`
			DelayedHandlingTime struct {
				Period string  `json:"period"`
				Rate   float64 `json:"rate"`
				Value  int     `json:"value"`
			} `json:"delayed_handling_time"`
			Cancellations struct {
				Period string `json:"period"`
				Rate   int    `json:"rate"`
				Value  int    `json:"value"`
			} `json:"cancellations"`
		} `json:"metrics"`
	} `json:"seller_reputation"`
	BuyerReputation struct {
		CanceledTransactions int           `json:"canceled_transactions"`
		Tags                 []interface{} `json:"tags"`
		Transactions         struct {
			Canceled struct {
				Paid  interface{} `json:"paid"`
				Total interface{} `json:"total"`
			} `json:"canceled"`
			Completed   interface{} `json:"completed"`
			NotYetRated struct {
				Paid  interface{} `json:"paid"`
				Total interface{} `json:"total"`
				Units interface{} `json:"units"`
			} `json:"not_yet_rated"`
			Period  string      `json:"period"`
			Total   interface{} `json:"total"`
			Unrated struct {
				Paid  interface{} `json:"paid"`
				Total interface{} `json:"total"`
			} `json:"unrated"`
		} `json:"transactions"`
	} `json:"buyer_reputation"`
	Status struct {
		Billing struct {
			Allow bool          `json:"allow"`
			Codes []interface{} `json:"codes"`
		} `json:"billing"`
		Buy struct {
			Allow            bool          `json:"allow"`
			Codes            []interface{} `json:"codes"`
			ImmediatePayment struct {
				Reasons  []interface{} `json:"reasons"`
				Required bool          `json:"required"`
			} `json:"immediate_payment"`
		} `json:"buy"`
		ConfirmedEmail bool `json:"confirmed_email"`
		ShoppingCart   struct {
			Buy  string `json:"buy"`
			Sell string `json:"sell"`
		} `json:"shopping_cart"`
		ImmediatePayment bool `json:"immediate_payment"`
		List             struct {
			Allow            bool          `json:"allow"`
			Codes            []interface{} `json:"codes"`
			ImmediatePayment struct {
				Reasons  []interface{} `json:"reasons"`
				Required bool          `json:"required"`
			} `json:"immediate_payment"`
		} `json:"list"`
		Mercadoenvios          string      `json:"mercadoenvios"`
		MercadopagoAccountType string      `json:"mercadopago_account_type"`
		MercadopagoTcAccepted  bool        `json:"mercadopago_tc_accepted"`
		RequiredAction         interface{} `json:"required_action"`
		Sell                   struct {
			Allow            bool          `json:"allow"`
			Codes            []interface{} `json:"codes"`
			ImmediatePayment struct {
				Reasons  []interface{} `json:"reasons"`
				Required bool          `json:"required"`
			} `json:"immediate_payment"`
		} `json:"sell"`
		SiteStatus string `json:"site_status"`
		UserType   string `json:"user_type"`
	} `json:"status"`
	SecureEmail string `json:"secure_email"`
	Company     struct {
		BrandName      string      `json:"brand_name"`
		CityTaxID      string      `json:"city_tax_id"`
		CorporateName  string      `json:"corporate_name"`
		Identification string      `json:"identification"`
		StateTaxID     string      `json:"state_tax_id"`
		CustTypeID     string      `json:"cust_type_id"`
		SoftDescriptor interface{} `json:"soft_descriptor"`
	} `json:"company"`
	Credit struct {
		Consumed      int    `json:"consumed"`
		CreditLevelID string `json:"credit_level_id"`
		Rank          string `json:"rank"`
	} `json:"credit"`
	Context struct {
		Device string `json:"device"`
		Flow   string `json:"flow"`
		Source string `json:"source"`
	} `json:"context"`
	RegistrationIdentifiers []interface{} `json:"registration_identifiers"`
}

// Site search seller id
// /sites/search/$seller_id
type MercadoLivreSiteSerachSellerID struct {
	// SiteID string `json:"site_id"`
	// Seller struct {
	// ID               int    `json:"id"`
	// Nickname         string `json:"nickname"`
	// Permalink        string `json:"permalink"`
	// RegistrationDate string `json:"registration_date"`
	// SellerReputation struct {
	// LevelID           string      `json:"level_id"`
	// PowerSellerStatus interface{} `json:"power_seller_status"`
	// Transactions      struct {
	// Total    int    `json:"total"`
	// Canceled int    `json:"canceled"`
	// Period   string `json:"period"`
	// Ratings  struct {
	// Negative float64 `json:"negative"`
	// Positive float64 `json:"positive"`
	// Neutral  float64 `json:"neutral"`
	// } `json:"ratings"`
	// Completed int `json:"completed"`
	// } `json:"transactions"`
	// Metrics struct {
	// Sales struct {
	// Period    string `json:"period"`
	// Completed int    `json:"completed"`
	// } `json:"sales"`
	// } `json:"metrics"`
	// } `json:"seller_reputation"`
	// RealEstateAgency bool        `json:"real_estate_agency"`
	// CarDealer        bool        `json:"car_dealer"`
	// Tags             []string    `json:"tags"`
	// Eshop            interface{} `json:"eshop"`
	// } `json:"seller"`
	Paging struct {
		Total          int `json:"total"`
		PrimaryResults int `json:"primary_results"`
		Offset         int `json:"offset"`
		Limit          int `json:"limit"`
	} `json:"paging"`
	Results []struct {
		ID     string `json:"id"`
		SiteID string `json:"site_id"`
		Title  string `json:"title"`
		Seller struct {
			ID               int      `json:"id"`
			Permalink        string   `json:"permalink"`
			RegistrationDate string   `json:"registration_date"`
			CarDealer        bool     `json:"car_dealer"`
			RealEstateAgency bool     `json:"real_estate_agency"`
			Tags             []string `json:"tags"`
			SellerReputation struct {
				Transactions struct {
					Total    int    `json:"total"`
					Canceled int    `json:"canceled"`
					Period   string `json:"period"`
					Ratings  struct {
						Negative float64 `json:"negative"`
						Positive float64 `json:"positive"`
						Neutral  float64 `json:"neutral"`
					} `json:"ratings"`
					Completed int `json:"completed"`
				} `json:"transactions"`
				PowerSellerStatus interface{} `json:"power_seller_status"`
				Metrics           struct {
					Claims struct {
						Rate   float64 `json:"rate"`
						Value  int     `json:"value"`
						Period string  `json:"period"`
					} `json:"claims"`
					DelayedHandlingTime struct {
						Rate   int    `json:"rate"`
						Value  int    `json:"value"`
						Period string `json:"period"`
					} `json:"delayed_handling_time"`
					Sales struct {
						Period    string `json:"period"`
						Completed int    `json:"completed"`
					} `json:"sales"`
					Cancellations struct {
						Rate   int    `json:"rate"`
						Value  int    `json:"value"`
						Period string `json:"period"`
					} `json:"cancellations"`
				} `json:"metrics"`
				LevelID string `json:"level_id"`
			} `json:"seller_reputation"`
			Nickname string `json:"nickname"`
		} `json:"seller"`
		Price  int `json:"price"`
		Prices struct {
			ID     string `json:"id"`
			Prices []struct {
				ID         string `json:"id"`
				Type       string `json:"type"`
				Conditions struct {
					ContextRestrictions []interface{} `json:"context_restrictions"`
					StartTime           interface{}   `json:"start_time"`
					EndTime             interface{}   `json:"end_time"`
					Eligible            bool          `json:"eligible"`
				} `json:"conditions"`
				Amount              int         `json:"amount"`
				RegularAmount       interface{} `json:"regular_amount"`
				CurrencyID          string      `json:"currency_id"`
				ExchangeRateContext string      `json:"exchange_rate_context"`
				Metadata            struct {
				} `json:"metadata"`
				LastUpdated time.Time `json:"last_updated"`
			} `json:"prices"`
			Presentation struct {
				DisplayCurrency string `json:"display_currency"`
			} `json:"presentation"`
			PaymentMethodPrices []interface{} `json:"payment_method_prices"`
		} `json:"prices"`
		SalePrice          interface{} `json:"sale_price"`
		CurrencyID         string      `json:"currency_id"`
		AvailableQuantity  int         `json:"available_quantity"`
		SoldQuantity       int         `json:"sold_quantity"`
		BuyingMode         string      `json:"buying_mode"`
		ListingTypeID      string      `json:"listing_type_id"`
		StopTime           time.Time   `json:"stop_time"`
		Condition          string      `json:"condition"`
		Permalink          string      `json:"permalink"`
		Thumbnail          string      `json:"thumbnail"`
		ThumbnailID        string      `json:"thumbnail_id"`
		AcceptsMercadopago bool        `json:"accepts_mercadopago"`
		Installments       interface{} `json:"installments"`
		Address            struct {
			StateID   string `json:"state_id"`
			StateName string `json:"state_name"`
			CityID    string `json:"city_id"`
			CityName  string `json:"city_name"`
		} `json:"address"`
		Shipping struct {
			FreeShipping bool          `json:"free_shipping"`
			Mode         string        `json:"mode"`
			Tags         []interface{} `json:"tags"`
			LogisticType string        `json:"logistic_type"`
			StorePickUp  bool          `json:"store_pick_up"`
		} `json:"shipping"`
		SellerAddress struct {
			ID          string `json:"id"`
			Comment     string `json:"comment"`
			AddressLine string `json:"address_line"`
			ZipCode     string `json:"zip_code"`
			Country     struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"country"`
			State struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"state"`
			City struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"city"`
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"seller_address"`
		Attributes []struct {
			Values []struct {
				Struct interface{} `json:"struct"`
				Source int         `json:"source"`
				ID     string      `json:"id"`
				Name   string      `json:"name"`
			} `json:"values"`
			AttributeGroupID   string      `json:"attribute_group_id"`
			AttributeGroupName string      `json:"attribute_group_name"`
			Source             int         `json:"source"`
			ID                 string      `json:"id"`
			ValueID            string      `json:"value_id"`
			ValueName          string      `json:"value_name"`
			ValueStruct        interface{} `json:"value_struct"`
			Name               string      `json:"name"`
		} `json:"attributes"`
		DifferentialPricing struct {
			ID int `json:"id"`
		} `json:"differential_pricing"`
		OriginalPrice    interface{} `json:"original_price"`
		CategoryID       string      `json:"category_id"`
		OfficialStoreID  interface{} `json:"official_store_id"`
		DomainID         string      `json:"domain_id"`
		CatalogProductID string      `json:"catalog_product_id"`
		Tags             []string    `json:"tags"`
		CatalogListing   bool        `json:"catalog_listing,omitempty"`
		UseThumbnailID   bool        `json:"use_thumbnail_id"`
		OrderBackend     int         `json:"order_backend"`
	} `json:"results"`
	SecondaryResults []interface{} `json:"secondary_results"`
	RelatedResults   []interface{} `json:"related_results"`
	Sort             struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"sort"`
	AvailableSorts []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"available_sorts"`
	Filters []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Type   string `json:"type"`
		Values []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			PathFromRoot []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"path_from_root"`
		} `json:"values"`
	} `json:"filters"`
	AvailableFilters []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Type   string `json:"type"`
		Values []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Results int    `json:"results"`
		} `json:"values"`
	} `json:"available_filters"`
}
