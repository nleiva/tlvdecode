package main

// Data represents the format of a Cisco-IOS-XR-clns-isis-oper
// :isis/instances/instance/levels/level/detailed-lsps/detailed-lsp
// message.
type Data struct {
	Rows []struct {
		Content struct {
			LspBody       string `json:"lsp_body,omitempty"`
			LspHeaderData struct {
				LocalLspFlag                   bool   `json:"local_lsp_flag,omitempty"`
				LspActiveFlag                  bool   `json:"lsp_active_flag,omitempty"`
				LspAttachedFlag                bool   `json:"lsp_attached_flag,omitempty"`
				LspChecksum                    int64  `json:"lsp_checksum,omitempty"`
				LspHoldtime                    int64  `json:"lsp_holdtime,omitempty"`
				LspID                          string `json:"lsp_id,omitempty"`
				LspLength                      int64  `json:"lsp_length,omitempty"`
				LspLevel                       string `json:"lsp_level,omitempty"`
				LspNonV1AFlag                  int64  `json:"lsp_non_v1_a_flag,omitempty"`
				LspOverloadedFlag              bool   `json:"lsp_overloaded_flag,omitempty"`
				LspParitionRepairSupportedFlag bool   `json:"lsp_parition_repair_supported_flag,omitempty"`
				LspSequenceNumber              int64  `json:"lsp_sequence_number,omitempty"`
			} `json:"lsp_header_data,omitempty"`
		} `json:"Content,omitempty"`
		Keys struct {
			InstanceName string `json:"instance_name,omitempty"`
			Level        int64  `json:"level,string,omitempty"`
			LspID        string `json:"lsp_id,omitempty"`
		} `json:"Keys,omitempty"`
		Timestamp int64 `json:"Timestamp,omitempty"`
	} `json:"Rows,omitempty"`
	Source    string `json:"Source,omitempty"`
	Telemetry struct {
		CollectionEndTime   int64  `json:"collection_end_time,omitempty"`
		CollectionID        int64  `json:"collection_id,omitempty"`
		CollectionStartTime int64  `json:"collection_start_time,omitempty"`
		EncodingPath        string `json:"encoding_path,omitempty"`
		MsgTimestamp        int64  `json:"msg_timestamp,omitempty"`
		NodeIDStr           string `json:"node_id_str,omitempty"`
		SubscriptionIDStr   string `json:"subscription_id_str,omitempty"`
	} `json:"Telemetry,omitempty"`
}
