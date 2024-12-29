export function isFunction(val: unknown): boolean {
    return typeof(val) === 'function';
}

export function isDefined(val: unknown): boolean {
    return typeof val !== 'undefined';
}

export function isObject(val: unknown): boolean {
    return val != null && typeof(val) === 'object' && !isArray(val);
}

export function isArray(val: unknown): boolean {
    if (isFunction(Array.isArray)) {
        return Array.isArray(val);
    }

    return Object.prototype.toString.call(val) === '[object Array]';
}

export function isString(val: unknown): boolean {
    return typeof(val) === 'string';
}

export function isNumber(val: unknown): boolean {
    return typeof(val) === 'number';
}

export function isInteger(val: unknown): boolean {
    return Number.isInteger(val);
}

export function isBoolean(val: unknown): boolean {
    return typeof(val) === 'boolean';
}

export function isYearMonth(val: unknown): boolean {
    if (typeof(val) !== 'string') {
        return false;
    }

    const items = val.split('-');

    if (items.length !== 2) {
        return false;
    }

    return !!parseInt(items[0]) && !!parseInt(items[1]);
}

export function isEquals(obj1: unknown, obj2: unknown): boolean {
    if (obj1 === obj2) {
        return true;
    }

    if (isArray(obj1) && isArray(obj2)) {
        const arr1 = obj1 as unknown[];
        const arr2 = obj2 as unknown[];

        if (arr1.length !== arr2.length) {
            return false;
        }

        for (let i = 0; i < arr1.length; i++) {
            if (!isEquals(arr1[i], arr2[i])) {
                return false;
            }
        }

        return true;
    } else if (isObject(obj1) && isObject(obj2)) {
        const keys1 = Object.keys(obj1 as Record<string, unknown>);
        const keys2 = Object.keys(obj2 as Record<string, unknown>);

        if (keys1.length !== keys2.length) {
            return false;
        }

        const keyExistsMap2 = new Map<string, boolean>();

        for (let i = 0; i < keys2.length; i++) {
            const key = keys2[i];

            keyExistsMap2.set(key, true);
        }

        for (let i = 0; i < keys1.length; i++) {
            const key = keys1[i];

            if (!keyExistsMap2.get(key)) {
                return false;
            }

            if (!isEquals((obj1 as Record<string, unknown>)[key], (obj2 as Record<string, unknown>)[key])) {
                return false;
            }
        }

        return true;
    } else {
        return obj1 === obj2;
    }
}

export function isYearMonthEquals(val1: unknown, val2: unknown): boolean {
    if (typeof(val1) !== 'string' || typeof(val2) !== 'string') {
        return false;
    }

    const items1 = val1.split('-');
    const items2 = val2.split('-');

    if (items1.length !== 2 || items2.length !== 2) {
        return false;
    }

    return (!!parseInt(items1[0]) && !!parseInt(items1[1])) && (parseInt(items1[0]) === parseInt(items2[0])) && (parseInt(items1[1]) === parseInt(items2[1]));
}

export function isObjectEmpty(obj: object): boolean {
    if (!obj) {
        return true;
    }

    for (const field in obj) {
        if (!Object.prototype.hasOwnProperty.call(obj, field)) {
            continue;
        }

        return false;
    }

    return true;
}

export function sortNumbersArray(array: number[]): number[] {
    return array.sort(function (num1, num2) {
        return num1 - num2;
    });
}

export function getObjectOwnFieldCount(object: object): number {
    let count = 0;

    if (!object || !isObject(object)) {
        return count;
    }

    for (const field in object) {
        if (!Object.prototype.hasOwnProperty.call(object, field)) {
            continue;
        }

        count++;
    }

    return count;
}

