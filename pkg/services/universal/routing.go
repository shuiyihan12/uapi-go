package universal

import "fmt"

// operationService follows the port addresses in
// wsdl/universal_v55_0/UniversalRecord.wsdl. The universal_v55_0 XML namespace
// does not imply UniversalRecordService: hotel/air/vehicle bookings and other
// operations in this namespace have their own physical service endpoints.
func operationService(operation string) (string, error) {
	switch operation {
	case "HotelCreateReservation", "HotelCancel":
		return "HotelService", nil
	case "AirCreateReservation", "AirCancel", "AirMerchandisingFulfillment":
		return "AirService", nil
	case "VehicleCreateReservation", "VehicleCancel":
		return "VehicleService", nil
	case "RailCreateReservation":
		return "RailService", nil
	case "PassiveCreateReservation", "PassiveCancel":
		return "PassiveService", nil
	case "UniversalRecordRetrieve", "UniversalRecordSearch", "UniversalRecordImport",
		"UniversalRecordModify", "UniversalRecordCancel", "ProviderReservationDivide":
		return "UniversalRecordService", nil
	case "AckScheduleChange":
		return "ScheduleChangeService", nil
	case "ProviderReservationDisplayDetails":
		return "ProviderReservationDisplayService", nil
	case "SavedTripCreate", "SavedTripModify", "SavedTripRetrieve", "SavedTripDelete", "SavedTripSearch":
		return operation + "Service", nil
	default:
		return "", fmt.Errorf("no WSDL service route for universal operation %q", operation)
	}
}
