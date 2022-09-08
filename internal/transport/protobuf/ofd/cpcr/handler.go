package cpcr

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/golang/protobuf/proto"

	"github.com/JohanVong/fiscal_v3_sample/internal/models"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/graylog"
	"github.com/JohanVong/fiscal_v3_sample/internal/transport/protobuf/ofd/kkm_proto"
)

var (
	appCode = uint16(0x81A2)
	version = uint16(124)
)

type OFD struct {
	tcpAddr *net.TCPAddr
	timeout time.Duration
}

type OfdHeader struct {
	Id     uint32
	Token  uint32
	ReqNum uint16
}

type OfdSystemPart struct {
	OfflinePeriod *OfflinePeriod
}

type OfflinePeriod struct {
	BeginTime time.Time
	EndTime   time.Time
}

type OfdResult struct {
	Header   *OfdHeader
	Response []interface{}
}

type TicketNumber struct {
	Number  uint64
	Offline bool
	OFDQR   string
}

type ofdRequest struct {
	request  *kkm_proto.Request
	header   *OfdHeader
	respChan chan *ofdResponse
}

type ofdResponse struct {
	resp *kkm_proto.Response
	h    *OfdHeader
	err  error
}

func (h *OFD) SendRaw(request []byte) (response []byte) {
	conn, err := net.DialTimeout("tcp", h.tcpAddr.String(), h.timeout)

	if err != nil {

		if conn != nil {
			conn.Close()
		}

		return
	}
	//conn.SetDeadline(time.Now().Add(h.timeout))
	reqHeader, reqPayload := extractHeader(request)

	p1 := ofdRequest{
		request: new(kkm_proto.Request),
	}
	_ = proto.Unmarshal(reqPayload, p1.request)
	m1, err := proto.Marshal(p1.request)

	m, _ := createMessage(*reqHeader, m1)
	_ = m
	jsonStr := base64.StdEncoding.EncodeToString(m)
	log.Println(jsonStr)
	_, err = conn.Write(m)

	if err != nil {

		conn.Close()
		return
	}

	response, err = h.receiveResponse(conn)
	if err != nil {
		return
	}
	respHeader, respPayload := extractHeader(response)
	log.Println(respHeader)
	resp := ofdResponse{
		resp: new(kkm_proto.Response),
	}
	err = proto.Unmarshal(respPayload, resp.resp)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp.resp)
	return
}

func (h *OFD) SetTimeout(timeout time.Duration) {
	h.timeout = timeout
}

func NewOFDHandler(connString string) (handler *OFD, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", connString)

	if err != nil {
		return nil, err
	}

	handler = new(OFD)
	handler.tcpAddr = tcpAddr
	handler.timeout = time.Second * 10

	return handler, nil
}

func (h *OFD) SendSystemMessage(args ...interface{}) (interface{}, error) {
	reqHeader, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	req := ofdRequest{
		request: &kkm_proto.Request{
			Command: kkm_proto.CommandTypeEnum_COMMAND_SYSTEM.Enum(),
		},
	}

	req.header = reqHeader
	req.respChan = make(chan *ofdResponse)
	resp := h.performRequest(&req)

	if resp.err != nil {
		return nil, resp.err
	}

	return resp.h, nil
}

func (h *OFD) SendOfflineModeInformation(start time.Time, args ...interface{}) (interface{}, error) {
	reqHeader, err := findHeader(args)
	if err != nil {
		return nil, err
	}
	endTime := time.Now()
	req := ofdRequest{
		request: &kkm_proto.Request{
			Command: kkm_proto.CommandTypeEnum_COMMAND_SYSTEM.Enum(),
			Service: &kkm_proto.ServiceRequest{
				OfflinePeriod: &kkm_proto.ServiceRequest_OfflinePeriod{
					BeginTime: NewDateTime(&start),
					EndTime:   NewDateTime(&endTime),
				},
			},
		},
	}

	req.header = reqHeader
	req.respChan = make(chan *ofdResponse)
	resp := h.performRequest(&req)

	if resp.err != nil {
		return nil, resp.err
	}

	return resp.h, nil
}

