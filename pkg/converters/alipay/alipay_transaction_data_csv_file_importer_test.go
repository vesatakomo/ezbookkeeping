package alipay

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestAlipayCsvFileImporterParseImportedData_MinimumValidData(t *testing.T) {
	converter := AlipayTransactionDataCsvImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01 01:23:45 ,0.12   ,交易成功    ,已收入      ,\n" +
		"2024-09-01 12:34:56 ,123.45  ,交易成功    ,已支出      ,\n" +
		"2024-09-01 23:59:59 ,0.05   ,交易成功    ,资金转移     ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, allNewAccounts, allNewSubCategories, allNewTags, err := converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[0].Type)
	assert.Equal(t, "2024-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(12), allNewTransactions[0].Amount)
	assert.Equal(t, "Alipay", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[1].Type)
	assert.Equal(t, "2024-09-01 12:34:56", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime), time.UTC))
	assert.Equal(t, int64(12345), allNewTransactions[1].Amount)
	assert.Equal(t, "", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[2].Type)
	assert.Equal(t, "2024-09-01 23:59:59", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime), time.UTC))
	assert.Equal(t, int64(5), allNewTransactions[2].Amount)
	assert.Equal(t, "", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[2].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Alipay", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewSubCategories[0].Uid)
	assert.Equal(t, "", allNewSubCategories[0].Name)
}

func TestAlipayCsvFileImporterParseImportedData_ParseRefundTransaction(t *testing.T) {
	converter := AlipayTransactionDataCsvImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01 01:23:45 ,0.12   ,退款成功    ,已收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, "2024-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(-12), allNewTransactions[0].Amount)
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	data2, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01 01:23:45 ,0.12   ,退税成功    ,已收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, "2024-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(-12), allNewTransactions[0].Amount)
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)
}

func TestAlipayCsvFileImporterParseImportedData_ParseInvalidTime(t *testing.T) {
	converter := AlipayTransactionDataCsvImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01T12:34:56 ,0.12   ,交易成功    ,已收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	data2, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,交易状态    ,资金状态     ,\n" +
		"09/01/2024 12:34:56 ,0.12   ,交易成功    ,已收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)
}

func TestAlipayCsvFileImporterParseImportedData_ParseInvalidType(t *testing.T) {
	converter := AlipayTransactionDataCsvImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,0.12   ,交易成功    ,          ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestAlipayCsvFileImporterParseImportedData_ParseAccountName(t *testing.T) {
	converter := AlipayTransactionDataCsvImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	// income to alipay wallet
	data1, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,test                ,0.12   ,交易成功    ,已收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "Alipay", allNewTransactions[0].OriginalSourceAccountName)

	// income to other account
	data2, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,test                ,0.12   ,退款成功    ,已收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)

	// transfer to alipay wallet
	data3, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,商品名称                ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,test                ,充值-普通充值             ,0.12   ,交易成功    ,资金转移     ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, err = converter.ParseImportedData(context, user, []byte(data3), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "Alipay", allNewTransactions[0].OriginalDestinationAccountName)

	// transfer from alipay wallet
	data4, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,商品名称                ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,test                ,提现-实时提现             ,0.12   ,交易成功    ,资金转移     ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, err = converter.ParseImportedData(context, user, []byte(data4), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "Alipay", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "test", allNewTransactions[0].OriginalDestinationAccountName)

	// transfer in
	data5, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,商品名称                ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,test                ,xx-转入             ,0.12   ,交易成功    ,资金转移     ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, err = converter.ParseImportedData(context, user, []byte(data5), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "test", allNewTransactions[0].OriginalDestinationAccountName)

	// transfer out
	data6, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,商品名称                ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,test                ,xx-转出             ,0.12   ,交易成功    ,资金转移     ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, err = converter.ParseImportedData(context, user, []byte(data6), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "test", allNewTransactions[0].OriginalDestinationAccountName)

	// repayment
	data7, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,商品名称                ,金额（元）,交易状态    ,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,test                ,xx还款             ,0.12   ,交易成功    ,资金转移     ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, err = converter.ParseImportedData(context, user, []byte(data7), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "test", allNewTransactions[0].OriginalDestinationAccountName)
}

func TestAlipayCsvFileImporterParseImportedData_ParseDescription(t *testing.T) {
	converter := AlipayTransactionDataCsvImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,商品名称                ,金额（元）,交易状态    ,备注                  ,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,test                ,0.12   ,交易成功    ,test2               ,已收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "test2", allNewTransactions[0].Comment)

	data2, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,商品名称                ,金额（元）,交易状态    ,备注                  ,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,test                ,0.12   ,交易成功    ,                    ,已收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "test", allNewTransactions[0].Comment)
}

func TestAlipayCsvFileImporterParseImportedData_MissingFileHeader(t *testing.T) {
	converter := AlipayTransactionDataCsvImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	data, err := simplifiedchinese.GB18030.NewEncoder().String(
		"交易创建时间              ,金额（元）,交易状态    ,资金状态     ,\n" +
			"2024-09-01 12:34:56 ,0.12   ,交易成功    ,Type      ,\n" +
			"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte(""), 0, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)
}

func TestAlipayCsvFileImporterParseImportedData_MissingRequiredColumn(t *testing.T) {
	converter := AlipayTransactionDataCsvImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Time Column
	data1, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"金额（元）,交易状态    ,资金状态     ,\n" +
		"0.12   ,交易成功    ,已收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	_, _, _, _, err = converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Amount Column
	data2, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易状态    ,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,交易成功    ,已收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	_, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Status Column
	data3, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,资金状态     ,\n" +
		"2024-09-01 12:34:56 ,0.12   ,已收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	_, _, _, _, err = converter.ParseImportedData(context, user, []byte(data3), 0, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Fund Status Column
	data4, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,交易状态    ,\n" +
		"2024-09-01 12:34:56 ,0.12   ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	_, _, _, _, err = converter.ParseImportedData(context, user, []byte(data4), 0, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)
}