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

// Active confidential tags.
// Used to increase parsing speed.
var tags = []string{
	def,
	cvv,
	cardNumber,
	cardHolder,
	phone,
	ip,
}
