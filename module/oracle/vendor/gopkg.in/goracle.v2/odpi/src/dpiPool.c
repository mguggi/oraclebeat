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
// dpiPool.c
//   Implementation of session pools.
//-----------------------------------------------------------------------------

#include "dpiImpl.h"

//-----------------------------------------------------------------------------
// dpiPool__acquireConnection() [INTERNAL]
//   Internal method used for acquiring a connection from a pool.
//-----------------------------------------------------------------------------
int dpiPool__acquireConnection(dpiPool *pool, const char *userName,
        uint32_t userNameLength, const char *password, uint32_t passwordLength,
        dpiConnCreateParams *params, dpiConn **conn, dpiError *error)
{
    dpiConn *tempConn;

    // allocate new connection
    if (dpiGen__allocate(DPI_HTYPE_CONN, pool->env, (void**) &tempConn,
            error) < 0)
        return DPI_FAILURE;

    // create the connection
    if (dpiConn__get(tempConn, userName, userNameLength, password,
            passwordLength, pool->name, pool->nameLength, params, pool,
            error) < 0) {
        dpiConn__free(tempConn, error);
        return DPI_FAILURE;
    }

    *conn = tempConn;
    return DPI_SUCCESS;
}


//-----------------------------------------------------------------------------
// dpiPool__checkConnected() [INTERNAL]
//   Determine if the session pool is connected to the database. If not, an
// error is raised.
//-----------------------------------------------------------------------------
static int dpiPool__checkConnected(dpiPool *pool, const char *fnName,
        dpiError *error)
{
    if (dpiGen__startPublicFn(pool, DPI_HTYPE_POOL, fnName, error) < 0)
        return DPI_FAILURE;
    if (!pool->handle)
        return dpiError__set(error, "check pool", DPI_ERR_NOT_CONNECTED);
    return DPI_SUCCESS;
}


//-----------------------------------------------------------------------------
// dpiPool__create() [INTERNAL]
//   Internal method for creating a session pool.
//-----------------------------------------------------------------------------
static int dpiPool__create(dpiPool *pool, const char *userName,
        uint32_t userNameLength, const char *password, uint32_t passwordLength,
        const char *connectString, uint32_t connectStringLength,
        const dpiCommonCreateParams *commonParams,
        dpiPoolCreateParams *createParams, dpiError *error)
{
    uint32_t poolMode;
    uint8_t getMode;
    void *authInfo;

    // validate parameters
    if (createParams->externalAuth &&
            ((userName && userNameLength > 0) ||
             (password && passwordLength > 0)))
        return dpiError__set(error, "check mixed credentials",
                DPI_ERR_EXT_AUTH_WITH_CREDENTIALS);

    // create the session pool handle
    if (dpiOci__handleAlloc(pool->env, &pool->handle, DPI_OCI_HTYPE_SPOOL,
            "allocate pool handle", error) < 0)
        return DPI_FAILURE;

    // prepare pool mode
    poolMode = DPI_OCI_SPC_STMTCACHE;
    if (createParams->homogeneous)
        poolMode |= DPI_OCI_SPC_HOMOGENEOUS;

    // create authorization handle
    if (dpiOci__handleAlloc(pool->env, &authInfo, DPI_OCI_HTYPE_AUTHINFO,
            "allocate authinfo handle", error) < 0)
        return DPI_FAILURE;

    // set context attributes
    if (dpiUtils__setAttributesFromCommonCreateParams(authInfo,
            DPI_OCI_HTYPE_AUTHINFO, commonParams, error) < 0)
        return DPI_FAILURE;

    // set authorization info on session pool
    if (dpiOci__attrSet(pool->handle, DPI_OCI_HTYPE_SPOOL, (void*) authInfo, 0,
            DPI_OCI_ATTR_SPOOL_AUTH, "set auth info", error) < 0)
        return DPI_FAILURE;

    // create pool
    if (dpiOci__sessionPoolCreate(pool, connectString, connectStringLength,
            createParams->minSessions, createParams->maxSessions,
            createParams->sessionIncrement, userName, userNameLength, password,
            passwordLength, poolMode, error) < 0)
        return DPI_FAILURE;

    // set the get mode on the pool
    getMode = (uint8_t) createParams->getMode;
    if (dpiOci__attrSet(pool->handle, DPI_OCI_HTYPE_SPOOL, (void*) &getMode, 0,
            DPI_OCI_ATTR_SPOOL_GETMODE, "set get mode", error) < 0)
        return DPI_FAILURE;

    // set reamining attributes directly
    pool->homogeneous = createParams->homogeneous;
    pool->externalAuth = createParams->externalAuth;
    pool->pingInterval = createParams->pingInterval;
    pool->pingTimeout = createParams->pingTimeout;
    return DPI_SUCCESS;
}