func (h *OFD) GetKKMInfo(args ...interface{}) (interface{}, error) {
	reqHeader, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	req := ofdRequest{
		request: &kkm_proto.Request{
			Command: kkm_proto.CommandTypeEnum_COMMAND_INFO.Enum(),
		},
	}

	req.header = reqHeader
	req.respChan = make(chan *ofdResponse)
	resp := h.performRequest(&req)

	if resp.err != nil {
		return nil, resp.err
	}
	result := OfdResult{
		Header: resp.h,
	}
	result.Response = append(result.Response, resp.resp.Report)
	result.Response = append(result.Response, resp.resp.Service)
	return result, resp.err
}

func (h *OFD) SendOperationPurchase(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	systemPart, err := findSystemPart(args)

	return h.performTicketOperation(kkm_proto.OperationTypeEnum_OPERATION_BUY, document, items, header, systemPart)
}

func (h *OFD) SendOperationPurchaseReturn(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	systemPart, err := findSystemPart(args)

	return h.performTicketOperation(kkm_proto.OperationTypeEnum_OPERATION_BUY_RETURN, document, items, header, systemPart)

}

func (h *OFD) SendOperationSale(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	systemPart, err := findSystemPart(args)

	return h.performTicketOperation(kkm_proto.OperationTypeEnum_OPERATION_SELL, document, items, header, systemPart)
}

func (h *OFD) SendOperationSaleReturn(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	systemPart, err := findSystemPart(args)

	return h.performTicketOperation(kkm_proto.OperationTypeEnum_OPERATION_SELL_RETURN, document, items, header, systemPart)
}

func (h *OFD) SendTicketRollback(args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	req := ofdRequest{
		request: &kkm_proto.Request{
			Command: kkm_proto.CommandTypeEnum_COMMAND_CANCEL_TICKET.Enum(),
		},
	}

	req.header = header
	req.respChan = make(chan *ofdResponse)
	resp := h.performRequest(&req)
	result := OfdResult{
		Header: resp.h,
	}

	return result, resp.err
}

func (h *OFD) CloseShift(document *models.Document, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	systemPart, err := findSystemPart(args)

	req := ofdRequest{
		request: &kkm_proto.Request{
			Command: kkm_proto.CommandTypeEnum_COMMAND_CLOSE_SHIFT.Enum(),
			CloseShift: &kkm_proto.CloseShiftRequest{
				CloseTime: NewDateTime(&document.ReceiptDate),
			},
			Service: createServiceRequest(systemPart),
		},
	}

	req.header = header
	req.respChan = make(chan *ofdResponse)
	resp := h.performRequest(&req)
	result := OfdResult{
		Header: resp.h,
	}
	if resp.resp != nil {
		if resp.resp.Report != nil {
			result.Response = append(result.Response, resp.resp.Report)
		}
	}

	return result, resp.err

}

func (h *OFD) RequestZReport(document *models.Document, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, nil
	}

	systemPart, err := findSystemPart(args)

	req := ofdRequest{
		request: &kkm_proto.Request{
			Command: kkm_proto.CommandTypeEnum_COMMAND_REPORT.Enum(),
			Report: &kkm_proto.ReportRequest{
				Report:   kkm_proto.ReportTypeEnum_REPORT_Z.Enum(),
				DateTime: NewDateTime(&document.ReceiptDate),
			},
			Service: createServiceRequest(systemPart),
		},
	}

	req.header = header
	req.respChan = make(chan *ofdResponse)
	resp := h.performRequest(&req)
	result := OfdResult{
		Header: resp.h,
	}
	if resp.resp != nil {
		if resp.resp.Report != nil {
			result.Response = append(result.Response, resp.resp.Report)
		}
	}
	return result, resp.err
}

