package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/api_client"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/payment-gateway/domain"
	handler_dto "code.evixo.ru/ncc/ncc-backend/services/payment-gateway/handlers/dto"
	"code.evixo.ru/ncc/ncc-backend/services/payment-gateway/usecases"
	usecase_dto "code.evixo.ru/ncc/ncc-backend/services/payment-gateway/usecases/dto"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	QueryTypeCheck = "check"
	QueryTypePay   = "pay"
)

const (
	ResponseOK                        = 0
	ResponseTemporaryError            = 1
	ResponseInternalError             = 2
	ResponseWrongCustomerIdFormat     = 3
	ResponseCustomerNotFound          = 21
	ResponsePaymentForbidden          = 22
	ResponsePaymentTemporaryForbidden = 23
	ResponseCustomerInactive          = 24
	ResponseUnableToCheck             = 25
	ResponsePaymentInProgress         = 100
	ResponseAmountTooLow              = 241
	ResponseAmountTooHigh             = 242
	ResponseUnknownError              = 299
)

type CitypayPayment struct {
	cfg *config.Config
	log logger.Logger
}

func NewCitypayPayment(
	cfg *config.Config,
	log logger.Logger,
) CitypayPayment {
	return CitypayPayment{
		cfg: cfg,
		log: log,
	}
}

type H map[string]any

type XMLPaymentResponse struct {
	XMLName        xml.Name `xml:"Response"`
	TransactionId  string   `xml:"TransactionId"`
	TransactionExt string   `xml:"TransactionExt"`
	Amount         float64  `xml:"Amount"`
	ResultCode     int      `xml:"ResultCode"`
	Fields         Fields   `xml:"Fields"`
	Comment        string   `xml:"Comment"`
}

type XMLCheckResponse struct {
	XMLName       xml.Name `xml:"Response"`
	TransactionId string   `xml:"TransactionId"`
	ResultCode    int      `xml:"ResultCode"`
	Fields        Fields   `xml:"Fields"`
	Comment       string   `xml:"Comment"`
}

type XMLErrorResponse struct {
	XMLName       xml.Name `xml:"Response"`
	TransactionId string   `xml:"TransactionId"`
	ResultCode    int      `xml:"ResultCode"`
	Comment       string   `xml:"Comment"`
}

type Fields struct {
	Fields []NameField `xml:"fields"`
}

type NameField struct {
	XMLName xml.Name `xml:"field1"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",chardata"`
}

