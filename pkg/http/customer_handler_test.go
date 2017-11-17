package http_test

import (
	"bytes"
	"datwire/pkg/apps/customer"
	"datwire/pkg/bolt"
	dwhttp "datwire/pkg/http"
	"datwire/pkg/shared"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	sysbolt "github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type getCustomerResponseTemplate struct {
	Message string
	Error   string
	Data    *customer.Customer
}

type CustomerHandlerTestSuite struct {
	suite.Suite
	customerService *bolt.CustomerService
	custID_1        uuid.UUID
	custID_2        uuid.UUID
	custID_3        uuid.UUID
}

func (suite *CustomerHandlerTestSuite) SetupSuite() {
	suite.customerService = &bolt.CustomerService{}
	suite.customerService = &bolt.CustomerService{}
	suite.custID_1 = uuid.FromStringOrNil("53d21e77-3556-47f6-872d-46c513b9566e")
	suite.custID_2 = uuid.FromStringOrNil("370433ba-dfb2-47f6-a7d9-d0f3dcb045a2")
	suite.custID_3 = uuid.FromStringOrNil("5b72cb5b-f9b0-42f8-955c-32b93c8c5290")

	// seed data
	suite.customerService.Open()
	defer suite.customerService.Close()

	suite.customerService.CreateCustomer(&customer.Customer{
		ID:                       suite.custID_1,
		Name:                     "PT Bungasari Flourmills",
		Address:                  "Wisma 46 Kota BNI Lantai 28 Suite 2801, Jl. Jend. Sudirman Kav. 1, RT.10/RW.11, Karet Tengsin, Kota Jakarta Pusat, Daerah Khusus Ibukota Jakarta 10220",
		Telephone:                "+6221678278",
		ProcurementPIC:           "Grant Lutz",
		ProcurementContactNumber: "+62817682791989",
		OperationsPIC:            "Listiana Putri",
		OperationsContactNumber:  "+628719882791",
		Industry:                 "flour mill",
		Notes:                    "lorem ipsum dolor sit amet",
	})

	suite.customerService.CreateCustomer(&customer.Customer{
		ID:                       suite.custID_2,
		Name:                     "Mayora",
		Address:                  "Gedung Mayora, Jalan Tomang Raya 21-23, RT.1/RW.1, Tomang, Grogol petamburan, Kota Jakarta Barat, Daerah Khusus Ibukota Jakarta 11530",
		Telephone:                "+6221678278",
		ProcurementPIC:           "Henry Atmadja",
		ProcurementContactNumber: "+62817682791989",
		OperationsPIC:            "Trisna W.",
		OperationsContactNumber:  "+628719882791",
		Industry:                 "flour mill",
		Notes:                    "lorem ipsum dolor sit amet",
	})
}

func (suite *CustomerHandlerTestSuite) TearDownSuite() {
	db, err := sysbolt.Open("customer.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *sysbolt.Tx) error {
		err := tx.DeleteBucket([]byte("customers"))
		if err != nil {
			return err
		}
		return nil
	})
}

func (suite *CustomerHandlerTestSuite) TestCustomerHandler_CreateUser() {
	mockData := []byte(`{
		"id": "` + suite.custID_3.String() + `",
		"name":                     "Cargill Feed",
		"address":                  "Infina Park Blok B 73 No. 45, Jl. Dr. Saharjo, RT.1/RW.7, Manggarai, Tebet, Kota Jakarta Selatan, Daerah Khusus Ibukota Jakarta 12850",
		"telephone":                "+6221678278",
		"procurement_pic":           "DK Lee",
		"procurement_contact_number": "+62817682791989",
		"operations_pic":            "Kim Li",
		"operations_contact_number":  "+628719882791",
		"industry":                 "feed mill",
		"notes":                    "lorem ipsum dolor sit amet"
	}`)

	request, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(mockData))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	h := dwhttp.NewCustomerHandler()
	h.CustomerService.Open()
	defer h.CustomerService.Close()
	h.ServeHTTP(response, request)

	var responseBody *shared.ResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	suite.Equal("", responseBody.Error, "error should be empty")
	suite.Equal("success", responseBody.Message, "message should match")
}

func (suite *CustomerHandlerTestSuite) TestCustomerHandler_FetchUser() {
	request, _ := http.NewRequest("GET", "/customers/"+suite.custID_3.String(), nil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	h := dwhttp.NewCustomerHandler()
	h.CustomerService.Open()
	defer h.CustomerService.Close()
	h.ServeHTTP(response, request)

	var responseBody *getCustomerResponseTemplate
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	suite.Equal("", responseBody.Error, "error should be empty")
	suite.Equal("success", responseBody.Message, "message should match")

	suite.Equal(suite.custID_3.String(), responseBody.Data.ID.String(), "id should match")
	suite.Equal("Cargill Feed", responseBody.Data.Name, "name should match")
	suite.Equal(
		"Infina Park Blok B 73 No. 45, Jl. Dr. Saharjo, RT.1/RW.7, Manggarai, Tebet, Kota Jakarta Selatan, Daerah Khusus Ibukota Jakarta 12850",
		responseBody.Data.Address,
		"address should match",
	)
	suite.Equal("+6221678278", responseBody.Data.Telephone, "telephone should match")
	suite.Equal("DK Lee", responseBody.Data.ProcurementPIC, "procurement PIC name")
	suite.Equal("+62817682791989", responseBody.Data.ProcurementContactNumber, "procurement contact number should match")
	suite.Equal("Kim Li", responseBody.Data.OperationsPIC, "operations PIC should match")
	suite.Equal("+628719882791", responseBody.Data.OperationsContactNumber, "operations contact number should match")
	suite.Equal("feed mill", responseBody.Data.Industry, "industry should match")
	suite.Equal("lorem ipsum dolor sit amet", responseBody.Data.Notes, "notes should match")
}

func TestCustomerHandlerSuite(t *testing.T) {
	suite.Run(t, new(CustomerHandlerTestSuite))
}