func (h *OFD) RequestXReport(document *models.Document, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	systemPart, err := findSystemPart(args)

	req := ofdRequest{
		request: &kkm_proto.Request{
			Command: kkm_proto.CommandTypeEnum_COMMAND_REPORT.Enum(),
			Report: &kkm_proto.ReportRequest{
				Report:   kkm_proto.ReportTypeEnum_REPORT_X.Enum(),
				DateTime: NewDateTime(&document.ReceiptDate),
			},
			Service: createServiceRequest(systemPart),
		},
	}

	req.header = header
	req.respChan = make(chan *ofdResponse)
	resp := h.performRequest(&req)
	result := OfdResult{
		Header: resp.h,
	}
	if resp.resp != nil {
		if resp.resp.Report != nil {
			result.Response = append(result.Response, resp.resp.Report)
		}
	}

	return result, resp.err
}

func (h *OFD) DepositMoney(document *models.Document, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	systemPart, err := findSystemPart(args)

	req := ofdRequest{
		request: &kkm_proto.Request{
			Command: kkm_proto.CommandTypeEnum_COMMAND_MONEY_PLACEMENT.Enum(),
			MoneyPlacement: &kkm_proto.MoneyPlacementRequest{
				Operation: kkm_proto.MoneyPlacementEnum_MONEY_PLACEMENT_DEPOSIT.Enum(),
				Datetime:  NewDateTime(&document.ReceiptDate),
				Sum:       NewMoneyDecimal(&document.Value.Big),
				Operator: &kkm_proto.TicketRequest_Operator{
					Code: proto.Uint32(uint32(document.IdUser)),
				},
			},
			Service: createServiceRequest(systemPart),
		},
	}

	if document.User != nil {
		req.request.MoneyPlacement.Operator.Name = proto.String(document.User.Name)
	}

	if document.Offline == true {
		req.request.MoneyPlacement.IsOffline = proto.Bool(true)
	}

	req.header = header
	req.respChan = make(chan *ofdResponse)
	resp := h.performRequest(&req)
	result := OfdResult{
		Header: resp.h,
	}

	return result, resp.err
}

func (h *OFD) WithdrawMoney(document *models.Document, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, nil
	}

	systemPart, err := findSystemPart(args)

	req := ofdRequest{
		request: &kkm_proto.Request{
			Command: kkm_proto.CommandTypeEnum_COMMAND_MONEY_PLACEMENT.Enum(),
			MoneyPlacement: &kkm_proto.MoneyPlacementRequest{
				Operation: kkm_proto.MoneyPlacementEnum_MONEY_PLACEMENT_WITHDRAWAL.Enum(),
				Datetime:  NewDateTime(&document.ReceiptDate),
				Sum:       NewMoneyDecimal(&document.Value.Big),
				Operator: &kkm_proto.TicketRequest_Operator{
					Code: proto.Uint32(uint32(document.IdUser)),
				},
			},
			Service: createServiceRequest(systemPart),
		},
	}

	if document.User != nil {
		req.request.MoneyPlacement.Operator.Name = proto.String(document.User.Name)
	}

	if document.Offline == true {
		req.request.MoneyPlacement.IsOffline = proto.Bool(true)
	}

	req.header = header
	req.respChan = make(chan *ofdResponse)
	resp := h.performRequest(&req)
	result := OfdResult{
		Header: resp.h,
	}

	return result, resp.err
}