// MarshalXML allows type H to be used with xml.Marshal.
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "Response",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func (h CitypayPayment) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := handler_dto.NewCitypayPaymentRequestDTO()
		err := request.Parse(ctx)
		if err != nil {
			h.log.Error("[%s] Can't parse params: %v: %s", request.TransactionID, err.Error(), ctx.Request.RequestURI)
			ctx.Status(http.StatusBadRequest)
			return
		}

		paymentClient := api_client.NewPaymentClient(h.cfg.API.URL, h.cfg.API.Token)
		customerClient := api_client.NewCustomerClient(h.cfg.API.URL, h.cfg.API.Token)

		paymentSystem, err := paymentClient.GetPaymentSystemByID(request.PaymentSystemID)
		if err != nil {
			h.log.Error("[%s] Can't get payment system [id=%s]: %v", request.TransactionID, request.PaymentSystemID, err.Error())
			ctx.Status(http.StatusUnauthorized)
			return
		}

		if paymentSystem.Token != request.Token {
			h.log.Error("[%s] Invalid token: %s", request.TransactionID, request.Token)
			ctx.Status(http.StatusUnauthorized)
			return
		}

		customer, err := customerClient.GetByUID(request.Account)
		if err != nil {
			h.log.Error("[%s] Can't get customer: [UID=%s] %v", request.TransactionID, request.Account, err.Error())
			ctx.XML(http.StatusOK, h.errorResponse(request.TransactionID, ResponseCustomerNotFound, "can't get customer"))
			return
		}

		switch request.QueryType {
		case QueryTypePay:
			paymentUsecase := usecases.NewPaymentUsecase(h.cfg, h.log, paymentClient)

			usecaseRequest := usecase_dto.NewPaymentRequestDTO(
				customer.ID,
				paymentSystem.UserName,
				paymentSystem.PaymentTypeID,
				paymentSystem.DstAccountID,
				request.Amount,
				paymentSystem.PaymentDescr,
				request.TransactionID,
				paymentSystem.ID,
				request.IP,
			)

			payment, err := paymentUsecase.Execute(usecaseRequest)
			if err != nil {
				h.log.Error("[%s] Can't execute: %v", request.TransactionID, err.Error())
				ctx.XML(http.StatusOK, h.errorResponse(request.TransactionID, ResponseUnknownError, "payment error"))
				return
			}

			if payment == nil {
				h.log.Error("[%s] payment is nil", request.TransactionID)
				ctx.XML(http.StatusOK, h.errorResponse(request.TransactionID, ResponseUnknownError, "payment error"))
				return
			}

			h.log.Info("[%s] Successful payment: account=%s amount=%0.2f pid=%d", request.TransactionID, request.Account, request.Amount, payment.Pid)

			ctx.XML(http.StatusOK, h.paymentResponse(request.TransactionID, strconv.Itoa(payment.Pid), payment.Amount, ResponseOK, customer.Name, customer.Name))
		case QueryTypeCheck:
			checkUsecase := usecases.NewCheckUsecase(h.cfg, h.log, paymentClient)

			usecaseRequest := usecase_dto.NewCheckRequestDTO(customer.ID)

			err = checkUsecase.Execute(usecaseRequest)
			if err != nil {
				if err == domain.ErrPaymentNotAllowed {
					h.log.Error("[%s] Payment not allowed for [id=%s]", request.TransactionID, customer.ID)
					ctx.XML(http.StatusOK, h.errorResponse(request.TransactionID, ResponsePaymentForbidden, "payment not allowed for customer"))
					return
				}
				h.log.Error("[%s] Can't execute: %v", request.TransactionID, err.Error())
				ctx.XML(http.StatusOK, h.errorResponse(request.TransactionID, ResponseUnknownError, "check error"))
				return
			}

			h.log.Info("[%s] Check OK [UID=%s]", request.TransactionID, request.Account)

			ctx.XML(http.StatusOK, h.checkResponse(request.TransactionID, ResponseOK, customer.Name, customer.Name))
		default:
			h.log.Error("[%s] Unknown QueryType: %s", request.TransactionID, request.QueryType)
			ctx.Status(http.StatusBadRequest)
			return
		}
	}
}

func (h CitypayPayment) paymentResponse(
	transactionId string,
	transactionExt string,
	amount float64,
	resultCode int,
	comment string,
	customerName string,
) XMLPaymentResponse {
	return XMLPaymentResponse{
		TransactionId:  transactionId,
		TransactionExt: transactionExt,
		Amount:         amount,
		ResultCode:     resultCode,
		Comment:        comment,
		Fields: Fields{
			[]NameField{
				{Name: "customerName", Value: customerName},
			},
		},
	}
}

func (h CitypayPayment) checkResponse(
	transactionId string,
	resultCode int,
	comment string,
	customerName string,
) XMLCheckResponse {
	return XMLCheckResponse{
		TransactionId: transactionId,
		ResultCode:    resultCode,
		Comment:       comment,
		Fields: Fields{
			[]NameField{
				{Name: "customerName", Value: customerName},
			},
		},
	}
}

func (h CitypayPayment) errorResponse(
	transactionId string,
	resultCode int,
	comment string,
) XMLErrorResponse {
	return XMLErrorResponse{
		TransactionId: transactionId,
		ResultCode:    resultCode,
		Comment:       comment,
	}
}
