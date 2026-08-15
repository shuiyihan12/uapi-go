package usecase

import (
	"strings"
	"testing"
	"time"

	hotelxsd "github.com/shuiyihan12/uapi-go/pkg/generated/hotel"
)

// TestValidateHotelStayRejectsPastCheckinDate verifies that a check-in date before today
// is rejected.
func TestValidateHotelStayRejectsPastCheckinDate(t *testing.T) {
	stay := hotelxsd.HotelStay{
		CheckinDate:  hotelxsd.TypeDate(time.Now().AddDate(0, 0, -1).Format("2006-01-02")),
		CheckoutDate: hotelxsd.TypeDate(time.Now().AddDate(0, 0, 1).Format("2006-01-02")),
	}

	err := validateHotelStay(stay)
	if err == nil {
		t.Fatalf("expected past checkin date to be rejected")
	}
	if !strings.Contains(err.Error(), "checkinDate: 不能早于今天") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestValidateHotelStayRejectsCheckoutBeforeCheckin verifies that a check-out date before
// the check-in date is rejected.
func TestValidateHotelStayRejectsCheckoutBeforeCheckin(t *testing.T) {
	stay := hotelxsd.HotelStay{
		CheckinDate:  hotelxsd.TypeDate(time.Now().AddDate(0, 0, 2).Format("2006-01-02")),
		CheckoutDate: hotelxsd.TypeDate(time.Now().AddDate(0, 0, 1).Format("2006-01-02")),
	}

	err := validateHotelStay(stay)
	if err == nil {
		t.Fatalf("expected checkout before checkin to be rejected")
	}
	if !strings.Contains(err.Error(), "checkoutDate: 必须晚于 checkinDate") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestValidateHotelStayAllowsTodayCheckin verifies that today is allowed as a check-in
// date.
func TestValidateHotelStayAllowsTodayCheckin(t *testing.T) {
	stay := hotelxsd.HotelStay{
		CheckinDate:  hotelxsd.TypeDate(time.Now().Format("2006-01-02")),
		CheckoutDate: hotelxsd.TypeDate(time.Now().AddDate(0, 0, 1).Format("2006-01-02")),
	}

	if err := validateHotelStay(stay); err != nil {
		t.Fatalf("expected today checkin to be allowed: %v", err)
	}
}

// TestNormalizeRateRuleDetailAcceptsAPIUppercase verifies that uppercase "COMPLETE" is
// normalized correctly.
func TestNormalizeRateRuleDetailAcceptsAPIUppercase(t *testing.T) {
	got, err := normalizeRateRuleDetail("COMPLETE")
	if err != nil {
		t.Fatalf("expected COMPLETE to be accepted: %v", err)
	}
	if got != "Complete" {
		t.Fatalf("unexpected normalized value: %s", got)
	}
}

// TestNormalizeRateRuleDetailRejectsInvalidValue verifies that an invalid rate rule value
// is rejected.
func TestNormalizeRateRuleDetailRejectsInvalidValue(t *testing.T) {
	if _, err := normalizeRateRuleDetail("FULL"); err == nil {
		t.Fatalf("expected invalid rateRuleDetail to be rejected")
	}
}