func (h *OFD) roundTrip(req *ofdRequest) {
	//defer recoverpanic.RecoverPanic()
	l := graylog.NewGELFMessage("OFD")
	l.Address = h.tcpAddr.String()
	l.IDOFD = req.header.Id
	if req.header.ReqNum < math.MaxUint16 {
		req.header.ReqNum++
	} else {
		req.header.ReqNum = 1
	}
	l.TokenReqnum = strconv.FormatUint(uint64(req.header.Token), 10) + "|" + strconv.FormatUint(uint64(req.header.
		ReqNum), 10)
	defer func() {
		logSendMessage, _ := json.Marshal(req.request)
		l.SendMessage = string(logSendMessage)
		_, _ = graylog.SendMessage(l)
	}()
	response := new(ofdResponse)

	payload, err := proto.Marshal(req.request)
	if err != nil {
		l.ReceiveMessage = err.Error()
		response.err = err
		req.respChan <- response
		return
	}

	message, err := createMessage(*req.header, payload)
	//log.Println(base64.StdEncoding.EncodeToString(message))
	if err != nil {
		l.ReceiveMessage = err.Error()
		response.err = err
		req.respChan <- response
		return
	}

	conn, err := net.DialTimeout("tcp", h.tcpAddr.String(), h.timeout)

	if err != nil {
		l.ReceiveMessage = err.Error()
		response.err = err
		req.respChan <- response
		if conn != nil {
			conn.Close()
		}

		return
	}
	conn.SetDeadline(time.Now().Add(h.timeout))
	*l.SendDate = time.Now()
	_, err = conn.Write(message)

	if err != nil {
		l.ReceiveMessage = err.Error()
		response.err = err
		req.respChan <- response
		conn.Close()
		return
	}

	message, err = h.receiveResponse(conn)

	*l.ReceiveDate = time.Now()
	if err != nil {
		l.ReceiveMessage = err.Error()
		response.err = err
		req.respChan <- response
		return
	}

	if message == nil {
		l.ReceiveMessage = err.Error()
		response.err = err
		req.respChan <- response
		return
	}

	respHeader, respPayload := extractHeader(message)

	resp := kkm_proto.Response{}

	err = proto.Unmarshal(respPayload, &resp)
	if err != nil {
		l.ReceiveMessage = err.Error()
		response.err = err
		req.respChan <- response
		return
	}

	//if respHeader.ReqNum < math.MaxUint16 {
	//	respHeader.ReqNum++
	//} else {
	//	respHeader.ReqNum = 0
	//}
	if resp.Result.ResultCode != nil {
		if *resp.Result.ResultCode != 0 {
			response.err = NewOfdError(*resp.Result.ResultCode, *resp.Result.ResultText)
			l.ErrorMessageOFD = response.err.Error()
			l.ErrorCodeOFD = *resp.Result.ResultCode
		}
	}
	response.resp = &resp
	response.h = respHeader

	req.respChan <- response

	respJson, _ := json.MarshalIndent(response.resp, "", "\t")
	l.ReceiveMessage = string(respJson)
}

func (h *OFD) performRequest(req *ofdRequest) (resp *ofdResponse) {
	req.respChan = make(chan *ofdResponse)
	go h.roundTrip(req)

	resp = <-req.respChan
	return
}

func createMessage(header OfdHeader, payload []byte) (message []byte, err error) {
	message = make([]byte, 18)

	var messageLen uint32

	if len(payload)+18 > math.MaxInt32 {
		return nil, errors.New("payload is too long")
	}
	messageLen = uint32(len(payload)) + 18

	binary.LittleEndian.PutUint16(message, appCode)
	binary.LittleEndian.PutUint16(message[2:], version)
	binary.LittleEndian.PutUint32(message[4:], messageLen)
	binary.LittleEndian.PutUint32(message[8:], header.Id)
	binary.LittleEndian.PutUint32(message[12:], header.Token)
	binary.LittleEndian.PutUint16(message[16:], header.ReqNum)
	message = append(message, payload...)

	return
}

