package converters

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

// DataTableColumn represents the data column type of data table
type DataTableColumn byte

// Data table columns
const (
	DATA_TABLE_TRANSACTION_TIME         DataTableColumn = 1
	DATA_TABLE_TRANSACTION_TIMEZONE     DataTableColumn = 2
	DATA_TABLE_TRANSACTION_TYPE         DataTableColumn = 3
	DATA_TABLE_CATEGORY                 DataTableColumn = 4
	DATA_TABLE_SUB_CATEGORY             DataTableColumn = 5
	DATA_TABLE_ACCOUNT_NAME             DataTableColumn = 6
	DATA_TABLE_ACCOUNT_CURRENCY         DataTableColumn = 7
	DATA_TABLE_AMOUNT                   DataTableColumn = 8
	DATA_TABLE_RELATED_ACCOUNT_NAME     DataTableColumn = 9
	DATA_TABLE_RELATED_ACCOUNT_CURRENCY DataTableColumn = 10
	DATA_TABLE_RELATED_AMOUNT           DataTableColumn = 11
	DATA_TABLE_GEOGRAPHIC_LOCATION      DataTableColumn = 12
	DATA_TABLE_TAGS                     DataTableColumn = 13
	DATA_TABLE_DESCRIPTION              DataTableColumn = 14
)

// DataTableTransactionDataConverter defines the structure of data table importer for transaction data
type DataTableTransactionDataConverter struct {
	dataColumnMapping          map[DataTableColumn]string
	transactionTypeMapping     map[models.TransactionDbType]string
	transactionTypeNameMapping map[string]models.TransactionDbType
	columnSeparator            string
	lineSeparator              string
	geoLocationSeparator       string
	transactionTagSeparator    string
}

func (c *DataTableTransactionDataConverter) buildExportedContent(ctx core.Context, dataTableBuilder DataTableBuilder, uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) error {
	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			continue
		}

		dataRowMap := make(map[DataTableColumn]string, 15)
		transactionTimeZone := time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)

		dataRowMap[DATA_TABLE_TRANSACTION_TIME] = utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime), transactionTimeZone)
		dataRowMap[DATA_TABLE_TRANSACTION_TIMEZONE] = utils.FormatTimezoneOffset(transactionTimeZone)
		dataRowMap[DATA_TABLE_TRANSACTION_TYPE] = c.replaceDelimiters(c.getDisplayTransactionTypeName(transaction.Type))
		dataRowMap[DATA_TABLE_CATEGORY] = c.getExportedTransactionCategoryName(transaction.CategoryId, categoryMap)
		dataRowMap[DATA_TABLE_SUB_CATEGORY] = c.getExportedTransactionSubCategoryName(transaction.CategoryId, categoryMap)
		dataRowMap[DATA_TABLE_ACCOUNT_NAME] = c.getExportedAccountName(transaction.AccountId, accountMap)
		dataRowMap[DATA_TABLE_ACCOUNT_CURRENCY] = c.getAccountCurrency(transaction.AccountId, accountMap)
		dataRowMap[DATA_TABLE_AMOUNT] = utils.FormatAmount(transaction.Amount)

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			dataRowMap[DATA_TABLE_RELATED_ACCOUNT_NAME] = c.getExportedAccountName(transaction.RelatedAccountId, accountMap)
			dataRowMap[DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = c.getAccountCurrency(transaction.RelatedAccountId, accountMap)
			dataRowMap[DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(transaction.RelatedAccountAmount)
		}

		dataRowMap[DATA_TABLE_GEOGRAPHIC_LOCATION] = c.getExportedGeographicLocation(transaction)
		dataRowMap[DATA_TABLE_TAGS] = c.getExportedTags(transaction.TransactionId, allTagIndexes, tagMap)
		dataRowMap[DATA_TABLE_DESCRIPTION] = c.replaceDelimiters(transaction.Comment)

		dataTableBuilder.AppendTransaction(dataRowMap)
	}

	return nil
}