//-----------------------------------------------------------------------------
// dpiPool__free() [INTERNAL]
//   Free any memory associated with the pool.
//-----------------------------------------------------------------------------
void dpiPool__free(dpiPool *pool, dpiError *error)
{
    if (pool->handle) {
        dpiOci__sessionPoolDestroy(pool, DPI_OCI_SPD_FORCE, 0, error);
        dpiOci__handleFree(pool->handle, DPI_OCI_HTYPE_SPOOL);
        pool->handle = NULL;
    }
    if (pool->env) {
        dpiEnv__free(pool->env, error);
        pool->env = NULL;
    }
    free(pool);
}


//-----------------------------------------------------------------------------
// dpiPool__getAttributeUint() [INTERNAL]
//   Return the value of the attribute as an unsigned integer.
//-----------------------------------------------------------------------------
static int dpiPool__getAttributeUint(dpiPool *pool, uint32_t attribute,
        uint32_t *value, const char *fnName)
{
    dpiError error;

    if (dpiPool__checkConnected(pool, fnName, &error) < 0)
        return DPI_FAILURE;
    DPI_CHECK_PTR_NOT_NULL(value)
    switch (attribute) {
        case DPI_OCI_ATTR_SPOOL_BUSY_COUNT:
        case DPI_OCI_ATTR_SPOOL_MAX_LIFETIME_SESSION:
        case DPI_OCI_ATTR_SPOOL_OPEN_COUNT:
        case DPI_OCI_ATTR_SPOOL_STMTCACHESIZE:
        case DPI_OCI_ATTR_SPOOL_TIMEOUT:
            return dpiOci__attrGet(pool->handle, DPI_OCI_HTYPE_SPOOL, value,
                    NULL, attribute, "get attribute value", &error);
    }
    return dpiError__set(&error, "get attribute", DPI_ERR_NOT_SUPPORTED);
}


//-----------------------------------------------------------------------------
// dpiPool__setAttributeUint() [INTERNAL]
//   Set the value of the OCI attribute as an unsigned integer.
//-----------------------------------------------------------------------------
static int dpiPool__setAttributeUint(dpiPool *pool, uint32_t attribute,
        uint32_t value, const char *fnName)
{
    uint8_t shortValue;
    void *ociValue;
    dpiError error;

    // make sure session pool is connected
    if (dpiPool__checkConnected(pool, fnName, &error) < 0)
        return DPI_FAILURE;

    // determine pointer to pass (OCI uses different sizes)
    switch (attribute) {
        case DPI_OCI_ATTR_SPOOL_GETMODE:
            shortValue = (uint8_t) value;
            ociValue = &shortValue;
            break;
        case DPI_OCI_ATTR_SPOOL_MAX_LIFETIME_SESSION:
        case DPI_OCI_ATTR_SPOOL_STMTCACHESIZE:
        case DPI_OCI_ATTR_SPOOL_TIMEOUT:
            ociValue = &value;
            break;
        default:
            return dpiError__set(&error, "set attribute",
                    DPI_ERR_NOT_SUPPORTED);
    }

    // set value in the OCI
    return dpiOci__attrSet(pool->handle, DPI_OCI_HTYPE_SPOOL, ociValue, 0,
            attribute, "set attribute value", &error);
}


//-----------------------------------------------------------------------------
// dpiPool_acquireConnection() [PUBLIC]
//   Acquire a connection from the pool.
//-----------------------------------------------------------------------------
int dpiPool_acquireConnection(dpiPool *pool, const char *userName,
        uint32_t userNameLength, const char *password, uint32_t passwordLength,
        dpiConnCreateParams *params, dpiConn **conn)
{
    dpiConnCreateParams localParams;
    dpiError error;

    // validate parameters
    if (dpiPool__checkConnected(pool, __func__, &error) < 0)
        return DPI_FAILURE;
    DPI_CHECK_PTR_AND_LENGTH(userName)
    DPI_CHECK_PTR_AND_LENGTH(password)
    DPI_CHECK_PTR_NOT_NULL(conn)

    // use default parameters if none provided
    if (!params) {
        if (dpiContext__initConnCreateParams(pool->env->context, &localParams,
                &error) < 0)
            return DPI_FAILURE;
        params = &localParams;
    }

    return dpiPool__acquireConnection(pool, userName, userNameLength, password,
            passwordLength, params, conn, &error);
}


