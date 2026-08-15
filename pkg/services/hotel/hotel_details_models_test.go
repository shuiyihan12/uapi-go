package hotel

import (
	"encoding/xml"
	"testing"

	common55 "github.com/shuiyihan12/uapi-go/pkg/generated/common55"
	hotelxsd "github.com/shuiyihan12/uapi-go/pkg/generated/hotel"
)

// TestHotelDetailsRspDecode is a regression test: a real GDS message (taken
// from production HotelDetails logs) must fully deserialize into the
// WSDL-generated types (after the generator fix for cross-element ref
// resolution, the details body no longer degrades to interface{}/null), and
// the paging token NextResultReference must be populated correctly.
func TestHotelDetailsRspDecode(t *testing.T) {
	const payload = `<HotelDetailsRsp xmlns="http://www.travelport.com/schema/hotel_v55_0" xmlns:common_v55_0="http://www.travelport.com/schema/common_v55_0" TraceId="trace-abc" TransactionId="txn-1" ResponseTime="123">
  <common_v55_0:ResponseMessage Code="0" Type="Info">OK</common_v55_0:ResponseMessage>
  <common_v55_0:NextResultReference ProviderCode="1G">token-page-2</common_v55_0:NextResultReference>
  <RequestedHotelDetails>
    <HotelProperty HotelChain="WV" HotelCode="76034" HotelLocation="TSN" Name="WANDA VISTA TIANJIN">
      <PropertyAddress>
        <Address>No 486 No 8 Rd.of Da Zhi Gu</Address>
        <Address>Tianjin 12 CN 300170 </Address>
      </PropertyAddress>
    </HotelProperty>
    <HotelRateDetail RatePlanType="DLXRAC" Base="CNY754.72" Total="CNY880.00">
      <RoomRateDescription Name="Room"><Text>Grand Deluxe Twin Room</Text></RoomRateDescription>
    </HotelRateDetail>
  </RequestedHotelDetails>
</HotelDetailsRsp>`

	var rsp hotelxsd.HotelDetailsRsp
	if err := xml.Unmarshal([]byte(payload), &rsp); err != nil {
		t.Fatalf("unmarshal HotelDetailsRsp: %v", err)
	}

	if rsp.TraceId == nil || *rsp.TraceId != "trace-abc" {
		t.Fatalf("BaseRsp.TraceId wrong: %v", rsp.TraceId)
	}
	if rsp.NextResultReference == nil || string(rsp.NextResultReference.Value) != "token-page-2" {
		t.Fatalf("NextResultReference wrong: %+v", rsp.NextResultReference)
	}
	if rsp.NextResultReference.ProviderCode == nil || string(*rsp.NextResultReference.ProviderCode) != "1G" {
		t.Fatalf("NextResultReference.ProviderCode wrong: %+v", rsp.NextResultReference.ProviderCode)
	}
	if rsp.RequestedHotelDetails == nil {
		t.Fatalf("RequestedHotelDetails must not be nil")
	}
	if len(rsp.RequestedHotelDetails.HotelRateDetail) != 1 || string(rsp.RequestedHotelDetails.HotelRateDetail[0].RatePlanType) != "DLXRAC" {
		t.Fatalf("HotelRateDetail wrong: %+v", rsp.RequestedHotelDetails.HotelRateDetail)
	}
}