func (c *DataTableTransactionDataConverter) parseImportedData(ctx core.Context, user *models.User, dataTable ImportedDataTable, defaultTimezoneOffset int16, accountMap map[string]*models.Account, categoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionTag, error) {
	if dataTable.DataRowCount() < 1 {
		log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse import data for user \"uid:%d\", because data table row count is less 1", user.Uid)
		return nil, nil, nil, nil, errs.ErrOperationFailed
	}

	headerLineItems := dataTable.HeaderLineColumnNames()
	headerItemMap := make(map[string]int)

	for i := 0; i < len(headerLineItems); i++ {
		headerItemMap[headerLineItems[i]] = i
	}

	timeColumnIdx, timeColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_TRANSACTION_TIME]]
	timezoneColumnIdx, timezoneColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_TRANSACTION_TIMEZONE]]
	typeColumnIdx, typeColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_TRANSACTION_TYPE]]
	subCategoryColumnIdx, subCategoryColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_SUB_CATEGORY]]
	accountColumnIdx, accountColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_ACCOUNT_NAME]]
	accountCurrencyColumnIdx, accountCurrencyColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_ACCOUNT_CURRENCY]]
	amountColumnIdx, amountColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_AMOUNT]]
	account2ColumnIdx, account2ColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_RELATED_ACCOUNT_NAME]]
	account2CurrencyColumnIdx, account2CurrencyColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_RELATED_ACCOUNT_CURRENCY]]
	amount2ColumnIdx, amount2ColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_RELATED_AMOUNT]]
	geoLocationIdx, geoLocationExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_GEOGRAPHIC_LOCATION]]
	tagsColumnIdx, tagsColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_TAGS]]
	descriptionColumnIdx, descriptionColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_DESCRIPTION]]

	if !timeColumnExists || !typeColumnExists || !subCategoryColumnExists ||
		!accountColumnExists || !amountColumnExists || !account2ColumnExists || !amount2ColumnExists {
		log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse import data for user \"uid:%d\", because missing essential columns in header row", user.Uid)
		return nil, nil, nil, nil, errs.ErrFormatInvalid
	}

	if accountMap == nil {
		accountMap = make(map[string]*models.Account)
	}

	if categoryMap == nil {
		categoryMap = make(map[string]*models.TransactionCategory)
	}

	if tagMap == nil {
		tagMap = make(map[string]*models.TransactionTag)
	}

	allNewTransactions := make(models.ImportedTransactionSlice, 0, dataTable.DataRowCount())
	allNewAccounts := make([]*models.Account, 0)
	allNewSubCategories := make([]*models.TransactionCategory, 0)
	allNewTags := make([]*models.TransactionTag, 0)

	dataRowIterator := dataTable.DataRowIterator()
	dataRowIndex := 0

	for dataRowIterator.HasNext() {
		dataRowIndex++
		dataRow := dataRowIterator.Next()
		columnCount := dataRow.ColumnCount()

		if columnCount < 1 || (columnCount == 1 && dataRow.GetData(0) == "") {
			continue
		}

		if columnCount < len(headerLineItems) {
			log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse data row \"index:%d\" for user \"uid:%d\", because may missing some columns (column count %d in data row is less than header column count %d)", dataRowIndex, user.Uid, columnCount, len(headerLineItems))
			return nil, nil, nil, nil, errs.ErrFormatInvalid
		}

		timezoneOffset := defaultTimezoneOffset

		if timezoneColumnExists {
			transactionTimezone, err := utils.ParseFromTimezoneOffset(dataRow.GetData(timezoneColumnIdx))

			if err != nil {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse time zone \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(timezoneColumnIdx), dataRowIndex, user.Uid, err.Error())
				return nil, nil, nil, nil, err
			}

			timezoneOffset = utils.GetTimezoneOffsetMinutes(transactionTimezone)
		}

		transactionTime, err := utils.ParseFromLongDateTime(dataRow.GetData(timeColumnIdx), timezoneOffset)

		if err != nil {
			log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse time \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(timeColumnIdx), dataRowIndex, user.Uid, err.Error())
			return nil, nil, nil, nil, err
		}

		transactionDbType, err := c.getTransactionDbType(dataRow.GetData(typeColumnIdx))

		if err != nil {
			log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse transaction type \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(typeColumnIdx), dataRowIndex, user.Uid, err.Error())
			return nil, nil, nil, nil, err
		}

		categoryId := int64(0)
		subCategoryName := ""

		if transactionDbType != models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			transactionCategoryType, err := c.getTransactionCategoryType(transactionDbType)

			if err != nil {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse transaction category type in data row \"index:%d\" for user \"uid:%d\", because %s", dataRowIndex, user.Uid, err.Error())
				return nil, nil, nil, nil, err
			}

			subCategoryName = dataRow.GetData(subCategoryColumnIdx)

			if subCategoryName == "" {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] sub category type is empty in data row \"index:%d\" for user \"uid:%d\"", dataRowIndex, user.Uid)
				return nil, nil, nil, nil, errs.ErrFormatInvalid
			}

			subCategory, exists := categoryMap[subCategoryName]

			if !exists {
				subCategory = c.createNewTransactionCategoryModel(user.Uid, subCategoryName, transactionCategoryType)
				allNewSubCategories = append(allNewSubCategories, subCategory)
				categoryMap[subCategoryName] = subCategory
			}

			categoryId = subCategory.CategoryId
		}

		accountName := dataRow.GetData(accountColumnIdx)

		if accountName == "" {
			log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] account name is empty in data row \"index:%d\" for user \"uid:%d\"", dataRowIndex, user.Uid)
			return nil, nil, nil, nil, errs.ErrFormatInvalid
		}

		accountCurrency := user.DefaultCurrency

		if accountCurrencyColumnExists {
			accountCurrency = dataRow.GetData(accountCurrencyColumnIdx)

			if _, ok := validators.AllCurrencyNames[accountCurrency]; !ok {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] account currency \"%s\" is not supported in data row \"index:%d\" for user \"uid:%d\"", accountCurrency, dataRowIndex, user.Uid)
				return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
			}
		}

		account, exists := accountMap[accountName]

		if !exists {
			account = c.createNewAccountModel(user.Uid, accountName, accountCurrency)
			allNewAccounts = append(allNewAccounts, account)
			accountMap[accountName] = account
		}

		if accountCurrencyColumnExists {
			if account.Currency != accountCurrency {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] currency \"%s\" in data row \"index:%d\" not equals currency \"%s\" of the account for user \"uid:%d\"", accountCurrency, dataRowIndex, account.Currency, user.Uid)
				return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
			}
		} else if exists {
			accountCurrency = account.Currency
		}

		amount, err := utils.ParseAmount(dataRow.GetData(amountColumnIdx))

		if err != nil {
			log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse acmount \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(amountColumnIdx), dataRowIndex, user.Uid, err.Error())
			return nil, nil, nil, nil, err
		}

		relatedAccountId := int64(0)
		relatedAccountAmount := int64(0)
		account2Name := ""
		account2Currency := ""

		if transactionDbType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			account2Name = dataRow.GetData(account2ColumnIdx)

			if account2Name == "" {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] account2 name is empty in data row \"index:%d\" for user \"uid:%d\"", dataRowIndex, user.Uid)
				return nil, nil, nil, nil, errs.ErrFormatInvalid
			}

			account2Currency = user.DefaultCurrency

			if account2CurrencyColumnExists {
				account2Currency = dataRow.GetData(account2CurrencyColumnIdx)

				if _, ok := validators.AllCurrencyNames[account2Currency]; !ok {
					log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] account2 currency \"%s\" is not supported in data row \"index:%d\" for user \"uid:%d\"", account2Currency, dataRowIndex, user.Uid)
					return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
				}
			}

			account2, exists := accountMap[account2Name]

			if !exists {
				account2 = c.createNewAccountModel(user.Uid, account2Name, account2Currency)
				allNewAccounts = append(allNewAccounts, account2)
				accountMap[account2Name] = account2
			}

			if account2CurrencyColumnExists {
				if account2.Currency != account2Currency {
					log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] currency \"%s\" in data row \"index:%d\" not equals currency \"%s\" of the account2 for user \"uid:%d\"", account2Currency, dataRowIndex, account2.Currency, user.Uid)
					return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
				}
			} else if exists {
				account2Currency = account2.Currency
			}

			relatedAccountId = account2.AccountId
			relatedAccountAmount, err = utils.ParseAmount(dataRow.GetData(amount2ColumnIdx))

			if err != nil {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse acmount2 \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(amount2ColumnIdx), dataRowIndex, user.Uid, err.Error())
				return nil, nil, nil, nil, err
			}
		}

		geoLongitude := float64(0)
		geoLatitude := float64(0)

		if geoLocationExists {
			geoLocationItems := strings.Split(dataRow.GetData(geoLocationIdx), c.geoLocationSeparator)

			if len(geoLocationItems) == 2 {
				geoLongitude, err = utils.StringToFloat64(geoLocationItems[0])

				if err != nil {
					log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse geographic location \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(geoLocationIdx), dataRowIndex, user.Uid, err.Error())
					return nil, nil, nil, nil, err
				}

				geoLatitude, err = utils.StringToFloat64(geoLocationItems[1])

				if err != nil {
					log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse geographic location \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(geoLocationIdx), dataRowIndex, user.Uid, err.Error())
					return nil, nil, nil, nil, err
				}
			}
		}

		var tagIds []string
		var tagNames []string

		if tagsColumnExists {
			tagNameItems := strings.Split(dataRow.GetData(tagsColumnIdx), c.transactionTagSeparator)

			for i := 0; i < len(tagNameItems); i++ {
				tagName := tagNameItems[i]

				if tagName == "" {
					continue
				}

				tag, exists := tagMap[tagName]

				if !exists {
					tag = c.createNewTransactionTagModel(user.Uid, tagName)
					allNewTags = append(allNewTags, tag)
					tagMap[tagName] = tag
				}

				if tag != nil {
					tagIds = append(tagIds, utils.Int64ToString(tag.TagId))
				}

				tagNames = append(tagNames, tagName)
			}
		}

		description := ""

		if descriptionColumnExists {
			description = dataRow.GetData(descriptionColumnIdx)
		}

		transaction := &models.ImportTransaction{
			Transaction: &models.Transaction{
				Uid:                  user.Uid,
				Type:                 transactionDbType,
				CategoryId:           categoryId,
				TransactionTime:      utils.GetMinTransactionTimeFromUnixTime(transactionTime.Unix()),
				TimezoneUtcOffset:    timezoneOffset,
				AccountId:            account.AccountId,
				Amount:               amount,
				HideAmount:           false,
				RelatedAccountId:     relatedAccountId,
				RelatedAccountAmount: relatedAccountAmount,
				Comment:              description,
				GeoLongitude:         geoLongitude,
				GeoLatitude:          geoLatitude,
				CreatedIp:            "127.0.0.1",
			},
			TagIds:                             tagIds,
			OriginalCategoryName:               subCategoryName,
			OriginalSourceAccountName:          accountName,
			OriginalSourceAccountCurrency:      accountCurrency,
			OriginalDestinationAccountName:     account2Name,
			OriginalDestinationAccountCurrency: account2Currency,
			OriginalTagNames:                   tagNames,
		}

		allNewTransactions = append(allNewTransactions, transaction)
	}

	sort.Sort(allNewTransactions)

	return allNewTransactions, allNewAccounts, allNewSubCategories, allNewTags, nil
}

