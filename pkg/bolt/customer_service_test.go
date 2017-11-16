package bolt_test

import (
	"datwire/pkg/apps/customer"
	"datwire/pkg/bolt"
	"log"
	"testing"

	sysbolt "github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type CustomerServiceTestSuite struct {
	suite.Suite
	customerService *bolt.CustomerService
	custID_1        uuid.UUID
	custID_2        uuid.UUID
	custID_3        uuid.UUID
}

func (suite *CustomerServiceTestSuite) SetupSuite() {
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

func (suite *CustomerServiceTestSuite) TearDownSuite() {
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

func (suite *CustomerServiceTestSuite) TestCustomerService_CreateCustomer() {
	suite.customerService.Open()
	defer suite.customerService.Close()
	err := suite.customerService.CreateCustomer(&customer.Customer{
		ID:                       suite.custID_3,
		Name:                     "Cargill Feed",
		Address:                  "Infina Park Blok B 73 No. 45, Jl. Dr. Saharjo, RT.1/RW.7, Manggarai, Tebet, Kota Jakarta Selatan, Daerah Khusus Ibukota Jakarta 12850",
		Telephone:                "+6221678278",
		ProcurementPIC:           "DK Lee",
		ProcurementContactNumber: "+62817682791989",
		OperationsPIC:            "Kim Li",
		OperationsContactNumber:  "+628719882791",
		Industry:                 "feed mill",
		Notes:                    "lorem ipsum dolor sit amet",
	})
	suite.Nil(err)
}

func (suite *CustomerServiceTestSuite) TestCustomerService_FetchCustomer() {
	suite.customerService.Open()
	defer suite.customerService.Close()

	c, err := suite.customerService.Customer(suite.custID_3)
	suite.Nil(err)
	suite.Equal(suite.custID_3.String(), c.ID.String(), "id should match")
	suite.Equal("Cargill Feed", c.Name, "name should match")
	suite.Equal(
		"Infina Park Blok B 73 No. 45, Jl. Dr. Saharjo, RT.1/RW.7, Manggarai, Tebet, Kota Jakarta Selatan, Daerah Khusus Ibukota Jakarta 12850",
		c.Address,
		"address should match",
	)
	suite.Equal("+6221678278", c.Telephone, "telephone should match")
	suite.Equal("DK Lee", c.ProcurementPIC, "procurement PIC name")
	suite.Equal("+62817682791989", c.ProcurementContactNumber, "procurement contact number should match")
	suite.Equal("Kim Li", c.OperationsPIC, "operations PIC should match")
	suite.Equal("+628719882791", c.OperationsContactNumber, "operations contact number should match")
	suite.Equal("feed mill", c.Industry, "industry should match")
	suite.Equal("lorem ipsum dolor sit amet", c.Notes, "notes should match")
}

func (suite *CustomerServiceTestSuite) TestCustomerService_FetchCustomers() {
	suite.customerService.Open()
	defer suite.customerService.Close()

	customers, err := suite.customerService.Customers()
	suite.Nil(err)
	suite.Equal(3, len(customers), "amount of customers should match")
}

func (suite *CustomerServiceTestSuite) TestCustomerService_UpdateCustomer() {
	suite.customerService.Open()
	defer suite.customerService.Close()

	err := suite.customerService.UpdateCustomer(suite.custID_3, "name", "ola")
	suite.Nil(err)
}

func (suite *CustomerServiceTestSuite) TestCustomerService_UpdateCustomer_VerifyUpdate() {
	suite.customerService.Open()
	defer suite.customerService.Close()

	c, err := suite.customerService.Customer(suite.custID_3)
	suite.Nil(err)
	suite.Equal("ola", c.Name, "name should be updated")
}

func (suite *CustomerServiceTestSuite) TestCustomerService_RemoveCustomer() {
	suite.customerService.Open()
	defer suite.customerService.Close()

	err := suite.customerService.DeleteCustomer(suite.custID_1)
	suite.Nil(err)
}

func (suite *CustomerServiceTestSuite) TestCustomerService_RemoveCustomer_VerifyRemoval() {
}

func TestCustomerServiceSuite(t *testing.T) {
	suite.Run(t, new(CustomerServiceTestSuite))
}