// TestRequestedHotelDetailsDecode verifies, with a full real message, that the
// generated types parse hotel static information and rate details.
func TestRequestedHotelDetailsDecode(t *testing.T) {
	const payload = `<RequestedHotelDetails xmlns="http://www.travelport.com/schema/hotel_v55_0" xmlns:common_v55_0="http://www.travelport.com/schema/common_v55_0">
  <HotelProperty HotelChain="WV" HotelCode="76034" HotelLocation="TSN" Name="WANDA VISTA TIANJIN">
    <PropertyAddress>
      <Address>No 486 No 8 Rd.of Da Zhi Gu</Address>
      <Address>Tianjin 12 CN 300170 </Address>
    </PropertyAddress>
    <common_v55_0:PhoneNumber Type="Business" Number="86-22-24626888"/>
    <common_v55_0:PhoneNumber Type="Fax" Number="86-22-2463 7112"/>
    <common_v55_0:Distance Units="KM" Value="11" Direction="W"/>
    <HotelRating RatingProvider="NTM"><Rating>4</Rating></HotelRating>
  </HotelProperty>
  <HotelDetailItem Name="CheckInTime"><Text>2PM</Text></HotelDetailItem>
  <HotelDetailItem Name="CheckOutTime"><Text>12N</Text></HotelDetailItem>
  <HotelRateDetail RatePlanType="DLXRAC" Base="CNY754.72" Total="CNY880.00" RateCategory="13" RateChangeIndicator="false" ExtraFeesIncluded="false">
    <RoomRateDescription Name="Room"><Text>Grand Deluxe Twin Room</Text><Text>Grand Deluxe Twin Bed Room</Text></RoomRateDescription>
    <RoomRateDescription Name="Rate"><Text>Best Available Rate</Text></RoomRateDescription>
    <HotelRateByDate EffectiveDate="2026-10-07" ExpireDate="2026-10-08" Base="CNY754.72"/>
    <Commission Indicator="true" Percent="10"/>
    <RateMatchIndicator Type="RoomCount" Status="Available" Value="1"/>
    <RateMatchIndicator Type="AdultCount" Status="Available" Value="1"/>
    <CancelInfo NonRefundableStayIndicator="false" CancelDeadline="2026-10-06T14:00:00.000+08:00" TaxInclusive="false" FeeInclusive="false"/>
    <GuaranteeInfo AbsoluteDeadline="2026-10-07T00:00:00.000+00:00" CredentialsRequired="false" GuaranteeType="Guarantee"><DepositAmount Amount="CNY0.00"/></GuaranteeInfo>
    <Inclusions SmokingRoomIndicator="unknown"><BedTypes Code="248"/><MealPlans Breakfast="false" Lunch="false" Dinner="false"><MealPlan Code="4"/></MealPlans><RoomView Code="20"/></Inclusions>
  </HotelRateDetail>
  <HotelType SourceLink="true"/>
</RequestedHotelDetails>`

	var details hotelxsd.RequestedHotelDetails
	if err := xml.Unmarshal([]byte(payload), &details); err != nil {
		t.Fatalf("unmarshal RequestedHotelDetails: %v", err)
	}

	hp := details.HotelProperty
	if string(hp.HotelChain) != "WV" || string(hp.HotelCode) != "76034" || hp.Name == nil || *hp.Name != "WANDA VISTA TIANJIN" {
		t.Fatalf("HotelProperty attrs wrong: %+v", hp)
	}
	if hp.PropertyAddress == nil || len(hp.PropertyAddress.Address) != 2 || hp.PropertyAddress.Address[0] != "No 486 No 8 Rd.of Da Zhi Gu" {
		t.Fatalf("PropertyAddress wrong: %+v", hp.PropertyAddress)
	}
	if len(hp.PhoneNumber) != 2 {
		t.Fatalf("PhoneNumber wrong: %+v", hp.PhoneNumber)
	}
	if hp.Distance == nil || hp.Distance.Value != 11 || hp.Distance.Direction == nil || *hp.Distance.Direction != "W" {
		t.Fatalf("Distance wrong: %+v", hp.Distance)
	}
	if len(hp.HotelRating) != 1 || hp.HotelRating[0].RatingProvider != "NTM" || len(hp.HotelRating[0].Rating) != 1 || hp.HotelRating[0].Rating[0] != 4 {
		t.Fatalf("HotelRating wrong: %+v", hp.HotelRating)
	}

	if len(details.HotelDetailItem) != 2 || details.HotelDetailItem[0].Name != "CheckInTime" || details.HotelDetailItem[0].Text[0] != "2PM" {
		t.Fatalf("HotelDetailItem wrong: %+v", details.HotelDetailItem)
	}

	if len(details.HotelRateDetail) != 1 {
		t.Fatalf("expected 1 HotelRateDetail, got %d", len(details.HotelRateDetail))
	}
	rd := details.HotelRateDetail[0]
	if string(rd.RatePlanType) != "DLXRAC" || rd.Base == nil || string(*rd.Base) != "CNY754.72" || rd.Total == nil || string(*rd.Total) != "CNY880.00" {
		t.Fatalf("HotelRateDetail attrs wrong: %+v", rd)
	}
	if len(rd.RoomRateDescription) != 2 || rd.RoomRateDescription[0].Name == nil || *rd.RoomRateDescription[0].Name != "Room" {
		t.Fatalf("RoomRateDescription wrong: %+v", rd.RoomRateDescription)
	}
	if rd.Commission == nil || rd.CancelInfo == nil || rd.GuaranteeInfo == nil || rd.Inclusions == nil {
		t.Fatalf("rate sub-structs must not be nil: %+v", rd)
	}
	if len(rd.RateMatchIndicator) != 2 || len(rd.HotelRateByDate) != 1 {
		t.Fatalf("RateMatchIndicator/HotelRateByDate wrong: %+v", rd)
	}

	if details.HotelType == nil {
		t.Fatalf("HotelType must not be nil")
	}
}

