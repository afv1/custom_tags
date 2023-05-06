package confidential

const confidentialTagKey = "confidential"

// tags
const (
	def        = "default"
	cvv        = "cvv"
	cardNumber = "cardnumber"
	cardHolder = "cardholder"
	phone      = "phone"
	ip         = "ip"
)

// Tags contains all active confidential tags.
// Used to increase parsing speed.
var Tags = []string{
	def,
	cvv,
	cardNumber,
	cardHolder,
	phone,
	ip,
}