export function replaceAll(value: string, originalValue: string, targetValue: string): string {
    // Escape special characters in originalValue to safely use it in a regex pattern.
    // This ensures that characters like . (dot), * (asterisk), +, ?, etc. are treated literally,
    // rather than as special regex symbols.
    const escapedOriginalValue = originalValue.replace(/([.*+?^=!:${}()|\-/\\])/g, '\\$1');

    return value.replaceAll(new RegExp(escapedOriginalValue, 'g'), targetValue);
}

export function removeAll(value: string, originalValue: string): string {
    return replaceAll(value, originalValue, '');
}

export function limitText(value: string, maxLength: number): string {
    let length = 0;

    for (let i = 0; i < value.length; i++) {
        const c = value.charCodeAt(i);

        if ((c >= 0x0001 && c <= 0x007e) || (0xff60 <= c && c <= 0xff9f)) {
            length++;
        } else {
            length += 2;
        }
    }

    if (length <= maxLength || maxLength <= 3) {
        return value;
    }

    return value.substring(0, maxLength - 3) + '...';
}

export function getTextBefore(fullText: string, text: string): string {
    if (!text) {
        return fullText;
    }

    const index = fullText.indexOf(text);

    if (index >= 0) {
        return fullText.substring(0, index);
    }

    return '';
}

export function getTextAfter(fullText: string, text: string): string {
    if (!text) {
        return fullText;
    }

    let index = fullText.indexOf(text);

    if (index >= 0) {
        index += text.length;
        return fullText.substring(index);
    }

    return '';
}

export function base64encode(arrayBuffer: ArrayBuffer): string | null {
    if (!arrayBuffer) {
        return null;
    }

    return btoa(String.fromCharCode.apply(null, Array.from(new Uint8Array(arrayBuffer))));
}

export function base64decode(str: string): string | null {
    if (!str) {
        return '';
    }

    return atob(str);
}

export function arrayBufferToString(arrayBuffer: ArrayBuffer): string {
    return String.fromCharCode.apply(null, Array.from(new Uint8Array(arrayBuffer)));
}

export function stringToArrayBuffer(str: string): ArrayBuffer {
    return Uint8Array.from(str, c => c.charCodeAt(0)).buffer;
}

export function getFirstVisibleItem(items: Record<string, unknown>[] | Record<string, Record<string, unknown>>, hiddenField: string): unknown {
    if (isArray(items) && (items as Record<string, unknown>[]).length > 0) {
        const arr = items as Record<string, unknown>[];

        for (let i = 0; i < arr.length; i++) {
            if (hiddenField && arr[i][hiddenField]) {
                continue;
            }

            return arr[i];
        }
    } else if (isObject(items)) {
        const obj = items as Record<string, Record<string, unknown>>;

        for (const field in obj) {
            if (!Object.prototype.hasOwnProperty.call(obj, field)) {
                continue;
            }

            if (hiddenField && obj[field][hiddenField]) {
                continue;
            }

            return obj[field];
        }
    }

    return null;
}

export function getItemByKeyValue(src: Record<string, unknown>[] | Record<string, Record<string, unknown>>, value: unknown, keyField: string) {
    if (isArray(src)) {
        const arr = src as Record<string, unknown>[];

        for (let i = 0; i < arr.length; i++) {
            const item = arr[i];

            if (item[keyField] === value) {
                return item;
            }
        }
    } else if (isObject(src)) {
        const obj = src as Record<string, Record<string, unknown>>;

        for (const field in obj) {
            if (!Object.prototype.hasOwnProperty.call(src, field)) {
                continue;
            }

            const item = obj[field];

            if (item[keyField] === value) {
                return item;
            }
        }
    }

    return null;
}

