//-----------------------------------------------------------------------------
// Copyright (c) 2016, 2017 Oracle and/or its affiliates.  All rights reserved.
// This program is free software: you can modify it and/or redistribute it
// under the terms of:
//
// (i)  the Universal Permissive License v 1.0 or at your option, any
//      later version (http://oss.oracle.com/licenses/upl); and/or
//
// (ii) the Apache License v 2.0. (http://www.apache.org/licenses/LICENSE-2.0)
//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------
// dpiUtils.c
//   Utility methods that aren't specific to a particular type.
//-----------------------------------------------------------------------------

#include "dpiImpl.h"

//-----------------------------------------------------------------------------
// dpiUtils__clearMemory() [INTERNAL]
//   Method for clearing memory that will not be optimised away by the
// compiler. Simple use of memset() can be optimised away. This routine makes
// use of a volatile pointer which most compilers will avoid optimising away,
// even if the pointer appears to be unused after the call.
//-----------------------------------------------------------------------------
void dpiUtils__clearMemory(void *ptr, size_t length)
{
    volatile unsigned char *temp = ptr;

    while (length--)
        *temp++ = '\0';
}


//-----------------------------------------------------------------------------
// dpiUtils__getAttrStringWithDup() [INTERNAL]
//   Get the string attribute from the OCI and duplicate its contents.
//-----------------------------------------------------------------------------
int dpiUtils__getAttrStringWithDup(const char *action, const void *ociHandle,
        uint32_t ociHandleType, uint32_t ociAttribute, const char **value,
        uint32_t *valueLength, dpiError *error)
{
    char *source, *temp;

    if (dpiOci__attrGet(ociHandle, ociHandleType, (void*) &source,
            valueLength, ociAttribute, action, error) < 0)
        return DPI_FAILURE;
    temp = malloc(*valueLength);
    if (!temp)
        return dpiError__set(error, action, DPI_ERR_NO_MEMORY);
    *value = memcpy(temp, source, *valueLength);
    return DPI_SUCCESS;
}


