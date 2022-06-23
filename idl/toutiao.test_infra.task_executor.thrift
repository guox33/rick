include "base.thrift"

namespace go toutiao.test_infra.task_executor
namespace py toutiao.test_infra.task_executor

typedef i64 AssertType
typedef i64 ServiceType

enum CheckMode {
    RESULT_UNKNOWN = 0
    Value_Equal
    Type_Equal
    Schema
    Value_Contain
    Type_Contain
}

enum RuleType {
    Unkmown = 0
    Value_Equal
    Type_Equal
    Value_UnEqual
    Value_Interval     // 值区间
    Len_Interval       // 长度区间
    In                 // 在某个数组
    NotIn              // 不在某个数组
    Regexp             // 正则
    null_normal        // 是否允许空值 不传未false
    Str_Not_Contain    // string不包含
    Str_Contain        // 包含
    Time_Filter        // 包含
    Struct_Contain     // 结构体包含
    Optional           // 可选择的
    Required_Value     // 要求的值
    Required_Type      // 要求的类型
    Max                // 要求的类型
    Min                // 要求的类型
}

enum TaskStatus {
    Unknown = 0
    Success
    Failed
    Interrupted
}

enum CaseType {
    Unknown = 0
    Http
    Rpc
}

struct CaseResult {
    1: bool assert_result
    2: bool analysis_result
    3: bool final_result
    4: optional string diff_info
    5: optional string fail_reason
}

struct BasicCase {
    2: string psm
    3: string route
    4: string method
    5: string api_id
    6: ServiceType service_type
    7: string expect_response
    8: bool usable_case
    9: string cluster
    10: string env
    11: string idc
    12: bool is_egress
    13: AssertType assert_type
    14: string task_id
}

struct HttpRequestInfo {
    1: string domain
    2: string path
    3: string req_method
    4: map<string, string> req_headers
    5: string req_body
    6: map<string, list<string>> req_query
    7: string req_content_type
    8: i32 resp_logic_code
}

struct ExpectRuleInfo {
    1: i64 id
    2: string psm
    3: string route
    4: string method
    5: string rule_key
    6: string prefix
    7: RuleType rule_type
    8: string expect_value_type
    9: string expect_value
    10: string operator
}

struct CustomizeAssertOption {
    1: optional bool NotLogicCodeCheck
    2: optional bool is_agw_psm
}

struct AssertRequestInfo {
    1: optional CheckMode check_mode
    2: optional CustomizeAssertOption customize_assert_option
    3: required map<string, ExpectRuleInfo> expect_rule_info_map
}

struct AssertResult {
    1: bool success
    2: string err_keys
    3: string err_msg
}

struct HttpCase {
    1: required string case_id // 业务方用于标识用例
    2: required BasicCase basic_case_info
    3: required HttpRequestInfo http_request_info
    4: required AssertRequestInfo assertion_info
}

struct HttpCaseExecutionResult {
    1: i64 status_code
    2: i64 logic_code
    3: string resp_body_str
    4: string log_id
    5: map<string, string> resp_header
}

struct HttpCaseExecutionDetail {
    1: string case_id
    2: HttpCaseExecutionResult execution_result
    3: AssertResult assert_result
    4: string err // 执行此case时发生的error
}

struct ExecutionMsg {
    1: string biz_need
    2: TaskStatus task_status
    3: CaseType case_type
    4: list<HttpCaseExecutionDetail> case_execution_detail_list
}

struct ExecuteHttpCaseBatchRequest {
    1: list<HttpCase> case_list
    2: string biz_need // 业务方执行case的目的，会包含在结果消息中投递

    255: base.Base base
}

struct ExecuteHttpCaseBatchResponse {
    1: required i32 code
    2: required string message

    255: base.BaseResp base_resp
}

struct ExecuteHttpCaseQpsRequest {
    1: required i32 qps
    2: required i32 duration // 超时时间
    3: required list<HttpCase> case_list
    4: bool is_random // 是否随机执行用例
    5: string biz_need
}

struct ExecuteHttpCaseQpsResponse {
    1: required i32 code
    2: required string message
}

service TaskExecutorService {
    ExecuteHttpCaseBatchResponse ExecuteHttpCaseBatch(1: ExecuteHttpCaseBatchRequest req)
    ExecuteHttpCaseQpsResponse ExecuteHttpCaseQps(1: ExecuteHttpCaseQpsRequest req)
}
