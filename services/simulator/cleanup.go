package simulator

func (ths *Simulator) Cleanup() error {
	ths.log.Info("Cleaning up...")

	r := ths.storage.GetDB().Exec(`
		DELETE FROM ncc_fdb;
		DELETE FROM ncc_iface;
		DELETE FROM ncc_iface_state;
		DELETE FROM ncc_pinger_status;
		DELETE FROM ncc_trigger;
		DELETE FROM ncc_dhcp_lease;
		DELETE FROM ncc_session;
		DELETE FROM ncc_pon_onu;
		DELETE FROM ncc_pon_onu_state;
		DELETE FROM ncc_dhcp_binding;
		DELETE FROM ncc_device;
		DELETE FROM ncc_device_state;
		DELETE FROM ncc_hardware_model;
		DELETE FROM ncc_vendor;
		DELETE FROM ncc_payment;
		DELETE FROM ncc_issue_action;
		DELETE FROM ncc_issue;
		DELETE FROM ncc_customer_fee_log;
		DELETE FROM ncc_customer;
		DELETE FROM ncc_contract;
		DELETE FROM ncc_map_node;
		DELETE FROM ncc_city;
		DELETE FROM ncc_street;
		DELETE FROM ncc_customer_group;
		DELETE FROM ncc_payment_system_allowed_source_link;
		DELETE FROM ncc_payment_system_allowed_source;
		DELETE FROM ncc_payment_system;
		DELETE FROM ncc_payment_type;
		DELETE FROM ncc_service_internet_custom_data;
		DELETE FROM ncc_service_internet;
		DELETE FROM ncc_nas_attribute;
		DELETE FROM ncc_nas;
		DELETE FROM ncc_nas_type;
	`)
	if r.Error != nil {
		return r.Error
	}
	return nil
}