//-----------------------------------------------------------------------------
// dpiUtils__parseNumberString() [INTERNAL]
//   Parse the contents of a string that is supposed to contain a number. The
// number is expected to be in the format (www.json.org):
//   - optional negative sign (-)
//   - any number of digits but at least one (0-9)
//   - an optional decimal point (.)
//   - any number of digits but at least one if decimal point specified (0-9)
//   - an optional exponent indicator (e or E)
//   - an optional exponent sign (+ or -)
//   - any number of digits, but at least one if exponent specified (0-9)
// What is returned is an indication of whether the number is negative, what
// the index of the decimal point in the string is and the list of digits
// without the decimal point. Note that OCI doesn't support more than 40 digits
// so if there are more than this amount an error is raised. OCI doesn't
// support larger than 1e126 so check for this value and raise a numeric
// overflow error if found. OCI also doesn't support smaller than 1E-130 so
// check for this value as well and if smaller than that value simply return
// zero.
//-----------------------------------------------------------------------------
int dpiUtils__parseNumberString(const char *value, uint32_t valueLength,
        uint16_t charsetId, int *isNegative, int16_t *decimalPointIndex,
        uint8_t *numDigits, uint8_t *digits, dpiError *error)
{
    char convertedValue[DPI_NUMBER_AS_TEXT_CHARS], exponentDigits[4];
    int exponentIsNegative, exponent;
    uint8_t numExponentDigits, digit;
    uint32_t convertedValueLength;
    uint16_t *utf16chars, i;
    const char *endValue;

    // empty strings are not valid numbers
    if (valueLength == 0)
        return dpiError__set(error, "zero length", DPI_ERR_INVALID_NUMBER);

    // strings longer than the maximum length of a valid number are also
    // excluded
    if ((charsetId == DPI_CHARSET_ID_UTF16 &&
                    valueLength > DPI_NUMBER_AS_TEXT_CHARS * 2) ||
            (charsetId != DPI_CHARSET_ID_UTF16 &&
                    valueLength > DPI_NUMBER_AS_TEXT_CHARS))
        return dpiError__set(error, "check length",
                DPI_ERR_NUMBER_STRING_TOO_LONG);

    // if value is encoded in UTF-16, convert to single byte encoding first
    // check for values that cannot be encoded in a single byte and are
    // obviously not part of a valid numeric string
    // also verify maximum length of number
    if (charsetId == DPI_CHARSET_ID_UTF16) {
        utf16chars = (uint16_t*) value;
        convertedValue[0] = '\0';
        convertedValueLength = valueLength / 2;
        for (i = 0; i < convertedValueLength; i++) {
            if (*utf16chars > 127)
                return dpiError__set(error, "convert from UTF-16",
                        DPI_ERR_INVALID_NUMBER);
            convertedValue[i] = (char) *utf16chars++;
        }
        value = convertedValue;
        valueLength = convertedValueLength;
    }

    // see if first character is a minus sign (number is negative)
    endValue = value + valueLength;
    *isNegative = (*value == '-');
    if (*isNegative)
        value++;

    // scan for digits until the decimal point or exponent indicator is found
    *numDigits = 0;
    while (value < endValue) {
        if (*value == '.' || *value == 'e' || *value == 'E')
            break;
        if (*value < '0' || *value > '9')
            return dpiError__set(error, "check digits before decimal point",
                    DPI_ERR_INVALID_NUMBER);
        digit = *value++ - '0';
        if (digit == 0 && *numDigits == 0)
            continue;
        *digits++ = digit;
        (*numDigits)++;
    }
    *decimalPointIndex = *numDigits;

    // scan for digits following the decimal point, if applicable
    if (value < endValue && *value == '.') {
        value++;
        while (value < endValue) {
            if (*value == 'e' || *value == 'E')
                break;
            if (*value < '0' || *value > '9')
                return dpiError__set(error, "check digits after decimal point",
                        DPI_ERR_INVALID_NUMBER);
            digit = *value++ - '0';
            if (digit == 0 && *numDigits == 0) {
                (*decimalPointIndex)--;
                continue;
            }
            *digits++ = digit;
            (*numDigits)++;
        }
    }

    // handle exponent, if applicable
    exponent = 0;
    if (value < endValue && (*value == 'e' || *value == 'E')) {
        value++;
        exponentIsNegative = 0;
        numExponentDigits = 0;
        if (value < endValue && (*value == '+' || *value == '-')) {
            exponentIsNegative = (*value == '-');
            value++;
        }
        while (value < endValue) {
            if (*value < '0' || *value > '9')
                return dpiError__set(error, "check digits in exponent",
                        DPI_ERR_INVALID_NUMBER);
            if (numExponentDigits == 3)
                return dpiError__set(error, "check exponent digits > 3",
                        DPI_ERR_NOT_SUPPORTED);
            exponentDigits[numExponentDigits] = *value++;
            numExponentDigits++;
        }
        if (numExponentDigits == 0)
            return dpiError__set(error, "no digits in exponent",
                    DPI_ERR_INVALID_NUMBER);
        exponentDigits[numExponentDigits] = '\0';
        exponent = (int) strtol(exponentDigits, NULL, 0);
        if (exponentIsNegative)
            exponent = -exponent;
        *decimalPointIndex += exponent;
    }

    // if there is anything left in the string, that indicates an invalid
    // number as well
    if (value < endValue)
        return dpiError__set(error, "check string used",
                DPI_ERR_INVALID_NUMBER);

    // strip trailing zeroes
    digits--;
    while (*numDigits > 0 && *digits-- == 0)
        (*numDigits)--;

    // only 40 digits are allowed in an OCI number
    if (*numDigits > DPI_NUMBER_MAX_DIGITS)
        return dpiError__set(error, "check number of digits > 40",
                DPI_ERR_NOT_SUPPORTED);

    // values must be less than 1e126
    if (*decimalPointIndex > 126)
        return dpiError__set(error, "check size of value",
                DPI_ERR_NUMBER_TOO_LARGE);

    // values smaller than 1e-130 are simply returned as zero
    if (*decimalPointIndex < -129) {
        *numDigits = 0;
        *isNegative = 0;
    }

    return DPI_SUCCESS;
}


