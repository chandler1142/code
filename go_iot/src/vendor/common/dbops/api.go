package dbops

func InsertMonitorRecord(record *MonitorRecord) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO monitor_record (type, value, ip, create_time, properties) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(record.Type, record.Value, record.IP, record.CreateTime, record.Properties)
	defer stmtIns.Close()
	if err != nil {
		return err
	}
	
	return nil

}