func (h *OFD) receiveResponse(conn net.Conn) (message []byte, err error) {
	defer conn.Close()
	headerBuffer := make([]byte, 18)
	timer := time.NewTimer(h.timeout)

	for {

		select {
		case <-timer.C:

			return nil, errors.New("connection timed out while awaiting response")
		default:
			size, err := conn.Read(headerBuffer)

			if err != nil && err != io.EOF {
				return nil, err
			}

			if size < len(headerBuffer) {
				continue
			}

			messageSize := binary.LittleEndian.Uint32(headerBuffer[4:])

			messageBuffer := make([]byte, messageSize-18)

			_, err = conn.Read(messageBuffer)
			if err != nil && err != io.EOF {
				return nil, err
			}

			message = append(headerBuffer, messageBuffer...)
			return message, nil

		}
	}

}

func extractHeader(message []byte) (h *OfdHeader, payload []byte) {
	h = new(OfdHeader)
	h.Id = binary.LittleEndian.Uint32(message[8:])
	h.Token = binary.LittleEndian.Uint32(message[12:])
	h.ReqNum = binary.LittleEndian.Uint16(message[16:])

	return h, message[18:]

}

func findHeader(args []interface{}) (*OfdHeader, error) {
	if len(args) == 0 {
		return nil, errors.New("header is missing")
	}
	var header *OfdHeader
	for _, arg := range args {
		switch a := arg.(type) {
		case *OfdHeader:
			header = a
			break
		}
	}

	if header == nil {
		return nil, errors.New("header is missing")
	}

	return header, nil
}

func findSystemPart(args []interface{}) (*OfdSystemPart, error) {
	if len(args) == 0 {
		return nil, errors.New("system part is missing")
	}
	var systemPart *OfdSystemPart
	for _, arg := range args {
		switch a := arg.(type) {
		case *OfdSystemPart:
			systemPart = a
			break
		}
	}

	if systemPart == nil {
		return nil, errors.New("system part is missing")
	}

	return systemPart, nil
}

func createServiceRequest(systemPart *OfdSystemPart) *kkm_proto.ServiceRequest {
	if systemPart == nil {
		return nil
	}

	var serviceRequest = new(kkm_proto.ServiceRequest)

	if systemPart.OfflinePeriod != nil {
		serviceRequest.OfflinePeriod = &kkm_proto.ServiceRequest_OfflinePeriod{
			BeginTime: NewDateTime(&systemPart.OfflinePeriod.BeginTime),
			EndTime:   NewDateTime(&systemPart.OfflinePeriod.EndTime),
		}

	}

	return serviceRequest
}