//-----------------------------------------------------------------------------
// dpiPool_addRef() [PUBLIC]
//   Add a reference to the pool.
//-----------------------------------------------------------------------------
int dpiPool_addRef(dpiPool *pool)
{
    return dpiGen__addRef(pool, DPI_HTYPE_POOL, __func__);
}


//-----------------------------------------------------------------------------
// dpiPool_close() [PUBLIC]
//   Destroy the pool now, not when the reference count reaches zero.
//-----------------------------------------------------------------------------
int dpiPool_close(dpiPool *pool, dpiPoolCloseMode mode)
{
    dpiError error;

    if (dpiPool__checkConnected(pool, __func__, &error) < 0)
        return DPI_FAILURE;
    if (dpiOci__sessionPoolDestroy(pool, mode, 1, &error) < 0)
        return DPI_FAILURE;
    dpiOci__handleFree(pool->handle, DPI_OCI_HTYPE_SPOOL);
    pool->handle = NULL;
    return DPI_SUCCESS;
}


//-----------------------------------------------------------------------------
// dpiPool_create() [PUBLIC]
//   Create a new session pool and return it.
//-----------------------------------------------------------------------------
int dpiPool_create(const dpiContext *context, const char *userName,
        uint32_t userNameLength, const char *password, uint32_t passwordLength,
        const char *connectString, uint32_t connectStringLength,
        const dpiCommonCreateParams *commonParams,
        dpiPoolCreateParams *createParams, dpiPool **pool)
{
    dpiCommonCreateParams localCommonParams;
    dpiPoolCreateParams localCreateParams;
    dpiPool *tempPool;
    dpiError error;

    // validate parameters
    if (dpiContext__startPublicFn(context, __func__, &error) < 0)
        return DPI_FAILURE;
    DPI_CHECK_PTR_AND_LENGTH(userName)
    DPI_CHECK_PTR_AND_LENGTH(password)
    DPI_CHECK_PTR_AND_LENGTH(connectString)
    DPI_CHECK_PTR_NOT_NULL(pool)

    // use default parameters if none provided
    if (!commonParams) {
        if (dpiContext__initCommonCreateParams(context, &localCommonParams,
                &error) < 0)
            return DPI_FAILURE;
        commonParams = &localCommonParams;
    }
    if (!createParams) {
        if (dpiContext__initPoolCreateParams(context, &localCreateParams,
                &error) < 0)
            return DPI_FAILURE;
        createParams = &localCreateParams;
    }

    // allocate memory for pool
    if (dpiGen__allocate(DPI_HTYPE_POOL, NULL, (void**) &tempPool, &error) < 0)
        return DPI_FAILURE;

    // initialize environment
    if (dpiEnv__init(tempPool->env, context, commonParams, &error) < 0) {
        dpiPool__free(tempPool, &error);
        return DPI_FAILURE;
    }

    // perform remaining steps required to create pool
    if (dpiPool__create(tempPool, userName, userNameLength, password,
            passwordLength, connectString, connectStringLength, commonParams,
            createParams, &error) < 0) {
        dpiPool__free(tempPool, &error);
        return DPI_FAILURE;
    }

    createParams->outPoolName = tempPool->name;
    createParams->outPoolNameLength = tempPool->nameLength;
    *pool = tempPool;
    return DPI_SUCCESS;
}


//-----------------------------------------------------------------------------
// dpiPool_getBusyCount() [PUBLIC]
//   Return the pool's busy count.
//-----------------------------------------------------------------------------
int dpiPool_getBusyCount(dpiPool *pool, uint32_t *value)
{
    return dpiPool__getAttributeUint(pool, DPI_OCI_ATTR_SPOOL_BUSY_COUNT,
            value, __func__);
}


//-----------------------------------------------------------------------------
// dpiPool_getEncodingInfo() [PUBLIC]
//   Get the encoding information from the pool.
//-----------------------------------------------------------------------------
int dpiPool_getEncodingInfo(dpiPool *pool, dpiEncodingInfo *info)
{
    dpiError error;

    if (dpiPool__checkConnected(pool, __func__, &error) < 0)
        return DPI_FAILURE;
    DPI_CHECK_PTR_NOT_NULL(info)
    return dpiEnv__getEncodingInfo(pool->env, info);
}


