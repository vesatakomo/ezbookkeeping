import { defineStore } from 'pinia';

import { useUserStore } from './user.js';
import { useAccountsStore } from './account.js';
import { useTransactionCategoriesStore } from './transactionCategory.js';
import { useExchangeRatesStore } from './exchangeRates.js';

import statisticsConstants from '@/consts/statistics.js';
import categoryConstants from '@/consts/category.js';
import iconConstants from '@/consts/icon.js';
import colorConstants from '@/consts/color.js';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';
import { isNumber, isObject } from '@/lib/common.js';

function loadTransactionStatistics(state, { statistics, defaultCurrency }) {
    if (statistics && statistics.items && statistics.items.length) {
        const accountsStore = useAccountsStore();
        const transactionCategoriesStore = useTransactionCategoriesStore();
        const exchangeRatesStore = useExchangeRatesStore();

        for (let i = 0; i < statistics.items.length; i++) {
            const item = statistics.items[i];

            if (item.accountId) {
                item.account = accountsStore.allAccountsMap[item.accountId];
            }

            if (item.account && item.account.parentId !== '0') {
                item.primaryAccount = accountsStore.allAccountsMap[item.account.parentId];
            } else {
                item.primaryAccount = item.account;
            }

            if (item.categoryId) {
                item.category = transactionCategoriesStore.allTransactionCategoriesMap[item.categoryId];
            }

            if (item.category && item.category.parentId !== '0') {
                item.primaryCategory = transactionCategoriesStore.allTransactionCategoriesMap[item.category.parentId];
            } else {
                item.primaryCategory = item.category;
            }

            if (item.account && item.account.currency !== defaultCurrency) {
                const amount = exchangeRatesStore.getExchangedAmount(item.amount, item.account.currency, defaultCurrency);

                if (isNumber(amount)) {
                    item.amountInDefaultCurrency = Math.floor(amount);
                }
            } else if (item.account && item.account.currency === defaultCurrency) {
                item.amountInDefaultCurrency = item.amount;
            } else {
                item.amountInDefaultCurrency = null;
            }
        }
    }

    state.transactionStatistics = statistics;
}