func (h *OFD) performTicketOperation(ticketType kkm_proto.OperationTypeEnum, document *models.Document,
	items []models.Position, header *OfdHeader, systemPart *OfdSystemPart) (result OfdResult, err error) {
	req := ofdRequest{
		request: &kkm_proto.Request{
			Command: kkm_proto.CommandTypeEnum_COMMAND_TICKET.Enum(),
			Ticket: &kkm_proto.TicketRequest{
				Operation: ticketType.Enum(),
				DateTime:  NewDateTime(&document.ReceiptDate),
				Operator: &kkm_proto.TicketRequest_Operator{
					Code: proto.Uint32(uint32(document.IdUser)),
				},
			},
			Service: createServiceRequest(systemPart),
		},
	}

	if document.AutonomousNumber > 0 {
		req.request.Ticket.OfflineTicketNumber = proto.Uint32(document.AutonomousNumber)
	}

	if document.Cash.Cmp(decimal.New(0, 0)) > 0 {
		req.request.Ticket.Payments = append(req.request.Ticket.Payments, &kkm_proto.TicketRequest_Payment{
			Type: kkm_proto.PaymentTypeEnum_PAYMENT_CASH.Enum(),
			Sum:  NewMoneyDecimal(new(decimal.Big).Sub(&document.Cash.Big, &document.Change.Big)),
		})
	}

	if document.NonCash.Cmp(decimal.New(0, 0)) > 0 {
		req.request.Ticket.Payments = append(req.request.Ticket.Payments, &kkm_proto.TicketRequest_Payment{
			Type: kkm_proto.PaymentTypeEnum_PAYMENT_CARD.Enum(),
			Sum:  NewMoneyDecimal(&document.NonCash.Big),
		})
	}

	if document.Mobile.Cmp(decimal.New(0, 0)) > 0 {
		req.request.Ticket.Payments = append(req.request.Ticket.Payments, &kkm_proto.TicketRequest_Payment{
			Type: kkm_proto.PaymentTypeEnum_PAYMENT_MOBILE.Enum(),
			Sum:  NewMoneyDecimal(&document.Mobile.Big),
		})
	}

	if document.User != nil {
		req.request.Ticket.Operator.Name = proto.String(document.User.Name)
	}

	switch document.IdDomain {
	case 1:
		req.request.Ticket.Domain = &kkm_proto.TicketRequest_Domain{
			Type: kkm_proto.DomainTypeEnum_DOMAIN_TRADING.Enum(),
		}

	case 2:
		req.request.Ticket.Domain = &kkm_proto.TicketRequest_Domain{
			Type: kkm_proto.DomainTypeEnum_DOMAIN_SERVICES.Enum(),
		}
	case 3:
		req.request.Ticket.Domain = &kkm_proto.TicketRequest_Domain{
			Type: kkm_proto.DomainTypeEnum_DOMAIN_GASOIL.Enum(),
		}
	case 4:
		req.request.Ticket.Domain = &kkm_proto.TicketRequest_Domain{
			Type: kkm_proto.DomainTypeEnum_DOMAIN_HOTELS.Enum(),
		}
	case 5:
		req.request.Ticket.Domain = &kkm_proto.TicketRequest_Domain{
			Type: kkm_proto.DomainTypeEnum_DOMAIN_TAXI.Enum(),
		}
	case 6:
		req.request.Ticket.Domain = &kkm_proto.TicketRequest_Domain{
			Type: kkm_proto.DomainTypeEnum_DOMAIN_PARKING.Enum(),
		}
	}
	var totalD = new(decimal.Big)
	var totalM = new(decimal.Big)
	var totalDNds = new(decimal.Big)
	var totalMNds = new(decimal.Big)
	var modifiers []*kkm_proto.TicketRequest_Item

	for _, i := range items {
		var commodityTotal = new(decimal.Big)
		var discountTotal = new(decimal.Big)
		var markupTotal = new(decimal.Big)
		var quantity = &i.Qty.Big
		commodityTotal.Mul(&i.Price.Big, quantity)
		//commodityTotal := i.Price * i.Qty
		var commodityNdsStruct kkm_proto.TicketRequest_Tax
		var discountNdsStruct kkm_proto.TicketRequest_Tax
		var markupNdsStruct kkm_proto.TicketRequest_Tax
		discountTotal.Mul(&i.Discount.Big, quantity)
		markupTotal.Mul(&i.Markup.Big, quantity)

		f64qty, ok := i.Qty.Big.Float64()
		if ok != true {
			return result, errors.New("Quantity conversion failed")
		}

		commodityItem := kkm_proto.TicketRequest_Item{
			Type: kkm_proto.ItemTypeEnum_ITEM_TYPE_COMMODITY.Enum(),
			Commodity: &kkm_proto.TicketRequest_Item_Commodity{
				SectionCode: proto.String(strconv.Itoa(i.IdSection)),
				Name:        proto.String(i.Name),
				Price:       NewMoneyDecimal(&i.Price.Big),
				Quantity:    proto.Uint32(uint32(f64qty * 1000)),
				Sum:         NewMoneyDecimal(commodityTotal),
				ExciseStamp: proto.String(strings.TrimSpace(i.ProductCode)),
			},
		}
		if i.Nds.Cmp(decimal.New(0, 0)) > 0 {
			//commodityNDS := new(decimal.Big).Mul(new(decimal.Big).Add(&i.Nds.Big, new(decimal.Big).Sub(&i.NdsMarkup.Big, &i.NdsDiscount.Big)), quantity)
			commodityNDS := new(decimal.Big).Mul(&i.Nds.Big, quantity)

			commodityNdsStruct = kkm_proto.TicketRequest_Tax{
				TaxType:      proto.Uint32(100),
				TaxationType: proto.Uint32(101),
				Percent:      proto.Uint32(uint32(i.Section.Nds * 1000)),
				Sum:          NewMoneyDecimal(commodityNDS),
				IsInTotalSum: proto.Bool(true),
			}
			commodityItem.Commodity.Taxes = append(commodityItem.Commodity.Taxes, &commodityNdsStruct)
		}

		req.request.Ticket.Items = append(req.request.Ticket.Items, &commodityItem)
		if i.Markup.Cmp(decimal.New(0, 0)) > 0 {
			if i.Storno == false {
				totalM.Add(totalM, &i.Markup.Big)
			}
			markup := kkm_proto.TicketRequest_Item{
				Type: kkm_proto.ItemTypeEnum_ITEM_TYPE_MARKUP.Enum(),
				Markup: &kkm_proto.TicketRequest_Modifier{
					Sum:  NewMoneyDecimal(markupTotal),
					Name: proto.String(i.Name),
				},
			}
			if i.Nds.Cmp(decimal.New(0, 0)) > 0 {
				markupNDS := new(decimal.Big).Mul(&i.NdsMarkup.Big, quantity)
				markupNdsStruct = kkm_proto.TicketRequest_Tax{
					TaxType:      proto.Uint32(100),
					TaxationType: proto.Uint32(101),
					Percent:      proto.Uint32(uint32(i.Section.Nds * 1000)),
					Sum:          NewMoneyDecimal(markupNDS),
					IsInTotalSum: proto.Bool(true),
				}
				markup.Markup.Taxes = append(markup.Markup.Taxes, &markupNdsStruct)
				if i.Storno == false {
					totalMNds.Add(totalMNds, markupNDS)
				}

			}
			modifiers = append(modifiers, &markup)

		}

		if i.Discount.Cmp(decimal.New(0, 0)) > 0 {
			if i.Storno == false {
				totalD.Add(totalD, &i.Discount.Big)
			}
			discount := kkm_proto.TicketRequest_Item{
				Type: kkm_proto.ItemTypeEnum_ITEM_TYPE_DISCOUNT.Enum(),
				Discount: &kkm_proto.TicketRequest_Modifier{
					Sum:  NewMoneyDecimal(discountTotal),
					Name: proto.String(i.Name),
				},
			}
			if i.Nds.Cmp(decimal.New(0, 0)) > 0 {
				discountNDS := new(decimal.Big).Mul(&i.NdsDiscount.Big, quantity)
				discountNdsStruct = kkm_proto.TicketRequest_Tax{
					TaxType:      proto.Uint32(100),
					TaxationType: proto.Uint32(101),
					Percent:      proto.Uint32(uint32(i.Section.Nds * 1000)),
					Sum:          NewMoneyDecimal(discountNDS),
					IsInTotalSum: proto.Bool(true),
				}
				discount.Discount.Taxes = append(discount.Discount.Taxes, &discountNdsStruct)
				if i.Storno == false {
					totalDNds.Add(totalDNds, discountNDS)
				}

			}
			modifiers = append(modifiers, &discount)
		}

		if i.Storno == true {
			stornoItem := kkm_proto.TicketRequest_Item{
				Type: kkm_proto.ItemTypeEnum_ITEM_TYPE_STORNO_COMMODITY.Enum(),
				StornoCommodity: &kkm_proto.TicketRequest_Item_StornoCommodity{
					SectionCode: proto.String(strconv.Itoa(i.IdSection)),
					Name:        proto.String(i.Name),
					Price:       NewMoneyDecimal(&i.Price.Big),
					Quantity:    proto.Uint32(uint32(f64qty) * 1000),
					Sum:         NewMoneyDecimal(commodityTotal),
				},
			}
			if i.Nds.Cmp(decimal.New(0, 0)) > 0 {
				stornoItem.StornoCommodity.Taxes = append(stornoItem.StornoCommodity.Taxes, &commodityNdsStruct)
			}

			req.request.Ticket.Items = append(req.request.Ticket.Items, &stornoItem)
			if i.Markup.Cmp(decimal.New(0, 0)) > 0 {
				stornoMarkup := kkm_proto.TicketRequest_Item{
					Type: kkm_proto.ItemTypeEnum_ITEM_TYPE_STORNO_MARKUP.Enum(),
					StornoMarkup: &kkm_proto.TicketRequest_Modifier{
						Sum:  NewMoneyDecimal(markupTotal),
						Name: proto.String(i.Name),
					},
				}
				if i.Nds.Cmp(decimal.New(0, 0)) > 0 {
					stornoMarkup.StornoMarkup.Taxes = append(stornoMarkup.StornoMarkup.Taxes, &markupNdsStruct)
				}
				modifiers = append(modifiers, &stornoMarkup)
			}

			if i.Discount.Cmp(decimal.New(0, 0)) > 0 {
				stornoDiscount := kkm_proto.TicketRequest_Item{
					Type: kkm_proto.ItemTypeEnum_ITEM_TYPE_STORNO_DISCOUNT.Enum(),
					StornoDiscount: &kkm_proto.TicketRequest_Modifier{
						Sum:  NewMoneyDecimal(discountTotal),
						Name: proto.String(i.Name),
					},
				}
				if i.Nds.Cmp(decimal.New(0, 0)) > 0 {
					stornoDiscount.StornoDiscount.Taxes = append(stornoDiscount.StornoDiscount.Taxes, &discountNdsStruct)
				}
				modifiers = append(modifiers, &stornoDiscount)
			}
		}
	}

	if len(modifiers) > 0 {
		req.request.Ticket.Items = append(req.request.Ticket.Items, modifiers...)
	}

	req.request.Ticket.Amounts = &kkm_proto.TicketRequest_Amounts{
		Total:  NewMoneyDecimal(&document.Value.Big),
		Taken:  NewMoneyDecimal(&document.Cash.Big),
		Change: NewMoneyDecimal(&document.Change.Big),
	}

	req.header = header
	req.respChan = make(chan *ofdResponse)
	resp := h.performRequest(&req)
	result = OfdResult{
		Header: resp.h,
	}
	if resp.resp != nil {
		result.Response = append(result.Response, resp.resp.Ticket)
	}

	return result, resp.err
}

