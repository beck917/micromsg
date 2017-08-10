package constant

const ERROR_UTILS_POST int = 10001 //请求错误
const ERROR_UTILS_POST_RESPONSE_READ int = 10002

const ERROR_VENDOR_MARSHAL int = 20001
const ERROR_VENDOR_POST int = 20002 //请求错误
const ERROR_VENDOR_RESPONSE_DECODE int = 20003

const ERROR_VENDOR_RESPONSE_MESSAGE_FORMAT int = 30001   //格式不对,放弃
const ERROR_VENDOR_RESPONSE_MESSAGE_EXPIRE int = 30002   //过期,需要重新打包
const ERROR_VENDOR_RESPONSE_MESSAGE_UNKNOWN int = 30003  //未知错误,重试n次后放弃
const ERROR_VENDOR_RESPONSE_MESSAGE_RESENDED int = 30004 //重复发送的票

/**
const ERROR_VENDOR_RESPONSE_MESSAGE_PROTOCAL int = 30001
const ERROR_VENDOR_RESPONSE_MESSAGE_ID int = 30001
const ERROR_VENDOR_RESPONSE_MESSAGE_TIMESTAMP int = 30001
const ERROR_VENDOR_RESPONSE_MESSAGE_MD5 int = 30001
const ERROR_VENDOR_RESPONSE_MESSAGE_UNSUPPUORT_TRANSTYPE int = 30001
const ERROR_VENDOR_RESPONSE_MESSAGE_ID_RETRY int = 30001
const ERROR_VENDOR_RESPONSE_MESSAGE_BINGFA int = 30001
*/