//-----------------------------------------------------------------------------
// dpiUtils__parseOracleNumber() [INTERNAL]
//   Parse the contents of an Oracle number and return its constituent parts
// so that a string can be generated from it easily.
//-----------------------------------------------------------------------------
int dpiUtils__parseOracleNumber(void *oracleValue, int *isNegative,
        int16_t *decimalPointIndex, uint8_t *numDigits, uint8_t *digits,
        dpiError *error)
{
    uint8_t *source, length, i, byte, digit;
    int8_t ociExponent;

    // the first byte of the structure is a length byte which includes the
    // exponent and the mantissa bytes
    source = (uint8_t*) oracleValue;
    length = *source++ - 1;

    // a mantissa length of 0 implies a value of 0
    if (length == 0) {
        *isNegative = 0;
        *decimalPointIndex = 1;
        *numDigits = 1;
        *digits = 0;
        return DPI_SUCCESS;
    }

    // a mantissa length longer than 20 signals corruption of some kind
    if (length > 20)
        return dpiError__set(error, "check mantissa length",
                DPI_ERR_INVALID_OCI_NUMBER);

    // the second byte of the structure is the exponent
    // positive numbers have the highest order bit set whereas negative numbers
    // have the highest order bit cleared and the bits inverted
    ociExponent = *source++;
    *isNegative = (ociExponent & 0x80) ? 0 : 1;
    if (*isNegative)
        ociExponent = ~ociExponent;
    ociExponent -= 193;
    *decimalPointIndex = ociExponent * 2 + 2;

    // check for the trailing 102 byte for negative numbers and if present,
    // reduce the number of mantissa digits
    if (*isNegative && source[length - 1] == 102)
        length--;

    // process the mantissa which are the remaining bytes
    // each mantissa byte is a base-100 digit
    *numDigits = length * 2;
    for (i = 0; i < length; i++) {
        byte = *source++;

        // positive numbers have 1 added to them; negative numbers are
        // subtracted from the value 101
        if (*isNegative)
            byte = 101 - byte;
        else byte--;

        // process the first digit; leading zeroes are ignored
        digit = (uint8_t) (byte / 10);
        if (digit == 0 && i == 0) {
            (*numDigits)--;
            (*decimalPointIndex)--;
        } else *digits++ = digit;

        // process the second digit; trailing zeroes are ignored
        digit = byte % 10;
        if (digit == 0 && i == length - 1)
            (*numDigits)--;
        else *digits++ = digit;

    }

    return DPI_SUCCESS;
}


//-----------------------------------------------------------------------------
// dpiUtils__setAttributesFromCommonCreateParams() [INTERNAL]
//   Set the attributes on the authorization info structure or session handle
// using the specified parameters.
//-----------------------------------------------------------------------------
int dpiUtils__setAttributesFromCommonCreateParams(void *handle,
        uint32_t handleType, const dpiCommonCreateParams *params,
        dpiError *error)
{
    uint32_t driverNameLength;
    const char *driverName;

    if (params->driverName && params->driverNameLength > 0) {
        driverName = params->driverName;
        driverNameLength = params->driverNameLength;
    } else {
        driverName = DPI_DEFAULT_DRIVER_NAME;
        driverNameLength = (uint32_t) strlen(driverName);
    }
    if (driverName && driverNameLength > 0 && dpiOci__attrSet(handle,
            handleType, (void*) driverName, driverNameLength,
            DPI_OCI_ATTR_DRIVER_NAME, "set driver name", error) < 0)
        return DPI_FAILURE;
    if (params->edition && params->editionLength > 0 &&
            dpiOci__attrSet(handle, handleType,
                    (void*) params->edition, params->editionLength,
                    DPI_OCI_ATTR_EDITION, "set edition", error) < 0)
        return DPI_FAILURE;

    return DPI_SUCCESS;
}