export function getNameByKeyValue(src: Record<string, unknown>[] | Record<string, Record<string, unknown>>, value: unknown, keyField: string, nameField: string, defaultName: string): unknown {
    if (isArray(src)) {
        const arr = src as Record<string, unknown>[];

        if (keyField) {
            for (let i = 0; i < arr.length; i++) {
                const option = arr[i];

                if (option[keyField] === value) {
                    return option[nameField];
                }
            }
        } else if (isNumber(value)) {
            const index = value as number;

            if (arr[index]) {
                const option = arr[index];

                return option[nameField];
            }
        }
    } else if (isObject(src)) {
        const obj = src as Record<string, Record<string, unknown>>;

        if (keyField) {
            for (const key in obj) {
                if (!Object.prototype.hasOwnProperty.call(src, key)) {
                    continue;
                }

                const option = obj[key];

                if (option[keyField] === value) {
                    return option[nameField];
                }
            }
        } else if (isNumber(value)) {
            const index = value as number;

            if (obj[index]) {
                const option = obj[index];

                return option[nameField];
            }
        }
    }

    return defaultName;
}

export function copyObjectTo(fromObject: Record<string, unknown>, toObject: Record<string, unknown>): object {
    if (!isObject(fromObject)) {
        return toObject;
    }

    if (!isObject(toObject)) {
        toObject = {};
    }

    for (const key in fromObject) {
        if (!Object.prototype.hasOwnProperty.call(fromObject, key)) {
            continue;
        }

        const fromValue = fromObject[key];
        const toValue = toObject[key];

        if (isArray(fromValue)) {
            toObject[key] = copyArrayTo(fromValue as unknown[], toValue as unknown[]);
        } else if (isObject(fromValue)) {
            toObject[key] = copyObjectTo(fromValue as Record<string, unknown>, toValue as Record<string, unknown>);
        } else {
            if (fromValue !== toValue) {
                toObject[key] = fromValue;
            }
        }
    }

    return toObject;
}

export function copyArrayTo(fromArray: unknown[], toArray: unknown[]) {
    if (!isArray(fromArray)) {
        return toArray;
    }

    if (!isArray(toArray)) {
        toArray = [];
    }

    for (let i = 0; i < fromArray.length; i++) {
        const fromValue = fromArray[i];

        if (toArray.length > i) {
            const toValue = toArray[i];

            if (isArray(fromValue)) {
                toArray[i] = copyArrayTo(fromValue as unknown[], toValue as unknown[]);
            } else if (isObject(fromValue)) {
                toArray[i] = copyObjectTo(fromValue as Record<string, unknown>, toValue as Record<string, unknown>);
            } else {
                if (fromValue !== toValue) {
                    toArray[i] = fromValue;
                }
            }
        } else {
            if (isArray(fromValue)) {
                toArray.push(copyArrayTo(fromValue as unknown[], []));
            } else if (isObject(fromValue)) {
                toArray.push(copyObjectTo(fromValue as Record<string, unknown>, {}));
            } else {
                toArray.push(fromValue);
            }
        }
    }

    return toArray;
}

export function arrayContainsFieldValue(array: Record<string, unknown>[], fieldName: string, value: unknown): boolean {
    if (!value || !array || !array.length) {
        return false;
    }

    for (let i = 0; i < array.length; i++) {
        if (array[i][fieldName] === value) {
            return true;
        }
    }

    return false;
}

export function objectFieldToArrayItem(object: object): string[] {
    const ret: string[] = [];

    for (const field in object) {
        if (!Object.prototype.hasOwnProperty.call(object, field)) {
            continue;
        }

        ret.push(field);
    }

    return ret;
}

export function arrayItemToObjectField(array: string[], value: unknown): Record<string, unknown> {
    const ret: Record<string, unknown> = {};

    for (let i = 0; i < array.length; i++) {
        ret[array[i]] = value;
    }

    return ret;
}

export function categorizedArrayToPlainArray(object: Record<string, unknown[]>): unknown[] {
    const ret: unknown[] = [];

    for (const field in object) {
        if (!Object.prototype.hasOwnProperty.call(object, field)) {
            continue;
        }

        const array = object[field];

        for (let i = 0; i < array.length; i++) {
            ret.push(array[i]);
        }
    }

    return ret;
}

