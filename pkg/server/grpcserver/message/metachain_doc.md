# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [message.proto](#message-proto)
    - [Tx](#message-Tx)
    - [req_balance](#message-req_balance)
    - [req_block_by_hash](#message-req_block_by_hash)
    - [req_block_by_number](#message-req_block_by_number)
    - [req_max_blockHeight](#message-req_max_blockHeight)
    - [req_nonce](#message-req_nonce)
    - [req_transaction](#message-req_transaction)
    - [req_tx_by_hash](#message-req_tx_by_hash)
    - [res_balance](#message-res_balance)
    - [res_max_blockHeight](#message-res_max_blockHeight)
    - [res_transaction](#message-res_transaction)
    - [resp_block](#message-resp_block)
    - [resp_block_data](#message-resp_block_data)
    - [resp_tx_by_hash](#message-resp_tx_by_hash)
    - [respose_nonce](#message-respose_nonce)
  
    - [Greeter](#message-Greeter)
  
- [Scalar Value Types](#scalar-value-types)



<a name="message-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## message.proto



<a name="message-Tx"></a>

### Tx
交易数据


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Version | [uint64](#uint64) |  | 版本号 |
| Type | [uint64](#uint64) |  | 交易类型 |
| Amount | [string](#string) |  | 交易数额 |
| From | [string](#string) |  | 交易发送者 |
| To | [string](#string) |  | 交易接收者 |
| GasPrice | [string](#string) |  | gas费单价 |
| GasFeeCap | [string](#string) |  | gas费容量 |
| GasLimit | [string](#string) |  | gas数量限制 |
| Input | [bytes](#bytes) |  | 附加数据 |
| Nonce | [uint64](#uint64) |  | 随机数 |






<a name="message-req_balance"></a>

### req_balance



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  | 地址 |






<a name="message-req_block_by_hash"></a>

### req_block_by_hash
通过块hash获取块数据的请求


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hash | [string](#string) |  | 块哈希 |






<a name="message-req_block_by_number"></a>

### req_block_by_number
通过块高获取块数据的请求


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [uint64](#uint64) |  | 块高 |






<a name="message-req_max_blockHeight"></a>

### req_max_blockHeight







<a name="message-req_nonce"></a>

### req_nonce
获取nonce接口的请求


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  | 地址 |






<a name="message-req_transaction"></a>

### req_transaction
发送交易接口的请求数据


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| From | [string](#string) |  | 交易的发送者 |
| To | [string](#string) |  | 交易的接收者 |
| Amount | [string](#string) |  | 交易数量，18位精度 |
| Nonce | [uint64](#uint64) |  | 递增的无符号整数，由From地址维护。调用GetAddressNonceAt接口获取当前发送交易所需的nonce。 |
| Sign | [bytes](#bytes) |  | 签名数据 |
| GasLimit | [string](#string) |  | gas费数量限制 |
| GasFeeCap | [string](#string) |  | gas费上限 |
| GasPrice | [string](#string) |  | gas费单价 |
| Input | [bytes](#bytes) |  | 附加数据，可以为nil |






<a name="message-req_tx_by_hash"></a>

### req_tx_by_hash
查询交易的请求


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hash | [string](#string) |  | 交易哈希 |






<a name="message-res_balance"></a>

### res_balance



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| balance | [string](#string) |  | 余额，18位精度 |






<a name="message-res_max_blockHeight"></a>

### res_max_blockHeight



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| max_height | [uint64](#uint64) |  | 当前最大的块高度 |






<a name="message-res_transaction"></a>

### res_transaction
发送交易接口的返回值


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Hash | [string](#string) |  | 交易哈希 |






<a name="message-resp_block"></a>

### resp_block
块数据


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  | 状态码，0为正常 |
| message | [string](#string) |  | 错误信息 |
| data | [bytes](#bytes) |  | 块数据 |






<a name="message-resp_block_data"></a>

### resp_block_data
获取块数据的响应


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  | 状态码，0为正常 |
| data | [bytes](#bytes) |  | 块数据 |
| message | [string](#string) |  | 错误信息 |






<a name="message-resp_tx_by_hash"></a>

### resp_tx_by_hash
查询交易的返回值


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  | 状态码，0为正常 |
| data | [bytes](#bytes) |  | 交易数据 |
| message | [string](#string) |  | 接口返回的错误信息 |






<a name="message-respose_nonce"></a>

### respose_nonce
获取nonce的响应


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nonce | [uint64](#uint64) |  | 随机数 |





 

 

 


<a name="message-Greeter"></a>

### Greeter


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetBalance | [req_balance](#message-req_balance) | [res_balance](#message-res_balance) | 获取地址对应的余额 |
| SendTransaction | [req_transaction](#message-req_transaction) | [res_transaction](#message-res_transaction) | 发送已经签名的交易 |
| GetBlockByNum | [req_block_by_number](#message-req_block_by_number) | [resp_block](#message-resp_block) | 通过块的高度获取到块数据 |
| GetTxByHash | [req_tx_by_hash](#message-req_tx_by_hash) | [resp_tx_by_hash](#message-resp_tx_by_hash) | 通过交易哈希获取交易数据 |
| GetAddressNonceAt | [req_nonce](#message-req_nonce) | [respose_nonce](#message-respose_nonce) | 获取该地址发送交易所需要的noce |
| GetBlockByHash | [req_block_by_hash](#message-req_block_by_hash) | [resp_block_data](#message-resp_block_data) | 通过块哈希获取到块数据 |
| GetMaxBlockHeight | [req_max_blockHeight](#message-req_max_blockHeight) | [res_max_blockHeight](#message-res_max_blockHeight) | 获取当前最大的块高 |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