func (c *DataTableTransactionDataConverter) getTransactionDbType(transactionTypeName string) (models.TransactionDbType, error) {
	transactionType, exists := c.transactionTypeNameMapping[transactionTypeName]

	if !exists {
		return 0, errs.ErrTransactionTypeInvalid
	}

	return transactionType, nil
}

func (c *DataTableTransactionDataConverter) getTransactionCategoryType(transactionType models.TransactionDbType) (models.TransactionCategoryType, error) {
	if transactionType == models.TRANSACTION_DB_TYPE_INCOME {
		return models.CATEGORY_TYPE_INCOME, nil
	} else if transactionType == models.TRANSACTION_DB_TYPE_EXPENSE {
		return models.CATEGORY_TYPE_EXPENSE, nil
	} else if transactionType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		return models.CATEGORY_TYPE_TRANSFER, nil
	} else {
		return 0, errs.ErrTransactionTypeInvalid
	}
}

func (c *DataTableTransactionDataConverter) getDisplayTransactionTypeName(transactionDbType models.TransactionDbType) string {
	transactionTypeName, exists := c.transactionTypeMapping[transactionDbType]

	if !exists {
		return ""
	}

	return transactionTypeName
}

func (c *DataTableTransactionDataConverter) getExportedTransactionCategoryName(categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if !exists {
		return ""
	}

	if category.ParentCategoryId == 0 {
		return c.replaceDelimiters(category.Name)
	}

	parentCategory, exists := categoryMap[category.ParentCategoryId]

	if !exists {
		return ""
	}

	return c.replaceDelimiters(parentCategory.Name)
}

