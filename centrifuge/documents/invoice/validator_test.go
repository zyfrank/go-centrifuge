// +build unit

package invoice

import (
	"fmt"
	"testing"

	"github.com/centrifuge/go-centrifuge/centrifuge/coredocument"
	"github.com/centrifuge/go-centrifuge/centrifuge/documents"
	"github.com/centrifuge/go-centrifuge/centrifuge/tools"
	"github.com/stretchr/testify/assert"
)

func TestFieldValidator_Validate(t *testing.T) {
	fv := fieldValidator()

	//  nil error
	err := fv.Validate(nil, nil)
	assert.Error(t, err)
	errs := documents.Errors(err)
	assert.Len(t, errs, 1, "errors length must be one")
	assert.Contains(t, errs[0].Error(), "nil document")

	// unknown type
	err = fv.Validate(nil, &mockModel{})
	assert.Error(t, err)
	errs = documents.Errors(err)
	assert.Len(t, errs, 1, "errors length must be one")
	assert.Contains(t, errs[0].Error(), "unknown document type")

	// fail
	err = fv.Validate(nil, new(InvoiceModel))
	assert.Error(t, err)
	errs = documents.Errors(err)
	assert.Len(t, errs, 1, "errors length must be 2")
	assert.Contains(t, errs[0].Error(), "currency is invalid")

	// success
	err = fv.Validate(nil, &InvoiceModel{
		Currency: "EUR",
	})
	assert.Nil(t, err)
}

func TestDataRootValidation_Validate(t *testing.T) {
	drv := dataRootValidator()

	// nil error
	err := drv.Validate(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil document")

	// pack coredoc failed
	model := &mockModel{}
	model.On("PackCoreDocument").Return(nil, fmt.Errorf("error")).Once()
	err = drv.Validate(nil, model)
	model.AssertExpectations(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to pack coredocument")

	// missing data root
	model = &mockModel{}
	model.On("PackCoreDocument").Return(coredocument.New(), nil).Once()
	err = drv.Validate(nil, model)
	model.AssertExpectations(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data root missing")

	// unknown doc type
	cd := coredocument.New()
	cd.DataRoot = tools.RandomSlice(32)
	model = &mockModel{}
	model.On("PackCoreDocument").Return(cd, nil).Once()
	err = drv.Validate(nil, model)
	model.AssertExpectations(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown document type")

	// mismatch
	inv := new(InvoiceModel)
	err = inv.InitInvoiceInput(createPayload())
	assert.Nil(t, err)
	inv.CoreDocument = cd
	err = drv.Validate(nil, inv)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mismatched data root")

	// success
	inv = new(InvoiceModel)
	err = inv.InitInvoiceInput(createPayload())
	assert.Nil(t, err)
	err = inv.calculateDataRoot()
	assert.Nil(t, err)
	err = drv.Validate(nil, inv)
	assert.Nil(t, err)
}

func TestCreateValidator(t *testing.T) {
	cv := CreateValidator()
	assert.Len(t, cv, 2)
}