export function selectAll(filterItemIds: Record<string, boolean>, allItemsMap: { [key: string]: { id: string } }): void {
    for (const itemId in filterItemIds) {
        if (!Object.prototype.hasOwnProperty.call(filterItemIds, itemId)) {
            continue;
        }

        const item = allItemsMap[itemId];

        if (item) {
            filterItemIds[item.id] = false;
        }
    }
}

export function selectNone(filterItemIds: Record<string, boolean>, allItemsMap: { [key: string]: { id: string } }): void {
    for (const itemId in filterItemIds) {
        if (!Object.prototype.hasOwnProperty.call(filterItemIds, itemId)) {
            continue;
        }

        const item = allItemsMap[itemId];

        if (item) {
            filterItemIds[item.id] = true;
        }
    }
}

export function selectInvert(filterItemIds: Record<string, boolean>, allItemsMap: { [key: string]: { id: string } }): void {
    for (const itemId in filterItemIds) {
        if (!Object.prototype.hasOwnProperty.call(filterItemIds, itemId)) {
            continue;
        }

        const item = allItemsMap[itemId];

        if (item) {
            filterItemIds[item.id] = !filterItemIds[item.id];
        }
    }
}

export function isPrimaryItemHasSecondaryValue(primaryItem: Record<string, Record<string, unknown>[]>, primarySubItemsField: string, secondaryValueField: string, secondaryHiddenField: string, secondaryValue: unknown): boolean {
    for (let i = 0; i < primaryItem[primarySubItemsField].length; i++) {
        const secondaryItem = primaryItem[primarySubItemsField][i];

        if (secondaryHiddenField && secondaryItem[secondaryHiddenField]) {
            continue;
        }

        if (secondaryValueField && secondaryItem[secondaryValueField] === secondaryValue) {
            return true;
        } else if (!secondaryValueField && secondaryItem === secondaryValue) {
            return true;
        }
    }

    return false;
}

export function getPrimaryValueBySecondaryValue(items: Record<string, Record<string, unknown>[]>[] | Record<string, Record<string, Record<string, unknown>[]>>, primarySubItemsField: string, primaryValueField: string, primaryHiddenField: string, secondaryValueField: string, secondaryHiddenField: string, secondaryValue: unknown): unknown {
    if (primarySubItemsField) {
        if (isArray(items)) {
            const arr = items as Record<string, Record<string, unknown>[]>[];

            for (let i = 0; i < arr.length; i++) {
                const primaryItem = arr[i];

                if (primaryHiddenField && primaryItem[primaryHiddenField]) {
                    continue;
                }

                if (isPrimaryItemHasSecondaryValue(primaryItem, primarySubItemsField, secondaryValueField, secondaryHiddenField, secondaryValue)) {
                    if (primaryValueField) {
                        return primaryItem[primaryValueField];
                    } else {
                        return primaryItem;
                    }
                }
            }
        } else {
            const obj = items as Record<string, Record<string, Record<string, unknown>[]>>;

            for (const field in obj) {
                if (!Object.prototype.hasOwnProperty.call(obj, field)) {
                    continue;
                }

                const primaryItem = obj[field];

                if (primaryHiddenField && primaryItem[primaryHiddenField]) {
                    continue;
                }

                if (isPrimaryItemHasSecondaryValue(primaryItem, primarySubItemsField, secondaryValueField, secondaryHiddenField, secondaryValue)) {
                    if (primaryValueField) {
                        return primaryItem[primaryValueField];
                    } else {
                        return primaryItem;
                    }
                }
            }
        }
    }

    return null;
}

export function arrangeArrayWithNewStartIndex(array: unknown[], startIndex: number): unknown[] {
    if (startIndex <= 0 || startIndex >= array.length) {
        return array;
    }

    const newArray = [];

    for (let i = startIndex; i < array.length; i++) {
        newArray.push(array[i]);
    }

    for (let i = 0; i < startIndex; i++) {
        newArray.push(array[i]);
    }

    return newArray;
}