func (c *DataTableTransactionDataConverter) getExportedTransactionSubCategoryName(categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if exists {
		return c.replaceDelimiters(category.Name)
	} else {
		return ""
	}
}

func (c *DataTableTransactionDataConverter) getExportedAccountName(accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return c.replaceDelimiters(account.Name)
	} else {
		return ""
	}
}

func (c *DataTableTransactionDataConverter) getAccountCurrency(accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return c.replaceDelimiters(account.Currency)
	} else {
		return ""
	}
}

func (c *DataTableTransactionDataConverter) getExportedGeographicLocation(transaction *models.Transaction) string {
	if transaction.GeoLongitude != 0 || transaction.GeoLatitude != 0 {
		return fmt.Sprintf("%f%s%f", transaction.GeoLongitude, c.geoLocationSeparator, transaction.GeoLatitude)
	}

	return ""
}

func (c *DataTableTransactionDataConverter) getExportedTags(transactionId int64, allTagIndexes map[int64][]int64, tagMap map[int64]*models.TransactionTag) string {
	tagIndexes, exists := allTagIndexes[transactionId]

	if !exists {
		return ""
	}

	var ret strings.Builder

	for i := 0; i < len(tagIndexes); i++ {
		tagIndex := tagIndexes[i]
		tag, exists := tagMap[tagIndex]

		if !exists {
			continue
		}

		if ret.Len() > 0 {
			ret.WriteString(c.transactionTagSeparator)
		}

		ret.WriteString(strings.Replace(tag.Name, c.transactionTagSeparator, " ", -1))
	}

	return c.replaceDelimiters(ret.String())
}

func (c *DataTableTransactionDataConverter) replaceDelimiters(text string) string {
	text = strings.Replace(text, "\r\n", " ", -1)
	text = strings.Replace(text, "\r", " ", -1)
	text = strings.Replace(text, "\n", " ", -1)
	text = strings.Replace(text, c.columnSeparator, " ", -1)
	text = strings.Replace(text, c.lineSeparator, " ", -1)

	return text
}

func (c *DataTableTransactionDataConverter) createNewAccountModel(uid int64, accountName string, currency string) *models.Account {
	return &models.Account{
		Uid:      uid,
		Name:     accountName,
		Currency: currency,
	}
}

func (c *DataTableTransactionDataConverter) createNewTransactionCategoryModel(uid int64, categoryName string, transactionCategoryType models.TransactionCategoryType) *models.TransactionCategory {
	return &models.TransactionCategory{
		Uid:  uid,
		Name: categoryName,
		Type: transactionCategoryType,
	}
}

func (c *DataTableTransactionDataConverter) createNewTransactionTagModel(uid int64, tagName string) *models.TransactionTag {
	return &models.TransactionTag{
		Uid:  uid,
		Name: tagName,
	}
}