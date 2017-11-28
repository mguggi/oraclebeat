.. _dpiPoolCreateParams:

ODPI-C Public Structure dpiPoolCreateParams
-------------------------------------------

This structure is used for creating session pools, which can in turn be used to
create connections that are acquired from that session pool. All members are
initialized to default values using the
:func:`dpiContext_initPoolCreateParams()` function.

.. member:: uint32_t dpiPoolCreateParams.minSessions

    Specifies the minimum number of sessions to be created by the session pool.
    This value is ignored if the :member:`dpiPoolCreateParams.homogeneous`
    member has a value of 0. The default value is 1.

.. member:: uint32_t dpiPoolCreateParams.maxSessions

    Specifies the maximum number of sessions that can be created by the session
    pool. Values of 1 and higher are acceptable. The default value is 1.

.. member:: uint32_t dpiPoolCreateParams.sessionIncrement

    Specifies the number of sessions that will be created by the session pool
    when more sessions are required and the number of sessions is less than
    the maximum allowed. This value is ignored if the
    :member:`dpiPoolCreateParams.homogeneous` member has a value of 0. This
    value added to the :member:`dpiPoolCreateParams.minSessions` member value
    must not exceed the :member:`dpiPoolCreateParams.maxSessions` member value.
    The default value is 0.

.. member:: int dpiPoolCreateParams.pingInterval

    Specifies the number of seconds since a connection has last been used
    before a ping will be performed to verify that the connection is still
    valid. A negative value disables this check. The default value is 60.
    This value is ignored in clients 12.2 and later since a much faster
    internal check is done by the Oracle client.

.. member:: int dpiPoolCreateParams.pingTimeout

    Specifies the number of milliseconds to wait when performing a ping to
    verify the connection is still valid before the connection is considered
    invalid and is dropped. The default value is 5000 (5 seconds).  This value
    is ignored in clients 12.2 and later since a much faster internal check is
    done by the Oracle client.

.. member:: int dpiPoolCreateParams.homogeneous

    Specifies whether the pool is homogeneous or not. In a homogeneous pool all
    connections use the same credentials whereas in a heterogeneous pool other
    credentials are permitted. The default value is 1.

.. member:: int dpiPoolCreateParams.externalAuth

    Specifies whether external authentication should be used to create the
    sessions in the pool. If this value is 0, the user name and password values
    must be specified in the call to :func:`dpiPool_create()`; otherwise, the
    user name and password values must be zero length or NULL. The default
    value is 0. External authentication cannot be used with homogeneous pools.

.. member:: dpiPoolGetMode dpiPoolCreateParams.getMode

    Specifies the mode to use when sessions are acquired from the pool. It is
    expected to be one of the values from the enumeration
    :ref:`dpiPoolGetMode<dpiPoolGetMode>`. The default value is
    DPI_MODE_POOL_GET_NOWAIT

.. member:: const char \*dpiPoolCreateParams.outPoolName

    This member is populated upon successful creation of a pool using the
    function :func:`dpiPool_create()`. It is a byte string in the encoding
    used for CHAR data. Any value specified prior to creating the session pool
    is ignored.

.. member:: uint32_t dpiPoolCreateParams.outPoolNameLength

    This member is populated upon successful creation of a pool using the
    function :func:`dpiPool_create()`. It is the length of the
    :member:`dpiPoolCreateParams.outPoolName` member, in bytes. Any value
    specified prior to creating the session pool is ignored.