//-----------------------------------------------------------------------------
// dpiPool_getGetMode() [PUBLIC]
//   Return the pool's "get" mode.
//-----------------------------------------------------------------------------
int dpiPool_getGetMode(dpiPool *pool, dpiPoolGetMode *value)
{
    uint8_t ociValue;
    dpiError error;

    if (dpiPool__checkConnected(pool, __func__, &error) < 0)
        return DPI_FAILURE;
    DPI_CHECK_PTR_NOT_NULL(value)
    if (dpiOci__attrGet(pool->handle, DPI_OCI_HTYPE_SPOOL, &ociValue, NULL,
            DPI_OCI_ATTR_SPOOL_GETMODE, "get attribute value", &error) < 0)
        return DPI_FAILURE;
    *value = ociValue;
    return DPI_SUCCESS;
}


//-----------------------------------------------------------------------------
// dpiPool_getMaxLifetimeSession() [PUBLIC]
//   Return the pool's maximum lifetime session.
//-----------------------------------------------------------------------------
int dpiPool_getMaxLifetimeSession(dpiPool *pool, uint32_t *value)
{
    return dpiPool__getAttributeUint(pool,
            DPI_OCI_ATTR_SPOOL_MAX_LIFETIME_SESSION, value, __func__);
}


//-----------------------------------------------------------------------------
// dpiPool_getOpenCount() [PUBLIC]
//   Return the pool's open count.
//-----------------------------------------------------------------------------
int dpiPool_getOpenCount(dpiPool *pool, uint32_t *value)
{
    return dpiPool__getAttributeUint(pool, DPI_OCI_ATTR_SPOOL_OPEN_COUNT,
            value, __func__);
}


//-----------------------------------------------------------------------------
// dpiPool_getStmtCacheSize() [PUBLIC]
//   Return the pool's default statement cache size.
//-----------------------------------------------------------------------------
int dpiPool_getStmtCacheSize(dpiPool *pool, uint32_t *value)
{
    return dpiPool__getAttributeUint(pool, DPI_OCI_ATTR_SPOOL_STMTCACHESIZE,
            value, __func__);
}


//-----------------------------------------------------------------------------
// dpiPool_getTimeout() [PUBLIC]
//   Return the pool's timeout value.
//-----------------------------------------------------------------------------
int dpiPool_getTimeout(dpiPool *pool, uint32_t *value)
{
    return dpiPool__getAttributeUint(pool, DPI_OCI_ATTR_SPOOL_TIMEOUT, value,
            __func__);
}


//-----------------------------------------------------------------------------
// dpiPool_release() [PUBLIC]
//   Release a reference to the pool.
//-----------------------------------------------------------------------------
int dpiPool_release(dpiPool *pool)
{
    return dpiGen__release(pool, DPI_HTYPE_POOL, __func__);
}


//-----------------------------------------------------------------------------
// dpiPool_setGetMode() [PUBLIC]
//   Set the pool's "get" mode.
//-----------------------------------------------------------------------------
int dpiPool_setGetMode(dpiPool *pool, dpiPoolGetMode value)
{
    return dpiPool__setAttributeUint(pool, DPI_OCI_ATTR_SPOOL_GETMODE, value,
            __func__);
}


//-----------------------------------------------------------------------------
// dpiPool_setMaxLifetimeSession() [PUBLIC]
//   Set the pool's maximum lifetime session.
//-----------------------------------------------------------------------------
int dpiPool_setMaxLifetimeSession(dpiPool *pool, uint32_t value)
{
    return dpiPool__setAttributeUint(pool,
            DPI_OCI_ATTR_SPOOL_MAX_LIFETIME_SESSION, value, __func__);
}


//-----------------------------------------------------------------------------
// dpiPool_setStmtCacheSize() [PUBLIC]
//   Set the pool's default statement cache size.
//-----------------------------------------------------------------------------
int dpiPool_setStmtCacheSize(dpiPool *pool, uint32_t value)
{
    return dpiPool__setAttributeUint(pool, DPI_OCI_ATTR_SPOOL_STMTCACHESIZE,
            value, __func__);
}


//-----------------------------------------------------------------------------
// dpiPool_setTimeout() [PUBLIC]
//   Set the pool's timeout value.
//-----------------------------------------------------------------------------
int dpiPool_setTimeout(dpiPool *pool, uint32_t value)
{
    return dpiPool__setAttributeUint(pool, DPI_OCI_ATTR_SPOOL_TIMEOUT, value,
            __func__);
}

