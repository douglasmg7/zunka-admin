package main

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

// Active products
type MLProductsTitles struct {
	Paging struct {
		Total          int `json:"total"`
		PrimaryResults int `json:"primary_results"`
		Offset         int `json:"offset"`
		Limit          int `json:"limit"`
	} `json:"paging"`
	Results []struct {
		ID                string `json:"id"`
		Title             string `json:"title"`
		Price             int    `json:"price"`
		AvailableQuantity int    `json:"available_quantity"`
	} `json:"results"`
}