func NewMoneyFloat64(sum float64) (money *kkm_proto.Money) {
	w, d := math.Modf(sum)
	money = &kkm_proto.Money{
		Bills: proto.Uint64(uint64(w)),
		Coins: proto.Uint32(uint32(math.Round(d * 100))),
	}

	return
}

func NewDateTime(t *time.Time) *kkm_proto.DateTime {
	return &kkm_proto.DateTime{
		Time: &kkm_proto.Time{
			Hour:   proto.Uint32(uint32(t.Hour())),
			Minute: proto.Uint32(uint32(t.Minute())),
			Second: proto.Uint32(uint32(t.Second())),
		},
		Date: &kkm_proto.Date{
			Day:   proto.Uint32(uint32(t.Day())),
			Month: proto.Uint32(uint32(t.Month())),
			Year:  proto.Uint32(uint32(t.Year())),
		},
	}
}

func NewMoneyDecimal(sum *decimal.Big) (money *kkm_proto.Money) {
	floatSum, _ := sum.Float64()
	w, d := math.Modf(floatSum)

	bs := proto.Uint64(uint64(w))
	cs := proto.Uint32(uint32(math.Round(d * 100)))
	if *cs > 99 {
		*bs++
		*cs = 0
	}

	money = &kkm_proto.Money{
		Bills: bs,
		Coins: cs,
	}

	return
}