export const useStatisticsStore = defineStore('statistics', {
    state: () => ({
        transactionStatisticsFilter: {
            dateType: statisticsConstants.defaultDataRangeType,
            startTime: 0,
            endTime: 0,
            chartType: statisticsConstants.defaultChartType,
            chartDataType: statisticsConstants.defaultChartDataType,
            filterAccountIds: {},
            filterCategoryIds: {}
        },
        transactionStatistics: [],
        transactionStatisticsStateInvalid: true
    }),
    getters: {
        statisticsItemsByTransactionStatisticsData(state) {
            if (!state.transactionStatistics || !state.transactionStatistics.items) {
                return null;
            }

            const allDataItems = {};
            let totalAmount = 0;
            let totalNonNegativeAmount = 0;

            for (let i = 0; i < state.transactionStatistics.items.length; i++) {
                const item = state.transactionStatistics.items[i];

                if (!item.primaryAccount || !item.account || !item.primaryCategory || !item.category) {
                    continue;
                }

                if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type ||
                    state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type ||
                    state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type) {
                    if (item.category.type !== categoryConstants.allCategoryTypes.Expense) {
                        continue;
                    }
                } else if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type ||
                    state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type ||
                    state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type) {
                    if (item.category.type !== categoryConstants.allCategoryTypes.Income) {
                        continue;
                    }
                } else {
                    continue;
                }

                if (state.transactionStatisticsFilter.filterAccountIds && state.transactionStatisticsFilter.filterAccountIds[item.account.id]) {
                    continue;
                }

                if (state.transactionStatisticsFilter.filterCategoryIds && state.transactionStatisticsFilter.filterCategoryIds[item.category.id]) {
                    continue;
                }

                if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type ||
                    state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type) {
                    if (isNumber(item.amountInDefaultCurrency)) {
                        let data = allDataItems[item.account.id];

                        if (data) {
                            data.totalAmount += item.amountInDefaultCurrency;
                        } else {
                            data = {
                                name: item.account.name,
                                type: 'account',
                                id: item.account.id,
                                icon: item.account.icon || iconConstants.defaultAccountIcon.icon,
                                color: item.account.color || colorConstants.defaultAccountColor,
                                hidden: item.primaryAccount.hidden || item.account.hidden,
                                displayOrders: [item.primaryAccount.category, item.primaryAccount.displayOrder, item.account.displayOrder],
                                totalAmount: item.amountInDefaultCurrency
                            }
                        }

                        totalAmount += item.amountInDefaultCurrency;

                        if (item.amountInDefaultCurrency > 0) {
                            totalNonNegativeAmount += item.amountInDefaultCurrency;
                        }

                        allDataItems[item.account.id] = data;
                    }
                } else if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type ||
                    state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type) {
                    if (isNumber(item.amountInDefaultCurrency)) {
                        let data = allDataItems[item.primaryCategory.id];

                        if (data) {
                            data.totalAmount += item.amountInDefaultCurrency;
                        } else {
                            data = {
                                name: item.primaryCategory.name,
                                type: 'category',
                                id: item.primaryCategory.id,
                                icon: item.primaryCategory.icon || iconConstants.defaultCategoryIcon.icon,
                                color: item.primaryCategory.color || colorConstants.defaultCategoryColor,
                                hidden: item.primaryCategory.hidden,
                                displayOrders: [item.primaryCategory.type, item.primaryCategory.displayOrder],
                                totalAmount: item.amountInDefaultCurrency
                            }
                        }

                        totalAmount += item.amountInDefaultCurrency;

                        if (item.amountInDefaultCurrency > 0) {
                            totalNonNegativeAmount += item.amountInDefaultCurrency;
                        }

                        allDataItems[item.primaryCategory.id] = data;
                    }
                } else if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type ||
                    state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type) {
                    if (isNumber(item.amountInDefaultCurrency)) {
                        let data = allDataItems[item.category.id];

                        if (data) {
                            data.totalAmount += item.amountInDefaultCurrency;
                        } else {
                            data = {
                                name: item.category.name,
                                type: 'category',
                                id: item.category.id,
                                icon: item.category.icon || iconConstants.defaultCategoryIcon.icon,
                                color: item.category.color || colorConstants.defaultCategoryColor,
                                hidden: item.primaryCategory.hidden || item.category.hidden,
                                displayOrders: [item.primaryCategory.type, item.primaryCategory.displayOrder, item.category.displayOrder],
                                totalAmount: item.amountInDefaultCurrency
                            }
                        }

                        totalAmount += item.amountInDefaultCurrency;

                        if (item.amountInDefaultCurrency > 0) {
                            totalNonNegativeAmount += item.amountInDefaultCurrency;
                        }

                        allDataItems[item.category.id] = data;
                    }
                }
            }

            return {
                totalAmount: totalAmount,
                totalNonNegativeAmount: totalNonNegativeAmount,
                items: allDataItems
            }
        },
        statisticsItemsByAccountsData(state) {
            const userStore = useUserStore();
            const accountsStore = useAccountsStore();
            const exchangeRatesStore = useExchangeRatesStore();

            if (!accountsStore.allPlainAccounts) {
                return null;
            }

            const allDataItems = {};
            let totalAmount = 0;
            let totalNonNegativeAmount = 0;

            for (let i = 0; i < accountsStore.allPlainAccounts.length; i++) {
                const account = accountsStore.allPlainAccounts[i];

                if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalAssets.type) {
                    if (!account.isAsset) {
                        continue;
                    }
                } else if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalLiabilities.type) {
                    if (!account.isLiability) {
                        continue;
                    }
                }

                if (state.transactionStatisticsFilter.filterAccountIds && state.transactionStatisticsFilter.filterAccountIds[account.id]) {
                    continue;
                }

                let primaryAccount = accountsStore.allAccountsMap[account.parentId];

                if (!primaryAccount) {
                    primaryAccount = account;
                }

                let amount = account.balance;

                if (account.currency !== userStore.currentUserDefaultCurrency) {
                    amount = Math.floor(exchangeRatesStore.getExchangedAmount(amount, account.currency, userStore.currentUserDefaultCurrency));

                    if (!isNumber(amount)) {
                        continue;
                    }
                }

                if (account.isLiability) {
                    amount = -amount;
                }

                const data = {
                    name: account.name,
                    type: 'account',
                    id: account.id,
                    icon: account.icon || iconConstants.defaultAccountIcon.icon,
                    color: account.color || colorConstants.defaultAccountColor,
                    hidden: primaryAccount.hidden || account.hidden,
                    displayOrders: [primaryAccount.category, primaryAccount.displayOrder, account.displayOrder],
                    totalAmount: amount
                };

                totalAmount += amount;

                if (amount > 0) {
                    totalNonNegativeAmount += amount;
                }

                allDataItems[account.id] = data;
            }

            return {
                totalAmount: totalAmount,
                totalNonNegativeAmount: totalNonNegativeAmount,
                items: allDataItems
            }
        }
    },
    actions: {
        updateTransactionStatisticsInvalidState(invalidState) {
            this.transactionStatisticsStateInvalid = invalidState;
        },
        resetTransactionStatistics() {
            this.transactionStatisticsFilter.dateType = statisticsConstants.defaultDataRangeType;
            this.transactionStatisticsFilter.startTime = 0;
            this.transactionStatisticsFilter.endTime = 0;
            this.transactionStatisticsFilter.chartType = statisticsConstants.defaultChartType;
            this.transactionStatisticsFilter.chartDataType = statisticsConstants.defaultChartDataType;
            this.transactionStatisticsFilter.filterAccountIds = {};
            this.transactionStatisticsFilter.filterCategoryIds = {};
            this.transactionStatistics = {};
            this.transactionStatisticsStateInvalid = true;
        },
        initTransactionStatisticsFilter(filter) {
            if (filter && isNumber(filter.dateType)) {
                this.transactionStatisticsFilter.dateType = filter.dateType;
            } else {
                this.transactionStatisticsFilter.dateType = statisticsConstants.defaultDataRangeType;
            }

            if (filter && isNumber(filter.startTime)) {
                this.transactionStatisticsFilter.startTime = filter.startTime;
            } else {
                this.transactionStatisticsFilter.startTime = 0;
            }

            if (filter && isNumber(filter.endTime)) {
                this.transactionStatisticsFilter.endTime = filter.endTime;
            } else {
                this.transactionStatisticsFilter.endTime = 0;
            }

            if (filter && isNumber(filter.chartType)) {
                this.transactionStatisticsFilter.chartType = filter.chartType;
            } else {
                this.transactionStatisticsFilter.chartType = statisticsConstants.defaultChartType;
            }

            if (filter && isNumber(filter.chartDataType)) {
                this.transactionStatisticsFilter.chartDataType = filter.chartDataType;
            } else {
                this.transactionStatisticsFilter.chartDataType = statisticsConstants.defaultChartDataType;
            }

            if (filter && isObject(filter.filterAccountIds)) {
                this.transactionStatisticsFilter.filterAccountIds = filter.filterAccountIds;
            } else {
                this.transactionStatisticsFilter.filterAccountIds = {};
            }

            if (filter && isObject(filter.filterCategoryIds)) {
                this.transactionStatisticsFilter.filterCategoryIds = filter.filterCategoryIds;
            } else {
                this.transactionStatisticsFilter.filterCategoryIds = {};
            }

            if (filter && isNumber(filter.sortingType)) {
                this.transactionStatisticsFilter.sortingType = filter.sortingType;
            } else {
                this.transactionStatisticsFilter.sortingType = statisticsConstants.defaultSortingType;
            }
        },
        updateTransactionStatisticsFilter(filter) {
            if (filter && isNumber(filter.dateType)) {
                this.transactionStatisticsFilter.dateType = filter.dateType;
            }

            if (filter && isNumber(filter.startTime)) {
                this.transactionStatisticsFilter.startTime = filter.startTime;
            }

            if (filter && isNumber(filter.endTime)) {
                this.transactionStatisticsFilter.endTime = filter.endTime;
            }

            if (filter && isNumber(filter.chartType)) {
                this.transactionStatisticsFilter.chartType = filter.chartType;
            }

            if (filter && isNumber(filter.chartDataType)) {
                this.transactionStatisticsFilter.chartDataType = filter.chartDataType;
            }

            if (filter && isObject(filter.filterAccountIds)) {
                this.transactionStatisticsFilter.filterAccountIds = filter.filterAccountIds;
            }

            if (filter && isObject(filter.filterCategoryIds)) {
                this.transactionStatisticsFilter.filterCategoryIds = filter.filterCategoryIds;
            }

            if (filter && isNumber(filter.sortingType)) {
                this.transactionStatisticsFilter.sortingType = filter.sortingType;
            }
        },
        loadTransactionStatistics({ defaultCurrency }) {
            const self = this;

            return new Promise((resolve, reject) => {
                services.getTransactionStatistics({
                    startTime: self.transactionStatisticsFilter.startTime,
                    endTime: self.transactionStatisticsFilter.endTime
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to get transaction statistics' });
                        return;
                    }

                    loadTransactionStatistics(self,  {
                        statistics: data.result,
                        defaultCurrency: defaultCurrency
                    });

                    if (self.transactionStatisticsStateInvalid) {
                        self.updateTransactionStatisticsInvalidState(false);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to get transaction statistics', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to get transaction statistics' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
    }
});