// TestHotelAlternatePropertiesDecode is a regression test: when the requested
// hotel is unbookable, the GDS returns a list of alternate hotels, and the
// generated types must parse the HotelProperty list correctly.
func TestHotelAlternatePropertiesDecode(t *testing.T) {
	const payload = `<HotelAlternateProperties xmlns="http://www.travelport.com/schema/hotel_v55_0">
  <HotelProperty HotelChain="WV" HotelCode="76034" HotelLocation="TSN" Name="WANDA VISTA TIANJIN"/>
  <HotelProperty HotelChain="SI" HotelCode="23397" Name="SHERATON"/>
</HotelAlternateProperties>`

	var alt hotelxsd.HotelAlternateProperties
	if err := xml.Unmarshal([]byte(payload), &alt); err != nil {
		t.Fatalf("unmarshal HotelAlternateProperties: %v", err)
	}
	if len(alt.HotelProperty) != 2 {
		t.Fatalf("expected 2 HotelProperty, got %d", len(alt.HotelProperty))
	}
	if string(alt.HotelProperty[0].HotelChain) != "WV" || string(alt.HotelProperty[1].HotelChain) != "SI" {
		t.Fatalf("HotelProperty attrs wrong: %+v", alt.HotelProperty)
	}
}

// TestDetailPagingHelpers covers the paging deduplication and merge helpers.
func TestDetailPagingHelpers(t *testing.T) {
	mkRsp := func(planTypes ...string) *hotelxsd.HotelDetailsRsp {
		rsp := &hotelxsd.HotelDetailsRsp{RequestedHotelDetails: &hotelxsd.RequestedHotelDetails{}}
		for _, pt := range planTypes {
			rsp.RequestedHotelDetails.HotelRateDetail = append(rsp.RequestedHotelDetails.HotelRateDetail,
				hotelxsd.HotelRateDetail{RatePlanType: common55.TypeRatePlanType(pt)})
		}
		return rsp
	}

	first := mkRsp("DLXRAC", "BAR")
	seen := collectRatePlanTypes(first)
	if len(seen) != 2 {
		t.Fatalf("collectRatePlanTypes expected 2, got %d", len(seen))
	}

	if !duplicatesRatePlanType(mkRsp("BAR", "NEW"), seen) {
		t.Fatalf("duplicatesRatePlanType should detect repeated BAR")
	}
	if duplicatesRatePlanType(mkRsp("NEW1", "NEW2"), seen) {
		t.Fatalf("duplicatesRatePlanType must not flag fresh plans")
	}

	mergeRateDetails(first, mkRsp("NEW"), seen)
	if len(first.RequestedHotelDetails.HotelRateDetail) != 3 {
		t.Fatalf("mergeRateDetails expected 3 rate details, got %d", len(first.RequestedHotelDetails.HotelRateDetail))
	}
	if _, ok := seen["NEW"]; !ok {
		t.Fatalf("mergeRateDetails must extend seen set")
	}

	provider := common55.TypeProviderCode("1G")
	if got := nextResultReferenceKey(&common55.NextResultReference{ProviderCode: &provider, Value: "tok"}); got != "1G|tok" {
		t.Fatalf("nextResultReferenceKey wrong: %q", got)
	}
	if got := nextResultReferenceKey(nil); got != "" {
		t.Fatalf("nextResultReferenceKey(nil) must be empty, got %q", got)
	}
}
